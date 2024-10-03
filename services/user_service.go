package services

import (
    "errors"
    "receipt-uploader-service/models"
    "receipt-uploader-service/storage"
)

// UserService provides methods for handling users
type UserService struct {
    storage storage.UserStorage 
}

// NewUserService creates a new instance of UserService
func NewUserService(storage storage.UserStorage) *UserService {
    return &UserService{storage: storage}
}

// SaveUser saves a new user
func (s *UserService) SaveUser(user models.User) error {
    existingUser, _ := s.storage.GetUserByUsername(user.Username)
    if existingUser != nil {
        return errors.New("username already exists")
    }
    return s.storage.SaveUser(user)
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(userID string) (*models.User, error) {
    return s.storage.GetUserByID(userID)
}

// DeleteUser removes a user
func (s *UserService) DeleteUser(userID string) error {
    return s.storage.DeleteUser(userID)
}
