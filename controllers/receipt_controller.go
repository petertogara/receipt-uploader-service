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

func (c *ReceiptController) UploadReceipt(ctx *gin.Context) {
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

    f, err := file.Open()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not open file"})
        return
    }
    defer f.Close()

    err = c.service.SaveReceipt(receipt, f, file)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, receipt)
}


func (c *ReceiptController) GetReceipt(ctx *gin.Context) {
    userID := ctx.Param("user_id")
    receiptID := ctx.Param("receipt_id")

    receipt, err := c.service.GetReceipt(userID, receiptID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, receipt)
}


func (c *ReceiptController) DeleteReceipt(ctx *gin.Context) {
    userID := ctx.Param("user_id")
    receiptID := ctx.Param("receipt_id")

    err := c.service.DeleteReceipt(userID, receiptID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusNoContent, nil)
}
