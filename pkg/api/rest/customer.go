package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type Customer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

var rdb *redis.Client
var ctx = context.Background()

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	// Fetch customer IDs from Redis Set
	customerIDs, err := rdb.SMembers(ctx, "customers").Result()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Fetch each customer by ID
	var customers []Customer
	for _, id := range customerIDs {
		customerJSON, err := rdb.Get(ctx, "customer:"+id).Result()
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		var customer Customer
		err = json.Unmarshal([]byte(customerJSON), &customer)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		customers = append(customers, customer)
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	customerJSON, err := rdb.Get(ctx, "customer:"+id).Result()
	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	var customer Customer
	err = json.Unmarshal([]byte(customerJSON), &customer)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&customer); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	customerJSON, err := json.Marshal(customer)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	id := strconv.Itoa(customer.ID)
	rdb.Set(ctx, "customer:"+id, customerJSON, 0)
	rdb.SAdd(ctx, "customers", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var updated Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updated); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updated.ID, _ = strconv.Atoi(id)
	updatedJSON, err := json.Marshal(updated)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	rdb.Set(ctx, "customer:"+id, updatedJSON, 0)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	rdb.Del(ctx, "customer:"+id)
	rdb.SRem(ctx, "customers", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
