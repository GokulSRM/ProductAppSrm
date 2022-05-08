package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"ProductApp/app"
)

func main() {

	route := mux.NewRouter()
	s := route.PathPrefix("/api").Subrouter() //Base Path

	// Product Routes

	// s.HandleFunc("/createProfile", app.CreateProduct).Methods("POST")
	// s.HandleFunc("/getAllUsers", app.GetAllProduct).Methods("GET")
	// s.HandleFunc("/getUserProfile", app.GetProduct).Methods("POST")
	// s.HandleFunc("/updateProfile", app.UpdateProduct).Methods("PUT")
	// s.HandleFunc("/deleteProfile/{id}", app.DeleteProduct).Methods("DELETE")

	// Category Routes
	s.HandleFunc("/CreateCategory", app.CreateCategory).Methods("POST")
	s.HandleFunc("/GetAllCategory", app.GetAllCategory).Methods("GET")
	s.HandleFunc("/GetCategory", app.GetCategory).Methods("GET")
	s.HandleFunc("/UpdateCategory", app.UpdateCategory).Methods("PUT")
	s.HandleFunc("/UpdateCategoryStatus", app.UpdateCategoryStatus).Methods("PUT")
	s.HandleFunc("/DeleteCategory/{id}", app.DeleteCategory).Methods("DELETE")



	log.Fatal(http.ListenAndServe(":8000", s)) // Run Server
}
