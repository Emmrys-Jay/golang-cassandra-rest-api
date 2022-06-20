package config

import (
	"fmt"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func ConfigCassandraDB() {
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "proposals_and_comments"
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("cassandra well initialized")
}
