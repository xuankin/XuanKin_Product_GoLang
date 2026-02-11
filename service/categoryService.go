package service

import (
	"Product_Mangement_Api/entity"
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/repository"
	"context"
	"github.com/google/uuid"
	"time"
)

type CategoryService interface {
	Create(ctx context.Context, req models.CreateCategoryRequest) (*models.CategoryResponse, error)
	ListAll(ctx context.Context) ([]models.CategoryResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, req models.UpdateCategoryRequest, Id uuid.UUID) (*models.CategoryResponse, error)
}
type categoryService struct {
	repo      repository.CategoryRepository
	cacheRepo repository.CacheRepository
}

func NewCategoryService(repo repository.CategoryRepository, cacheRepo repository.CacheRepository) CategoryService {
	return &categoryService{repo: repo,
		cacheRepo: cacheRepo}
}
func (s *categoryService) Create(ctx context.Context, req models.CreateCategoryRequest) (*models.CategoryResponse, error) {
	cat := &entity.Category{Name: toJson(req.Name)}
	if req.ParentID != nil {
		cat.ParentId = req.ParentID
	}
	res, err := s.repo.Create(ctx, cat)
	if err != nil {
		return nil, err
	}
	s.cacheRepo.Delete(ctx, models.CacheKeyCategoryAll)
	return &models.CategoryResponse{ID: res.ID, Name: toMap(res.Name), ParentID: res.ParentId}, nil
}
func (s *categoryService) ListAll(ctx context.Context) ([]models.CategoryResponse, error) {
	cacheKey := models.CacheKeyCategoryAll
	var resCache []models.CategoryResponse
	if err := s.cacheRepo.Get(ctx, cacheKey, &resCache); err == nil {
		return resCache, err
	}

	cats, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	var res []models.CategoryResponse
	for _, cat := range cats {
		res = append(res, models.CategoryResponse{ID: cat.ID, Name: toMap(cat.Name), ParentID: cat.ParentId})
	}
	s.cacheRepo.Set(ctx, cacheKey, res, 24*time.Hour)
	return res, nil
}
func (s *categoryService) Update(ctx context.Context, req models.UpdateCategoryRequest, Id uuid.UUID) (*models.CategoryResponse, error) {
	cat, err := s.repo.GetById(ctx, Id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		cat.Name = toJson(req.Name)
	}
	if req.ParentID != nil {
		cat.ParentId = req.ParentID
	}
	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}
	s.cacheRepo.Delete(ctx, models.CacheKeyCategoryAll)
	return &models.CategoryResponse{ID: cat.ID, Name: toMap(cat.Name), ParentID: cat.ParentId}, nil
}
func (s *categoryService) Delete(ctx context.Context, id uuid.UUID) error {

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.cacheRepo.Delete(ctx, models.CacheKeyCategoryAll)
	return nil
}
