package storage

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "os"
    "receipt-uploader-service/models"
    "sync"
)

const receiptFile = "receipts.json"

type ReceiptStorage interface {
    SaveReceipt(receipt models.Receipt) error
    GetReceiptByID(userID, receiptID string) (*models.Receipt, error)
    DeleteReceipt(userID, receiptID string) error
}


type FileReceiptStorage struct {
    receipts []models.Receipt
    lock     sync.Mutex
}

func NewReceiptStorage() (ReceiptStorage, error) {
    storage := &FileReceiptStorage{}
    err := storage.loadReceipts()
    return storage, err
}

func (s *FileReceiptStorage) SaveReceipt(receipt models.Receipt) error {
    s.lock.Lock()
    defer s.lock.Unlock()

    s.receipts = append(s.receipts, receipt)
    return s.saveReceipts()
}


func (s *FileReceiptStorage) GetReceiptByID(userID, receiptID string) (*models.Receipt, error) {
    s.lock.Lock()
    defer s.lock.Unlock()

    for _, receipt := range s.receipts {
        if receipt.ID == receiptID && receipt.UserID == userID {
            return &receipt, nil
        }
    }
    return nil, errors.New("receipt not found")
}

func (s *FileReceiptStorage) DeleteReceipt(userID, receiptID string) error {
    s.lock.Lock()
    defer s.lock.Unlock()

    for i, receipt := range s.receipts {
        if receipt.ID == receiptID && receipt.UserID == userID {
            s.receipts = append(s.receipts[:i], s.receipts[i+1:]...)
            return s.saveReceipts()
        }
    }
    return errors.New("receipt not found")
}


func (s *FileReceiptStorage) loadReceipts() error {
    file, err := os.Open(receiptFile)
    if err != nil {
        if os.IsNotExist(err) {
            return nil 
        }
        return err
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return err
    }

    return json.Unmarshal(data, &s.receipts)
}


func (s *FileReceiptStorage) saveReceipts() error {
    data, err := json.Marshal(s.receipts)
    if err != nil {
        return err
    }

    return ioutil.WriteFile(receiptFile, data, 0644)
}
