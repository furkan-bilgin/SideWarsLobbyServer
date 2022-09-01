package main

import (
	"sidewarslobby/pkg/rest"
	"sidewarslobby/platform/cache"
	"sidewarslobby/platform/database"
)

func main() {
	rest := rest.Create()
	InitDatabases()

	rest.Listen(":3000")
}

func InitDatabases() {
	dbq, dbc, err := database.MysqlConnection()
	database.DBQueries, database.DBConnection = dbq, dbc
	if err != nil {
		panic(err)
	}

	rdc, err := cache.RedisConnection()
	cache.RedisClient = rdc
	if err != nil {
		panic(err)
	}
}
