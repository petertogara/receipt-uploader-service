package services

import (
    "encoding/json"
    "errors"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "receipt-uploader-service/models"
    "receipt-uploader-service/storage"
)


var allowedImageFileTypes = map[string]bool{
    ".jpg":  true,
    ".jpeg": true,
    ".png":  true,
}

type ReceiptService struct {
    storage storage.ReceiptStorage 
}

/
func NewReceiptService(storage storage.ReceiptStorage) *ReceiptService {
    return &ReceiptService{storage: storage}
}


func (s *ReceiptService) SaveReceipt(receipt models.Receipt, file multipart.File, header *multipart.FileHeader) error {
   
    if !allowedImageFileTypes[filepath.Ext(header.Filename)] {
        return errors.New("invalid file type; only images are allowed")
    }

    dir := filepath.Join("uploads")
    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        return err
    }

    filePath := filepath.Join(dir, header.Filename)

    if _, err := os.Stat(filePath); err == nil {
        return errors.New("file already exists")
    }

    outFile, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer outFile.Close()

    _, err = io.Copy(outFile, file)
    if err != nil {
        return err
    }

    // Set the receipt path and save the receipt
    receipt.Path = filePath
    if err := s.storage.SaveReceipt(receipt); err != nil {
        return err
    }

    // Update receipts.json after saving
    return s.UpdateReceiptsJSON()
}

// UpdateReceiptsJSON updates the receipts.json file with the current receipts
func (s *ReceiptService) UpdateReceiptsJSON() error {
    receipts, err := s.storage.GetAllReceipts() // Implement this method in your storage layer
    if err != nil {
        return err
    }

    receiptsFilePath := filepath.Join("storage", "receipts.json")
    file, err := os.Create(receiptsFilePath)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(receipts)
}

// GetReceipt retrieves a receipt by user ID and receipt ID
func (s *ReceiptService) GetReceipt(userID, receiptID string) (*models.Receipt, error) {
    receipt, err := s.storage.GetReceiptByID(userID, receiptID)
    if err != nil {
        return nil, err
    }

    // Check if receipt file exists
    if _, err := os.Stat(receipt.Path); os.IsNotExist(err) {
        return nil, errors.New("receipt file does not exist")
    }

    return receipt, nil
}

// DownloadReceipt serves the receipt file for download
func (s *ReceiptService) DownloadReceipt(userID, receiptID string, w http.ResponseWriter) error {
    receipt, err := s.GetReceipt(userID, receiptID)
    if err != nil {
        return err
    }

    // Set the content type based on the file type
    w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(receipt.Path))
    w.Header().Set("Content-Type", "application/octet-stream")
    
    file, err := os.Open(receipt.Path)
    if err != nil {
        return err
    }
    defer file.Close()

    // Serve the file for download
    _, err = io.Copy(w, file)
    return err
}

// DeleteReceipt removes a receipt and its image file
func (s *ReceiptService) DeleteReceipt(userID, receiptID string) error {
    receipt, err := s.storage.GetReceiptByID(userID, receiptID)
    if err != nil {
        return err // Receipt not found
    }

    // Check if the physical file exists before deletion
    if _, err := os.Stat(receipt.Path); os.IsNotExist(err) {
        return errors.New("receipt file does not exist")
    }

  
    if err := os.Remove(receipt.Path); err != nil {
        return err // Error deleting the file
    }

    if err := s.storage.DeleteReceipt(userID, receiptID); err != nil {
        return err
    }

    // Update receipts.json after deletion
    return s.UpdateReceiptsJSON()
}
