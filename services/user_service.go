package services

import (
    "errors"
    "time"
    "github.com/dgrijalva/jwt-go"
    "receipt_uploader_service/models"
    "receipt_uploader_service/storage"
    "github.com/gin-gonic/gin"
)

// CreateUser handles user registration
func CreateUser(c *gin.Context) (map[string]interface{}, error) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        return nil, errors.New("invalid input")
    }

    // Store the user
    userID, err := storage.StoreUser(&user)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "user_id": userID,
        "message": "User created successfully",
    }, nil
}

// LoginUser authenticates the user and returns a JWT
func LoginUser(c *gin.Context) (map[string]interface{}, error) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        return nil, errors.New("invalid input")
    }

    // Validate user credentials
    storedUser, err := storage.GetUserByUsername(user.Username)
    if err != nil || !checkPasswordHash(user.Password, storedUser.Password) {
        return nil, errors.New("invalid credentials")
    }

    token, err := generateJWT(storedUser.ID)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "token": token,
    }, nil
}

// checkPasswordHash compares a hash with a password
func checkPasswordHash(password, hash string) bool {
    // Implement password hash comparison logic
    return password == hash // Placeholder for actual hash comparison
}

// generateJWT generates a JWT token
func generateJWT(userID string) (string, error) {
    claims := jwt.MapClaims{}
    claims["user_id"] = userID
    claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // Token expires in 72 hours

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("Zimbabwe"))
}
