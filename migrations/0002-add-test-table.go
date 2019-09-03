package migrations

import (
	"github.com/austinkline/go-migrate/db"
)

func init(){
	query := `CREATE TABLE test_migrations.new_table (
				idnew_table INT NOT NULL,
				new_tablecol VARCHAR(45) NULL,
				PRIMARY KEY (idnew_table));`
	name := "0002-add-test-table"

	m := db.Migration{
		Version: 2,
		Query: query,
		Name: name,
		Project: "test-migrations",
	}

	db.RegisterMigration(m)
}