package tools

import (
	"mime/multipart"
)

type FileUploader interface {
	Upload(string, *multipart.FileHeader) (string, error)
	Delete(string) error
}
