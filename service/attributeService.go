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
	AddValue(ctx context.Context, req models.CreateAttributeValueRequest) (*models.AttributeValueResponse, error)
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
	res := &models.AttributeResponse{
		ID:   attr.ID,
		Name: toMap(attr.Name),
	}

	for _, v := range attr.Values {
		res.Values = append(res.Values, models.AttributeValueResponse{
			ID:            v.ID,
			AttributeID:   v.AttributeID,
			AttributeName: toMap(attr.Name),
			Value:         toMap(v.Value),
		})
	}
	return res
}
func (s *attributeService) CreateAttribute(ctx context.Context, req models.CreateAttributeRequest) (*models.AttributeResponse, error) {
	attr := &entity.Attribute{
		Name: toJson(req.Name),
	}
	if err := s.repo.CreateAttribute(ctx, attr); err != nil {
		return nil, err
	}
	return s.mapToResponse(attr), nil
}
func (s *attributeService) AddValue(ctx context.Context, req models.CreateAttributeValueRequest) (*models.AttributeValueResponse, error) {
	val := &entity.AttributeValue{
		AttributeID: req.AttributeID,
		Value:       toJson(req.Value),
	}
	if err := s.repo.CreateValue(ctx, val); err != nil {
		return nil, err
	}
	return &models.AttributeValueResponse{
		ID:          val.ID,
		AttributeID: val.AttributeID,
		Value:       toMap(val.Value),
	}, nil
}
func (s *attributeService) ListAttributes(ctx context.Context) ([]models.AttributeResponse, error) {
	attrs, err := s.repo.ListAttributes(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]models.AttributeResponse, 0, len(attrs))
	for _, v := range attrs {
		res = append(res, *s.mapToResponse(&v))
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
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
func (s *attributeService) UpdatedAttribute(ctx context.Context, id uuid.UUID, req models.UpdateAttributeRequest) (*models.AttributeResponse, error) {
	attr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		attr.Name = toJson(req.Name)
	}
	attrUpdated, err := s.repo.Update(ctx, attr)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(attrUpdated), nil
}
