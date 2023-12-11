package service

import (
	"iam-service/model"
	"iam-service/vault"
	"context"
	"iam-service/client"
	"log"
	"strings"
	oort "github.com/c12s/oort/pkg/api"
)

type AuthService struct {
	repo model.UserRepo 
	v *vault.VaultClientService
}

// init
func NewAuthService(repo model.UserRepo, v *vault.VaultClientService) (*AuthService, error) {
	return &AuthService{
		repo: repo,
		v: v,
	}, nil
}

func (h AuthService) RegisterUser(ctx context.Context, req model.User) model.RegisterResp {
	refClient := *h.v
	registerResp := h.repo.CreateUser(ctx, req)

	if registerResp.Error == nil {
		err:= client.CreateOrgUserRelationship(registerResp.User.Org, registerResp.User.Username)
		if err != nil {
			log.Printf("Error while creating inheritance rel: %v", err)
			return model.RegisterResp{User: model.User{}, Error: err}
		}

		client.CreatePolicyAsync(registerResp.User.Org, 
								registerResp.User.Username, 
								getPermissionsForOort(registerResp.User.Permissions))
							
		
		refClient.RegisterUser(req.Username, req.Password, []string{"org.add"})
	}

	return registerResp
}

func (h AuthService) LoginUser(req model.LoginReq) model.LoginResp {
	refClient := *h.v
	token:= refClient.LoginUser(req)
	return model.LoginResp{Token: token, Error: nil}
}

func (h AuthService) Autorize(req model.AuthorizationReq) model.AuthorizationResp {
	return model.AuthorizationResp{Authorized: true}
}

func (h AuthService) VerifyToken(req model.Token) model.InternalToken {
	refClient := *h.v
	response := refClient.VerifyToken(req.Token)

	if !response.Verified {
		return model.InternalToken{Verified: response.Verified, Jwt: ""}
	}

	// proveriti da li ima nekih promena na oort-u
	permissions := client.GetGrantedPermissions(response.Username)

	// create jwt with permissions inside
	token, err := CreateToken(response.Username, transformPermissions(response.Username, permissions))
	if err != nil {
		return model.InternalToken{Verified: response.Verified, Jwt: ""}
	}

	return model.InternalToken{Verified: response.Verified, Jwt: token}
}

func (h AuthService) DecodeJwt(req model.Token) []string {
	return GetClaimsFromJwt(req)
}

func transformPermissions(username string, permissions []*oort.GrantedPermission) string {
	// format: perm_org, perm2_org, ...
	var transformed string

	if len(permissions) > 0 {
		for _, perm := range permissions {
			if !strings.Contains(perm.Object.Id, username) {
				transformed = transformed + perm.Name + "|" + perm.Object.Id + ","
			}
		}
		return transformed[:len(transformed)-1]
	}
	
	return transformed
}

func getPermissionsForOort(permissions []string) []*oort.Permission{
	var oortPermissions []*oort.Permission

	for _,perm := range permissions {
		oortPerm := &oort.Permission{
			Name:      perm,
			Kind:      oort.Permission_ALLOW,
			Condition: &oort.Condition{Expression: ""},
		}
		oortPermissions = append(oortPermissions, oortPerm)
	}

	return oortPermissions
}





		
		
	