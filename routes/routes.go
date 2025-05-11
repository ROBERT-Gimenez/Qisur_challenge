package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"qisur-challenge/controllers"
	"qisur-challenge/middlewares"
	ws "qisur-challenge/webSocket"
)

func ApplyMiddlewareRoute(router *mux.Router, route string, handler http.HandlerFunc, methods ...string) {
	router.Handle(route, middlewares.AuthMiddleware(handler)).Methods(methods...)
}

func RegisterRoutes(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/login", controllers.Login(db)).Methods("POST")
    r.Handle("/ws", middlewares.AuthMiddleware(http.HandlerFunc(ws.HandleWebSocket)))

	api := r.PathPrefix("/api").Subrouter()

	ProductRoutes(db, api)
	CategoriesRoutes(db, api)

	return r
}
