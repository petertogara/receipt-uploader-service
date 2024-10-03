package storage

import (
    "os"
    "receipt_uploader_service/models"
    "encoding/json"
    "github.com/google/uuid"
)

// StoreUser saves a new user and returns the user ID
func StoreUser(user *models.User) (string, error) {
    user.ID = uuid.New().String() // Generate a unique ID for the user

    // Here, you should hash the user's password before storing
    user.Password = hashPassword(user.Password) // Implement hashPassword

    // Store user info in a JSON file
    users := []models.User{user}
    // Load existing users from file, if any

    // Save updated user list back to file
    return user.ID, nil
}

// GetUserByUsername retrieves a user by their username
func GetUserByUsername(username string) (*models.User, error) {
    // Implement logic to read from JSON and return the user
    return nil, nil
}
