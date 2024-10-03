package services_test

import (
    "io/ioutil"
    "mime/multipart"
    "os"
    "path/filepath"
    "receipt-uploader-service/models"
    "receipt-uploader-service/services"
    "receipt-uploader-service/storage"
    "testing"
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

func createTestFile() (multipart.File, *multipart.FileHeader, error) {
    // Create a temporary file
    tmpFile, err := ioutil.TempFile(os.TempDir(), "test_receipt_*.png")
    if err != nil {
        return nil, nil, err
    }

    // Write some test data to the file
    if _, err := tmpFile.Write([]byte("test")); err != nil {
        return nil, nil, err
    }

    // Rewind the file
    if _, err := tmpFile.Seek(0, 0); err != nil {
        return nil, nil, err
    }

    // Create a FileHeader
    fileHeader := &multipart.FileHeader{
        Filename: filepath.Base(tmpFile.Name()),
        Size:     4, // size of "test"
        Header:   make(map[string][]string),
    }

    return tmpFile, fileHeader, nil
}

func TestSaveReceipt(t *testing.T) {
    mockStorage := &mockStorage{}
    receiptService := services.NewReceiptService(mockStorage)

    file, header, err := createTestFile()
    if err != nil {
        t.Fatalf("Failed to create test file: %v", err)
    }
    defer os.Remove(header.Filename) // clean up

    receipt := models.Receipt{
        ID:     "1",
        UserID: "user1",
        Path:   header.Filename,
    }

    err = receiptService.SaveReceipt(receipt, file)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    if len(mockStorage.receipts) != 1 {
        t.Errorf("Expected 1 receipt, got %d", len(mockStorage.receipts))
    }
}

func TestDownloadReceipt(t *testing.T) {
    mockStorage := &mockStorage{}
    receiptService := services.NewReceiptService(mockStorage)

    receipt := models.Receipt{
        ID:     "1",
        UserID: "user1",
        Path:   "test.png",
    }

    mockStorage.receipts = append(mockStorage.receipts, receipt)

    var writer io.Writer

    err := receiptService.DownloadReceipt("user1", "1", writer, "")
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}
