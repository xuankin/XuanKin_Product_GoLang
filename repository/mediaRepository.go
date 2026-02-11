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
	GetByProductID(ctx context.Context, productID uuid.UUID) (*entity.Media, error)
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
	err := repo.db.WithContext(ctx).Delete(entity).Error
	return err
}
func (repo *mediaRepository) GetByProductID(ctx context.Context, productID uuid.UUID) (*entity.Media, error) {
	var media entity.Media
	err := repo.db.WithContext(ctx).First(&media, "product_id = ?", productID).Error
	return &media, err
}
