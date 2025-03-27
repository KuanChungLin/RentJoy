package interfaces

import "mime/multipart"

type CloudinaryService interface {
	UploadImages(file []*multipart.FileHeader) ([]struct {
		PublicID string
		ImageURL string
	}, error)
}
