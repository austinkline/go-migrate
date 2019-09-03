package migrations

import (
	"github.com/austinkline/go-migrate/db"
)

func init(){
	query := "CREATE SCHEMA IF NOT EXISTS test_migrations;"
	name := "0001-create-schema"

	m := db.Migration{
		Version: 1,
		Query: query,
		Name: name,
		Project: "test-migrations",
	}	

	db.RegisterMigration(m)
}