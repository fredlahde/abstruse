package security

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/bleenco/abstruse/config"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSecret []byte
)

// UserJWT defines structure for generating JWT based on user.
type UserJWT struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

// InitSecurity initializes security tokens with values from config file.
func InitSecurity(cfg config.Security) {
	jwtSecret = []byte(cfg.JWTSecret)
}

// HashPassword generates encrypted password from password string
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// CheckPasswordHash compares password string with encrypted hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates JSON Web Token
func GenerateJWT(user UserJWT) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     user.ID,
		"email":  user.Email,
		"avatar": user.Avatar,
	})

	return token.SignedString(jwtSecret)
}

// GetUserIDFromJWT returns users ID from token.
func GetUserIDFromJWT(tokenString string) (int, error) {
	var userID int
	if tokenString == "" {
		return userID, errors.New("invalid token")
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Configuration.Security.JWTSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID, _ = strconv.Atoi(claims["id"].(string))
		return userID, nil
	}

	return userID, nil
}