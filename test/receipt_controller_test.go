package controllers_test

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "receipt-uploader-service/controllers"
    "receipt-uploader-service/models"
    "receipt-uploader-service/services"
    "receipt-uploader-service/storage"
    "testing"

    "github.com/gin-gonic/gin"
)

type mockStorage struct {
    receipts []models.Receipt
}

func (m *mockStorage) SaveReceipt(receipt models.Receipt) error {
    m.receipts = append(m.receipts, receipt)
    return nil
}

func (m *mockStorage) GetReceiptByID(userID, receiptID string) (*models.Receipt, error) {
    for _, receipt := range m.receipts {
        if receipt.UserID == userID && receipt.ID == receiptID {
            return &receipt, nil
        }
    }
    return nil, errors.New("receipt not found")
}

func (m *mockStorage) DeleteReceipt(userID, receiptID string) error {
    for i, receipt := range m.receipts {
        if receipt.UserID == userID && receipt.ID == receiptID {
            m.receipts = append(m.receipts[:i], m.receipts[i+1:]...)
            return nil
        }
    }
    return errors.New("receipt not found")
}

func TestSaveReceiptHandler(t *testing.T) {
    mockStorage := &mockStorage{}
    receiptService := services.NewReceiptService(mockStorage)
    receiptController := controllers.NewReceiptController(receiptService)

    router := gin.Default()
    router.POST("/api/receipts", receiptController.SaveReceipt)

    receipt := models.Receipt{
        ID:     "1",
        UserID: "user1",
        Path:   "test.png",
    }
    
    // Create a multipart form file
    body := new(bytes.Buffer)
    writer := multipart.NewWriter(body)
    file, err := writer.CreateFormFile("file", "test.png")
    if err != nil {
        t.Fatalf("Unable to create form file: %v", err)
    }
    file.Write([]byte("dummy content")) // Dummy content for the test

    // Close the writer to set the correct content type
    writer.Close()

    req, _ := http.NewRequest("POST", "/api/receipts", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)

    if resp.Code != http.StatusCreated {
        t.Errorf("Expected status code 201, got %d", resp.Code)
    }

    if len(mockStorage.receipts) != 1 {
        t.Errorf("Expected 1 receipt in storage, got %d", len(mockStorage.receipts))
    }
}

func TestDownloadReceiptHandler(t *testing.T) {
    mockStorage := &mockStorage{}
    receiptService := services.NewReceiptService(mockStorage)
    receiptController := controllers.NewReceiptController(receiptService)

    router := gin.Default()
    router.GET("/api/receipts/:userId/:id", receiptController.DownloadReceipt)

    receipt := models.Receipt{
        ID:     "1",
        UserID: "user1",
        Path:   "test.png",
    }
    mockStorage.receipts = append(mockStorage.receipts, receipt)

    req, _ := http.NewRequest("GET", "/api/receipts/user1/1", nil)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)

    if resp.Code != http.StatusOK {
        t.Errorf("Expected status code 200, got %d", resp.Code)
    }
}
