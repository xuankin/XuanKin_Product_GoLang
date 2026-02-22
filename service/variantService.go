package service

import (
	"Product_Mangement_Api/entity"
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/repository"
	"context"
	"errors"
	"github.com/google/uuid"
)

type VariantService interface {
	AddVariant(ctx context.Context, req models.CreateVariantRequest) (*models.VariantResponse, error)
	GetVariantByID(ctx context.Context, id uuid.UUID) (*models.VariantResponse, error)
	UpdateVariant(ctx context.Context, id uuid.UUID, req models.UpdateVariantRequest) (*models.VariantResponse, error)
	DeleteVariant(ctx context.Context, id uuid.UUID) error
	AddOption(ctx context.Context, variantID uuid.UUID, req models.VariantOptionRequest) error
	UpdateOption(ctx context.Context, optionID uuid.UUID, req models.UpdateVariantOptionRequest) error
}

type variantService struct {
	repo        repository.VariantRepository
	productRepo repository.ProductRepository
	esRepo      repository.ElasticsearchRepository
	cacheRepo   repository.CacheRepository
}

func NewVariantService(
	repo repository.VariantRepository,
	productRepo repository.ProductRepository,
	esRepo repository.ElasticsearchRepository,
	cacheRepo repository.CacheRepository,
) VariantService {
	return &variantService{
		repo:        repo,
		productRepo: productRepo,
		esRepo:      esRepo,
		cacheRepo:   cacheRepo,
	}
}

func (s *variantService) AddVariant(ctx context.Context, req models.CreateVariantRequest) (*models.VariantResponse, error) {
	_, err := s.productRepo.GetById(ctx, req.ProductID)
	if err != nil {
		return nil, errors.New("product does not exist")
	}

	variant := &entity.ProductVariant{
		ProductID: req.ProductID,
		Code:      req.Code,
		Name:      toJson(req.Name),
		Status:    models.StatusActive,
	}

	for _, optReq := range req.Options {
		option := entity.VariantOption{
			SKU:       optReq.SKU,
			Price:     optReq.Price,
			SalePrice: optReq.SalePrice,
			Weight:    optReq.Weight,
			Status:    models.StatusActive,
		}

		for _, valReq := range optReq.Values {
			option.Values = append(option.Values, entity.VariantOptionValue{
				Name:      valReq.Name,
				Value:     valReq.Value,
				SortOrder: valReq.SortOrder,
			})
		}

		for _, stockReq := range optReq.Inventories {
			option.Inventories = append(option.Inventories, entity.Inventory{
				WarehouseID: stockReq.WarehouseID,
				Quantity:    stockReq.Quantity,
			})
		}

		variant.Options = append(variant.Options, option)
	}

	if err := s.repo.Create(ctx, variant); err != nil {
		return nil, err
	}

	s.cacheRepo.Delete(ctx, models.CacheKeyProductDetail+req.ProductID.String())
	s.cacheRepo.DeleteByPrefix(ctx, models.CacheKeyProductList)

	savedVariant, _ := s.repo.GetByID(ctx, variant.ID)
	return s.mapToResponse(savedVariant), nil
}

func (s *variantService) GetVariantByID(ctx context.Context, id uuid.UUID) (*models.VariantResponse, error) {
	variant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("variant does not exist")
	}
	return s.mapToResponse(variant), nil
}

func (s *variantService) UpdateVariant(ctx context.Context, id uuid.UUID, req models.UpdateVariantRequest) (*models.VariantResponse, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("variant does not exist")
	}

	if req.Code != "" {
		v.Code = req.Code
	}
	if req.Name != nil {
		v.Name = toJson(req.Name)
	}
	if req.Status != "" {
		v.Status = req.Status
	}

	if err := s.repo.Update(ctx, v); err != nil {
		return nil, err
	}

	s.cacheRepo.Delete(ctx, models.CacheKeyProductDetail+v.ProductID.String())
	s.cacheRepo.DeleteByPrefix(ctx, models.CacheKeyProductList)

	updatedVariant, _ := s.repo.GetByID(ctx, id)
	return s.mapToResponse(updatedVariant), nil
}

func (s *variantService) DeleteVariant(ctx context.Context, id uuid.UUID) error {

	variant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("variant not found")
	}

	for _, opt := range variant.Options {
		for _, inv := range opt.Inventories {
			if inv.Quantity > 0 {
				return errors.New("cannot delete variant. Some options still have stock")
			}
		}
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.cacheRepo.Delete(ctx, models.CacheKeyProductDetail+variant.ProductID.String())
	s.cacheRepo.DeleteByPrefix(ctx, models.CacheKeyProductList)

	return nil
}

func (s *variantService) mapToResponse(v *entity.ProductVariant) *models.VariantResponse {
	res := &models.VariantResponse{
		ID:     v.ID,
		Code:   v.Code,
		Name:   toMap(v.Name),
		Status: v.Status,
	}

	for _, opt := range v.Options {
		optRes := models.VariantOptionResponse{
			ID:        opt.ID,
			SKU:       opt.SKU,
			Price:     opt.Price,
			SalePrice: opt.SalePrice,
			Weight:    opt.Weight,
			Status:    opt.Status,
		}
		for _, val := range opt.Values {
			optRes.Values = append(optRes.Values, models.VariantOptionValueResponse{
				ID:        val.ID,
				Name:      val.Name,
				Value:     val.Value,
				SortOrder: val.SortOrder,
			})
		}

		for _, inv := range opt.Inventories {
			optRes.Inventories = append(optRes.Inventories, models.InventoryResponse{
				ID:       inv.ID,
				OptionID: inv.OptionID,
				Quantity: inv.Quantity,
				Warehouse: models.WarehouseResponse{
					ID:   inv.Warehouse.ID,
					Name: toMap(inv.Warehouse.Name),
				},
			})
		}
	}
	return res
}
func (s *variantService) AddOption(ctx context.Context, variantID uuid.UUID, req models.VariantOptionRequest) error {

	v, err := s.repo.GetByID(ctx, variantID)
	if err != nil || v == nil {
		return errors.New("variant does not exist")
	}

	option := &entity.VariantOption{
		VariantID: variantID,
		SKU:       req.SKU,
		Price:     req.Price,
		SalePrice: req.SalePrice,
		Weight:    req.Weight,
		Status:    models.StatusActive,
	}

	for _, valReq := range req.Values {
		option.Values = append(option.Values, entity.VariantOptionValue{
			Name:      valReq.Name,
			Value:     valReq.Value,
			SortOrder: valReq.SortOrder,
		})
	}

	for _, stockReq := range req.Inventories {
		option.Inventories = append(option.Inventories, entity.Inventory{
			WarehouseID: stockReq.WarehouseID,
			Quantity:    stockReq.Quantity,
		})
	}

	if err := s.repo.CreateOption(ctx, option); err != nil {
		return err
	}

	s.cacheRepo.Delete(ctx, models.CacheKeyProductDetail+v.ProductID.String())
	s.cacheRepo.DeleteByPrefix(ctx, models.CacheKeyProductList)
	return nil
}

func (s *variantService) UpdateOption(ctx context.Context, optionID uuid.UUID, req models.UpdateVariantOptionRequest) error {

	option, err := s.repo.GetOptionByID(ctx, optionID)
	if err != nil {
		return errors.New("option does not exist")
	}

	if req.SKU != "" {
		option.SKU = req.SKU
	}
	if req.Price > 0 {
		option.Price = req.Price
	}
	if req.SalePrice >= 0 {
		option.SalePrice = req.SalePrice
	}
	if req.Weight > 0 {
		option.Weight = req.Weight
	}
	if req.Status != "" {
		option.Status = req.Status
	}

	if err := s.repo.UpdateOption(ctx, option); err != nil {
		return err
	}

	s.cacheRepo.DeleteByPrefix(ctx, models.CacheKeyProductList)
	return nil
}
