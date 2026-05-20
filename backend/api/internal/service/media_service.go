package service

import (
	"Server/internal/helpers"
	"Server/internal/model"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mediaService struct {
	cld *cloudinary.Cloudinary
}

func NewMediaService() *mediaService {
	service := &mediaService{}

	err := service.configureCloudService()
	if err != nil {
		panic(err)
	}

	return service
}

func (m *mediaService) Upload(file *multipart.FileHeader) (*model.Media, error) {
	if !isValidImage(file) {
		return nil, fmt.Errorf("invalid image type")
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	ctx, cancel := helpers.GenerateContext()
	defer cancel()
	unique := true
	uploadResult, err := m.cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder:         "posts",
		UniqueFilename: &unique,
		ResourceType:   "image",
	})
	if err != nil {
		return nil, err
	}

	media := &model.Media{
		ID:       primitive.NewObjectID(),
		URL:      uploadResult.SecureURL,
		PublicID: uploadResult.PublicID,
	}

	return media, nil
}

func (m *mediaService) UploadMany(files []*multipart.FileHeader) ([]*model.Media, error) {
	media := make([]*model.Media, len(files))

	for i, file := range files {
		f, err := m.Upload(file)
		if err != nil {
			for _, f := range media {
				if f != nil {
					m.Delete(f.PublicID)
				}
			}
			return nil, fmt.Errorf("failed to upload some files")
		}
		media[i] = f
	}

	return media, nil
}

func (m *mediaService) Delete(publicId string) error {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	_, err := m.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicId,
		ResourceType: "image",
	})

	if err != nil {
		return fmt.Errorf("error deleting image")
	}

	return nil
}

func (m *mediaService) DeleteMany(publicIds []string) error {
	for _, publicId := range publicIds {
		if err := m.Delete(publicId); err != nil {
			return err
		}
	}

	return nil
}

func (m *mediaService) DeleteManyByMedia(media []*model.Media) error {
	for _, file := range media {
		if err := m.Delete(file.PublicID); err != nil {
			return err
		}
	}
	return nil
}

func (m *mediaService) configureCloudService() error {
	cldApiKey := os.Getenv("CLD_API_KEY")
	cldApiSecret := os.Getenv("CLD_API_SECRET")
	cldName := os.Getenv("CLD_CLOUD_NAME")

	if cldApiKey == "" || cldApiSecret == "" || cldName == "" {
		return fmt.Errorf("no cloudinary credentials")
	}

	cloudinaryUrl := fmt.Sprintf("cloudinary://%s:%s@%s", cldApiKey, cldApiSecret, cldName)
	cld, err := cloudinary.NewFromURL(cloudinaryUrl)
	if err != nil {
		return err
	}

	m.cld = cld

	return nil
}

func isValidImage(file *multipart.FileHeader) bool {
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	return allowed[file.Header.Get("Content-Type")]
}
