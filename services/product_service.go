package services

import (
	"qisur-challenge/models"
	"qisur-challenge/repository"
	"time"

	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(product *models.Product) error
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
	ConvertToProductDTO(product *models.Product) models.ProductDTO
	ConvertToProductDTOs(products []models.Product) []models.ProductDTO
	UpdateProduct(id uint, req *models.UpdateProductRequest) (*models.Product, error)
	DeleteProduct(product *models.Product) error
	GetProductHistory(id uint, start, end *time.Time) ([]models.ProductHistory, error)
	SearchProducts(name, sort string, page, limit int) ([]models.Product, error)
	SearchCategories(name, sort string, page, limit int) ([]models.Category, error)
}

type productService struct {
	productRepo repository.ProductRepository
	db          *gorm.DB
}

func NewProductService(db *gorm.DB) *productService {
	return &productService{
		productRepo: repository.NewProductRepository(db),
		db:          db,
	}
}

func (ps *productService) GetAllProducts() ([]models.Product, error) {
	return ps.productRepo.GetAll()
}

func (ps *productService) GetProductByID(id uint) (*models.Product, error) {
	return ps.productRepo.GetByID(id)
}

func (ps *productService) ConvertToProductDTO(product *models.Product) models.ProductDTO {
	categories := make([]models.CategoryDTO, len(product.Categories))
	for i, c := range product.Categories {
		categories[i] = models.CategoryDTO{
			ID:   c.ID,
			Name: c.Name,
		}
	}
	return models.ProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Categories:  categories,
	}
}

func (s *productService) ConvertToProductDTOs(products []models.Product) []models.ProductDTO {
	dtos := make([]models.ProductDTO, len(products))
	for i, product := range products {
		dtos[i] = s.ConvertToProductDTO(&product)
	}
	return dtos
}

func (ps *productService) CreateProduct(product *models.Product) error {
	return ps.productRepo.Create(product)
}

func (ps *productService) UpdateProduct(id uint, req *models.UpdateProductRequest) (*models.Product, error) {
	product, err := ps.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	original := *product

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if req.Categories != nil {
		if err := ps.productRepo.UpdateCategories(product, *req.Categories); err != nil {
			return nil, err
		}
	}

	if err := ps.db.Save(&product).Error; err != nil {
		return nil, err
	}

	if err := ps.productRepo.SaveHistory(&original); err != nil {
		return nil, err
	}

	return product, nil
}

func (ps *productService) DeleteProduct(product *models.Product) error {
	return ps.productRepo.Delete(product)
}

func (ps *productService) GetProductHistory(id uint, start, end *time.Time) ([]models.ProductHistory, error) {
	var history []models.ProductHistory
	query := ps.db.Where("product_id = ?", id)

	if start != nil {
		query = query.Where("changed_at >= ?", *start)
	}
	if end != nil {
		query = query.Where("changed_at <= ?", *end)
	}

	err := query.Find(&history).Error
	return history, err
}

func (ps *productService) SearchProducts(name, sort string, page, limit int) ([]models.Product, error) {
	db := ps.db.Model(&models.Product{})

	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}

	switch sort {
	case "price_asc":
		db = db.Order("price ASC")
	case "price_desc":
		db = db.Order("price DESC")
	}

	offset := (page - 1) * limit
	var products []models.Product
	err := db.Offset(offset).Limit(limit).Find(&products).Error
	return products, err
}

func (ps *productService) SearchCategories(name, sort string, page, limit int) ([]models.Category, error) {
	db := ps.db.Model(&models.Category{})

	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}

	switch sort {
	case "name_asc":
		db = db.Order("name ASC")
	case "name_desc":
		db = db.Order("name DESC")
	}

	offset := (page - 1) * limit
	var categories []models.Category
	err := db.Offset(offset).Limit(limit).Find(&categories).Error
	return categories, err
}
