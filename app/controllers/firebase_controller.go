package controllers

import (
	"os"
	"sidewarslobby/pkg/utils"
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
		FirebaseToken string
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	token, err := FirebaseAuth.VerifyIDTokenAndCheckRevoked(c.Context(), payload.FirebaseToken)
	if err != nil {
		return utils.RESTError(c, "Hesap doğrulanamadı")
	}

	u, err := FirebaseAuth.GetUser(c.Context(), token.UID)
	if err != nil {
		return utils.RESTError(c, "Kullanıcı bulunamadı")
	}

	db_user := database.DBQueries.CreateOrUpdateUser(u)

	return c.JSON(fiber.Map{
		"Token": db_user.Token,
	})
}
