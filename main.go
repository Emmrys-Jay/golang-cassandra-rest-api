package main

import (
	config "github.com/test/gocql/config"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func main() {
	Session = config.ConfigCassandraDB()

}
