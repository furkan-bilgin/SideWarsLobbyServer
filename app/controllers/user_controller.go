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
	lastMatch := database.DBQueries.GetMatch(int(lastUserMatch.MatchID))
	teams := make(map[string][]interface{})

	// Add users to teams dict
	for _, v := range lastMatch.UserMatches {
		team := "BlueTeam"
		if v.TeamID == repository.TeamRed {
			team = "RedTeam"
		}

		// BlueTeam/RedTeam struct
		data := struct {
			Elo      int
			Username string
		}{Elo: v.User.CachedElo, Username: v.User.Username}
		teams[team] = append(teams[team], data)
	}

	return c.JSON(fiber.Map{
		"CurrentElo": user.CachedElo,
		"ScoreDiff":  lastUserMatch.ScoreDiff,
		"ShowRank":   len(user.UserMatches) >= repository.LerpKGameCount,
		"BlueTeam":   teams["BlueTeam"],
		"RedTeam":    teams["RedTeam"],
	})
}

func SetUserChampion(c *fiber.Ctx) error {
	user := validateUserToken(c)
	if user == nil {
		return nil
	}

	payload := struct {
		SelectedChampion int
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	// TODO: Validate champion ID
	err := database.DBQueries.UpdateUserInfo(user.UserInfo, models.UserInfo{
		SelectedChampion: uint8(payload.SelectedChampion),
	})

	if err != nil {
		return utils.RESTError(c, "Failed to update")
	}

	return c.JSON(fiber.Map{
		"Success":  true,
		"UserInfo": user.UserInfo.Sanitize(),
	})
}
