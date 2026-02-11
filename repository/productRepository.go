package repository

import (
	"Product_Mangement_Api/entity"
	"Product_Mangement_Api/models"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, p *entity.Product) (*entity.Product, error)
	Update(ctx context.Context, p *entity.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	List(ctx context.Context, offset, limit int) ([]entity.Product, int64, error)
}
type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}
func (r *productRepository) Create(ctx context.Context, p *entity.Product) (*entity.Product, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(p).Error; err != nil {
			return err
		}

		for _, variant := range p.Variants {
			for _, inv := range variant.Inventories {

				if inv.Quantity > 0 {
					movement := entity.StockMovement{
						InventoryID:  inv.ID,
						ChangeAmount: inv.Quantity,
						Type:         models.StockIn,
						Reason:       "Initialize opening balance",
					}
					if err := tx.Create(&movement).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
	return p, err
}
func (r *productRepository) Update(ctx context.Context, p *entity.Product) error {
	return r.db.WithContext(ctx).Model(p).Save(p).Error
}
func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, "id = ?", id).Error
}
func (r *productRepository) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	var p entity.Product

	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
		Preload("Media").
		Preload("Variants.Attributes").
		Preload("Variants.Inventories.Warehouse").
		Preload("Variants.Attributes.Attribute").
		First(&p, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}
func (r *productRepository) List(ctx context.Context, offset, limit int) ([]entity.Product, int64, error) {
	var products []entity.Product
	var count int64

	db := r.db.WithContext(ctx).Model(&entity.Product{})

	db.Count(&count)

	err := db.Order("created_at DESC").
		Offset(offset).Limit(limit).
		Preload("Category").
		Preload("Brand").
		Preload("Media", "is_primary=?", true).
		Preload("Variants.Attributes.Attribute").
		Find(&products).Error

	return products, count, err
}
