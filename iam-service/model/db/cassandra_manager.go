package db

import (
	"github.com/gocql/gocql"
	"log"
	"iam-service/model"
	"context"
)

type CassandraManager struct {
	session *gocql.Session
}

func NewCassandraManager() *CassandraManager {
	return &CassandraManager{
		session: Connect(),
	}
}

func Connect() *gocql.Session {
	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "apollo"
	
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	//defer session.Close()

	return session
}

func (cm CassandraManager) InitDb() {
	err := cm.session.Query("CREATE KEYSPACE IF NOT EXISTS apollo WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };").Exec()
	if err != nil {
		log.Println(err)
	}

	err = cm.session.Query("USE apollo;").Exec()
	if err != nil {
		log.Println(err)
	}

	err = cm.session.Query("CREATE TABLE IF NOT EXISTS org (id UUID PRIMARY KEY, name TEXT );").Exec()
	if err != nil {
		log.Println(err)
	}

	err = cm.session.Query("CREATE INDEX ON org (name);").Exec()
	if err != nil {
		log.Println(err)
	}

	err = cm.session.Query("CREATE TABLE IF NOT EXISTS permission (id UUID PRIMARY KEY, name TEXT);").Exec()
	if err != nil {
		log.Println(err)
	}

	err = cm.session.Query("CREATE TABLE IF NOT EXISTS user (id UUID PRIMARY KEY, name TEXT, surname TEXT, email TEXT, username TEXT, password TEXT, created_at DATE, updated_at DATE);").Exec()
	if err != nil {
		log.Println(err)
	}

	err = cm.session.Query("CREATE TABLE IF NOT EXISTS org_user (org_id UUID, user_id UUID, permissions SET<TEXT>, PRIMARY KEY (org_id, user_id));").Exec()
	if err != nil {
		log.Println(err)
	}
}

func (cm CassandraManager) SeedDb() {
	err := cm.session.Query("INSERT INTO permission (id, name) VALUES (uuid(), 'config.get') IF NOT EXISTS;").Exec()
	if err != nil {
		log.Println(err)
	}

	err = cm.session.Query("INSERT INTO permission (id, name) VALUES (uuid(), 'config.put') IF NOT EXISTS;").Exec()
	if err != nil {
		log.Println(err)
	}

	err = cm.session.Query("INSERT INTO permission (id, name) VALUES (uuid(), 'namespace.putconfig') IF NOT EXISTS;").Exec()
	if err != nil {
		log.Println(err)
	}

	err = cm.session.Query("INSERT INTO permission (id, name) VALUES (uuid(), 'node.get') IF NOT EXISTS;").Exec()
	if err != nil {
		log.Println(err)
	}

	err = cm.session.Query("INSERT INTO permission (id, name) VALUES (uuid(), 'node.put') IF NOT EXISTS;").Exec()
	if err != nil {
		log.Println(err)
	}

	err = cm.session.Query("INSERT INTO permission (id, name) VALUES (uuid(), 'node.label.put') IF NOT EXISTS;").Exec()
	if err != nil {
		log.Println(err)
	}

	err = cm.session.Query("INSERT INTO permission (id, name) VALUES (uuid(), 'node.label.get') IF NOT EXISTS;").Exec()
	if err != nil {
		log.Println(err)
	}

	err = cm.session.Query("INSERT INTO permission (id, name) VALUES (uuid(), 'node.label.delete') IF NOT EXISTS;").Exec()
	if err != nil {
		log.Println(err)
	}
}

const insertUserQuery = `
INSERT INTO user (id, name, surname, email, password, username)
VALUES (?, ?, ?, ?, ?, ?)`
func (cm CassandraManager) InsertUser(ctx context.Context, user model.User) (string, error) {
	id:=gocql.UUID{}
	query := cm.session.Query(insertUserQuery,
        id, 
		user.Name, 
		user.Surname, 
		user.Email,
		user.Password,
		user.Username, 
		)

    if err := query.Exec(); err != nil {
        return "", err
    }

	return id.String(), nil
}

const findOrgQuery = `SELECT id, name FROM org WHERE name = ?`
func (cm CassandraManager) FindOrgByName(ctx context.Context, orgName string) (model.Org, error) {
	query := cm.session.Query(findOrgQuery,
        orgName, 
		)

	var id gocql.UUID
	var name string
    if err := query.WithContext(ctx).Consistency(gocql.One).Scan(&id, &name); err != nil {
		log.Printf("FindOrgByName error")
		log.Printf("Error: %v", err)
        return model.Org{}, err
    }

	return model.Org{Id: id.String(), Name: name}, nil
}

const findAllPermQuery = `SELECT id, name FROM permission`
func (cm CassandraManager) GetAllPermissions() ([]string, error) {
	query := cm.session.Query(findAllPermQuery)

	var id, name string
	var permissions []string
	iter := query.Iter()
	
	for iter.Scan(&id, &name) {
		log.Printf("ID: %s, Name: %s", id, name)
		permissions = append(permissions, name)
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	return permissions, nil
}

const findUserPermQuery = `SELECT permissions FROM org_user WHERE org_id = ? AND user_id = ?;`
func (cm CassandraManager) GetUserPermissions(org_id string, user_id string) ([]string, error) {
	query := cm.session.Query(findUserPermQuery, org_id, user_id)

	var permissions map[string]struct{}
	var foundPermissions []string
	
	if err := query.Scan(&permissions); err != nil {
		log.Fatal(err)
	}

	log.Print("Permissions: ")
	for permission := range permissions {
		log.Printf("%s ", permission)
		foundPermissions = append(foundPermissions, permission)
	}

	return foundPermissions, nil
}

const insertOrgQuery = `
INSERT INTO org (id, name)
VALUES (?, ?)`
func (cm CassandraManager) InsertOrg(ctx context.Context, org model.Org) (model.Org, error) {
	orgUuid:=gocql.UUID{}
	query := cm.session.Query(insertOrgQuery,
        orgUuid, 
		org.Name, 
		)

    if err := query.Exec(); err != nil {
		log.Printf("InsertOrg error")
		log.Fatal(err)
        return model.Org{}, err
    }

	return model.Org{Id: orgUuid.String(), Name: org.Name}, nil
}

const createOrgUserQuery = `
INSERT INTO org_user (org_id, user_id, permissions)
VALUES (?, ?, ?)`
func (cm CassandraManager) CreateOrgUser(org_uuid string, user_uuid string) (bool, error) {
	permissions, err:= cm.GetAllPermissions()
	log.Printf("Array: %v", permissions)

	if err != nil {
		log.Fatal("Cannot find permissions")
	}

	mapPermissions := arrayToSet(permissions)
	query := cm.session.Query(createOrgUserQuery,
        org_uuid, 
		user_uuid, 
		mapPermissions,
		)

    if err := query.Exec(); err != nil {
		log.Println(err)
        return false, err
    }

	return true, nil
}

func arrayToSet(arr []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, value := range arr {
		set[value] = struct{}{}
	}
	return set
}