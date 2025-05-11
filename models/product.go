package models

import "time"

type Product struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Categories  []Category `json:"categories" gorm:"many2many:product_categories"`
}

type UpdateProductRequest struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Price       *float64  `json:"price"`
	Stock       *int      `json:"stock"`
	Categories  *[]uint   `json:"categories"` 
}

type ProductDTO struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Stock       int           `json:"stock"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Categories  []CategoryDTO `json:"categories"`
}

type CategoryDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}


type ProductCategory struct {
	ProductID  uint `gorm:"primaryKey"`
	CategoryID uint `gorm:"primaryKey"`
}
