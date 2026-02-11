package repository

import (
	"Product_Mangement_Api/entity"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WarehouseRepository interface {
	Create(ctx context.Context, w *entity.Warehouse) error
	List(ctx context.Context) ([]entity.Warehouse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error)
	Update(ctx context.Context, w *entity.Warehouse) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type warehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &warehouseRepository{db: db}
}
func (repo *warehouseRepository) Create(ctx context.Context, w *entity.Warehouse) error {
	return repo.db.WithContext(ctx).Create(w).Error
}
func (repo *warehouseRepository) List(ctx context.Context) ([]entity.Warehouse, error) {
	var warehouses []entity.Warehouse
	err := repo.db.WithContext(ctx).Find(&warehouses).Error
	return warehouses, err
}
func (repo *warehouseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error) {
	var warehouse entity.Warehouse
	err := repo.db.WithContext(ctx).Preload("Inventories").First(&warehouse, "id = ?", id).Error
	return &warehouse, err
}
func (repo *warehouseRepository) Update(ctx context.Context, w *entity.Warehouse) error {
	return repo.db.WithContext(ctx).Save(w).Error
}
func (repo *warehouseRepository) Delete(ctx context.Context, id uuid.UUID) error {

	return repo.db.WithContext(ctx).Delete(&entity.Warehouse{}, "id = ?", id).Error
}
