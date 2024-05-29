package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/Neniel/gotennis/categories/usecase"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/telemetry/grafana"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Usecases struct {
	CreateCategoryUsecase usecase.CreateCategory
	ListCategories        usecase.ListCategories
	GetCategory           usecase.GetCategory
	UpdateCategory        usecase.UpdateCategory
	DeleteCategory        usecase.DeleteCategory
}

type CategoryMicroservice struct {
	App app.IApp
	//Usecases *Usecases
}

type APIServer struct {
	CategoryMicroservice *CategoryMicroservice
}

func (ms *CategoryMicroservice) NewAPIServer() *APIServer {
	return &APIServer{
		CategoryMicroservice: &CategoryMicroservice{
			App: ms.App,
			//Usecases: ms.Usecases,
		},
	}

}

func (api *APIServer) Run() {
	log.Println("Starting API Server")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/ping", api.pingHandler)
	mux.HandleFunc("GET /api/categories", api.listCategories)
	mux.HandleFunc("GET /api/categories/{id}", api.getCategory)
	mux.HandleFunc("POST /api/categories", api.addCategory)
	mux.HandleFunc("PUT /api/categories/{id}", api.updateCategory)
	mux.HandleFunc("DELETE /api/categories/{id}", api.deleteCategory)
	mux.Handle("/metrics", promhttp.Handler())

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
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")

	/*
	   1. recibir el token
	   2. validar el token
	   3. obtener datos del token
	*/

	tenantID := r.Header.Get("X-Tenant-ID")

	client, ok := api.CategoryMicroservice.App.GetMongoDBClients()[tenantID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid value in header X-Tenant-ID"))
		return
	}

	listCategories := usecase.NewListCategories(database.NewDatabaseReader(client.MongoDBClient, client.DatabaseName))

	categories, err := listCategories.Do(r.Context())
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
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	if categoryId := r.PathValue("id"); categoryId != "" {

		/*
		   1. recibir el token
		   2. validar el token
		   3. obtener datos del token
		*/

		tenantID := r.Header.Get("X-Tenant-ID")

		client, ok := api.CategoryMicroservice.App.GetMongoDBClients()[tenantID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid value in header X-Tenant-ID"))
			return
		}

		getCategory := usecase.NewGetCategory(database.NewDatabaseReader(client.MongoDBClient, client.DatabaseName))

		categories, err := getCategory.Do(r.Context(), categoryId)
		if errors.Is(err, primitive.ErrInvalidHex) {
			w.WriteHeader(http.StatusBadRequest)
			grafana.SendMetric("get.category", 1, 1, map[string]interface{}{
				"status_code": http.StatusBadRequest,
			})
			return
		}

		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			grafana.SendMetric("get.category", 1, 1, map[string]interface{}{
				"status_code": http.NotFound,
			})
			return
		}

		err = json.NewEncoder(w).Encode(&categories)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			grafana.SendMetric("get.category", 1, 1, map[string]interface{}{
				"status_code": http.StatusInternalServerError,
			})
			return
		}
		grafana.SendMetric("get.category", 1, 1, map[string]interface{}{
			"status_code": http.StatusOK,
		})
	}
}

func (api *APIServer) addCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")

	var request usecase.CreateCategoryRequest

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	/*
	   1. recibir el token
	   2. validar el token
	   3. obtener datos del token
	*/

	tenantID := r.Header.Get("X-Tenant-ID")

	client, ok := api.CategoryMicroservice.App.GetMongoDBClients()[tenantID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid value in header X-Tenant-ID"))
		return
	}

	createCategory := usecase.NewCreateCategory(database.NewDatabaseWriter(client.MongoDBClient, client.DatabaseName))

	category, err := createCategory.Do(r.Context(), &request)
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
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	if id := r.PathValue("id"); id != "" {

		var request usecase.UpdateCategoryRequest

		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		/*
		   1. recibir el token
		   2. validar el token
		   3. obtener datos del token
		*/

		tenantID := r.Header.Get("X-Tenant-ID")

		client, ok := api.CategoryMicroservice.App.GetMongoDBClients()[tenantID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid value in header X-Tenant-ID"))
			return
		}

		updateCategory := usecase.NewUpdateCategory(database.NewDatabaseReader(client.MongoDBClient, client.DatabaseName), database.NewDatabaseWriter(client.MongoDBClient, client.DatabaseName))

		category, err := updateCategory.Do(r.Context(), id, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&category)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
}

func (api *APIServer) deleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	if id := r.PathValue("id"); id != "" {

		/*
		   1. recibir el token
		   2. validar el token
		   3. obtener datos del token
		*/

		tenantID := r.Header.Get("X-Tenant-ID")

		client, ok := api.CategoryMicroservice.App.GetMongoDBClients()[tenantID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid value in header X-Tenant-ID"))
			return
		}

		deleteCategory := usecase.NewDeleteCategory(database.NewDatabaseWriter(client.MongoDBClient, client.DatabaseName))

		err := deleteCategory.Do(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
