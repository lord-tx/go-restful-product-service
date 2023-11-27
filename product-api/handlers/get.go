package handlers

import (
	"go-restful-product-service/product-api/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
