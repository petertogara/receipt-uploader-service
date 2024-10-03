package models

type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
}

type LoginRequest struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}