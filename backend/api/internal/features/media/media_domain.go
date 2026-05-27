package media

import (
	"mime/multipart"
)

type MediaService interface {
	Upload(file *multipart.FileHeader) (*Media, error)
	UploadMany(files []*multipart.FileHeader) ([]*Media, error)
	Delete(publicId string) error
	DeleteMany(publicIds []string) error
	DeleteManyByMedia(media []*Media) error
}
