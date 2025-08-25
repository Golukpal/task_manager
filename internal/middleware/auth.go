package middleware

import (
    "os"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
)

func Protected() fiber.Handler {
    return func(c *fiber.Ctx) error {
        cookie := c.Cookies("jwt")

        token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Unauthorized",
            })
        }

        claims := token.Claims.(jwt.MapClaims)
        c.Locals("user", claims)

        return c.Next()
    }
}
