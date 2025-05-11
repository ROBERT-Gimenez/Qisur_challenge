package main

import (
	"log"
	"net/http"
	"os"

	"qisur-challenge/config"
	"qisur-challenge/routes"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar .env, se usarán variables de entorno existentes.")
	}

	config.LoadConfig()

	db, err := config.CONNECTDB()
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	config.AutoMigrate(db)

	r := routes.RegisterRoutes(db)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en puerto %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
