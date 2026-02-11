package repository

import (
	"Product_Mangement_Api/entity"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BrandRepository interface {
	Create(ctx context.Context, b *entity.Brand) error
	Update(ctx context.Context, b *entity.Brand) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Brand, error)
	List(ctx context.Context) ([]entity.Brand, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CountProductByBrand(ctx context.Context, brandId uuid.UUID) (int64, error)
}

type brandRepository struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) BrandRepository {
	return &brandRepository{db: db}
}
func (r *brandRepository) Create(ctx context.Context, b *entity.Brand) error {
	return r.db.WithContext(ctx).Create(b).Error
}
func (r *brandRepository) Update(ctx context.Context, b *entity.Brand) error {
	return r.db.WithContext(ctx).Save(b).Error
}
func (r *brandRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Brand, error) {
	var b entity.Brand
	err := r.db.WithContext(ctx).First(&b, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}
func (r *brandRepository) List(ctx context.Context) ([]entity.Brand, error) {
	var brands []entity.Brand
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&brands).Error
	return brands, err
}
func (r *brandRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Brand{}, "id = ?", id).Error
}
func (r *brandRepository) CountProductByBrand(ctx context.Context, brandId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Product{}).Where("brand_id = ?", brandId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
