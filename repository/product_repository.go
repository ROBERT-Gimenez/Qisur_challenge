package repository

import (
	"fmt"
	"log"
	"qisur-challenge/models"
	"time"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	GetByID(id uint) (*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(product *models.Product) error
	SaveHistory(product *models.Product) error
	UpdateCategories(product *models.Product, categoryIDs []uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Categories").Find(&products).Error
	if err != nil {
		log.Printf("ERROR: Falló al obtener productos: %v", err)
		return nil, err
	}
	return products, nil
}

func (r *productRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Categories", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).First(&product, id).Error; err != nil {
		log.Printf("ERROR: Falló al obtener producto con ID %d: %v", id, err)
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(product *models.Product) error {
	var existingProduct models.Product
	err := r.db.Where("name = ?", product.Name).First(&existingProduct).Error

	if err == nil {
		return fmt.Errorf("producto con nombre '%s' ya existe", product.Name)
	}

	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *models.Product) error {
	var existingProduct models.Product
	err := r.db.Where("name = ? AND id <> ?", product.Name, product.ID).First(&existingProduct).Error
	if err == nil {
		return fmt.Errorf("producto con nombre '%s' ya existe", product.Name)
	}
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(product *models.Product) error {
	if err := r.db.Model(product).Association("Categories").Clear(); err != nil {
		log.Printf("error al desasociar categorías: %v", err)
		return err
	}
	return r.db.Delete(product).Error
}

func (r *productRepository) UpdateCategories(product *models.Product, categoryIDs []uint) error {
	var categories []models.Category
	if len(categoryIDs) > 0 {
		if err := r.db.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
			return err
		}
	}
	return r.db.Model(product).Association("Categories").Replace(&categories)
}

func (r *productRepository) SaveHistory(product *models.Product) error {
	history := models.ProductHistory{
		ProductID: product.ID,
		Price:     product.Price,
		Stock:     product.Stock,
		ChangedAt: time.Now(),
	}
	return r.db.Create(&history).Error
}
