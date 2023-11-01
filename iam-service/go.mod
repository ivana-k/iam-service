module iam-service

go 1.21.3

require (
	github.com/c12s/oort v0.0.0
	github.com/neo4j/neo4j-go-driver/v4 v4.4.1
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/c12s/magnetar v1.0.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/nats-io/nats.go v1.28.0 // indirect
	github.com/nats-io/nkeys v0.4.4 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.9.0 // indirect
)

require (
	//github.com/c12s/oort v0.0.0-20231024081602-18cfc5f3c7c1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
)

replace github.com/c12s/oort => ../oort

replace github.com/c12s/magnetar => ../magnetar
