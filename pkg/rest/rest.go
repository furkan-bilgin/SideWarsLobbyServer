package rest

import (
	"sidewarslobby/app/controllers"
	"sidewarslobby/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

/*
Create a Rest API server and add routes to it.
*/
func Create() *fiber.App {
	app := fiber.New()

	middleware.FiberMiddleware(app)

	serveStaticFiles(app)
	serveV1Api(app)

	return app
}

func serveStaticFiles(app *fiber.App) {
	app.Static("/robots.txt", "static/robots.txt")
}

func serveV1Api(app *fiber.App) {
	v1 := app.Group("/api/v1")

	v1.Post("/auth/firebase", controllers.AuthViaFirebase)
	v1.Post("/server/confirm-user-match", controllers.ConfirmUserMatch)
}
