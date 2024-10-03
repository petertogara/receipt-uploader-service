package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
)

// AuthMiddleware checks the JWT token for protected routes
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.Request.Header.Get("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        // Validate the token
        claims := &jwt.MapClaims{}
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte("Zimbabwe"), nil 
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Set user ID in the context
        c.Set("userID", claims["user_id"])
        c.Next()
    }
}
