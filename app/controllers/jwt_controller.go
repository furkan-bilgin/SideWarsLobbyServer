package controllers

import (
	"fmt"
	"os"
	"sidewarslobby/app/models"

	"github.com/golang-jwt/jwt/v4"
)

func JWTGetKey() []byte {
	return []byte(os.Getenv("JWT_KEY"))
}

// Create UserMatchToken for clients to use
func JWTCreateUserMatchToken(userMatch *models.UserMatch) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"MatchID": userMatch.ID,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString(JWTGetKey())
	return tokenString
}

// Valide the UserMatchToken and return MatchID
func JWTValidateUserMatchToken(jwtToken string) (int, error) {
	claims, err := parseJWTToken(jwtToken)
	if err != nil {
		return -1, err
	}

	return int(claims["MatchID"].(float64)), nil
}

func parseJWTToken(jwtToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return JWTGetKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("token not valid")
	}
}
