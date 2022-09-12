package main

import (
	"sidewarslobby/pkg/rest"
	"sidewarslobby/pkg/websocket"
	"sidewarslobby/platform/cache"
	"sidewarslobby/platform/database"

	"sidewarslobby/app/controllers"
)

func main() {
	println("Creating rest")
	rest := rest.Create()
	println("Init database")
	InitDatabases()
	println("Init firebase")
	InitFirebase()
	println("Create websocket")
	websocket.Create(rest)
	
	println("Done")
	rest.Listen(":3000")
}

// Initialize MySQL and Redis
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

	controllers.InitRedisController()
}

func InitFirebase() {
	fb, au, err := controllers.InitFirebase()
	if err != nil {
		panic(err)
	}

	controllers.FirebaseApp, controllers.FirebaseAuth = fb, au
}
