package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/swaggo/gin-swagger" 
    "github.com/swaggo/gin-swagger/swaggerFiles"
    _ "receipt-uploader-service/docs" 
    "receipt-uploader-service/middleware"
    "receipt-uploader-service/api"
    "receipt-uploader-service/controllers"
    "receipt-uploader-service/services"
    "receipt-uploader-service/storage"
)

func main() {
    r := gin.Default()


    r.Use(middleware.CORSMiddleware()) 

    // API routes
    userStorage := storage.NewUserStorage()
    userService := services.NewUserService(userStorage)
    userController := controllers.NewUserController(userService)
    userHandler := api.NewUserAPI(userController)

    receiptStorage := storage.NewReceiptStorage()
    receiptService := services.NewReceiptService(receiptStorage)
    receiptController := controllers.NewReceiptController(receiptService)
    receiptHandler := api.NewReceiptAPI(receiptController)

    // Apply JWT middleware to protected routes
    r.Use(middleware.JWTMiddleware)

    // User routes
    r.POST("/api/users", userHandler.CreateUser)
    r.POST("/api/users/login", userHandler.Login)
    r.DELETE("/api/users/:user_id", userHandler.DeleteUser)
    r.GET("/api/users/:user_id", userHandler.GetUser) 

    // Receipt routes
    r.POST("/api/receipts", receiptHandler.UploadReceipt)
    r.GET("/api/receipts/:user_id/:receipt_id", receiptHandler.DownloadReceipt)
    r.DELETE("/api/receipts/:user_id/:receipt_id", receiptHandler.DeleteReceipt)

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    log.Println("Server started on :9090")
    log.Fatal(r.Run(":9090")) 
}
