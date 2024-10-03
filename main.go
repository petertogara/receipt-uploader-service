package main

import (
    "log"
    "net/http"
    
    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"
    "github.com/swaggo/http-swagger"
    _ "receipt-uploader-service/docs"
    "receipt-uploader-service/middleware"
    "receipt-uploader-service/controllers"
    "receipt-uploader-service/storage"
    "receipt-uploader-service/services"
)

func main() {
    r := mux.NewRouter()

    cors := handlers.AllowedOrigins([]string{"*"})

    receiptStorage, err := storage.NewReceiptStorage()
    if err != nil {
        log.Fatalf("could not initialize storage: %v", err)
    }
    receiptService := services.NewReceiptService(receiptStorage)

    receiptController := controllers.NewReceiptController(receiptService)

    // Apply JWT middleware to protected routes
    r.Use(middleware.JWTMiddleware)

    // API routes
    r.HandleFunc("/api/receipts", receiptController.SaveReceipt).Methods("POST")              
    r.HandleFunc("/api/receipts/{userId}/{id}", receiptController.DownloadReceipt).Methods("GET") 
    r.HandleFunc("/api/receipts/{userId}/{id}", receiptController.DeleteReceipt).Methods("DELETE")

    // Swagger documentation route
    r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

    log.Println("Server started on :9090")
    log.Fatal(http.ListenAndServe(":9090", handlers.CORS(cors)(r)))
}
