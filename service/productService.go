package service

import (
	"Product_Mangement_Api/entity"
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/repository"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type ProductService interface {
	Create(ctx context.Context, req models.CreateProductRequest) (*models.ProductResponse, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.ProductResponse, error)
	UpdateById(ctx context.Context, id uuid.UUID, req models.UpdateProductRequest) (*models.ProductResponse, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, limit int) (*models.PaginationResponse, error)
	Search(ctx context.Context, params models.FilterParams) (*models.PaginationResponse, error)
	SyncProductAttributes(ctx context.Context, productID uuid.UUID, attributes []models.ProductAttributeRequest) error
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

func (s *productService) generateSlug(name map[string]interface{}) string {
	var str string
	if val, ok := name["vi"]; ok && val != nil {
		str = fmt.Sprintf("%v", val)
	} else if val, ok := name["en"]; ok && val != nil {
		str = fmt.Sprintf("%v", val)
	} else {
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

	for _, pa := range p.ProductAttributes {
		var values []interface{}
		for _, pav := range pa.Values {
			values = append(values, toMap(pav.Value))
		}
		attrName := map[string]interface{}{"vi": "Unknown"}
		if pa.Attribute.ID != uuid.Nil {
			attrName = toMap(pa.Attribute.Name)
		}

		res.Attributes = append(res.Attributes, models.ProductAttributeDetail{
			AttributeID:   pa.AttributeID,
			AttributeName: attrName,
			Values:        values,
		})
	}

	for _, v := range p.Variants {
		vRes := models.VariantResponse{
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
			vRes.Options = append(vRes.Options, optRes)
		}
		res.Variants = append(res.Variants, vRes)
	}
	return res
}

func (s *productService) syncToElasticsearch(p *entity.Product) {
	go func(product entity.Product) {
		ctx := context.Background()
		minPrice := 0.0
		maxPrice := 0.0
		first := true

		for _, v := range product.Variants {
			for _, opt := range v.Options {
				if first {
					minPrice = opt.Price
					maxPrice = opt.Price
					first = false
				} else {
					if opt.Price < minPrice {
						minPrice = opt.Price
					}
					if opt.Price > maxPrice {
						maxPrice = opt.Price
					}
				}
			}
		}

		primaryImg := ""
		if len(product.Media) > 0 {
			primaryImg = product.Media[0].URL
		}

		var attributesSummary []models.EsAttributeSummary
		for _, pa := range product.ProductAttributes {

			attrNameMap := toMap(pa.Attribute.Name)
			attrName := ""

			if val, ok := attrNameMap["vi"]; ok && val != nil {
				attrName = fmt.Sprintf("%v", val)
			} else if val, ok := attrNameMap["en"]; ok && val != nil {
				attrName = fmt.Sprintf("%v", val)
			}

			var attrValues []string
			for _, pav := range pa.Values {
				valMap := toMap(pav.Value)
				if val, ok := valMap["vi"]; ok && val != nil {
					attrValues = append(attrValues, fmt.Sprintf("%v", val))
				} else if val, ok := valMap["en"]; ok && val != nil {
					attrValues = append(attrValues, fmt.Sprintf("%v", val))
				}
			}

			if attrName != "" || len(attrValues) > 0 {
				attributesSummary = append(attributesSummary, models.EsAttributeSummary{
					Name:   attrName,
					Values: attrValues,
				})
			}
		}

		esProduct := &models.EsProductIndex{
			ID:           product.ID,
			Name:         toMap(product.Name),
			Slug:         product.Slug,
			Description:  toMap(product.Description),
			CategoryID:   product.CategoryID,
			CategoryName: toMap(product.Category.Name),
			BrandID:      product.BrandID,
			BrandName:    toMap(product.Brand.Name),
			Status:       product.Status,
			MinPrice:     minPrice,
			MaxPrice:     maxPrice,
			PrimaryImage: primaryImg,
			CreatedAt:    product.CreatedAt,

			// Gán biến attributesSummary đúng chuẩn
			AttributesSummary: attributesSummary,
		}

		_ = s.esRepo.IndexProduct(ctx, esProduct)
	}(*p)
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
			Code:   vReq.Code,
			Name:   toJson(vReq.Name),
			Status: models.StatusActive,
		}

		for _, optReq := range vReq.Options {
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
		product.Variants = append(product.Variants, variant)
	}

	for _, attrReq := range req.Attributes {
		pa := entity.ProductAttribute{
			AttributeID: attrReq.AttributeID,
		}

		for _, val := range attrReq.Values {
			pa.Values = append(pa.Values, entity.ProductAttributeValue{
				Value: toJson(val),
			})
		}

		product.ProductAttributes = append(product.ProductAttributes, pa)
	}

	res, err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	s.cacheRepo.DeleteByPrefix(ctx, models.CacheKeyProductList)
	fullProduct, _ := s.repo.GetById(ctx, res.ID)
	if fullProduct != nil {
		s.syncToElasticsearch(fullProduct)
	}

	return s.mapToProductResponse(res), nil
}

func (s *productService) GetById(ctx context.Context, id uuid.UUID) (*models.ProductResponse, error) {
	cacheKey := models.CacheKeyProductDetail + id.String()
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
	if req.CategoryID != uuid.Nil {
		p.CategoryID = req.CategoryID
	}
	if req.BrandID != uuid.Nil {
		p.BrandID = req.BrandID
	}
	if req.Status != "" {
		p.Status = req.Status
	}
	if req.Name != nil {
		p.Name = toJson(req.Name)
	}
	if req.Description != nil {
		p.Description = toJson(req.Description)
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}

	s.cacheRepo.Delete(ctx, models.CacheKeyProductDetail+id.String())
	s.cacheRepo.DeleteByPrefix(ctx, models.CacheKeyProductList)
	fullProduct, _ := s.repo.GetById(ctx, id)
	if fullProduct != nil {
		s.syncToElasticsearch(fullProduct)
	}

	return s.mapToProductResponse(p), nil
}

func (s *productService) DeleteById(ctx context.Context, id uuid.UUID) error {

	p, err := s.repo.GetById(ctx, id)
	if err != nil {
		return err
	}

	for _, variant := range p.Variants {
		for _, option := range variant.Options {
			for _, inv := range option.Inventories {
				if inv.Quantity > 0 {
					return errors.New("cannot delete product. Some options still have stock in warehouses")
				}
			}
		}
	}

	s.cacheRepo.Delete(ctx, models.CacheKeyProductDetail+id.String())
	s.cacheRepo.DeleteByPrefix(ctx, models.CacheKeyProductList)

	err = s.repo.Delete(ctx, id)
	if err == nil {
		go func() { s.esRepo.DeleteProduct(context.Background(), id.String()) }()
	}
	return err
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
	for i := range products {
		data = append(data, *s.mapToProductResponse(&products[i]))
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

func (s *productService) Search(ctx context.Context, params models.FilterParams) (*models.PaginationResponse, error) {
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
func (s *productService) SyncProductAttributes(ctx context.Context, productID uuid.UUID, reqAttributes []models.ProductAttributeRequest) error {

	_, err := s.repo.GetById(ctx, productID)
	if err != nil {
		return errors.New("product does not exist")
	}

	var newAttributes []entity.ProductAttribute
	for _, attrReq := range reqAttributes {
		pa := entity.ProductAttribute{
			ProductID:   productID,
			AttributeID: attrReq.AttributeID,
		}

		for _, val := range attrReq.Values {
			pa.Values = append(pa.Values, entity.ProductAttributeValue{
				Value: toJson(val),
			})
		}
		newAttributes = append(newAttributes, pa)
	}

	if err := s.repo.SyncAttributes(ctx, productID, newAttributes); err != nil {
		return fmt.Errorf("failed to sync attributes: %w", err)
	}

	s.cacheRepo.Delete(ctx, models.CacheKeyProductDetail+productID.String())
	s.cacheRepo.DeleteByPrefix(ctx, models.CacheKeyProductList)

	fullProduct, _ := s.repo.GetById(ctx, productID)
	if fullProduct != nil {
		s.syncToElasticsearch(fullProduct)
	}

	return nil
}
