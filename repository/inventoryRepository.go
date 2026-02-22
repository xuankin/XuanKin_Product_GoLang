package repository

import (
	"Product_Mangement_Api/entity"
	"Product_Mangement_Api/models"
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InventoryRepository interface {
	GetByVariantId(ctx context.Context, variantId uuid.UUID) ([]entity.Inventory, error)
	GetByOptionId(ctx context.Context, optionId uuid.UUID) ([]entity.Inventory, error)
	UpdateStock(ctx context.Context, optionID uuid.UUID, warehouseID uuid.UUID, amount int, moveType string, reason string) error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *inventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetByVariantId(ctx context.Context, variantId uuid.UUID) ([]entity.Inventory, error) {
	var inventories []entity.Inventory
	err := r.db.WithContext(ctx).
		Joins("JOIN variant_options ON variant_options.id = inventories.option_id").
		Where("variant_options.variant_id = ?", variantId).
		Preload("Warehouse").
		Find(&inventories).Error
	return inventories, err
}

func (r *inventoryRepository) GetByOptionId(ctx context.Context, optionId uuid.UUID) ([]entity.Inventory, error) {
	var inventories []entity.Inventory
	err := r.db.WithContext(ctx).
		Where("option_id = ?", optionId).
		Preload("Warehouse").
		Find(&inventories).Error
	return inventories, err
}

func (r *inventoryRepository) UpdateStock(ctx context.Context, optionID uuid.UUID, warehouseID uuid.UUID, amount int, moveType string, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inv entity.Inventory

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("option_id = ? AND warehouse_id = ?", optionID, warehouseID).
			First(&inv).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {

				if moveType == models.StockOut {
					return errors.New("cannot out-stock: inventory record not found")
				}

				inv = entity.Inventory{
					OptionID:    optionID,
					WarehouseID: warehouseID,
					Quantity:    0,
				}
				if err := tx.Create(&inv).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		change := amount
		if moveType == models.StockOut {
			change = -amount
		}

		newQuantity := inv.Quantity + change
		if newQuantity < 0 {
			return errors.New("not enough inventory to complete the transaction")
		}

		if err := tx.Model(&inv).Update("quantity", newQuantity).Error; err != nil {
			return err
		}

		movement := entity.StockMovement{
			InventoryID:  inv.ID,
			ChangeAmount: change,
			Type:         moveType,
			Reason:       reason,
		}
		return tx.Create(&movement).Error
	})
}
