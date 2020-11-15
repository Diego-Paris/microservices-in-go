package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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

	// handle a creation
	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}

	// handle a total update
	if r.Method == http.MethodPut {
		p.l.Println("PUT")
		// expect the id in the URI

		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URI, more than one id")
			http.Error(w, "Invalid URI1", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI, more than one capture group", g)
			http.Error(w, "Invalid URI2", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI, unable to convert to number", idString)
			http.Error(w, "Invalid URI3", http.StatusBadRequest)
			return
		}

		p.updateProducts(id, w, r)
		return
	}

	// catch all
	w.WriteHeader(http.StatusMethodNotAllowed)

}

func (p *Products) getProducts(w http.ResponseWriter, h *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()

	err := lp.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}

func (p *Products) updateProducts(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Products")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	//p.l.Printf("Prod: %#v", prod)
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	
	if err != nil {
		http.Error(w, "Product update failed", http.StatusInternalServerError)
		return
	}


}