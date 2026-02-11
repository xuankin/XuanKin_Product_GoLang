package service

import (
	"Product_Mangement_Api/entity"
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/repository"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type MediaService interface {
	Upload(ctx context.Context, file io.Reader, filename string, req models.CreateMediaRequest) (*models.MediaResponse, error)
}

type mediaService struct {
	repo      repository.MediaRepository
	uploadDir string
	baseUrl   string
}

func NewMediaService(repo repository.MediaRepository, baseUrl string) MediaService {
	uploadDir := "./uploads"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	return &mediaService{
		repo:      repo,
		uploadDir: uploadDir,
		baseUrl:   baseUrl,
	}
}
func (s *mediaService) Upload(ctx context.Context, file io.Reader, filename string, req models.CreateMediaRequest) (*models.MediaResponse, error) {

	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return nil, errors.New("error reading file to check format")
	}

	fileType := http.DetectContentType(buff)

	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	} else {
		return nil, errors.New("File upload does not support seek (internal server error)")
	}

	isValidImage := strings.HasPrefix(fileType, "image/")
	isValidVideo := strings.HasPrefix(fileType, "video/")

	if !isValidImage && !isValidVideo {
		return nil, fmt.Errorf("invalid file format: %s (only accepts images or videos)", fileType)
	}

	ext := filepath.Ext(filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(s.uploadDir, newFileName)

	out, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return nil, err
	}

	mediaType := models.MediaTypeImage
	if isValidVideo {
		mediaType = models.MediaTypeVideo
	}

	fileUrl := fmt.Sprintf("%s/uploads/%s", s.baseUrl, newFileName)

	mediaEntity := &entity.Media{
		ProductID:    req.ProductID,
		VariantID:    req.VariantID,
		Type:         mediaType,
		URL:          fileUrl,
		ThumbnailURL: req.ThumbnailURL,
		IsPrimary:    req.IsPrimary,
		SortOrder:    req.SortOrder,
	}

	res, err := s.repo.Create(ctx, mediaEntity)
	if err != nil {
		return nil, err
	}

	return &models.MediaResponse{
		ID:           res.ID,
		Type:         res.Type,
		URL:          res.URL,
		ThumbnailURL: res.ThumbnailURL,
		IsPrimary:    res.IsPrimary,
		SortOrder:    res.SortOrder,
	}, nil
}
