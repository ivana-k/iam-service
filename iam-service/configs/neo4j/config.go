package neo4j

import (
	"fmt"
	"os"
)

type Config interface {
	Uri() string
	Username() string
	Password() string
	DbName() string
}

type config struct {
	hostname string
	port     string
	username string
	password string
	dbName   string
}

func NewConfig() Config {
	return config{
		hostname: os.Getenv("NEO4J_HOSTNAME"),
		port:     os.Getenv("NEO4J_BOLT_PORT"),
		username: os.Getenv("NEO4J_USERNAME"),
		password: os.Getenv("NEO4J_PASSWORD"),
		dbName:   os.Getenv("NEO4J_DBNAME"),
	}
}

func (c config) Uri() string {
	return fmt.Sprintf("bolt://%s:%s", c.hostname, c.port)
}

func (c config) Username() string {
	return c.username
}

func (c config) Password() string {
	return c.password
}

func (c config) DbName() string {
	return c.dbName
}
