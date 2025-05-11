package routes

import (
	"qisur-challenge/controllers"
	"qisur-challenge/services"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CategoriesRoutes(db *gorm.DB, api *mux.Router) {
	categorieService := services.NewCategoryService(db)
	categoriesController := controllers.NewCategoriesController(db, categorieService)
	//rutas publicas
	api.HandleFunc("/categories", categoriesController.GetCategories).Methods("GET")

	//rutas protegidas
	ApplyMiddlewareRoute(api, "/categories/{id}", categoriesController.GetCategory, "GET")
	ApplyMiddlewareRoute(api, "/categories", categoriesController.CreateCategory, "POST")
	ApplyMiddlewareRoute(api, "/categories/{id}", categoriesController.UpdateCategory, "PUT")
	ApplyMiddlewareRoute(api, "/categories/{id}", categoriesController.DeleteCategory, "DELETE")

}
