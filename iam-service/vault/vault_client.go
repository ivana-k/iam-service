package vault

import (
	"context"
	"log"
	"time"
	"os"
	vault "github.com/hashicorp/vault-client-go"
	schema "github.com/hashicorp/vault-client-go/schema"
	"iam-service/model"
)

type VaultClientService struct {
	client *vault.Client 
}

// init
func NewVaultClientService() (*VaultClientService, error) {
	return &VaultClientService{
		client: initClient(),
	}, nil
}

func initClient() *vault.Client {
	// prepare a client with the given base address
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

	initStatus := initResp.Data["initialized"].(bool)

	if initStatus {
		log.Println("Already init")
		log.Println(os.Getenv("VAULT_DEV_ROOT_TOKEN_ID"))
		rootToken:=os.Getenv("VAULT_DEV_ROOT_TOKEN_ID")
		if err := client.SetToken(rootToken); err != nil {
			log.Println(err)
		}

		respAuth, err := client.System.AuthEnableMethod(
			context.Background(),
			"userpass",
			schema.AuthEnableMethodRequest{
				Description: "Mount for user identity",
				Type: "userpass",
	
			},
		)
		if err != nil {
			log.Println(err)
		}
	
		log.Println(respAuth)

		return client
	}

	// init
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
		log.Println(err)
	}

	keysArray, ok := resp.Data["keys"].([]interface{})
	if !ok || len(keysArray) == 0 {
		log.Println("Error: Unable to access the 'keys' array")
		return nil
	}

	firstKey, ok := keysArray[0].(string)
	rootToken := resp.Data["root_token"].(string)
	
	log.Println("root token:")
	log.Println(rootToken)
	os.Setenv("VAULT_DEV_ROOT_TOKEN_ID", rootToken)

	// auth
	if err := client.SetToken(rootToken); err != nil {
		log.Fatal(err)
	}

	if ok {
		respUseal, err := client.System.Unseal(
			context.Background(),
			schema.UnsealRequest{
				Key: firstKey,		// first key in array
			},
		)
		if err != nil {
			log.Println(err)
		}

		log.Println(respUseal.Data)
	}

	// mount new secret engine
	respAuth, err := client.System.AuthEnableMethod(
		context.Background(),
		"userpass",
		schema.AuthEnableMethodRequest{
			Description: "Mount for user identity",
			Type: "userpass",

		},
	)
	if err != nil {
		log.Println(err)
	}

	log.Println(respAuth)

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
	}

	log.Println("vault registration finished")
	log.Println(resp)
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
		log.Printf("%v", err)
	}

	log.Println("vault login finished")
	log.Println(resp)

	authToken := resp.Auth.ClientToken
	return authToken
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