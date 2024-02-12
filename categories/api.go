package main

import (
	"categories/database"
	"categories/handler"
	"categories/helper"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/Neniel/gotennis/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIServer struct {
	Handler *handler.CategoriesHandler
}

func (a *App) NewAPIServer() *APIServer {
	return &APIServer{
		Handler: handler.NewCategoriesHandler(
			database.NewDatabaseReader(a.DBClients.Redis),
			database.NewDatabaseWriter(a.DBClients.Redis),
			database.NewDatabaseReader(a.DBClients.MongoDB),
			database.NewDatabaseWriter(a.DBClients.MongoDB)),
	}

}

func (api *APIServer) Run() {
	log.Println("Starting API Server")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", api.pingHandler)
	mux.HandleFunc("GET /categories", api.getAllCategories)
	mux.HandleFunc("GET /categories/{id}", api.getCategoryByID)
	mux.HandleFunc("POST /categories", api.addCategory)
	mux.HandleFunc("PUT /categories/{id}", api.updateCategory)
	mux.HandleFunc("DELETE /categories/{id}", api.deleteCategory)

	log.Fatal(http.ListenAndServe(os.Getenv("APP_PORT"), mux))
}

func (api *APIServer) pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("Ok"))
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (api *APIServer) getAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	categories, err := api.Handler.GetAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&categories)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (api *APIServer) getCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if categoryId := r.PathValue("id"); categoryId != "" {
		categories, err := api.Handler.GetByID(r.Context(), categoryId)
		if errors.Is(err, primitive.ErrInvalidHex) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(&categories)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
}

func (api *APIServer) addCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var categoryToCreate entity.Category
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&categoryToCreate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	createdCategory, err := api.Handler.Add(r.Context(), &categoryToCreate)
	if err != nil {
		if errors.Is(err, helper.ErrCategoryNameIsEmpty) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&createdCategory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
func (api *APIServer) updateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
}

func (api *APIServer) deleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if id := r.PathValue("id"); id != "" {
		err := api.Handler.Delete(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
