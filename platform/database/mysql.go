package database

import (
	"fmt"
	"os"

	"sidewarslobby/app/queries"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DBQueries    *Queries
	DBConnection *gorm.DB
)

type Queries struct {
	*queries.UserQueries
}

func MysqlConnection() (*Queries, *gorm.DB, error) {
	mysqlConnURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(mysqlConnURL), &gorm.Config{})

	if err != nil {
		return nil, nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	return &Queries{
		UserQueries: &queries.UserQueries{DB: db},
	}, db, nil
}
