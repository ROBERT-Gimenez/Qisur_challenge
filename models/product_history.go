package models

import "time"

type ProductHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	ChangedAt time.Time `json:"changed_at"`
}
