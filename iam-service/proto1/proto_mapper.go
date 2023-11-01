package proto1

import (
	"iam-service/model"
)

func UserToModel(req *User) (*model.User, error) {
	return &model.User{
		Name:       		 req.Name,
		Surname:         	 req.Surname,
		Email: 				 req.Email,
		Password:            req.Password,
		Org:            	 req.Org,
		Permission:          req.Permission,
	}, nil
}

func LoginToModel(req *LoginReq) (*model.LoginReq, error) {
	return &model.LoginReq{
		Email: 				 req.Email,
		Password:            req.Password,
	}, nil
}
