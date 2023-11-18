package handlers

import (
	"go-restful-product-service/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l: l}
}

func (p *Product) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()

	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Product) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle ADD Products")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	data.AddProduct(prod)
}

func (p *Product) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Unable to convert ID", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Product", id)

	prod := &data.Product{}

	err = prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProdNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
