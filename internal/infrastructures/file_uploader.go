package infrastructures

import (
	"mime/multipart"
)

type FileUploaderImpl struct{}

func (r *FileUploaderImpl) Upload(fileName string, fh *multipart.FileHeader) error {
	return nil
}

func (r *FileUploaderImpl) Delete(fileID string) error {
	return nil
}
