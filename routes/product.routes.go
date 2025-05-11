package routes

import (
	"qisur-challenge/controllers"
	"qisur-challenge/services"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ProductRoutes(db *gorm.DB, api *mux.Router) {
	productService := services.NewProductService(db)
	productController := controllers.NewProductController(db, productService)
	//rutas publicas
	api.HandleFunc("/products", productController.GetProducts).Methods("GET")
	api.HandleFunc("/search", productController.SearchHandler).Methods("GET")

	//rutas protegidas
	ApplyMiddlewareRoute(api, "/products/{id}", productController.GetProduct, "GET")
	ApplyMiddlewareRoute(api, "/products", productController.CreateProduct, "POST")
	ApplyMiddlewareRoute(api, "/products/{id}", productController.UpdateProduct, "PUT")
	ApplyMiddlewareRoute(api, "/products/{id}", productController.DeleteProduct, "DELETE")
	ApplyMiddlewareRoute(api, "/products/{id}/history", productController.GetProductHistory, "GET")

}
