package service

import (
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/repository"
	"context"
	"github.com/google/uuid"
)

type InventoryService interface {
	AdjustStock(ctx context.Context, req models.UpdateInventoryRequest) error
	GetStockByVariant(ctx context.Context, variantId uuid.UUID) ([]models.InventoryResponse, error)
}
type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) AdjustStock(ctx context.Context, req models.UpdateInventoryRequest) error {
	return s.repo.UpdateStock(ctx, req.VariantID, req.WarehouseID, req.Amount, req.Type, req.Reason)
}
func (s *inventoryService) GetStockByVariant(ctx context.Context, variantId uuid.UUID) ([]models.InventoryResponse, error) {
	invs, err := s.repo.GetByVariantId(ctx, variantId)
	if err != nil {
		return nil, err
	}
	var res []models.InventoryResponse
	for _, i := range invs {
		res = append(res, models.InventoryResponse{
			ID:               i.ID,
			VariantID:        i.VariantID,
			Quantity:         i.Quantity,
			ReservedQuantity: i.ReservedQuantity,
			Warehouse: models.WarehouseResponse{
				ID:   i.Warehouse.ID,
				Name: toMap(i.Warehouse.Name), Address: i.Warehouse.Address, Phone: i.Warehouse.Phone,
			},
		})
	}
	return res, nil
}
