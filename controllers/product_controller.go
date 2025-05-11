package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"qisur-challenge/models"
	"qisur-challenge/services"
	websocket "qisur-challenge/webSocket"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ProductController struct {
	ProductService services.ProductService
	DB             *gorm.DB
}

func NewProductController(db *gorm.DB, productService services.ProductService) *ProductController {
	return &ProductController{DB: db, ProductService: productService}
}

func (pc *ProductController) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := pc.ProductService.GetAllProducts()
	if err != nil {
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}
	productDtos := pc.ProductService.ConvertToProductDTOs(products)
	json.NewEncoder(w).Encode(productDtos)

}

func (pc *ProductController) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	product, err := pc.ProductService.GetProductByID(uint(id))
	if err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	productDTO := pc.ProductService.ConvertToProductDTO(product)
	json.NewEncoder(w).Encode(productDTO)
}

func (pc *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	if err := pc.ProductService.CreateProduct(&product); err != nil {
		if strings.Contains(err.Error(), "ya existe") {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Error al crear producto", http.StatusInternalServerError)
		}
		return
	}
	websocket.GetEventManager().BroadcastMessage(websocket.Message{
		Type: "product_created",
		Data: websocket.ProductData{
			ID:   int(product.ID),
			Name: product.Name,
		},
	})

	json.NewEncoder(w).Encode(product)
}

func (pc *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	updatedProduct, err := pc.ProductService.UpdateProduct(uint(id), &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("UpdateProduct: Producto no encontrado ID=%d", id)
			http.Error(w, "Producto no encontrado", http.StatusNotFound)
		} else {
			log.Printf("UpdateProduct: Error actualizando producto ID=%d, error=%v", id, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	websocket.GetEventManager().BroadcastMessage(websocket.Message{
		Type: "product_upgraded",
		Data: websocket.ProductData{
			ID:   int(updatedProduct.ID),
			Name: updatedProduct.Name,
		},
	})
	json.NewEncoder(w).Encode(pc.ProductService.ConvertToProductDTO(updatedProduct))
}

func (pc *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	product, err := pc.ProductService.GetProductByID(uint(id))
	if err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}
	if err := pc.ProductService.DeleteProduct(product); err != nil {
		http.Error(w, "Error al eliminar producto", http.StatusInternalServerError)
		return
	}
	websocket.GetEventManager().BroadcastMessage(websocket.Message{
		Type: "product_delete",
		Data: websocket.ProductData{
			ID:   int(product.ID),
			Name: product.Name,
		},
	})
	w.WriteHeader(http.StatusNoContent)

}

func (pc *ProductController) GetProductHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var startTime, endTime *time.Time
	const layout = "2006-01-02"

	if startStr != "" {
		t, err := time.Parse(layout, startStr)
		if err != nil {
			http.Error(w, "Fecha 'start' inválida. Formato esperado: YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		startTime = &t
	}

	if endStr != "" {
		t, err := time.Parse(layout, endStr)
		if err != nil {
			http.Error(w, "Fecha 'end' inválida. Formato esperado: YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		endTime = &t
	}

	if startTime != nil && endTime != nil && startTime.After(*endTime) {
		http.Error(w, "El parámetro 'start' no puede ser posterior a 'end'", http.StatusBadRequest)
		return
	}

	history, err := pc.ProductService.GetProductHistory(uint(id), startTime, endTime)
	if err != nil {
		http.Error(w, "Error al obtener historial", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(history)
}

func (pc *ProductController) SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	searchType := query.Get("type")
	name := query.Get("name")
	sort := query.Get("sort")
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 10
	}

	switch searchType {
	case "product":
		results, err := pc.ProductService.SearchProducts(name, sort, page, limit)
		if err != nil {
			http.Error(w, "Error al buscar productos", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(results)
	case "category":
		results, err := pc.ProductService.SearchCategories(name, sort, page, limit)
		if err != nil {
			http.Error(w, "Error al buscar categorías", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(results)
	default:
		http.Error(w, "Tipo de búsqueda inválido", http.StatusBadRequest)
	}
}
