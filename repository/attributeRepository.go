package repository

import (
	"Product_Mangement_Api/entity"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttributeRepository interface {
	CreateAttribute(ctx context.Context, attr *entity.Attribute) error
	ListAttributes(ctx context.Context) ([]entity.Attribute, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Attribute, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, attr *entity.Attribute) (*entity.Attribute, error)
}
type attributeRepository struct {
	db *gorm.DB
}

func NewAttributeRepository(db *gorm.DB) AttributeRepository {
	return &attributeRepository{db: db}
}

func (ar *attributeRepository) CreateAttribute(ctx context.Context, attr *entity.Attribute) error {
	return ar.db.WithContext(ctx).Create(attr).Error
}

func (ar *attributeRepository) ListAttributes(ctx context.Context) ([]entity.Attribute, error) {
	var attrs []entity.Attribute
	err := ar.db.WithContext(ctx).Find(&attrs).Error
	if err != nil {
		return nil, err
	}
	return attrs, nil
}

func (ar *attributeRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Attribute, error) {
	var attr entity.Attribute
	// Bỏ Preload("Values")
	err := ar.db.WithContext(ctx).First(&attr, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &attr, nil
}

func (ar *attributeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return ar.db.WithContext(ctx).Delete(&entity.Attribute{}, "id = ?", id).Error
}

func (ar *attributeRepository) Update(ctx context.Context, attr *entity.Attribute) (*entity.Attribute, error) {
	err := ar.db.WithContext(ctx).Save(attr).Error
	if err != nil {
		return nil, err
	}
	return attr, nil
}
