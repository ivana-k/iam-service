package configs

import (
	"iam-service/configs/neo4j"
	"iam-service/configs/server"
	"iam-service/configs/nats"
)

type Config interface {
	Neo4j() neo4j.Config
	Server() server.Config
	Nats()	nats.Config
}

type config struct {
	neo4j  neo4j.Config
	server server.Config
	nats   nats.Config
}

func NewConfig() (Config, error) {
	return &config{
		neo4j:  neo4j.NewConfig(),
		server: server.NewConfig(),
		nats: 	nats.NewConfig(),
	}, nil
}

func (c config) Neo4j() neo4j.Config {
	return c.neo4j
}

func (c config) Server() server.Config {
	return c.server
}

func (c config) Nats() nats.Config {
	return c.nats
}
