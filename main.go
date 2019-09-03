package main

import (
	"flag"

	"github.com/austinkline/go-migrate/db"
	_ "github.com/austinkline/go-migrate/migrations"
)

const (
	project = "test-migrations"
)

func main() {
	flag.Parse()

	err := db.SetupMultistatementDBWithEnv()
	if err != nil {
		print(err)
		return
	}

	err = db.DoMigrations(db.DBMultiStatement, project)
	if err != nil{
		print(err)
		return
	}
}

func init(){
	flag.Parse()
}
