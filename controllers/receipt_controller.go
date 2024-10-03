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

func NewReceiptController(service *services.ReceiptService) *ReceiptController {
    return &ReceiptController{service: service}
}

func (c *ReceiptController) SaveReceipt(ctx *gin.Context) {
    var receipt models.Receipt
    if err := ctx.ShouldBindJSON(&receipt); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
    }

    file, err := ctx.FormFile("file")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
        return
    }

    err = c.service.SaveReceipt(receipt, file)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, receipt)
}

func (c *ReceiptController) DownloadReceipt(ctx *gin.Context) {
    userID := ctx.Param("userId")
    receiptID := ctx.Param("id")
    resolution := ctx.Query("resolution") 

    
    if err := c.service.DownloadReceipt(userID, receiptID, ctx.Writer, resolution); err != nil {
        if err.Error() == "receipt not found" {
            ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        } else {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
}

func (c *ReceiptController) DeleteReceipt(ctx *gin.Context) {
    userID := ctx.Param("userId")
    receiptID := ctx.Param("id")

    err := c.service.DeleteReceipt(userID, receiptID)
    if err != nil {
        if err.Error() == "receipt not found" {
            ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        } else {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    ctx.JSON(http.StatusNoContent, nil)
}
