package service

import (
	"iam-service/model"
)

type AuthService struct {
	repo model.UserRepo 
}

// init
func NewAuthService(repo model.UserRepo) (*AuthService, error) {
	return &AuthService{
		repo: repo,
	}, nil
}

func (h AuthService) RegisterUser(req model.User) model.RegisterResp {
	return h.repo.CreateUser(req)
}

func (h AuthService) LoginUser(req model.LoginReq) model.LoginResp {
	return h.repo.LoginUser(req)
}

func (h AuthService) Autorize(req model.AuthorizationReq) model.AuthorizationResp {
	return model.AuthorizationResp{Authorized: true}
}



		
		
	