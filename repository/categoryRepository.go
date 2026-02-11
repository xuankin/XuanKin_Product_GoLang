package repository

import (
	"Product_Mangement_Api/entity"
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, c *entity.Category) (*entity.Category, error)
	Update(ctx context.Context, c *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	List(ctx context.Context) ([]entity.Category, error)
}
type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}
func (r *categoryRepository) Create(ctx context.Context, c *entity.Category) (*entity.Category, error) {
	err := r.db.WithContext(ctx).Create(c).Error
	if err != nil {
		return nil, err
	}
	return c, nil

}
func (r *categoryRepository) Update(ctx context.Context, c *entity.Category) error {
	return r.db.WithContext(ctx).Model(c).Select("*").Updates(c).Error
}
func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var subCatCount int64
		if err := tx.Model(&entity.Category{}).Where("parent_id = ?", id).Count(&subCatCount).Error; err != nil {
			return err
		}

		if subCatCount > 0 {

			return errors.New("Can't delete a category that contains subcategories (please delete or move the subcategories first)")
		}

		var productCount int64
		if err := tx.Model(&entity.Product{}).Where("category_id = ?", id).Count(&productCount).Error; err != nil {
			return err
		}

		if productCount > 0 {
			return errors.New("Cannot delete categories that have linked products")
		}

		if err := tx.Delete(&entity.Category{}, "id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}
func (r *categoryRepository) GetById(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	var c entity.Category
	err := r.db.WithContext(ctx).First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *categoryRepository) List(ctx context.Context) ([]entity.Category, error) {
	var cats []entity.Category
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&cats).Error
	return cats, err
}
