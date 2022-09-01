package rest

import (
	"sidewarslobby/app/controllers"

	"github.com/gofiber/fiber/v2"
)

/*
Create a Rest API server and add routes to it.
*/
func Create() *fiber.App {
	app := fiber.New()
	serveStaticFiles(app)
	serveV1Api(app)

	return app
}

func serveStaticFiles(app *fiber.App) {
	app.Static("robots.txt", "files/robots.txt")
}

func serveV1Api(app *fiber.App) {
	v1 := app.Group("/api/v1")

	v1.Post("/firebase-auth", controllers.AuthViaFirebase)
}
