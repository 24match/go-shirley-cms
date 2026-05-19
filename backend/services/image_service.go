package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
)

type ImageService struct{}

func NewImageService() *ImageService {
	return &ImageService{}
}

func generateFilename(originalName string) string {
	return fmt.Sprintf("%d_%s", time.Now().Unix(), originalName)
}

func (s *ImageService) GetImages(category string) ([]models.Image, error) {
	var images []models.Image
	query := config.DB.Order("created_at DESC")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Find(&images).Error
	return images, err
}

func (s *ImageService) UploadImage(filename, filepathStr, description, category string, fileSize int64, sortOrder int) (*models.Image, error) {
	os.MkdirAll("uploads", 0755)

	image := models.Image{
		Filename:    filename,
		FilePath:    filepathStr,
		FileSize:    fileSize,
		Description: description,
		Category:    category,
		SortOrder:   sortOrder,
	}
	err := config.DB.Create(&image).Error
	return &image, err
}

func (s *ImageService) UploadMultipleImages(files []ImageUpload) ([]models.Image, error) {
	os.MkdirAll("uploads", 0755)
	var uploadedImages []models.Image

	for i, file := range files {
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		filepathStr := filepath.Join("uploads", filename)

		if err := os.WriteFile(filepathStr, file.Data, 0644); err != nil {
			return nil, err
		}

		description := ""
		if i < len(file.Descriptions) {
			description = file.Descriptions[i]
		}

		image := models.Image{
			Filename:    filename,
			FilePath:    filepathStr,
			FileSize:    file.Size,
			Description: description,
			Category:    file.Category,
			SortOrder:   i,
		}
		config.DB.Create(&image)
		uploadedImages = append(uploadedImages, image)
	}

	return uploadedImages, nil
}

type ImageUpload struct {
	Filename     string
	Data         []byte
	Size         int64
	Category     string
	Descriptions []string
}

func (s *ImageService) DeleteImage(id uint) error {
	var image models.Image
	if err := config.DB.First(&image, id).Error; err != nil {
		return err
	}

	os.Remove(image.FilePath)
	return config.DB.Delete(&image).Error
}

func (s *ImageService) DeleteImageByFilename(filename string) error {
	var image models.Image
	if err := config.DB.Where("filename = ?", filename).First(&image).Error; err != nil {
		return err
	}

	os.Remove(image.FilePath)
	return config.DB.Delete(&image).Error
}

func (s *ImageService) UpdateImage(id uint, updates map[string]interface{}) (*models.Image, error) {
	var image models.Image
	if err := config.DB.First(&image, id).Error; err != nil {
		return nil, err
	}

	contentType := updates["contentType"].(string)
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if newFile, ok := updates["file"].(interface{}); ok && newFile != nil {
			oldPath := "./uploads/" + image.Filename
			os.Remove(oldPath)

			file := newFile.(map[string]interface{})
			filename := generateFilename(file["filename"].(string))
			fileData := file["data"].([]byte)
			os.WriteFile("./uploads/"+filename, fileData, 0644)
			image.Filename = filename
		}

		if description, ok := updates["description"].(string); ok {
			image.Description = description
		}
		if category, ok := updates["category"].(string); ok {
			image.Category = category
		}
		if sortOrderStr, ok := updates["sort_order"].(string); ok && sortOrderStr != "" {
			image.SortOrder, _ = strconv.Atoi(sortOrderStr)
		}
	} else {
		if description, ok := updates["description"].(string); ok && description != "" {
			image.Description = description
		}
		if longDescription, ok := updates["long_description"].(string); ok && longDescription != "" {
			image.LongDescription = longDescription
		}
		if category, ok := updates["category"].(string); ok && category != "" {
			image.Category = category
		}
		if sortOrder, ok := updates["sort_order"].(int); ok {
			image.SortOrder = sortOrder
		}
	}

	err := config.DB.Save(&image).Error
	return &image, err
}
