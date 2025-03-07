package store

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

type mockError struct {
	msg string
}

func (m *mockError) Error() string {
	return m.msg
}

func assertAnError() error {
	return &mockError{"some db error"}
}

func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Close()
	})
	sqlxDB := sqlx.NewDb(db, "mysql")
	return sqlxDB, mock
}
