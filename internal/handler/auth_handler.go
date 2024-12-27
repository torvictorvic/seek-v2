package handler

import (
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

// GenerateToken is responsible for creating and returning the JWT token
func GenerateToken(c *gin.Context) {

    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = ""
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "authorized": true,
        "user":       "usuario_demo",
        "exp":        time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Token not generated"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
