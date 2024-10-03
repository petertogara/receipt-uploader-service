package controllers

import (
    "net/http"
    "receipt-uploader-service/models"
    "receipt-uploader-service/services"
    "receipt-uploader-service/auth" // Import the auth package
    "github.com/gin-gonic/gin"
)

type UserController struct {
    service *services.UserService 
}

// NewUserController creates a new instance of UserController
func NewUserController(service *services.UserService) *UserController {
    return &UserController{service: service}
}

// RegisterUser handles user registration
func (c *UserController) RegisterUser(ctx *gin.Context) {
    var user models.User
    if err := ctx.ShouldBindJSON(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
    }

    err := c.service.SaveUser(user)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, user)
}

// Login handles user authentication
func (c *UserController) Login(ctx *gin.Context) {
    var loginRequest models.LoginRequest
    if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
    }

    token, err := auth.Authenticate(loginRequest.Username, loginRequest.Password, c.service.storage)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// GetUser handles retrieving a user
func (c *UserController) GetUser(ctx *gin.Context) {
    userID := ctx.Param("user_id")

    user, err := c.service.GetUser(userID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, user)
}

// DeleteUser handles deleting a user
func (c *UserController) DeleteUser(ctx *gin.Context) {
    userID := ctx.Param("user_id")

    err := c.service.DeleteUser(userID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusNoContent, nil)
}
