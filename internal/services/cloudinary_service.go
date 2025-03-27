package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	interfaces "rentjoy/internal/interfaces/services"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cloudinary *cloudinary.Cloudinary
}

func NewCloudinaryService() (interfaces.CloudinaryService, error) {
	// 從.env 取得 CLOUDINARY_URL
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		return nil, fmt.Errorf("CLOUDINARY_URL environment variable is not set")
	}

	// 初始化 Cloudinary
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %v", err)
	}
	// 設置 Secure 為 true
	cld.Config.URL.Secure = true

	return &CloudinaryService{cloudinary: cld}, nil
}

// 上傳圖片
func (s *CloudinaryService) UploadImages(files []*multipart.FileHeader) ([]struct {
	PublicID string
	ImageURL string
}, error) {
	log.Println("UploadImages", files)
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	var results []struct {
		PublicID string
		ImageURL string
	}

	for _, file := range files {
		// 檢查文件
		if err := s.checkImage(file); err != nil {
			return nil, err
		}

		// 打開文件
		src, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %v", err)
		}
		defer src.Close()

		// 確定上傳文件夾
		folder := "ManagementImg"
		if len(files) >= 3 {
			folder = "VenuesImg"
		}
		// 上傳到 Cloudinary
		uploadResult, err := s.cloudinary.Upload.Upload(context.Background(), src, uploader.UploadParams{
			Folder:       folder,
			ResourceType: "image",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to upload to cloudinary: %v", err)
		}
		// 檢查上傳狀態是否成功
		if uploadResult == nil {
			return nil, fmt.Errorf("upload failed: no result returned")
		}

		results = append(results, struct {
			PublicID string
			ImageURL string
		}{
			PublicID: uploadResult.PublicID,
			ImageURL: uploadResult.SecureURL,
		})
	}

	return results, nil
}

// 檢查圖片
func (s *CloudinaryService) checkImage(file *multipart.FileHeader) error {
	// 檢查文件大小
	if file.Size > 5*1024*1024 {
		return fmt.Errorf("file size exceeds 5MB: %s", file.Filename)
	}

	// 檢查文件類型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, validExt := range validExtensions {
		if ext == validExt {
			return nil
		}
	}

	return fmt.Errorf("invalid file type: %s", file.Filename)
}
