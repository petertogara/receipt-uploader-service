package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "receipt-uploader-service/services"
)

// UploadReceipt uploads a receipt image
// @Summary Upload a receipt
// @Description Uploads a receipt image
// @Tags receipts
// @Accept multipart/form-data
// @Produce json
// @Param receipt formData file true "Receipt Image"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/receipts [post]
func UploadReceipt(c *gin.Context) {
    // Call service to handle the upload logic
    response, err := services.UploadReceipt(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, response)
}

// DownloadReceipt downloads a receipt image
// @Summary Download a receipt
// @Description Downloads a receipt image by ID
// @Tags receipts
// @Produce json
// @Param id path string true "Receipt ID"
// @Param userId path string true "User ID"
// @Param resolution query string false "Image Resolution"
// @Success 200 {string} string "File downloaded"
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/users/{userId}/receipts/{id} [get]
func DownloadReceipt(c *gin.Context) {
    // Call service to handle the download logic
    services.DownloadReceipt(c)
}

// DeleteReceipt deletes a receipt image
// @Summary Delete a receipt
// @Description Deletes a receipt image by ID
// @Tags receipts
// @Param id path string true "Receipt ID"
// @Param userId path string true "User ID"
// @Success 204 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/users/{userId}/receipts/{id} [delete]
func DeleteReceipt(c *gin.Context) {
    services.DeleteReceipt(c)
}
