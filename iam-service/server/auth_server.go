package server

import (
	"context"
	"iam-service/proto1"
	"iam-service/service"
	"log"
)

type AuthServiceServer struct {
	service service.AuthService
	proto1.UnimplementedAuthServiceServer	
}

func NewAuthServiceServer(service service.AuthService) (proto1.AuthServiceServer, error) {
	return &AuthServiceServer{
		service: service,
	}, nil
}

// videti kako poslati gRPC 
func (o *AuthServiceServer) Authorize(ctx context.Context, req *proto1.AuthorizationReq) (*proto1.AuthorizationResp, error) {
	/*reqDomain, err := proto.AuthorizationReqToDomain(req)
	if err != nil {
		return nil, err
	}
	resp := o.service.Authorize(*reqDomain)
	log.Println(resp)
	return &proto.AuthorizationResp{Authorized: resp.Authorized}, resp.Error*/
	return &proto1.AuthorizationResp{Authorized: true}, nil
}

func (o *AuthServiceServer) RegisterUser(ctx context.Context, req *proto1.User) (*proto1.RegisterResp, error) {
	user, err := proto1.UserToModel(req)
	if err != nil {
		return nil, err
	}
	resp := o.service.RegisterUser(*user)
	log.Println(resp.User)
	return &proto1.RegisterResp{User: &proto1.User{Id: int64(resp.User.Id), 
		Name: resp.User.Name,
		Surname: resp.User.Surname,
		Password: resp.User.Password,
		Email: resp.User.Email}}, resp.Error
}

func (o *AuthServiceServer) LoginUser(ctx context.Context, req *proto1.LoginReq) (*proto1.LoginResp, error) {
	user, err := proto1.LoginToModel(req)

	if err != nil {
		return nil, err
	}
	
	resp := o.service.LoginUser(*user)
	log.Println(resp)
	return &proto1.LoginResp{Token: resp.Token}, nil
}