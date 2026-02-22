package repository

import (
	"Product_Mangement_Api/entity"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaRepository interface {
	Create(ctx context.Context, entity *entity.Media) (*entity.Media, error)
	Delete(ctx context.Context, entity *entity.Media) error
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]entity.Media, error)
	GetByOptionID(ctx context.Context, optionID uuid.UUID) ([]entity.Media, error)
}

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) MediaRepository {
	return &mediaRepository{db: db}
}

func (repo *mediaRepository) Create(ctx context.Context, entity *entity.Media) (*entity.Media, error) {
	err := repo.db.WithContext(ctx).Create(entity).Error
	return entity, err
}

func (repo *mediaRepository) Delete(ctx context.Context, entity *entity.Media) error {
	return repo.db.WithContext(ctx).Delete(entity).Error
}

func (repo *mediaRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]entity.Media, error) {
	var media []entity.Media
	// Lấy tất cả media của product, sắp xếp theo sort_order
	err := repo.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("sort_order ASC").
		Find(&media).Error
	return media, err
}

func (repo *mediaRepository) GetByOptionID(ctx context.Context, optionID uuid.UUID) ([]entity.Media, error) {
	var media []entity.Media
	err := repo.db.WithContext(ctx).
		Where("option_id = ?", optionID).
		Order("sort_order ASC").
		Find(&media).Error
	return media, err
}
