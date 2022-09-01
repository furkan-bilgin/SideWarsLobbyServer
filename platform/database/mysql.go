package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"sidewarslobby/app/queries"

	_ "github.com/go-sql-driver/mysql" // load driver for Mysql
	"github.com/jmoiron/sqlx"
)

var (
	DBQueries    *Queries
	DBConnection *sqlx.DB
)

type Queries struct {
	*queries.UserQueries
}

func MysqlConnection() (*Queries, *sqlx.DB, error) {
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	// Build Mysql connection URL.
	mysqlConnURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Define database connection for Mysql.
	db, err := sqlx.Connect("mysql", mysqlConnURL)

	if err != nil {
		return nil, nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// Set database connection settings:
	// 	- SetMaxOpenConns: the default is 0 (unlimited)
	// 	- SetMaxIdleConns: defaultMaxIdleConns = 2
	// 	- SetConnMaxLifetime: 0, connections are reused forever
	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close() // close database connection
		return nil, nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return &Queries{
		UserQueries: &queries.UserQueries{DB: db},
	}, db, nil
}
