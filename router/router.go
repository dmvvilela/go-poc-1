package router

import (
	"go-postgres/middleware"

	"github.com/gorilla/mux"
)

// Router retuens the router
func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/contacts", middleware.GetAllContacts).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/contacts/{id}", middleware.GetContact).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/contacts", middleware.CreateContact).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/contacts/{id}", middleware.UpdateContact).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/contacts/{id}", middleware.DeleteContact).Methods("DELETE", "OPTIONS")

	return router
}
