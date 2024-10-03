package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "receipt-uploader-service/services"
)

// CreateUser creates a new user
// @Summary Create a user
// @Description Creates a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User Info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/users [post]
func CreateUser(c *gin.Context) {
    response, err := services.CreateUser(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, response)
}

// LoginUser logs in a user and generates a JWT
// @Summary User Login
// @Description Logs in a user and returns a JWT
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User Login Info"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/login [post]
func LoginUser(c *gin.Context) {
    response, err := services.LoginUser(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, response)
}
