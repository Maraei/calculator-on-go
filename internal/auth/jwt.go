package auth

import (
	"fmt"
	"log"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

const (
	HmacSampleSecret = "an7DkUH?L8iClxbVj5JZdbRVO2M$1Jc~D6CXsL@4"
)

func GenerateToken(id int) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"nbf": now.Unix(),
		"exp": now.Add(1 * time.Hour).Unix(),
		"iat": now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(HmacSampleSecret))
	if err != nil {
		log.Printf("jwt generation error: %v", err)
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(HmacSampleSecret), nil
	})

	if err != nil {
		log.Println("token parse error:", err)
		return 0, err
	}

	if !token.Valid {
		log.Println("invalid token")
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("invalid token claims")
		return 0, fmt.Errorf("invalid token claims")
	}

	log.Printf("Claims: %+v", claims)

	id, ok := claims["id"].(float64)
	if !ok {
		log.Println("invalid user ID type in token")
		return 0, fmt.Errorf("invalid user ID type")
	}

	userID := int(id)

	if userID <= 0 {
		log.Println("user ID is zero or negative")
		return 0, fmt.Errorf("invalid user ID in token")
	}

	return userID, nil
}
