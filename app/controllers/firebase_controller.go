package controllers

import (
	"os"
	"sidewarslobby/platform/database"

	"github.com/gofiber/fiber/v2"

	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"google.golang.org/api/option"
)

var (
	FirebaseApp  *firebase.App
	FirebaseAuth *auth.Client
)

func InitFirebase() (*firebase.App, *auth.Client, error) {
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_ADMIN_JSON_PATH"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing app: %v", err)
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing app: %v", err)
	}

	return app, auth, err
}

// Validate Firebase token and authenticate the user with it
func AuthViaFirebase(c *fiber.Ctx) error {
	payload := struct {
		IdToken string
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	token, err := FirebaseAuth.VerifyIDTokenAndCheckRevoked(c.Context(), payload.IdToken)
	if err != nil {
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "Hatalı token",
		})
	}

	u, err := FirebaseAuth.GetUser(c.Context(), token.UID)
	if err != nil {
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "Kullanıcı bulunamadı",
		})
	}

	db_user := database.DBQueries.CreateOrUpdateUser(u)

	return c.JSON(fiber.Map{
		"token": db_user.Token,
	})
}
