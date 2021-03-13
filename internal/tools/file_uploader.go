package tools

import (
	"mime/multipart"
)

type FileUploader interface {
	Upload(string, *multipart.FileHeader) error
	Delete(string) error
}
