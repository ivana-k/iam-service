package client

import (
	"context"
	//"iam-service/proto1"
	"fmt"
	"log"
	oort "github.com/c12s/oort/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AuthorizeUser(permission string, subjectId string) bool {
	conn, err := grpc.Dial("oort:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	evaluatorClient := oort.NewOortEvaluatorClient(conn)

	getResp, err := evaluatorClient.Authorize(context.Background(), &oort.AuthorizationReq{
			Subject:        &oort.Resource{
				Id:   subjectId,
				Kind: "user",
			},
			Object:         &oort.Resource{
				Id:   "idk",
				Kind: "user",
			},
			PermissionName: permission,
	}) 
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(getResp.Authorized)
	}

	return getResp.Authorized
}

func CreateOrgUserRelationship(org_id string, user_id string) error {
	conn, err := grpc.Dial("oort:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	administratorClient := oort.NewOortAdministratorClient(conn)

	_, err = administratorClient.CreateInheritanceRel(context.TODO(), &oort.CreateInheritanceRelReq{
		From: &oort.Resource{
			Id:   org_id,
			Kind: "user-group",
		},
		To:   &oort.Resource{
			Id:   user_id,
			Kind: "user",
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func GetGrantedPermissions(user string) []*oort.GrantedPermission {
	conn, err := grpc.Dial("oort:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	evaluatorClient := oort.NewOortEvaluatorClient(conn)

	resp, err := evaluatorClient.GetGrantedPermissions(context.TODO(), &oort.GetGrantedPermissionsReq{
		Subject: &oort.Resource{
			Id:   user,
			Kind: "user",
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("permissions of user")
	for _, perm := range resp.Permissions {
		log.Printf("%s - %s/%s", perm.Name, perm.Object.Kind, perm.Object.Id)
	}

	return resp.Permissions
}

/*func CreatePolicyAsync() {
	err := administratorAsync.SendRequest(&oort.CreatePolicyReq{
		SubjectScope: group,
		ObjectScope:  parentNamespace,
		Permission:   getConfigPerm,
	}, func(resp *oort.AdministrationAsyncResp) {
		if len(resp.Error) > 0 {
			log.Println(resp.Error)
		}
	})
	if err != nil {
		log.Fatalln(err)
	}
}*/
