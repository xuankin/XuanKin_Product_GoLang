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
	UpdateStock(ctx context.Context, variantID uuid.UUID, warehouseID uuid.UUID, amount int, moveType string, reason string) error
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
		Joins("Warehouse").
		Where("variant_id = ?", variantId).
		Find(&inventories).Error
	return inventories, err
}
func (r *inventoryRepository) UpdateStock(ctx context.Context, variantID uuid.UUID, warehouseID uuid.UUID, amount int, moveType string, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inv entity.Inventory

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("variant_id = ? AND warehouse_id = ?", variantID, warehouseID).
			First(&inv).Error; err != nil {
			return err
		}

		change := amount
		if moveType == models.StockOut {
			change = -amount
		}

		newQuantity := inv.Quantity + change
		if newQuantity < 0 {
			return errors.New("Not enough inventory to complete the transaction")
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
