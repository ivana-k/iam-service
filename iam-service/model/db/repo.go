package db

import (
	"iam-service/model"
	"iam-service/service"
	"iam-service/client"
	"log"
	"fmt"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type UserRepo struct {
	manager *TransactionManager
	factory CypherFactory
}

func NewUserRepo(manager *TransactionManager, factory CypherFactory) model.UserRepo {
	return UserRepo{
		manager: manager,
		factory: factory,
	}
}

func (store UserRepo) CreateUser(req model.User) model.RegisterResp {
	cypher, params := store.factory.createResource(req)

	fmt.Printf("Cypher Query Parameters:\n")
        for key, value := range params {
            fmt.Printf("%s: %v\n", key, value)
        }

	err := store.manager.WriteTransaction(cypher, params)

	if err != nil {
		fmt.Printf("ipak vraca gresku")
		// test user da se proveri proto1 mapiranje
		return model.RegisterResp{User: model.User{Email: "nostrud sed velit reprehenderit", Id: 11, Name: "velit irure culpa ex", Password: "dolore ad incididunt ut eu", Surname: "ullamco culpa nostrud mollit"}, Error: err}
	}

	return model.RegisterResp{User: req, Error: nil}		// plus id

}

func (store UserRepo) LoginUser(req model.LoginReq) model.LoginResp {
	cypher, params := store.factory.findUser(req)

	result, err := store.manager.ReadTransaction(cypher, params)
	if err != nil {
		return model.LoginResp{Token: "", Error: err}
	}

	records, ok := result.([]*neo4j.Record)
	log.Println(len(records))
	if !ok {
		return model.LoginResp{Token: "", Error: errors.New("invalid resp format")}
	}


	for _, record := range records {
		userProps, found := record.Get("user")
		fmt.Printf("Record userprops: %+v\n", userProps.(neo4j.Node).Props)
		if !found {
			fmt.Println("User not found in record")
			return model.LoginResp{Token: "", Error: errors.New("User not found in record")}
		}

		if userNode, ok := userProps.(neo4j.Node); ok {
			fmt.Printf("Record usernode: %+v\n", userNode)
			userMap := userNode.Props
			name := userMap["name"].(string)
			permission := userMap["permission"].(string)
			isAuthorized := client.AuthorizeUser(permission, name)
			if isAuthorized {
				token, _ := service.CreateToken(name, "ALLOW", userMap["permission"].(string))
				return model.LoginResp{Token: token, Error: nil}
			} 
				
		} else {
			fmt.Println("invalid mapping")
		}
		
	}

	return model.LoginResp{Token: "", Error: errors.New("Invalid mapping")}		

}

/*func (store RHABACRepo) GetResource(req domain.GetResourceReq) domain.GetResourceResp {
	cypher, params := store.factory.getResource(req)
	records, err := store.manager.ReadTransaction(cypher, params)
	if err != nil {
		return domain.GetResourceResp{Resource: nil, Error: err}
	}

	recordList, ok := records.([]*neo4j.Record)
	if !ok {
		return domain.GetResourceResp{Error: errors.New("invalid resp format")}
	}
	if len(recordList) == 0 {
		return domain.GetResourceResp{Error: errors.New("resource not found")}
	}
	return domain.GetResourceResp{Resource: getResource(records), Error: nil}
}*/