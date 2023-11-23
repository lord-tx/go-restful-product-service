package handlers

import (
	"context"
	"fmt"
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

func (p *Product) GetProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert ID", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle GET Product", id)

	prod, err := data.GetProduct(id)
	if err != nil {
		http.Error(rw, "Unable to find product", http.StatusBadRequest)
		return
	}

	err = prod.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}
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
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

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

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

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

type KeyProduct struct{}

func (p *Product) MiddleWareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product")
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		/// validatee the input
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product")
			http.Error(rw, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
