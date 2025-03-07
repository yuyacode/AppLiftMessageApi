package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/yuyacode/AppLiftMessageApi/config"
)

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type Queryer interface {
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
}

var (
	_ Execer  = (*sqlx.DB)(nil)
	_ Queryer = (*sqlx.DB)(nil)
)

func New(ctx context.Context, cfg *config.Config, targetDB string) (*sqlx.DB, func(), error) {
	dbName, err := selectDB(cfg, targetDB)
	if err != nil {
		return nil, func() {}, err
	}
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Asia%%2FTokyo",
		cfg.DBUserName,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		dbName,
	))
	if err != nil {
		return nil, func() {}, err
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, func() { _ = db.Close() }, err
	}
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, func() { _ = db.Close() }, nil
}

func selectDB(cfg *config.Config, targetDB string) (string, error) {
	if targetDB == "company" {
		return cfg.DBCompany, nil
	} else if targetDB == "student" {
		return cfg.DBStudent, nil
	} else if targetDB == "common" {
		return cfg.DBCommon, nil
	} else {
		return "", fmt.Errorf("invalid database: %s", targetDB)
	}
}
