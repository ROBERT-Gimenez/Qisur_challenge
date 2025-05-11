package models

import "time"

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Products    []Product `json:"products" gorm:"many2many:product_categories"`
}


type ProductSummaryDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CategoryWithProductsDTO struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Products    []ProductSummaryDTO `json:"products"`
}
