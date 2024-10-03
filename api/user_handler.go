package api

import (
    "net/http"
    "receipt-uploader-service/controllers"
    "github.com/gin-gonic/gin"
)

var UserHandler = &UserAPI{}

type UserAPI struct {
    controller *controllers.UserController
}

func NewUserAPI(controller *controllers.UserController) *UserAPI {
    return &UserAPI{controller: controller}
}

func (h *UserAPI) CreateUser(c *gin.Context) {
    h.controller.RegisterUser(c)
}

func (h *UserAPI) Login(c *gin.Context) {
    h.controller.Login(c)
}

func (h *UserAPI) DeleteUser(c *gin.Context) {
    h.controller.DeleteUser(c)
}

func (h *UserAPI) GetUser(c *gin.Context) {
    h.controller.GetUser(c)
}
