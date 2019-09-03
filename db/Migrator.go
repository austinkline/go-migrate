package db

import (
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/golang/glog"
)

const (
	queryEnsureMigrationsSchema = "CREATE SCHEMA IF NOT EXISTS migrations;"
	queryEnsureMigrationsTable  = `
		CREATE TABLE IF NOT EXISTS migrations.applied_migrations (
			id int(11) NOT NULL,
			name varchar(90) NOT NULL,
			hash varchar(90) NOT NULL,
			project varchar(90) NOT NULL,
			PRIMARY KEY (id, project));`
	queryMaxID        = `SELECT MAX(id) FROM migrations.applied_migrations
						 WHERE project='%s';`
	queryAddMigration = "INSERT INTO `migrations`.`applied_migrations` (`id`, `name`, `hash`, `project`) VALUES ('%d', '%s', '%s', '%s');"
)

var (
	// the latest migration we've applied
	latestVersion int

	// map of version to it's migration func.
	migrations map[int]Migration
)

// Migration - The struct representing a migration
type Migration struct {
	Version int
	Name    string
	Query   string
	Project string
}

// RegisterMigration - Checks if the migration has been applied, and will queued if it hasnn't been applied yet.
func RegisterMigration(m Migration) {
	if migrations == nil {
		migrations = make(map[int]Migration)
	}

	migrations[m.Version] = m
	return
}

// getLatestAppliedVersion - Checks the db for the latest applied version of
func getLatestAppliedVersion(project string) (id int, err error) {
	maxIDQuery := fmt.Sprintf(queryMaxID, project)
	row := DBMultiStatement.QueryRow(maxIDQuery)
	row.Scan(&id)
	glog.Infof("processing latest version %d", id)
	return
}

// DoMigrations - performs all the migrations registered to the migrations map
func DoMigrations(db *sql.DB, project string) (err error) {
	err = ensureCreatedMigrations()
	if err != nil {
		glog.Fatal(err)
		return
	}

	id, err := getLatestAppliedVersion(project)
	if err != nil {
		return
	}

	nextVersion := id + 1
	for true {
		migration, found := migrations[nextVersion]
		if !found {
			glog.Infof("Done applying migrations, caught up to version %d\n", nextVersion-1)
			break
		}

		err = apply(db, migration)
		if err != nil {
			glog.Errorf("error while applying migration version %d\n%s\n", nextVersion, err)
			break
		}
		nextVersion++
	}

	return
}

func apply(db *sql.DB, m Migration) (err error) {
	glog.Infof("processing migration %s", m.Name)
	// generate a hash for the query
	h := sha1.New()
	h.Write([]byte(m.Query))
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))

	_, err = db.Exec(m.Query)
	if err != nil {
		return
	}

	updateVersionQuery := fmt.Sprintf(queryAddMigration, m.Version, m.Name, sha, m.Project)
	// update the last processed version to this one
	_, err = db.Exec(updateVersionQuery)
	if err != nil {
		glog.Infof("failed to apply migration version %d\n", m.Version)
	} else {
		glog.Infof("applied migration version %d\n", m.Version)
	}
	return
}

func ensureCreatedMigrations() (err error) {
	// ensure we have created the migrations schema and table
	_, err = DBMultiStatement.Exec(queryEnsureMigrationsSchema)
	if err != nil {
		glog.Fatal(err)
	}

	_, err = DBMultiStatement.Exec(queryEnsureMigrationsTable)
	if err != nil {
		glog.Fatal(err)
	}
	return
}