package repository

import (
	"Product_Mangement_Api/entity"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VariantRepository interface {
	Create(ctx context.Context, v *entity.ProductVariant) error
	Update(ctx context.Context, v *entity.ProductVariant) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductVariant, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]entity.ProductVariant, error)
}
type variantRepository struct {
	db *gorm.DB
}

func NewVariantRepository(db *gorm.DB) VariantRepository {
	return &variantRepository{db: db}
}
func (s *variantRepository) Create(ctx context.Context, v *entity.ProductVariant) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(v).Error; err != nil {
			return err
		}
		for _, inv := range v.Inventories {
			if inv.Quantity > 0 {
				if err := tx.Create(&entity.StockMovement{
					InventoryID:  inv.ID,
					ChangeAmount: inv.Quantity,
					Type:         "IN",
					Reason:       "Initial stock for new variant",
				}).Error; err != nil {
					return err
				}

			}
		}
		return nil
	})
}
func (s *variantRepository) Update(ctx context.Context, v *entity.ProductVariant) error {
	return s.db.WithContext(ctx).Model(v).Updates(v).Error
}
func (s *variantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return s.db.WithContext(ctx).Delete(&entity.ProductVariant{}, "id = ?", id).Error
}
func (s *variantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductVariant, error) {
	var v entity.ProductVariant
	err := s.db.WithContext(ctx).
		Preload("Attributes").
		Preload("Inventories").
		First(&v, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}
func (s *variantRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]entity.ProductVariant, error) {
	var v []entity.ProductVariant
	err := s.db.WithContext(ctx).
		Preload("Attributes").
		Where("product_id = ?", productID).
		Find(&v).Error
	if err != nil {
		return nil, err
	}
	return v, nil
}
