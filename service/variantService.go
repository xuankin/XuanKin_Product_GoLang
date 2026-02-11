package service

import (
	"Product_Mangement_Api/entity"
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/repository"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
)

type VariantService interface {
	AddVariant(ctx context.Context, req models.CreateVariantRequest) (*models.VariantResponse, error)
	UpdateVariant(ctx context.Context, id uuid.UUID, req models.UpdateVariantRequest) (*models.VariantResponse, error)
	DeleteVariant(ctx context.Context, id uuid.UUID) error
	GetVariantByID(ctx context.Context, id uuid.UUID) (*models.VariantResponse, error)
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

func (s *variantService) syncParentProduct(ctx context.Context, productID uuid.UUID) {

	s.cacheRepo.Delete(ctx, models.CacheKeyProductDetail+productID.String())
	s.cacheRepo.DeleteByPrefix(ctx, models.CacheKeyProductList)

	go func(pID uuid.UUID) {

		bgCtx := context.Background()

		fullProduct, err := s.productRepo.GetById(bgCtx, pID)
		if err != nil {
			log.Printf("[SyncES] Failed to get product %s: %v", pID, err)
			return
		}

		minPrice := 0.0
		maxPrice := 0.0
		if len(fullProduct.Variants) > 0 {
			minPrice = fullProduct.Variants[0].Price
			maxPrice = fullProduct.Variants[0].Price
			for _, v := range fullProduct.Variants {
				if v.Price < minPrice {
					minPrice = v.Price
				}
				if v.Price > maxPrice {
					maxPrice = v.Price
				}
			}
		}

		attrMap := make(map[string]map[string]bool)
		for _, v := range fullProduct.Variants {
			for _, attrVal := range v.Attributes {

				nameMap := toMap(attrVal.Attribute.Name)
				valMap := toMap(attrVal.Value)

				attrName := fmt.Sprintf("%v", nameMap["vi"])
				attrValue := fmt.Sprintf("%v", valMap["vi"])

				if attrName == "<nil>" || attrValue == "<nil>" {
					continue
				}

				if attrMap[attrName] == nil {
					attrMap[attrName] = make(map[string]bool)
				}
				attrMap[attrName][attrValue] = true
			}
		}

		var esAttrs []models.EsAttributeSummary
		for k, vMap := range attrMap {
			var vals []string
			for v := range vMap {
				vals = append(vals, v)
			}
			esAttrs = append(esAttrs, models.EsAttributeSummary{Name: k, Values: vals})
		}

		primaryImg := ""
		if len(fullProduct.Media) > 0 {
			primaryImg = fullProduct.Media[0].URL
		}

		esProduct := &models.EsProductIndex{
			ID:                fullProduct.ID,
			Name:              toMap(fullProduct.Name),
			Slug:              fullProduct.Slug,
			Description:       toMap(fullProduct.Description),
			CategoryID:        fullProduct.CategoryID,
			CategoryName:      toMap(fullProduct.Category.Name),
			BrandID:           fullProduct.BrandID,
			BrandName:         toMap(fullProduct.Brand.Name),
			Status:            fullProduct.Status,
			MinPrice:          minPrice,
			MaxPrice:          maxPrice,
			AttributesSummary: esAttrs,
			PrimaryImage:      primaryImg,
			CreatedAt:         fullProduct.CreatedAt,
		}

		if err := s.esRepo.IndexProduct(bgCtx, esProduct); err != nil {
			log.Printf("[SyncES] Failed to index product %s: %v", pID, err)
		} else {
			log.Printf("[SyncES] Indexed product %s successfully", pID)
		}

	}(productID)
}

func (s *variantService) AddVariant(ctx context.Context, req models.CreateVariantRequest) (*models.VariantResponse, error) {
	_, err := s.productRepo.GetById(ctx, req.ProductID)
	if err != nil {
		return nil, errors.New("product does not exist")
	}

	variant := &entity.ProductVariant{
		ProductID: req.ProductID,
		SKU:       req.SKU,
		Price:     req.Price,
		SalePrice: req.SalePrice,
		Weight:    req.Weight,
		Status:    models.StatusActive,
	}

	for _, attrID := range req.AttributeIds {

		variant.Attributes = append(variant.Attributes, entity.AttributeValue{
			Base: entity.Base{ID: attrID},
		})
	}

	for _, stockReq := range req.InitialStocks {
		variant.Inventories = append(variant.Inventories, entity.Inventory{
			WarehouseID: stockReq.WarehouseID,
			Quantity:    stockReq.Quantity,
		})
	}

	if err := s.repo.Create(ctx, variant); err != nil {
		return nil, err
	}

	s.syncParentProduct(ctx, req.ProductID)

	savedVariant, _ := s.repo.GetByID(ctx, variant.ID)
	return s.mapToResponse(savedVariant), nil
}

func (s *variantService) UpdateVariant(ctx context.Context, id uuid.UUID, req models.UpdateVariantRequest) (*models.VariantResponse, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("variant does not exist")
	}

	if req.Price != nil {
		v.Price = *req.Price
	}
	if req.SalePrice != nil {
		v.SalePrice = *req.SalePrice
	}
	if req.Weight != nil {
		v.Weight = *req.Weight
	}
	if req.Status != "" {
		v.Status = req.Status
	}

	if err := s.repo.Update(ctx, v); err != nil {
		return nil, err
	}

	s.syncParentProduct(ctx, v.ProductID)

	updatedV, _ := s.repo.GetByID(ctx, id)
	return s.mapToResponse(updatedV), nil
}

func (s *variantService) DeleteVariant(ctx context.Context, id uuid.UUID) error {
	variant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("variant does not exist")
	}

	totalStock := 0
	for _, inv := range variant.Inventories {
		totalStock += inv.Quantity
	}

	if totalStock > 0 {
		return errors.New("variations cannot be removed because they are still in stock")
	}

	productID := variant.ProductID

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Trigger Sync
	s.syncParentProduct(ctx, productID)

	return nil
}

func (s *variantService) GetVariantByID(ctx context.Context, id uuid.UUID) (*models.VariantResponse, error) {
	variant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("variant does not exist")
	}
	return s.mapToResponse(variant), nil
}

func (s *variantService) mapToResponse(v *entity.ProductVariant) *models.VariantResponse {
	res := &models.VariantResponse{
		ID:        v.ID,
		SKU:       v.SKU,
		Price:     v.Price,
		SalePrice: v.SalePrice,
		Weight:    v.Weight,
		Status:    v.Status,
	}

	for _, attr := range v.Attributes {
		res.Attributes = append(res.Attributes, models.AttributeValueResponse{
			ID:    attr.ID,
			Value: toMap(attr.Value),
		})
	}
	return res
}
