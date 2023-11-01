package configs

import (
	"iam-service/configs/neo4j"
	"iam-service/configs/server"
)

type Config interface {
	Neo4j() neo4j.Config
	Server() server.Config
}

type config struct {
	neo4j  neo4j.Config
	server server.Config
}

func NewConfig() (Config, error) {
	return &config{
		neo4j:  neo4j.NewConfig(),
		server: server.NewConfig(),
	}, nil
}

func (c config) Neo4j() neo4j.Config {
	return c.neo4j
}

func (c config) Server() server.Config {
	return c.server
}
