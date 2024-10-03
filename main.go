package main

import (
    "github.com/gin-gonic/gin"
    "receipt_uploader-service/docs"
    "receipt_uploader-service/middleware"
    "receipt_uploader-service/api"
)

func main() {
   
    r := gin.Default()

    // Middleware for JWT authentication
    r.Use(middleware.AuthMiddleware())

    // API routes
    r.POST("/api/users", api.CreateUser)             
    r.POST("/api/login", api.LoginUser)               
    r.POST("/api/receipts", api.UploadReceipt)          
    r.GET("/api/users/:userId/receipts/:id", api.DownloadReceipt) 
    r.DELETE("/api/users/:userId/receipts/:id", api.DeleteReceipt) 

   
    docs.SwaggerInfo.Title = "Receipt Uploader API"
    docs.SwaggerInfo.Description = "API for uploading and managing receipt images."
    docs.SwaggerInfo.Version = "1.0"
    
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


    r.Run(":9090") 
}

