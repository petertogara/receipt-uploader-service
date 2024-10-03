package services

import (
    "errors"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "receipt_uploader_service/models"
    "receipt_uploader_service/storage"
    "github.com/gin-gonic/gin"
)

// UploadReceipt handles the receipt upload logic
func UploadReceipt(c *gin.Context) (map[string]interface{}, error) {
    file, err := c.FormFile("receipt")
    if err != nil {
        return nil, err
    }

    // Validate file format and size
    if err := validateFile(file); err != nil {
        return nil, err
    }

    userID := c.Param("userId")
    receiptID, err := storage.StoreReceipt(file, userID)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "receipt_id": receiptID,
        "message":    "Receipt uploaded successfully",
    }, nil
}

// DownloadReceipt downloads a receipt image
func DownloadReceipt(c *gin.Context) {
    receiptID := c.Param("id")
    userID := c.Param("userId")
    resolution := c.Query("resolution")

    receipt, err := storage.GetReceiptByID(receiptID)
    if err != nil || receipt == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
        return
    }

    // Check ownership
    if receipt.UserID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this receipt"})
        return
    }

    imgPath := filepath.Join("uploads", receipt.Path)
    if resolution != "" {
        imgPath = adjustPathForResolution(imgPath, resolution)
    }

    c.File(imgPath) // Send the file to the client
}

// DeleteReceipt deletes a receipt
func DeleteReceipt(c *gin.Context) {
    receiptID := c.Param("id")
    userID := c.Param("userId")

    receipt, err := storage.GetReceiptByID(receiptID)
    if err != nil || receipt == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
        return
    }

    // Check ownership
    if receipt.UserID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this receipt"})
        return
    }

    err = storage.DeleteReceipt(receiptID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete receipt"})
        return
    }

    c.Status(http.StatusNoContent)
}

// validateFile checks if the uploaded file is valid
func validateFile(file *multipart.FileHeader) error {
    if file.Size > 5*1024*1024 { // Limit file size to 5MB
        return errors.New("file size exceeds limit")
    }

    validExtensions := map[string]bool{
        ".jpg":  true,
        ".jpeg": true,
        ".png":  true,
    }

    ext := filepath.Ext(file.Filename)
    if !validExtensions[ext] {
        return errors.New("invalid file type")
    }

    return nil
}

// adjustPathForResolution modifies the file path to point to the requested resolution
func adjustPathForResolution(originalPath string, resolution string) string {
    base := originalPath[:len(originalPath)-len(filepath.Ext(originalPath))]
    ext := filepath.Ext(originalPath)
    return base + "_" + resolution + ext
}
