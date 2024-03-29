package beans

import (
	"mime/multipart"

	"github.com/edy4c7/works-uploader/internal/common/constants"
)

type WorksFormBean struct {
	Type        constants.WorkType    `form:"type" binding:"required"`
	Title       string                `form:"title" binding:"required,max=40"`
	Description string                `form:"description" binding:"max=200"`
	ContentURL  string                `form:"url" binding:"required_if=Type 1,omitempty,url"`
	Thumbnail   *multipart.FileHeader `form:"thumbnail" binding:"required_if=Type 2"`
	Content     *multipart.FileHeader `form:"content" binding:"required_if=Type 2"`
}
