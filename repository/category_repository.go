package repository

import (
	"fmt"
	"qisur-challenge/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetAll() ([]models.Category, error)
	GetByID(id uint) (*models.Category, error)
	Create(category *models.Category) error
	Update(category *models.Category) error
	Delete(category *models.Category) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Preload("Products").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.Preload("Products").First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Create(category *models.Category) error {
	var existingCategory models.Category
	err := r.db.Where("name = ?", category.Name).First(&existingCategory).Error

	if err == nil {
		return fmt.Errorf("la categoria con nombre '%s' ya existe", category.Name)
	}

	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *models.Category) error {
	var existingCategory models.Category
	err := r.db.Where("name = ? AND id <> ?", category.Name, category.ID).First(&existingCategory).Error
	if err == nil {
		return fmt.Errorf("la categoria con nombre '%s' ya existe", category.Name)
	}

	if err := r.db.Preload("Products").First(&category, category.ID).Error; err != nil {
		return err
	}
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(category *models.Category) error {
	if err := r.db.Model(category).Association("Products").Clear(); err != nil {
		return err
	}
	return r.db.Delete(category).Error
}
