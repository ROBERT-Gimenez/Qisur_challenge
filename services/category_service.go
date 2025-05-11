package services

import (
	"qisur-challenge/models"
	"qisur-challenge/repository"

	"gorm.io/gorm"
)

type CategoryService interface {
	GetAllCategories() ([]models.Category, error)
	GetCategoryByID(id uint) (*models.Category, error)
	ConvertToCategoryDTO(category *models.Category) models.CategoryWithProductsDTO
	ConvertToCategoryWithProductsDTOs(categories []models.Category) []models.CategoryWithProductsDTO
	CreateCategory(category *models.Category) error
	UpdateCategory(category *models.Category) error
	DeleteCategory(category *models.Category) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
	db           *gorm.DB
}

func NewCategoryService(db *gorm.DB) CategoryService {
	return &categoryService{
		categoryRepo: repository.NewCategoryRepository(db),
		db:           db,
	}
}

func (s *categoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *categoryService) GetCategoryByID(id uint) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *categoryService) ConvertToCategoryDTO(category *models.Category) models.CategoryWithProductsDTO {
	products := make([]models.ProductSummaryDTO, len(category.Products))
	for i, p := range category.Products {
		products[i] = models.ProductSummaryDTO{
			ID:   p.ID,
			Name: p.Name,
		}
	}
	return models.CategoryWithProductsDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Products:    products,
	}
}

func (s *categoryService) ConvertToCategoryWithProductsDTOs(categories []models.Category) []models.CategoryWithProductsDTO {
	dtos := make([]models.CategoryWithProductsDTO, len(categories))
	for i, c := range categories {
		dtos[i] = s.ConvertToCategoryDTO(&c)
	}
	return dtos
}

func (s *categoryService) CreateCategory(category *models.Category) error {
	return s.categoryRepo.Create(category)
}

func (s *categoryService) UpdateCategory(category *models.Category) error {
	return s.categoryRepo.Update(category)
}

func (s *categoryService) DeleteCategory(category *models.Category) error {
	return s.categoryRepo.Delete(category)
}
