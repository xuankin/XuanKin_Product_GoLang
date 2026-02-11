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
	"strings"
	"time"
)

type ProductService interface {
	Create(ctx context.Context, req models.CreateProductRequest) (*models.ProductResponse, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.ProductResponse, error)
	UpdateById(ctx context.Context, id uuid.UUID, req models.UpdateProductRequest) (*models.ProductResponse, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, limit int) (*models.PaginationResponse, error)
	Search(ctx context.Context, params models.FilterParams) (*models.PaginationResponse, error) // Thêm method Search
}

type productService struct {
	repo         repository.ProductRepository
	cacheRepo    repository.CacheRepository
	categoryRepo repository.CategoryRepository
	brandRepo    repository.BrandRepository
	esRepo       repository.ElasticsearchRepository
}

func NewProductService(
	repo repository.ProductRepository,
	brandRepo repository.BrandRepository,
	cacheRepo repository.CacheRepository,
	categoryRepo repository.CategoryRepository,
	esRepo repository.ElasticsearchRepository,
) ProductService {
	return &productService{
		repo:         repo,
		cacheRepo:    cacheRepo,
		categoryRepo: categoryRepo,
		brandRepo:    brandRepo,
		esRepo:       esRepo,
	}
}

func (s *productService) syncToElasticsearch(p *entity.Product) {

	go func(product entity.Product) {
		ctx := context.Background()

		minPrice := 0.0
		maxPrice := 0.0
		if len(product.Variants) > 0 {
			minPrice = product.Variants[0].Price
			maxPrice = product.Variants[0].Price
			for _, v := range product.Variants {
				if v.Price < minPrice {
					minPrice = v.Price
				}
				if v.Price > maxPrice {
					maxPrice = v.Price
				}
			}
		}

		attrMap := make(map[string]map[string]bool)
		for _, v := range product.Variants {
			for _, attrVal := range v.Attributes {

				nameMap := toMap(attrVal.Attribute.Name)
				valMap := toMap(attrVal.Value)

				attrName := fmt.Sprintf("%v", nameMap["vi"])
				attrVal := fmt.Sprintf("%v", valMap["vi"])

				if attrName == "<nil>" || attrVal == "<nil>" {
					continue
				}

				if attrMap[attrName] == nil {
					attrMap[attrName] = make(map[string]bool)
				}
				attrMap[attrName][attrVal] = true
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
		if len(product.Media) > 0 {
			primaryImg = product.Media[0].URL
		}

		esProduct := &models.EsProductIndex{
			ID:                product.ID,
			Name:              toMap(product.Name),
			Slug:              product.Slug,
			Description:       toMap(product.Description),
			CategoryID:        product.CategoryID,
			CategoryName:      toMap(product.Category.Name),
			BrandID:           product.BrandID,
			BrandName:         toMap(product.Brand.Name),
			Status:            product.Status,
			MinPrice:          minPrice,
			MaxPrice:          maxPrice,
			AttributesSummary: esAttrs,
			PrimaryImage:      primaryImg,
			CreatedAt:         product.CreatedAt,
		}

		if err := s.esRepo.IndexProduct(ctx, esProduct); err != nil {
			log.Printf("Failed to index product %s to ES: %v\n", product.ID, err)
		} else {
			log.Printf("Indexed product %s to ES successfully\n", product.ID)
		}
	}(*p)
}

func (s *productService) generateSlug(name map[string]interface{}) string {
	var str string
	if val, ok := name["en"]; ok && val != nil {
		str = fmt.Sprintf("%v", val)
	} else if val, ok := name["vi"]; ok && val != nil {
		str = fmt.Sprintf("%v", val)
	} else {
		for _, v := range name {
			if v != nil {
				str = fmt.Sprintf("%v", v)
				break
			}
		}
	}

	if str == "" {
		str = "product"
	}
	return strings.ToLower(strings.ReplaceAll(str, " ", "-")) + "-" + uuid.New().String()[:8]
}

func (s *productService) mapToProductResponse(p *entity.Product) *models.ProductResponse {
	res := &models.ProductResponse{
		ID:          p.ID,
		Name:        toMap(p.Name),
		Description: toMap(p.Description),
		Slug:        p.Slug,
		Status:      p.Status,
		Category: models.CategoryResponse{
			ID:   p.CategoryID,
			Name: toMap(p.Category.Name),
		},
		Brand: models.BrandResponse{
			ID:   p.BrandID,
			Name: toMap(p.Brand.Name),
			Logo: p.Brand.Logo,
		},
	}
	for _, v := range p.Variants {
		vRes := models.VariantResponse{
			ID:        v.ID,
			SKU:       v.SKU,
			Price:     v.Price,
			SalePrice: v.SalePrice,
			Weight:    v.Weight,
			Status:    v.Status,
		}
		for _, attr := range v.Attributes {
			vRes.Attributes = append(vRes.Attributes, models.AttributeValueResponse{
				ID:    attr.ID,
				Value: toMap(attr.Value),
			})
		}
		res.Variants = append(res.Variants, vRes)
	}
	return res
}

func (s *productService) Create(ctx context.Context, req models.CreateProductRequest) (*models.ProductResponse, error) {
	product := &entity.Product{
		Name:        toJson(req.Name),
		Description: toJson(req.Description),
		CategoryID:  req.CategoryID,
		BrandID:     req.BrandID,
		Status:      req.Status,
		Slug:        s.generateSlug(req.Name),
	}

	for _, vReq := range req.Variants {
		variant := entity.ProductVariant{
			SKU:       vReq.SKU,
			Price:     vReq.Price,
			SalePrice: vReq.SalePrice,
			Weight:    vReq.Weight,
			Status:    models.StatusActive,
		}
		for _, attrId := range vReq.AttributeIds {
			variant.Attributes = append(variant.Attributes, entity.AttributeValue{Base: entity.Base{ID: attrId}})
		}
		for _, sReq := range vReq.InitialStocks {
			variant.Inventories = append(variant.Inventories, entity.Inventory{
				WarehouseID: sReq.WarehouseID,
				Quantity:    sReq.Quantity,
			})
		}
		product.Variants = append(product.Variants, variant)
	}

	res, err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	// Invalidate Cache
	s.cacheRepo.Delete(ctx, models.CacheKeyProductDetail+res.ID.String())
	s.cacheRepo.DeleteByPrefix(ctx, models.CacheKeyProductList)

	fullProduct, err := s.repo.GetById(ctx, res.ID)
	if err == nil {
		s.syncToElasticsearch(fullProduct)
	}

	return s.mapToProductResponse(res), nil
}

func (s *productService) GetById(ctx context.Context, id uuid.UUID) (*models.ProductResponse, error) {
	cacheKey := "product:detail" + id.String()
	var res models.ProductResponse
	if err := s.cacheRepo.Get(ctx, cacheKey, &res); err == nil {
		return &res, nil
	}
	p, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	response := s.mapToProductResponse(p)
	s.cacheRepo.Set(ctx, cacheKey, response, 30*time.Minute)
	return response, nil
}

func (s *productService) UpdateById(ctx context.Context, id uuid.UUID, req models.UpdateProductRequest) (*models.ProductResponse, error) {
	p, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		p.Name = toJson(req.Name)
	}
	if req.Description != nil {
		p.Description = toJson(req.Description)
	}
	if req.CategoryID != uuid.Nil {
		if _, err := s.categoryRepo.GetById(ctx, req.CategoryID); err != nil {
			return nil, errors.New("Category ID does not exist")
		}
		p.CategoryID = req.CategoryID
	}
	if req.BrandID != uuid.Nil {

		if _, err := s.brandRepo.GetByID(ctx, req.BrandID); err != nil {
			return nil, errors.New("Brand ID does not exist")
		}
		p.BrandID = req.BrandID
	}
	if req.Status != "" {
		p.Status = req.Status
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}

	s.cacheRepo.Delete(ctx, models.CacheKeyProductDetail+id.String())
	s.cacheRepo.DeleteByPrefix(ctx, models.CacheKeyProductList)

	fullProduct, err := s.repo.GetById(ctx, id)
	if err == nil {
		s.syncToElasticsearch(fullProduct)
	}

	return s.mapToProductResponse(p), nil
}

func (s *productService) DeleteById(ctx context.Context, id uuid.UUID) error {

	s.cacheRepo.Delete(ctx, models.CacheKeyProductDetail+id.String())
	s.cacheRepo.DeleteByPrefix(ctx, models.CacheKeyProductList)

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	go func() {
		if err := s.esRepo.DeleteProduct(context.Background(), id.String()); err != nil {
			log.Printf("Failed to delete product %s from ES: %v", id.String(), err)
		}
	}()

	return nil
}

func (s *productService) List(ctx context.Context, page, limit int) (*models.PaginationResponse, error) {
	cacheKey := fmt.Sprintf("products:list:page:%d:limit:%d", page, limit)
	var res models.PaginationResponse
	if err := s.cacheRepo.Get(ctx, cacheKey, &res); err == nil {
		return &res, nil
	}

	offset := (page - 1) * limit
	products, total, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	var data []models.ProductResponse
	for _, p := range products {
		data = append(data, *s.mapToProductResponse(&p))
	}

	finalRes := &models.PaginationResponse{
		Total:       total,
		Data:        data,
		CurrentPage: page,
		LastPage:    int((total + int64(limit) - 1) / int64(limit)),
	}
	s.cacheRepo.Set(ctx, cacheKey, finalRes, 10*time.Minute)
	return finalRes, nil
}

// Search implements ProductService.
func (s *productService) Search(ctx context.Context, params models.FilterParams) (*models.PaginationResponse, error) {
	// Gọi xuống ES Repo
	products, total, err := s.esRepo.SearchProducts(ctx, params)
	if err != nil {
		return nil, err
	}

	return &models.PaginationResponse{
		Total:       total,
		Data:        products,
		CurrentPage: params.Page,
		LastPage:    int((total + int64(params.Limit) - 1) / int64(params.Limit)),
	}, nil
}
