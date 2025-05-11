package config

import (
	"fmt"
	"log"
	"os"
	"qisur-challenge/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CONNECTDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	log.Printf("DEBUG: host=%s user=%s pass=%s dbname=%s port=%s\n", host, user, pass, dbname, port)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Conexión a la base de datos exitosa")
	return db, nil
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Product{},
		&models.Category{},
		&models.ProductHistory{},
	)
	if err != nil {
		log.Printf("Error al migrar modelos: %v\n", err)
	} else {
		log.Println("Migraciones completadas.")
	}
}
