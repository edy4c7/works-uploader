package beans

import "mime/multipart"

type WorksFormBean struct {
	Title       string                `form:"title" binding:"required"`
	Description string                `form:"description" binding:"max=200"`
	Thumbnail   *multipart.FileHeader `form:"thumbnail" binding:"required"`
	Content     *multipart.FileHeader `form:"content" binding:"required"`
}
