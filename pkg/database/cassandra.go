package database

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/jklaw90/shinfo/pkg/config"
)

func NewCassandra(ctx context.Context, cfg config.Provider) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cfg.GetStringSlice("cassandra.addresses")...)
	cluster.Keyspace = cfg.GetString("cassandra.keyspace")
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.GetString("cassandra.username"),
		Password: cfg.GetString("cassandra.password"),
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}
