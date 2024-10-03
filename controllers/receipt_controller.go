package controllers

import (
    "net/http"
    "receipt-uploader-service/models"
    "receipt-uploader-service/services"
    "github.com/gin-gonic/gin"
)

type ReceiptController struct {
    service *services.ReceiptService
}

// NewReceiptController creates a new instance of ReceiptController
func NewReceiptController(service *services.ReceiptService) *ReceiptController {
    return &ReceiptController{service: service}
}

// SaveReceipt handles the uploading of a receipt
func (c *ReceiptController) SaveReceipt(ctx *gin.Context) {
    var receipt models.Receipt
    file, err := ctx.FormFile("file")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
        return
    }

    if err := ctx.ShouldBindJSON(&receipt); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
    }

    // Call the service to save the receipt
    if err := c.service.SaveReceipt(receipt, file); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, receipt)
}

// GetReceipt handles downloading a receipt
func (c *ReceiptController) GetReceipt(ctx *gin.Context) {
    userID := ctx.Param("user_id")            // Get user ID from the URL
    receiptID := ctx.Param("receipt_id")      // Get receipt ID from the URL
    resolution := ctx.Query("resolution")     // Get resolution from query parameters

    // Serve the file for download
    if err := c.service.DownloadReceipt(userID, receiptID, ctx.Writer, resolution); err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
}

// DeleteReceipt handles deleting a receipt
func (c *ReceiptController) DeleteReceipt(ctx *gin.Context) {
    userID := ctx.Param("user_id")
    receiptID := ctx.Param("receipt_id")

    if err := c.service.DeleteReceipt(userID, receiptID); err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusNoContent, nil)
}
