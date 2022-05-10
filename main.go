package main

import (
	"ProductApp/app"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()
	s := route.PathPrefix("/api").Subrouter() //Base Path

	// Product Routes
	s.HandleFunc("/CreateProduct", app.CreateProduct).Methods("POST")
	s.HandleFunc("/GetAllProduct", app.GetAllProduct).Methods("GET")
	s.HandleFunc("/GetProduct/{id}", app.GetProduct).Methods("GET")
	s.HandleFunc("/UpdateProduct", app.UpdateProduct).Methods("PUT")
	s.HandleFunc("/UpdateProductStatus", app.UpdateProductStatus).Methods("PUT")
	s.HandleFunc("/DeleteProduct/{id}", app.DeleteProduct).Methods("DELETE")

	// Category Routes
	s.HandleFunc("/CreateCategory", app.CreateCategory).Methods("POST")
	s.HandleFunc("/GetAllCategory", app.GetAllCategory).Methods("GET")
	s.HandleFunc("/GetCategory/{id}", app.GetCategory).Methods("GET")
	s.HandleFunc("/UpdateCategory", app.UpdateCategory).Methods("PUT")
	s.HandleFunc("/UpdateCategoryStatus", app.UpdateCategoryStatus).Methods("PUT")
	s.HandleFunc("/DeleteCategory/{id}", app.DeleteCategory).Methods("DELETE")

	// SubCategory Routes
	s.HandleFunc("/CreateSubCategory", app.CreateSubCategory).Methods("POST")
	s.HandleFunc("/GetAllSubCategory", app.GetAllSubCategory).Methods("GET")
	s.HandleFunc("/GetSubCategory/{id}", app.GetSubCategory).Methods("GET")
	s.HandleFunc("/UpdateSubCategory", app.UpdateSubCategory).Methods("PUT")
	s.HandleFunc("/UpdateSubCategoryStatus", app.UpdateSubCategoryStatus).Methods("PUT")
	s.HandleFunc("/DeleteSubCategory/{id}", app.DeleteSubCategory).Methods("DELETE")

	// Brand Routes
	s.HandleFunc("/CreateBrand", app.CreateBrand).Methods("POST")
	s.HandleFunc("/GetAllBrand", app.GetAllBrand).Methods("GET")
	s.HandleFunc("/GetBrand/{id}", app.GetBrand).Methods("GET")
	s.HandleFunc("/UpdateBrand", app.UpdateBrand).Methods("PUT")
	s.HandleFunc("/UpdateBrandStatus", app.UpdateBrandStatus).Methods("PUT")
	s.HandleFunc("/DeleteBrand/{id}", app.DeleteBrand).Methods("DELETE")

	// Varient Routes
	s.HandleFunc("/CreateVarient", app.CreateVarient).Methods("POST")
	s.HandleFunc("/GetAllVarient", app.GetAllVarient).Methods("GET")
	s.HandleFunc("/GetVarient/{id}", app.GetVarient).Methods("GET")
	s.HandleFunc("/UpdateVarient", app.UpdateVarient).Methods("PUT")
	s.HandleFunc("/UpdateVarientStatus", app.UpdateVarientStatus).Methods("PUT")
	s.HandleFunc("/DeleteVarient/{id}", app.DeleteVarient).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", s)) // Run Server
}
