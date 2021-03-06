package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// Product defines the data we are managing
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// FromJSON parses from json to struct
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// ToJSON encodes the contents of the type directly into a write
// does not need to allocate memory into a buffer, it is direct
func (p *Product) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// Validate validates
func (p *Product) Validate() error {
	validate := validator.New()

	validate.RegisterValidation("sku", validateSKU)

	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	//sku is of format abc-abcd-efgh
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

// Products is a slice containing the address of multiple products
type Products []*Product

// ToJSON encodes the contents of the type directly into a write
// does not need to allocate memory into a buffer, it is direct
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

//! Below this is pseudo-database logic

// GetProducts abstracts the logic from accessing a database
func GetProducts() Products {
	return productList
}

// GetProduct abstracts the logic from getting one item from database
func GetProduct(id int) (*Product, error) {
	prod, _, err := findProduct(id)
	if err != nil {
		return nil, err
	}

	return prod, nil
}

// AddProduct abstracts the logic of adding a product to a database
func AddProduct(p *Product) {

	p.ID = getNextID()

	productList = append(productList, p)
}

// UpdateProduct updates product in database
func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p
	return nil
}

// DeleteProduct deletes a product by id in database
func DeleteProduct(id int) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	// swap items
	productList[pos], productList[len(productList) - 1] = productList[len(productList) - 1], productList[pos] 

	// truncate list
	productList = productList[:len(productList) - 1]

	return nil
}

// ErrProductNotFound , the product wasn't found
var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
