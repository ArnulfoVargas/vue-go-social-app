package domain

import (
	"Server/internal/model"
	"mime/multipart"
)

type MediaService interface {
	Upload(file *multipart.FileHeader) (*model.Media, error)
	UploadMany(files []*multipart.FileHeader) ([]*model.Media, error)
	Delete(publicId string) error
	DeleteMany(publicIds []string) error
	DeleteManyByMedia(media []*model.Media) error
}
