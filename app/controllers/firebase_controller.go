package controllers

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"context"
	"fmt"

	firebase "firebase.google.com/go"
	//"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func FirebaseApp() (*firebase.App, error) {
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_ADMIN_JSON_PATH"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	return app, err
}

// Validate Firebase token and authenticate the user with it
func AuthViaFirebase(c *fiber.Ctx) error {
	return nil
}
