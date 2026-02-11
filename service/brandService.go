package service

import (
	"Product_Mangement_Api/entity"
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/repository"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type BrandService interface {
	Create(ctx context.Context, req models.CreateBrandRequest) (*models.BrandResponse, error)
	List(ctx context.Context) ([]models.BrandResponse, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.BrandResponse, error)
	UpdateById(ctx context.Context, id uuid.UUID, req models.UpdateBrandRequest) (*models.BrandResponse, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type brandService struct {
	repo      repository.BrandRepository
	cacheRepo repository.CacheRepository
}

func NewBrandService(repo repository.BrandRepository, cacheRepo repository.CacheRepository) BrandService {
	return &brandService{repo: repo, cacheRepo: cacheRepo}
}
func (r *brandService) Create(ctx context.Context, req models.CreateBrandRequest) (*models.BrandResponse, error) {

	brand := &entity.Brand{Name: toJson(req.Name),
		Logo: req.Logo}
	err := r.repo.Create(ctx, brand)
	if err != nil {
		return nil, err
	}
	if err == nil {
		r.cacheRepo.Delete(ctx, models.CacheKeyBrandAll)
	}
	return &models.BrandResponse{
		ID:   brand.ID,
		Name: toMap(brand.Name),
		Logo: brand.Logo,
	}, nil
}
func (r *brandService) List(ctx context.Context) ([]models.BrandResponse, error) {
	cacheKey := models.CacheKeyBrandAll
	var resCache []models.BrandResponse
	if err := r.cacheRepo.Get(ctx, cacheKey, &resCache); err == nil {
		return resCache, nil
	}

	brands, err := r.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	var res []models.BrandResponse
	for _, brand := range brands {
		res = append(res, models.BrandResponse{
			ID:   brand.ID,
			Name: toMap(brand.Name),
			Logo: brand.Logo,
		})
	}
	r.cacheRepo.Set(ctx, cacheKey, res, 24*time.Hour)
	return res, nil
}
func (r *brandService) GetById(ctx context.Context, id uuid.UUID) (*models.BrandResponse, error) {
	brand, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &models.BrandResponse{
		ID:   brand.ID,
		Name: toMap(brand.Name),
		Logo: brand.Logo,
	}, nil
}
func (r *brandService) UpdateById(ctx context.Context, id uuid.UUID, req models.UpdateBrandRequest) (*models.BrandResponse, error) {
	br, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		br.Name = toJson(req.Name)
	}
	if req.Logo != "" {
		br.Logo = req.Logo
	}
	if err := r.repo.Update(ctx, br); err != nil {
		return nil, err
	}
	if err == nil {
		r.cacheRepo.Delete(ctx, models.CacheKeyBrandAll)
	}

	return &models.BrandResponse{
		ID:   br.ID,
		Name: toMap(br.Name),
		Logo: br.Logo,
	}, nil
}
func (r *brandService) DeleteById(ctx context.Context, id uuid.UUID) error {
	_, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("Brand not found")
	}
	count, er := r.repo.CountProductByBrand(ctx, id)
	if er != nil {
		return er
	}
	if count > 0 {
		return errors.New("Cannot delete brand with existing products")
	}
	if err == nil {
		r.cacheRepo.Delete(ctx, models.CacheKeyBrandAll)
	}
	return r.repo.Delete(ctx, id)
}
