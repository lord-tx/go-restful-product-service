package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gte=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {

	/// sku is of the format abc-bcd-exef

	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProduct(id int) (*Product, error) {
	prod, _, err := findOneProduct(id)

	if err != nil {
		return nil, err
	}

	return prod, nil
}

func GetProducts() Products {
	return productsList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productsList = append(productsList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findOneProduct(id)

	if err != nil {
		return err
	}

	p.ID = id
	productsList[pos] = p

	return nil
}

var ErrProdNotFound = fmt.Errorf("Product not found")

func findOneProduct(id int) (*Product, int, error) {

	for pos, prod := range productsList {
		if prod.ID == id {
			return prod, pos, nil
		}
	}

	return nil, -1, ErrProdNotFound
}

func getNextID() int {
	lp := productsList[len(productsList)-1]
	return lp.ID + 1
}

var productsList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc133",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fdj34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
