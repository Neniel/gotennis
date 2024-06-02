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
	CreateTenant usecase.CreateTenant
	ListTenants  usecase.ListTenants
	GetTenant    usecase.GetTenant
	//UpdateCustomer usecase.UpdateTenant
	DeleteTenant usecase.DeleteTenant
}

type CustomerMicroservice struct {
	App      app.IApp
	Usecases *Usecases
}

type APIServer struct {
	CustomerMicroservice *CustomerMicroservice
}

func (ms *CustomerMicroservice) NewAPIServer() *APIServer {
	return &APIServer{
		CustomerMicroservice: &CustomerMicroservice{
			App:      ms.App,
			Usecases: ms.Usecases,
		},
	}

}

func (api *APIServer) Run() {
	log.Println("Starting API Server")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", api.pingHandler)
	mux.HandleFunc("GET /tenants", api.listTenants)
	mux.HandleFunc("GET /tenants/{id}", api.getTenant)
	mux.HandleFunc("POST /tenants", api.addTenant)
	//mux.HandleFunc("PUT /api/tenants/{id}", api.updateCustomer)
	mux.HandleFunc("DELETE /tenants/{id}", api.deleteTenant)
	mux.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(os.Getenv("APP_PORT"), mux))
}

func (api *APIServer) pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("Tenants is ok"))
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (api *APIServer) listTenants(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	customers, err := api.CustomerMicroservice.Usecases.ListTenants.Do(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&customers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (api *APIServer) getTenant(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	if id := r.PathValue("id"); id != "" {
		tenants, err := api.CustomerMicroservice.Usecases.GetTenant.Do(r.Context(), id)
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

		err = json.NewEncoder(w).Encode(&tenants)
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

func (api *APIServer) addTenant(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	var request usecase.CreateTenantRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	customer, err := api.CustomerMicroservice.Usecases.CreateTenant.Do(r.Context(), &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&customer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (api *APIServer) deleteTenant(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	if id := r.PathValue("id"); id != "" {
		err := api.CustomerMicroservice.Usecases.DeleteTenant.Do(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
