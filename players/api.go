package main

import (
	"encoding/json"
	"errors"

	"net/http"
	"os"

	"github.com/Neniel/gotennis/players/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/log"
	"github.com/Neniel/gotennis/lib/telemetry/grafana"
)

type Usecases struct {
	CreatePlayer          usecase.CreatePlayer
	SavePlayer            interface{}
	ListPlayers           usecase.ListPlayers
	GetPlayer             usecase.GetPlayer
	UpdatePlayer          usecase.UpdatePlayer
	PartiallyUpdatePlayer usecase.PartialltUpdatePlayer
	DeletePlayer          usecase.DeletePlayer
}

type PlayerMicroservice struct {
	App app.IApp
	//Usecases *Usecases
}

type APIServer struct {
	PlayerMicroservice *PlayerMicroservice
}

func (ms *PlayerMicroservice) NewAPIServer() *APIServer {
	return &APIServer{
		PlayerMicroservice: &PlayerMicroservice{
			App: ms.App,
			//Usecases: ms.Usecases,
		},
	}

}

func (api *APIServer) Run() {
	log.Logger.Info("Starting API Server")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", api.pingHandler)
	mux.HandleFunc("GET /api/players", api.listPlayers)
	mux.HandleFunc("GET /api/players/{id}", api.getPlayer)
	mux.HandleFunc("POST /api/players", api.addPlayer)
	mux.HandleFunc("PUT /api/players/{id}", api.updatePlayer)
	mux.HandleFunc("PATCH /api/players/{id}", api.partiallyUpdatePlayer)
	mux.HandleFunc("DELETE /api/players/{id}", api.deletePlayer)

	log.Logger.Error(http.ListenAndServe(os.Getenv("APP_PORT"), mux).Error())
}

func (api *APIServer) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodGet {
		w.Write([]byte("Ok"))
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (api *APIServer) listPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")

	/*
	   1. recibir el token
	   2. validar el token
	   3. obtener datos del token
	*/

	customerID := "" // viene del token

	client, ok := api.PlayerMicroservice.App.GetMongoDBClients()[customerID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	listPlayers := usecase.NewListPlayers(database.NewDatabaseReader(client.MongoDBClient, client.DatabaseName))

	categories, err := listPlayers.Do(r.Context())
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

func (api *APIServer) getPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	if categoryId := r.PathValue("id"); categoryId != "" {

		/*
		   1. recibir el token
		   2. validar el token
		   3. obtener datos del token
		*/

		customerID := "" // viene del token

		client, ok := api.PlayerMicroservice.App.GetMongoDBClients()[customerID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		getPlayer := usecase.NewGetPlayer(database.NewDatabaseReader(client.MongoDBClient, client.DatabaseName))

		categories, err := getPlayer.Do(r.Context(), categoryId)
		if errors.Is(err, primitive.ErrInvalidHex) {
			w.WriteHeader(http.StatusBadRequest)
			grafana.SendMetric("get.player", 1, 1, map[string]interface{}{
				"status_code": http.StatusBadRequest,
			})
			return
		}

		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			grafana.SendMetric("get.player", 1, 1, map[string]interface{}{
				"status_code": http.NotFound,
			})
			return
		}

		err = json.NewEncoder(w).Encode(&categories)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			grafana.SendMetric("get.player", 1, 1, map[string]interface{}{
				"status_code": http.StatusInternalServerError,
			})
			return
		}

		grafana.SendMetric("get.player", 1, 1, map[string]interface{}{
			"status_code": http.StatusOK,
		})
	}
}

func (api *APIServer) addPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	var request usecase.CreatePlayerRequest
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

	client, ok := api.PlayerMicroservice.App.GetMongoDBClients()[tenantID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid value in header X-Tenant-ID"))
		return
	}

	createPlayer := usecase.NewCreatePlayer(
		database.NewDatabaseWriter(client.MongoDBClient, client.DatabaseName),
		database.NewDatabaseReader(client.MongoDBClient, client.DatabaseName),
	)

	player, err := createPlayer.Do(r.Context(), &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&player)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (api *APIServer) updatePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	if id := r.PathValue("id"); id != "" {
		var request usecase.UpdatePlayerRequest
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

		customerID := "" // viene del token

		client, ok := api.PlayerMicroservice.App.GetMongoDBClients()[customerID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updatePlayer := usecase.NewUpdatePlayer(database.NewDatabaseWriter(client.MongoDBClient, client.DatabaseName), database.NewDatabaseReader(client.MongoDBClient, client.DatabaseName))

		category, err := updatePlayer.Do(r.Context(), id, &request)
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

func (api *APIServer) partiallyUpdatePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if id := r.PathValue("id"); id != "" {
		var request usecase.PartiallyUpdatePlayerRequest
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

		customerID := "" // viene del token

		client, ok := api.PlayerMicroservice.App.GetMongoDBClients()[customerID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		partiallyUpdatePlayer := usecase.NewPartiallyUpdatePlayer(database.NewDatabaseWriter(client.MongoDBClient, client.DatabaseName), database.NewDatabaseReader(client.MongoDBClient, client.DatabaseName))

		player, err := partiallyUpdatePlayer.Do(r.Context(), id, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&player)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
}

func (api *APIServer) deletePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	if id := r.PathValue("id"); id != "" {

		/*
		   1. recibir el token
		   2. validar el token
		   3. obtener datos del token
		*/

		customerID := "" // viene del token

		client, ok := api.PlayerMicroservice.App.GetMongoDBClients()[customerID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		deletePlayer := usecase.NewDeletePlayer(database.NewDatabaseWriter(client.MongoDBClient, client.DatabaseName))
		err := deletePlayer.Do(r.Context(), id)
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
