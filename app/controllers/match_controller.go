package controllers

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"sidewarslobby/app/models"
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

	// Parse JWT token
	matchID, err := JWTValidateUserMatchToken(payload.UserMatchToken)
	if err != nil {
		return utils.RESTError(c, "Token hatalı")
	}

	// Get UserMatch from database
	userMatch, err := database.DBQueries.GetUserMatch(matchID)
	if err != nil {
		return utils.RESTError(c, "Maç bulunamadı")
	}

	return c.JSON(fiber.Map{
		"UserID":       userMatch.User.ID,
		"RoomID":       userMatch.Match.MatchmakingID,
		"Username":     userMatch.User.Username,
		"UserChampion": userMatch.UserChampion,
		"UserMatchID":  userMatch.ID,
		"TeamID":       userMatch.TeamID,
	})
}

func FinishUserMatches(c *fiber.Ctx) error {
	if !validateServerToken(c) {
		return nil
	}

	payload := struct {
		UserMatchIDs   string
		WinnerMatchIDs string
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	// Decode json lists
	var userMatchIDs []int
	var winnerMatchIDs []int
	json.Unmarshal([]byte(payload.UserMatchIDs), &userMatchIDs)
	json.Unmarshal([]byte(payload.WinnerMatchIDs), &winnerMatchIDs)

	for _, v := range userMatchIDs {
		userMatch, err := database.DBQueries.GetUserMatch(v)
		if err != nil {
			return utils.RESTError(c, "Maç bulunamadı, "+err.Error())
		}

		// Find enemy team
		enemyTeam := models.TeamRed
		if userMatch.TeamID == models.TeamRed {
			enemyTeam = models.TeamBlue
		}

		// Calculate average enemy elo
		enemies, err := database.DBQueries.GetMatchUsersByTeamID(&userMatch.Match, enemyTeam)
		if err != nil {
			return utils.RESTError(c, "Hata, "+err.Error())
		}

		enemySum := 0
		for _, v := range enemies {
			enemySum += v.CachedElo
		}
		averageEnemyElo := enemySum / len(enemies)

		// Lerp: K_Beginner -> K_Default over the course of LerpKGameCount games
		t := float64(repository.LerpKGameCount-len(userMatch.User.UserMatches)) / float64(repository.LerpKGameCount)
		t = math.Min(0, t)
		kValue := int(utils.LinearInterp(repository.DefaultEloK, repository.BeginnerEloK, t))

		elo := utils.NewEloWithFactors(kValue, utils.NewElo().D)

		// Update UserMatch
		gameResult := 0
		userMatch.UpdatedAt = time.Now()

		if !userMatch.Match.Finished {
			userMatch.Match.Finished = true
			database.DBQueries.UpdateMatch(&userMatch.Match)
		}

		// If we won, change these vars accordingly
		if utils.Contains(winnerMatchIDs, v) {
			gameResult = 1
			userMatch.UserWon = true
		}

		// gameResult = 1 -> user wins, gameResult = 0 -> enemy wins
		userMatch.ScoreDiff = elo.RatingDelta(userMatch.User.CachedElo, averageEnemyElo, float64(gameResult))

		// Finally, update UserMatch info, and re-cache user elo
		database.DBQueries.UpdateUserMatch(userMatch)
		database.DBQueries.CacheUserElo(&userMatch.User)
	}

	return c.JSON(fiber.Map{
		"Success": true,
	})
}
