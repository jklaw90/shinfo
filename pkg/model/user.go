package model

import "github.com/gocql/gocql"

type User struct {
	ID     gocql.UUID `cql:"id"`
	Name   string     `cql:"name,omitempty"`
	Avatar string     `cql:"avatar,omitempty"`
}
