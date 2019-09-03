# go-migrate

Go package to help make migrations to a mysql database. This package doesn't
support up and down migrations it only applies a set of registered migrations in 
version order. 

If you want to use this package's DB setup, you can use `SetupMultistatementDBWithEnv` under the `go-migrate/db`
package. It uses environment variables to setup a connection to mysql with the following environment variables:

```
MIGRATE_DB_USER=[user]
MIGRATE_DB_HOST=localhost
MIGRATE_DB_PASSWORD=password
MIGRATE_DB_PORT=3306
```

The migrations are applied in the order of their associated versions. You can register 
migrations by adding Migration structs. One way to do this is to have a package for migrations, 
with one file per migration like below:


```
# main.go
package main

import (
	"flag"

	"github.com/austinkline/go-migrate/db"
	_ "github.com/austinkline/go-migrate/testmigrations"
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
```

```
# migrations/0001-create-schema.go
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
```

```
# migrations/0002-add-test-table.go
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
```