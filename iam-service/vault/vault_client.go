package vault

import (
	"context"
	"log"
	"time"
	"os"
	"io/ioutil"
	"encoding/json"
	vault "github.com/hashicorp/vault-client-go"
	schema "github.com/hashicorp/vault-client-go/schema"
	"iam-service/model"
)

type VaultClientService struct {
	client *vault.Client 
}

type VaultKey struct {
	RootKey string `json:"root_key"`
}

// init
func NewVaultClientService() (*VaultClientService, error) {
	return &VaultClientService{
		client: initClient(),
	}, nil
}

func initClient() *vault.Client {
	client, err := vault.New(
		vault.WithAddress("http://vault:8200"),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		log.Println(err)
	}

	// check if its initialized
	initResp, err := client.System.ReadInitializationStatus(
		context.Background(),
	)
	if err != nil {
		log.Println(err)
	}

	log.Println(initResp)
	initStatus := initResp.Data["initialized"].(bool)

	if initStatus {
		log.Println("Vault already initialized.")
		vaultKey := loadKeyFromJson()
		log.Println(vaultKey.RootKey)
		if err := client.SetToken(vaultKey.RootKey); err != nil {
			log.Printf("Error while trying to set vault token: %v" , err)
		}

		Unseal(client, "firstKey")

		return client
	}

	// init
	initializedVault := Initialize(client)

	log.Println("root token:")
	log.Println(initializedVault.rootKey)
	var vaultKey VaultKey
	vaultKey.RootKey = initializedVault.rootKey
	saveKeyToJson(vaultKey)

	// auth
	if err := client.SetToken(initializedVault.rootKey); err != nil {
		log.Fatal(err)
	}

	Unseal(client, initializedVault.keysArray[0].(string))
	MountSecretEngine(client)

	return client
}

func (v VaultClientService) RegisterUser(username string, password string, policies []string) {
	resp, err := v.client.Auth.UserpassWriteUser(
		context.Background(),
		username,
		schema.UserpassWriteUserRequest{
			Password: password,
			Policies: policies,
			TokenPeriod: "0.5h",
		},
		vault.WithMountPath("userpass"),
	)
	if err != nil {
		log.Println("vault registration failed")
		log.Printf("Error: %v", err)
	} else {
		log.Println("vault registration finished")
		log.Println(resp)
	}

	
}

func (v VaultClientService) LoginUser(req model.LoginReq) string {
	resp, err := v.client.Auth.UserpassLogin(
		context.Background(),
		req.Username,
		schema.UserpassLoginRequest{
			Password: req.Password,
		},
		vault.WithMountPath("userpass"),
	)
	if err != nil {
		log.Printf("VaultLogin error: %v", err)
		return ""
	}

	log.Println("vault login finished")
	log.Println(resp)

	return resp.Auth.ClientToken
}

func (v VaultClientService) VerifyToken(token string) model.VerificationResp {
	resp, err := v.client.Auth.TokenLookUp(
		context.Background(),
		schema.TokenLookUpRequest{
			Token: token,
		},
	)

	if err != nil {
		log.Printf("%v", err)
		return model.VerificationResp{Verified: false, Username: ""}
	}

	expTime:= resp.Data["expire_time"].(string)
	metaMap:= resp.Data["meta"].(map[string]interface{})
	username:=metaMap["username"].(string)
	timestamp, err := time.Parse(time.RFC3339Nano, expTime)

	if err != nil {
		log.Printf("Error parsing timestamp: %v", err)
		return model.VerificationResp{Verified: false, Username: username}
	}

	currentTime := time.Now()
	isBefore := timestamp.Before(currentTime)

	return model.VerificationResp{Verified: !isBefore, Username: username}
}

func Initialize(client *vault.Client) VaultClient {
	resp, err := client.System.Initialize(
		context.Background(),
		schema.InitializeRequest{
			PgpKeys: nil,
			RootTokenPgpKey: "",
			SecretShares: 1,
			SecretThreshold: 1,
		},
	)
	if err != nil {
		log.Printf("Vault failed to initialize %v", err)
	}

	keysArray, ok := resp.Data["keys"].([]interface{})
	if !ok || len(keysArray) == 0 {
		log.Println("Error: Unable to access the 'keys' array")
		return VaultClient{}
	}

	return VaultClient{keysArray: keysArray, rootKey: resp.Data["root_token"].(string)}
}

func Unseal(client *vault.Client, firstKey string) {
	_, err := client.System.Unseal(
		context.Background(),
		schema.UnsealRequest{
			Key: firstKey,		// first key in array
		},
	)
	if err != nil {
		log.Printf("Vault failed to unseal: %v", err)
	}
}

func MountSecretEngine(client *vault.Client) {
	_, err := client.System.AuthEnableMethod(
		context.Background(),
		"userpass",
		schema.AuthEnableMethodRequest{
			Description: "Mount for user identity",
			Type: "userpass",
		},
	)
	if err != nil {
		log.Printf("Vault failed to mount secret engine %v", err)
	}
}


func loadKeyFromJson() VaultKey {
	path := "api_key.json"

	jsonFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("%s", err)
	}

	var vaultKey VaultKey
	err = json.Unmarshal(jsonFile, &vaultKey)
	if err != nil {
		log.Println("Error:", err)
	}

	return vaultKey
}

func saveKeyToJson(vaultKey VaultKey) {
	path := "api_key.json"

	updatedJSON, err := json.MarshalIndent(vaultKey, "", "  ")
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	// Write the updated JSON back to the file
	err = ioutil.WriteFile(path, updatedJSON, 0644)
	if err != nil {
		log.Println("Error writing JSON file:", err)
		return
	}

}