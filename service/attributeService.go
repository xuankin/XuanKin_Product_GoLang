package service

import (
	"Product_Mangement_Api/entity"
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/repository"
	"context"

	"github.com/google/uuid"
)

type AttributeService interface {
	CreateAttribute(ctx context.Context, req models.CreateAttributeRequest) (*models.AttributeResponse, error)
	ListAttributes(ctx context.Context) ([]models.AttributeResponse, error)
	GetAttributesById(ctx context.Context, id uuid.UUID) (*models.AttributeResponse, error)
	DeleteAttribute(ctx context.Context, id uuid.UUID) error
	UpdatedAttribute(ctx context.Context, id uuid.UUID, req models.UpdateAttributeRequest) (*models.AttributeResponse, error)
}

type attributeService struct {
	repo repository.AttributeRepository
}

func NewAttributeService(repo repository.AttributeRepository) AttributeService {
	return &attributeService{repo: repo}
}

func (s *attributeService) mapToResponse(attr *entity.Attribute) *models.AttributeResponse {
	return &models.AttributeResponse{
		ID:           attr.ID,
		Name:         toMap(attr.Name),
		Type:         attr.Type,
		IsFilterable: attr.IsFilterable,
		IsRequired:   attr.IsRequired,
	}
}

func (s *attributeService) CreateAttribute(ctx context.Context, req models.CreateAttributeRequest) (*models.AttributeResponse, error) {
	attr := &entity.Attribute{
		Name:         toJson(req.Name),
		Type:         req.Type,
		IsFilterable: req.IsFilterable,
		IsRequired:   req.IsRequired,
	}
	if err := s.repo.CreateAttribute(ctx, attr); err != nil {
		return nil, err
	}
	return s.mapToResponse(attr), nil
}

func (s *attributeService) UpdatedAttribute(ctx context.Context, id uuid.UUID, req models.UpdateAttributeRequest) (*models.AttributeResponse, error) {
	attr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		attr.Name = toJson(req.Name)
	}

	if req.Type != "" {
		attr.Type = req.Type
	}
	if req.IsFilterable != nil {
		attr.IsFilterable = *req.IsFilterable
	}
	if req.IsRequired != nil {
		attr.IsRequired = *req.IsRequired
	}

	attrUpdated, err := s.repo.Update(ctx, attr)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(attrUpdated), nil
}

func (s *attributeService) ListAttributes(ctx context.Context) ([]models.AttributeResponse, error) {
	attributes, err := s.repo.ListAttributes(ctx)
	if err != nil {
		return nil, err
	}

	var res []models.AttributeResponse
	for _, attr := range attributes {
		res = append(res, *s.mapToResponse(&attr))
	}
	return res, nil
}

func (s *attributeService) GetAttributesById(ctx context.Context, id uuid.UUID) (*models.AttributeResponse, error) {
	attr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(attr), nil
}

func (s *attributeService) DeleteAttribute(ctx context.Context, id uuid.UUID) error {

	return s.repo.Delete(ctx, id)
}
