package storage

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "os"
    "receipt-uploader-service/models"
)

const userFile = "users.json"

type UserStorage interface {
    SaveUser(user models.User) error
    GetUserByID(userID string) (*models.User, error)
    DeleteUser(userID string) error
}

type FileUserStorage struct {
    users map[string]models.User
}

func NewUserStorage() (UserStorage, error) {
    storage := &FileUserStorage{
        users: make(map[string]models.User),
    }
    err := storage.loadUsers()
    return storage, err
}


func (s *FileUserStorage) SaveUser(user models.User) error {
    s.users[user.ID] = user
    return s.saveUsers()
}

func (s *FileUserStorage) GetUserByID(userID string) (*models.User, error) {
    user, exists := s.users[userID]
    if !exists {
        return nil, errors.New("user not found")
    }
    return &user, nil
}


func (s *FileUserStorage) DeleteUser(userID string) error {
    delete(s.users, userID)
    return s.saveUsers()
}

func (s *FileUserStorage) loadUsers() error {
    file, err := os.Open(userFile)
    if err != nil {
        if os.IsNotExist(err) {
            return nil 
        }
        return err
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return err
    }

    return json.Unmarshal(data, &s.users)
}

// saveUsers saves the users to the JSON file
func (s *FileUserStorage) saveUsers() error {
    data, err := json.Marshal(s.users)
    if err != nil {
        return err
    }

    return ioutil.WriteFile(userFile, data, 0644)
}
