package lib

import "mime/multipart"

type StorageClient interface {
	Upload(string, *multipart.FileHeader) (string, error)
}
