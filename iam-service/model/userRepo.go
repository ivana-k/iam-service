package model

/*import (
	"errors"
)*/

//In this file we localize all the operations on our database (currently in-memory database with productList)
//When we decIde to migrate this to a genuine solution we would have to make changes just to this file

type UserRepo interface {
	CreateUser(req User) RegisterResp
	LoginUser(req LoginReq) LoginResp
}




