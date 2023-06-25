package main

import (
	"go-commerce/pkg/api/rest"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize Redis
	rest.InitRedis()

	// Initialize Router
	r := mux.NewRouter()

	// Route Handlers
	r.HandleFunc("/api/customers", rest.GetCustomers).Methods("GET")
	r.HandleFunc("/api/customers/{id}", rest.GetCustomer).Methods("GET")
	r.HandleFunc("/api/customers", rest.CreateCustomer).Methods("POST")
	r.HandleFunc("/api/customers/{id}", rest.UpdateCustomer).Methods("PUT")
	r.HandleFunc("/api/customers/{id}", rest.DeleteCustomer).Methods("DELETE")

	// Start Server
	log.Fatal(http.ListenAndServe(":8000", r))
}
