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
