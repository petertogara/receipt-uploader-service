package auth

import (
    "errors"
    "time"
    "github.com/dgrijalva/jwt-go"
    "receipt-uploader-service/models"
    "receipt-uploader-service/storage"
)

// Secret key used for signing tokens
var secretKey = []byte("your_secret_key")

// Claims structure for JWT
type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

// GenerateToken generates a JWT token for the authenticated user
func GenerateToken(username string) (string, error) {
    claims := Claims{
        Username: username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}

// Authenticate validates user credentials and generates a JWT token
func Authenticate(username, password string, userStorage storage.UserStorage) (string, error) {
    user, err := userStorage.GetUserByUsername(username)
    if err != nil {
        return "", errors.New("user not found")
    }

    if user.Password != password {
        return "", errors.New("invalid password")
    }

    return GenerateToken(user.Username)
}

// Middleware for JWT authentication
func Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.Request.Header.Get("Authorization")

        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, errors.New("unexpected signing method")
            }
            return secretKey, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        c.Next()
    }
}
