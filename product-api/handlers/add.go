package handlers

import (
	"go-restful-product-service/product-api/data"
	"net/http"
)

func (p *Product) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle ADD Products")
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	data.AddProduct(prod)
}
