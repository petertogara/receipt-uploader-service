package services

import (
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
    "errors"
    "receipt-uploader-service/models"
    "receipt-uploader-service/storage"
)

// ReceiptService provides methods for handling receipts
type ReceiptService struct {
    storage storage.ReceiptStorage 
}

// NewReceiptService creates a new instance of ReceiptService
func NewReceiptService(storage storage.ReceiptStorage) *ReceiptService {
    return &ReceiptService{storage: storage}
}

// SaveReceipt saves a receipt and its file
func (s *ReceiptService) SaveReceipt(receipt models.Receipt, file multipart.File, header *multipart.FileHeader) error {
    dir := filepath.Join("storage", "uploads")
    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        return err
    }

    filePath := filepath.Join(dir, header.Filename)
    outFile, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer outFile.Close()

    _, err = io.Copy(outFile, file)
    if err != nil {
        return err
    }

    receipt.Path = filePath
    return s.storage.SaveReceipt(receipt)
}

// GetReceipt retrieves a receipt by ID
func (s *ReceiptService) GetReceipt(userID, receiptID string) (*models.Receipt, error) {
    return s.storage.GetReceiptByID(userID, receiptID)
}

// DeleteReceipt removes a receipt
func (s *ReceiptService) DeleteReceipt(userID, receiptID string) error {
    return s.storage.DeleteReceipt(userID, receiptID)
}
