package main

import (
	"sidewarslobby/pkg/rest"
	"sidewarslobby/platform/database"
)

func main() {
	rest := rest.Create()

	dbq, dbc, err := database.MysqlConnection()
	database.DBQueries, database.DBConnection = dbq, dbc

	if err != nil {
		panic(err)
	}

	rest.Listen(":3000")
}
