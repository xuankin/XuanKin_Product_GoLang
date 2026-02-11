package service

import (
	"Product_Mangement_Api/entity"
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/repository"
	"context"
	"errors"
	"github.com/google/uuid"
)

type WarehouseService interface {
	Create(ctx context.Context, req models.CreateWarehouseRequest) (*models.WarehouseResponse, error)
	Update(ctx context.Context, req models.UpdateWarehouseRequest, Id uuid.UUID) (*models.WarehouseResponse, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.WarehouseResponse, error)
	List(ctx context.Context) ([]models.WarehouseResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
type warehouseService struct {
	repo repository.WarehouseRepository
}

func NewWarehouseService(repo repository.WarehouseRepository) WarehouseService {
	return &warehouseService{repo: repo}
}
func (s *warehouseService) Create(ctx context.Context, req models.CreateWarehouseRequest) (*models.WarehouseResponse, error) {
	w := &entity.Warehouse{
		Name:    toJson(req.Name),
		Address: req.Address,
		Phone:   req.Phone,
		Status:  models.StatusActive,
	}
	if err := s.repo.Create(ctx, w); err != nil {
		return nil, err
	}
	return &models.WarehouseResponse{
		ID:      w.ID,
		Name:    toMap(w.Name),
		Address: w.Address,
		Phone:   w.Phone,
		Status:  w.Status,
	}, nil
}
func (s *warehouseService) List(ctx context.Context) ([]models.WarehouseResponse, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	var res []models.WarehouseResponse
	for _, w := range list {
		res = append(res, models.WarehouseResponse{
			ID:      w.ID,
			Name:    toMap(w.Name),
			Address: w.Address,
			Phone:   w.Phone,
			Status:  w.Status,
		})
	}
	return res, nil
}
func (s *warehouseService) Update(ctx context.Context, req models.UpdateWarehouseRequest, Id uuid.UUID) (*models.WarehouseResponse, error) {
	w, err := s.repo.GetByID(ctx, Id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		w.Name = toJson(req.Name)
	}
	if req.Address != "" {
		w.Address = req.Address
	}
	if req.Status != "" {
		w.Status = req.Status
	}
	if req.Phone != "" {
		w.Phone = req.Phone
	}
	if err := s.repo.Update(ctx, w); err != nil {
		return nil, err
	}
	return &models.WarehouseResponse{
		ID:      w.ID,
		Name:    toMap(w.Name),
		Address: w.Address,
		Phone:   w.Phone,
		Status:  w.Status,
	}, nil
}
func (s *warehouseService) GetById(ctx context.Context, id uuid.UUID) (*models.WarehouseResponse, error) {
	w, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &models.WarehouseResponse{
		ID:      w.ID,
		Name:    toMap(w.Name),
		Address: w.Address,
		Phone:   w.Phone,
		Status:  w.Status,
	}, nil
}
func (s *warehouseService) Delete(ctx context.Context, id uuid.UUID) error {
	w, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if len(w.Inventories) > 0 {
		return errors.New("Cannot delete warehouse containing products")
	}
	return s.repo.Delete(ctx, w.ID)
}
