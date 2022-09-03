package controllers

import (
	"bytes"
	"os"
	"sidewarslobby/pkg/repository"
	"sidewarslobby/pkg/utils"
	"sidewarslobby/platform/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

func validateServerToken(c *fiber.Ctx) bool {
	serverToken := bytes.NewBuffer(c.Request().Header.Peek("SW-ServerToken")).String()
	if serverToken == os.Getenv("SW_SERVER_TOKEN") {
		return true
	}

	utils.RESTError(c, "Invalid server token")
	return false
}

func ConfirmUserMatch(c *fiber.Ctx) error {
	if !validateServerToken(c) {
		return nil
	}

	payload := struct {
		UserMatchToken string
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	userMatch, err := database.DBQueries.GetUserMatch(payload.UserMatchToken)
	if err != nil {
		return utils.RESTError(c, "Maç bulunamadı")
	}

	return c.JSON(fiber.Map{
		"Username":     userMatch.User.Username,
		"UserChampion": userMatch.UserChampion,
	})
}

func FinishUserMatches(c *fiber.Ctx) error {
	if !validateServerToken(c) {
		return nil
	}

	payload := struct {
		UserMatchTokens   []string
		WinnerMatchTokens []string
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	for _, v := range payload.UserMatchTokens {
		userMatch, err := database.DBQueries.GetUserMatch(v)
		if err != nil {
			return utils.RESTError(c, "Maç bulunamadı")
		}

		userMatch.UpdatedAt = time.Now()
		userMatch.Finished = true
		userMatch.ScoreDiff = repository.LoseScoreDiff
		if utils.Contains(payload.WinnerMatchTokens, v) {
			// Change ScoreDiff If the token is also in the WinnerMatchTokens
			userMatch.ScoreDiff = repository.WinScoreDiff
			userMatch.UserWon = true
		}

		database.DBQueries.UpdateUserMatch(userMatch)
	}

	return c.JSON(fiber.Map{
		"Success": true,
	})
}
