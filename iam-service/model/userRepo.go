package model

import (
	"context"
)

type UserRepo interface {
	CreateUser(ctx context.Context, req User) RegisterResp
	LoginUser(ctx context.Context, req LoginReq) LoginResp
}




