package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"qisur-challenge/models"
	"qisur-challenge/services"
	websocket "qisur-challenge/webSocket"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type CategoriesController struct {
	CategoriesService services.CategoryService
	DB                *gorm.DB
}

func NewCategoriesController(db *gorm.DB, categoriesService services.CategoryService) *CategoriesController {
	return &CategoriesController{DB: db, CategoriesService: categoriesService}
}

func (sc *CategoriesController) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := sc.CategoriesService.GetAllCategories()
	if err != nil {
		http.Error(w, "Error al obtener categorías", http.StatusInternalServerError)
		return
	}
	categoryDTOs := sc.CategoriesService.ConvertToCategoryWithProductsDTOs(categories)
	json.NewEncoder(w).Encode(categoryDTOs)

}

func (sc *CategoriesController) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	category, err := sc.CategoriesService.GetCategoryByID(uint(id))
	if err != nil {
		http.Error(w, "Categoría no encontrada", http.StatusNotFound)
		return
	}
	categoryDTO := sc.CategoriesService.ConvertToCategoryDTO(category)
	json.NewEncoder(w).Encode(categoryDTO)

}

func (sc *CategoriesController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	if err := sc.CategoriesService.CreateCategory(&category); err != nil {
		if strings.Contains(err.Error(), "ya existe") {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Error al crear la categoria", http.StatusInternalServerError)
		}
		return
	}

	websocket.GetEventManager().BroadcastMessage(websocket.Message{
		Type: "category_created",
		Data: websocket.ProductData{
			ID:   int(category.ID),
			Name: category.Name,
		},
	})
	json.NewEncoder(w).Encode(category)
}

func (sc *CategoriesController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	category.ID = uint(id)
	if err := sc.CategoriesService.UpdateCategory(&category); err != nil {
		if strings.Contains(err.Error(), "ya existe") {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Error al actualizar la categoría", http.StatusInternalServerError)
		}
		return
	}
	categoryDTO := sc.CategoriesService.ConvertToCategoryDTO(&category)

	websocket.GetEventManager().BroadcastMessage(websocket.Message{
		Type: "category_updated",
		Data: websocket.ProductData{
			ID:   int(categoryDTO.ID),
			Name: categoryDTO.Name,
		},
	})
	json.NewEncoder(w).Encode(categoryDTO)
}

func (sc *CategoriesController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	category, err := sc.CategoriesService.GetCategoryByID(uint(id))
	if err != nil {
		http.Error(w, "Categoría no encontrada", http.StatusNotFound)
		return
	}
	if err := sc.CategoriesService.DeleteCategory(category); err != nil {
		http.Error(w, "Error al eliminar categoría", http.StatusInternalServerError)
		return
	}
	websocket.GetEventManager().BroadcastMessage(websocket.Message{
		Type: "category_deleted",
		Data: websocket.ProductData{
			ID:  int(category.ID),
			Name: category.Name,
		},
	})
	w.WriteHeader(http.StatusNoContent)
}
