package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Subject  string `json:"subject"`
	Username string `json:"username"`
	PhotoUrl string `json:"photo_url"`
}

func NewMiddleware(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		auth := c.Request().Header.Peek("Authorization")

		if len(auth) < 8 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Bearer <token>
		// 01234567
		tokenString := string(auth)[7:]

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("user", claims)

		return c.Next()
	}
}

func MustSubject(c fiber.Ctx) string {
	claims := c.Locals("user").(jwt.MapClaims)
	return claims["subject"].(string)
}

func MustUser(c fiber.Ctx) *UserClaims {
	claims := c.Locals("user").(jwt.MapClaims)
	user := &UserClaims{
		Subject:  claims["subject"].(string),
		Username: claims["username"].(string),
		PhotoUrl: claims["photo_url"].(string),
	}
	return user
}
