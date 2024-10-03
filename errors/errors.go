package errors

import "errors"

var (
    ErrNotFound         = errors.New("resource not found")
    ErrUnauthorized     = errors.New("unauthorized access")
    ErrInvalidInput     = errors.New("invalid input data")
    ErrFileTooLarge     = errors.New("file size exceeds limit")
    ErrInvalidFileType  = errors.New("invalid file type")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrUserNotFound     = errors.New("user not found")
    ErrReceiptNotFound   = errors.New("receipt not found")
    ErrTokenExpired      = errors.New("token has expired")
    ErrTokenInvalid      = errors.New("invalid token")
)
