package service

import (
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/repository"
	"context"
	"github.com/google/uuid"
)

type InventoryService interface {
	AdjustStock(ctx context.Context, req models.UpdateInventoryRequest) error
	GetStockByOption(ctx context.Context, optionId uuid.UUID) ([]models.InventoryResponse, error)
}
type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) AdjustStock(ctx context.Context, req models.UpdateInventoryRequest) error {
	return s.repo.UpdateStock(ctx, req.OptionID, req.WarehouseID, req.Amount, req.Type, req.Reason)
}

func (s *inventoryService) GetStockByOption(ctx context.Context, optionId uuid.UUID) ([]models.InventoryResponse, error) {
	invs, err := s.repo.GetByOptionId(ctx, optionId)
	if err != nil {
		return nil, err
	}
	var res []models.InventoryResponse
	for _, i := range invs {
		res = append(res, models.InventoryResponse{
			ID:               i.ID,
			OptionID:         i.OptionID,
			Quantity:         i.Quantity,
			ReservedQuantity: i.ReservedQuantity,
			Warehouse: models.WarehouseResponse{
				ID:      i.Warehouse.ID,
				Name:    toMap(i.Warehouse.Name),
				Address: i.Warehouse.Address,
				Phone:   i.Warehouse.Phone,
			},
		})
	}
	return res, nil
}
