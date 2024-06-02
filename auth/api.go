package main

import (
	"encoding/json"
	"os"

	"net/http"

	"github.com/Neniel/gotennis/auth/usecase"
	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/log"
	"github.com/Neniel/gotennis/lib/middleware"
	"github.com/Neniel/gotennis/lib/telemetry/grafana"
)

type Usecases struct {
	Login usecase.Login
}

type AuthMicroservice struct {
	App app.IApp
	//Usecases *Usecases
}

type APIServer struct {
	AuthMicroservice *AuthMicroservice
}

func (ms *AuthMicroservice) NewAPIServer() *APIServer {
	return &APIServer{
		AuthMicroservice: &AuthMicroservice{
			App: ms.App,
			//Usecases: ms.Usecases,
		},
	}

}

func (api *APIServer) Run() {
	log.Logger.Info("Starting API Server")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", api.pingHandler)
	mux.HandleFunc("POST /login", api.login)

	log.Logger.Error(http.ListenAndServe(os.Getenv("APP_PORT"), middleware.CORSMiddleware(mux)).Error())
}

func (api *APIServer) pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("Auth is ok"))
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (api *APIServer) login(w http.ResponseWriter, r *http.Request) {

	var request usecase.LoginRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		grafana.SendMetric("login", 1, 1, map[string]interface{}{
			"status_code": http.StatusBadRequest,
		})
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	/*
	   1. recibir el token
	   2. validar el token
	   3. obtener datos del token
	*/

	tenantName := r.Header.Get("X-Tenant-Name")

	tenant, err := api.AuthMicroservice.App.GetTenantByName(tenantName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("tenant not found"))
		return
	}

	client, err := api.AuthMicroservice.App.GetTenantMongoDBClient(tenant.ID.Hex())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid value in header X-Tenant-Name"))
		return
	}

	login := usecase.NewLogin(database.NewDatabaseReader(client.MongoDBClient, client.DatabaseName))

	err = login.Do(r.Context(), &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		grafana.SendMetric("login", 1, 1, map[string]interface{}{
			"status_code": http.StatusInternalServerError,
		})
		return
	}

	w.Header().Add("X-Tenant-ID", tenant.ID.Hex())

	grafana.SendMetric("login", 1, 1, map[string]interface{}{
		"status_code": http.StatusOK,
	})
}
