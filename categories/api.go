package main

import (
	"categories/usecase"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/Neniel/gotennis/app"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Usecases struct {
	CreateCategoryUsecase usecase.CreateCategoryUsecase
	ListCategories        usecase.ListCategoriesUsecase
	GetCategory           usecase.GetCategoryUsecase
	UpdateCategory        interface{}
	DeleteCategory        usecase.DeleteCategoryUsecase
}

type CategoryMicroservice struct {
	App      app.IApp
	Usecases *Usecases
}

type APIServer struct {
	CategoryMicroservice *CategoryMicroservice
}

func (ms *CategoryMicroservice) NewAPIServer() *APIServer {
	return &APIServer{
		CategoryMicroservice: &CategoryMicroservice{
			App:      ms.App,
			Usecases: ms.Usecases,
		},
	}

}

func (api *APIServer) Run() {
	log.Println("Starting API Server")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", api.pingHandler)
	mux.HandleFunc("GET /categories", api.listCategories)
	mux.HandleFunc("GET /categories/{id}", api.getCategory)
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

func (api *APIServer) listCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	categories, err := api.CategoryMicroservice.Usecases.ListCategories.List(r.Context())
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

func (api *APIServer) getCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if categoryId := r.PathValue("id"); categoryId != "" {
		categories, err := api.CategoryMicroservice.Usecases.GetCategory.Get(r.Context(), categoryId)
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
	}
}

func (api *APIServer) addCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var request usecase.CreateCategoryRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	category, err := api.CategoryMicroservice.Usecases.CreateCategoryUsecase.CreateCategory(r.Context(), &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&category)
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
		err := api.CategoryMicroservice.Usecases.DeleteCategory.Delete(r.Context(), id)
		if errors.Is(err, primitive.ErrInvalidHex) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
