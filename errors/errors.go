package errors

import "errors"

var (
    ErrNotFound         = errors.New("resource not found")
    ErrUnauthorized     = errors.New("unauthorized access")
    ErrInvalidInput     = errors.New("invalid input data")
    ErrFileTooLarge     = errors.New("file size exceeds limit")
    ErrInvalidFileType  = errors.New("invalid file type")
)
