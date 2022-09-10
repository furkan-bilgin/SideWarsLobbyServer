package database

import (
	"fmt"
	"os"

	"sidewarslobby/app/models"
	"sidewarslobby/app/queries"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DBQueries    *Queries
	DBConnection *gorm.DB
)

type Queries struct {
	*queries.UserQueries
	*queries.MatchQueries
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

	mysqlConnURL += "?parseTime=true"

	db, err := gorm.Open(mysql.Open(mysqlConnURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	AutoMigrateDatabase(db)

	return &Queries{
		UserQueries:  &queries.UserQueries{DB: db},
		MatchQueries: &queries.MatchQueries{DB: db},
	}, db, nil
}

func AutoMigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&models.Match{}, &models.User{}, &models.UserMatch{}, &models.UserInfo{})
	if err != nil {
		panic(err)
	}
}
