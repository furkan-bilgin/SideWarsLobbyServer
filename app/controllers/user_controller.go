package controllers

import (
	"bytes"
	"sidewarslobby/app/models"
	"sidewarslobby/pkg/repository"
	"sidewarslobby/pkg/utils"
	"sidewarslobby/platform/database"

	"github.com/gofiber/fiber/v2"
)

func validateUserToken(c *fiber.Ctx) *models.User {
	userToken := bytes.NewBuffer(c.Request().Header.Peek("SW-ClientToken")).String()
	user, err := database.DBQueries.GetUserByToken(userToken)
	if err != nil {
		utils.RESTError(c, "Invalid client token")
		return nil
	}

	return user
}

func GetLastFinishedUserMatch(c *fiber.Ctx) error {
	user := validateUserToken(c)
	if user == nil {
		return nil
	}
	if len(user.UserMatches) == 0 {
		return utils.RESTError(c, "No match found")
	}

	lastUserMatch := user.UserMatches[len(user.UserMatches)-1]
	return c.JSON(fiber.Map{
		"CurrentElo": user.CachedElo,
		"ScoreDiff":  lastUserMatch.ScoreDiff,
		"ShowRank":   len(user.UserMatches) >= repository.LerpKGameCount,
	})
}
