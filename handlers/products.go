package handlers

import (
	"log"
	"net/http"

	"github.com/Diego-Paris/microservices-in-go/data"
)

// Products does
type Products struct {
	l *log.Logger
}

// NewProducts returns address to product handler
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}

	// handle an update, PUT creates or replaces
	if r.Method ==  http.MethodPut {

	}

	// catch all
	w.WriteHeader(http.StatusMethodNotAllowed)

}

func (p *Products) getProducts(w http.ResponseWriter, h *http.Request) {
	p.l.Println("GET METHOD")
	lp := data.GetProducts()

	err := lp.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}
