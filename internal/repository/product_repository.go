package repository

import (
	"errors"
	"github.com/mdedealf/go-api/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type ProductRepository struct {
	Log *logrus.Logger
}

func NewProductRepository(log *logrus.Logger) *ProductRepository {
	return &ProductRepository{
		Log: log,
	}
}

func (p *ProductRepository) Save(db *gorm.DB, product *entity.Product) (*entity.Product, error) {
	err := db.Create(product).Error
	if err != nil {
		p.Log.Error(err)
		return nil, err
	}

	return product, nil
}

// GetByID implements ProductRepository
func (p *ProductRepository) GetByID(db *gorm.DB, id int64) (*entity.Product, error) {
	var product entity.Product
	err := db.Where("id = ? AND deleted_at IS NULL", id).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Log.WithField("id", id).Warn("product not found")
			return nil, errors.New("product not found")
		}
		p.Log.Error(err)
		return nil, err
	}
	return &product, nil
}

// Update implements ProductRepository
func (p *ProductRepository) Update(db *gorm.DB, product *entity.Product) (*entity.Product, error) {
	err := db.Where("id = ? AND deleted_at IS NULL", product.ID).Updates(product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Log.WithField("id", product.ID).Warn("product not found")
			return nil, errors.New("product not found")
		}
		p.Log.Error(err)
		return nil, err
	}

	// Fetch the updated product
	return p.GetByID(db, product.ID)
}

// Delete implements ProductRepository
func (p *ProductRepository) Delete(db *gorm.DB, id int64) error {
	result := db.Model(&entity.Product{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		p.Log.Error(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
