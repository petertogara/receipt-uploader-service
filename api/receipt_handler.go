package api

import (
    "net/http"
    "receipt-uploader-service/controllers"
    "github.com/gin-gonic/gin"
)

var ReceiptHandler = &ReceiptAPI{}

type ReceiptAPI struct {
    controller *controllers.ReceiptController
}

func NewReceiptAPI(controller *controllers.ReceiptController) *ReceiptAPI {
    return &ReceiptAPI{controller: controller}
}

func (h *ReceiptAPI) UploadReceipt(c *gin.Context) {
    h.controller.UploadReceipt(c)
}

func (h *ReceiptAPI) DownloadReceipt(c *gin.Context) {
    h.controller.GetReceipt(c)
}

func (h *ReceiptAPI) DeleteReceipt(c *gin.Context) {
    h.controller.DeleteReceipt(c)
}
