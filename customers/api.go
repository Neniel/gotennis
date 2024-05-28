package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/Neniel/gotennis/customers/usecase"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/telemetry/grafana"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Usecases struct {
	CreateCustomer usecase.CreateCustomer
	//ListCustomers  usecase.ListCategories
	//GetCustomer    usecase.GetCategory
	//UpdateCustomer usecase.UpdateCategory
	//DeleteCustomer usecase.DeleteCategory
}

type CustomerMicroservice struct {
	App app.IApp
}

type APIServer struct {
	CustomerMicroservice *CustomerMicroservice
}

func (ms *CustomerMicroservice) NewAPIServer() *APIServer {
	return &APIServer{
		CustomerMicroservice: &CustomerMicroservice{
			App: ms.App,
		},
	}

}

func (api *APIServer) Run() {
	log.Println("Starting API Server")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", api.pingHandler)
	mux.HandleFunc("GET /customers", api.listCustomers)
	mux.HandleFunc("GET /customers/{id}", api.getCustomer)
	mux.HandleFunc("POST /customers", api.addCustomer)
	mux.HandleFunc("PUT /customers/{id}", api.updateCustomer)
	mux.HandleFunc("DELETE /customers/{id}", api.deleteCustomer)
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

func (api *APIServer) listCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	createCustomer, err := usecase.NewCreateCustomer(api.CustomerMicroservice.App, "CUSTOMERIDGOESHERE")
	if err != nil {

	}

	categories, err := createCustomer.Do(r.Context())
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

func (api *APIServer) getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	if categoryId := r.PathValue("id"); categoryId != "" {
		categories, err := api.CustomerMicroservice.Usecases.GetCustomer.Do(r.Context(), categoryId)
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

func (api *APIServer) addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	var request usecase.CreateCustomerRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	category, err := api.CustomerMicroservice.Usecases.CreateCustomer.CreateCategory(r.Context(), &request)
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
func (api *APIServer) updateCustomer(w http.ResponseWriter, r *http.Request) {
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

		category, err := api.CustomerMicroservice.Usecases.UpdateCustomer.Do(r.Context(), id, &request)
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

func (api *APIServer) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	if id := r.PathValue("id"); id != "" {
		err := api.CustomerMicroservice.Usecases.DeleteCustomer.Do(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
