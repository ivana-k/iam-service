-- Create keyspace
CREATE KEYSPACE IF NOT EXISTS apollo
WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

-- Use keyspace
USE apollo;

-- Create table
CREATE TABLE IF NOT EXISTS org (
    id UUID PRIMARY KEY,
    name TEXT
);
CREATE INDEX ON org (name);

CREATE TABLE IF NOT EXISTS permission (
    id UUID PRIMARY KEY,
    name TEXT
);

CREATE TABLE IF NOT EXISTS user (
    id UUID PRIMARY KEY,
    name TEXT,
    surname TEXT,
    email TEXT,
    username TEXT,
    password TEXT,
    created_at DATE,
    updated_at DATE
);

CREATE TABLE IF NOT EXISTS org_user (
    org_id UUID,
    user_id UUID,
    permissions SET<TEXT>, 
    PRIMARY KEY (org_id, user_id)
);

INSERT INTO permission (id, name) VALUES (uuid(), 'config.get') IF NOT EXISTS;
INSERT INTO permission (id, name) VALUES (uuid(), 'config.put') IF NOT EXISTS;
INSERT INTO permission (id, name) VALUES (uuid(), 'namespace.putconfig') IF NOT EXISTS;
INSERT INTO permission (id, name) VALUES (uuid(), 'node.get') IF NOT EXISTS;
INSERT INTO permission (id, name) VALUES (uuid(), 'node.put') IF NOT EXISTS;
INSERT INTO permission (id, name) VALUES (uuid(), 'node.label.put') IF NOT EXISTS;
INSERT INTO permission (id, name) VALUES (uuid(), 'node.label.get') IF NOT EXISTS;
INSERT INTO permission (id, name) VALUES (uuid(), 'node.label.delete') IF NOT EXISTS;