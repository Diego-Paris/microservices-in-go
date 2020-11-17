// Package handlers classification of Product API
//
// Documentation for Product API
//
//  Schemes: http
//  BasePath: /
//  Version: 1.0.0
//
//  Consumes:
//  - application/json
// 
//  Produces:
//  - application/json
//  swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Diego-Paris/microservices-in-go/data"
	"github.com/gorilla/mux"
)

// A list of products returns in the response
// swagger:response productsResponse
type productsResponse struct {
	// All products in the system
	// in: body
	Body []data.Product
}


// Products does
type Products struct {
	l *log.Logger
}

// NewProducts returns address to product handler
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//  200: productsResponse

// GetProducts returns the products from the data store
func (p *Products) GetProducts(w http.ResponseWriter, h *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()

	err := lp.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

// AddProduct a product to the db
func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(&prod)
}

// UpdateProducts updates in db
func (p *Products) UpdateProducts(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Could not parse id to num", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Products", id)

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	//p.l.Printf("Prod: %#v", prod)
	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product update failed", http.StatusInternalServerError)
		return
	}

}

// KeyProduct used to identify product in context
type KeyProduct struct{}

// MiddlewareProductValidation does
func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		// validate the product
		err = prod.Validate()
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
