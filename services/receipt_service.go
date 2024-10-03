package services

import (
    "bytes"
    "errors"
    "image"
    "image/jpeg"
    "image/png"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
    "receipt-uploader-service/models"
    "receipt-uploader-service/storage"
    "golang.org/x/image/draw" // Import for image resizing
)

type ReceiptService struct {
    storage storage.ReceiptStorage 
}

func NewReceiptService(storage storage.ReceiptStorage) *ReceiptService {
    return &ReceiptService{storage: storage}
}

func (s *ReceiptService) SaveReceipt(receipt models.Receipt, file multipart.File) error {
    if err := validateImageFile(file); err != nil {
        return err
    }

    dir := filepath.Join("uploads") 
    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        return err
    }

    filePath := filepath.Join(dir, filepath.Base(receipt.Path))
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

func (s *ReceiptService) DeleteReceipt(userID, receiptID string) error {
    return s.storage.DeleteReceipt(userID, receiptID)
}

func (s *ReceiptService) DownloadReceipt(userID, receiptID string, writer io.Writer, resolution string) error {
    receipt, err := s.storage.GetReceiptByID(userID, receiptID)
    if err != nil {
        return err
    }

    
    return serveFileWithResolution(receipt.Path, writer, resolution)
}

func validateImageFile(file multipart.File) error {
    // Create a buffer to hold the file data
    buf := make([]byte, 512)
    _, err := file.Read(buf)
    if err != nil {
        return errors.New("unable to read file")
    }

    // Reset the file cursor
    _, err = file.Seek(0, 0)
    if err != nil {
        return errors.New("unable to reset file pointer")
    }

    // Check the MIME type
    fileType := http.DetectContentType(buf)
    if fileType != "image/jpeg" && fileType != "image/png" {
        return errors.New("uploaded file is not a valid image type (JPEG or PNG)")
    }
    return nil
}

func serveFileWithResolution(filePath string, writer io.Writer, resolution string) error {
    // Open the file
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    // Decode the image
    img, _, err := image.Decode(file)
    if err != nil {
        return err
    }

    // Resize the image based on the requested resolution
    var newImg image.Image
    switch resolution {
    case "small":
        newImg = resizeImage(img, 100, 100) // Resize to 100x100 for small
    case "medium":
        newImg = resizeImage(img, 400, 400) // Resize to 400x400 for medium
    case "large":
        newImg = resizeImage(img, 800, 800) // Resize to 800x800 for large
    default:
        newImg = img // Serve original image if resolution is invalid
    }

    // Encode the new image to the writer
    switch img.(type) {
    case *image.NRGBA:
        err = png.Encode(writer, newImg) // Encode as PNG
    case *image.YCbCr:
        err = jpeg.Encode(writer, newImg, nil) // Encode as JPEG
    default:
        err = jpeg.Encode(writer, newImg, nil) // Default to JPEG
    }

    return err
}

// resizeImage resizes the given image to the specified width and height
func resizeImage(img image.Image, width, height int) image.Image {
    // Create a new blank image with the desired size
    newImg := image.NewNRGBA(image.Rect(0, 0, width, height))
    // Use the draw package to resize the original image into the new image
    draw.BiLinear.Scale(newImg, newImg.Rect, img, img.Bounds(), draw.Over, nil)
    return newImg
}
