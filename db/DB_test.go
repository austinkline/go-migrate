package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	err := SetupMultistatementDBWithEnv()
	assert.Nil(t, err)

	err = DBMultiStatement.Ping()
	assert.Nil(t, err)
}
