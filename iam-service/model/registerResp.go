package model

type RegisterResp struct {
	User User
	Error error
}

type LoginReq struct {
	Email string
	Password string
}

type LoginResp struct {
	Token string	//
	Error error
}


type AuthorizationReq struct {
	Subject,
	Object Resource
	PermissionName string
	Env            []Attribute
}

type AuthorizationResp struct {
	Authorized bool
	Error      error
}