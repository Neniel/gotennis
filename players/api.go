package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"players/usecase"

	"github.com/Neniel/gotennis/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Usecases struct {
	CreatePlayerUsecase usecase.CreatePlayerUsecase
	SavePlayerUsecase   interface{}
	ListPlayersUsecase  usecase.ListPlayersUsecase
	GetPlayerUsecase    usecase.GetPlayerUsecase
	UpdatePlayerUsecase interface{}
	DeletePlayerUsecase usecase.DeletePlayerUsecase
}

type PlayerMicroservice struct {
	App      *app.App
	Usecases *Usecases
}

type APIServer struct {
	PlayerMicroservice *PlayerMicroservice
}

func (ms *PlayerMicroservice) NewAPIServer() *APIServer {
	return &APIServer{
		PlayerMicroservice: &PlayerMicroservice{
			App:      ms.App,
			Usecases: ms.Usecases,
		},
	}

}

func (api *APIServer) Run() {
	log.Println("Starting API Server")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", api.pingHandler)
	mux.HandleFunc("GET /players", api.listPlayers)
	mux.HandleFunc("GET /players/{id}", api.getPlayer)
	mux.HandleFunc("POST /players", api.addPlayer)
	mux.HandleFunc("PUT /players/{id}", api.updatePlayer)
	mux.HandleFunc("DELETE /players/{id}", api.deletePlayer)

	log.Fatal(http.ListenAndServe(os.Getenv("APP_PORT"), mux))
}

func (api *APIServer) pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("Ok"))
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (api *APIServer) listPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	categories, err := api.PlayerMicroservice.Usecases.ListPlayersUsecase.List(r.Context())
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
	w.Header().Add("Content-Type", "application/json")
	if categoryId := r.PathValue("id"); categoryId != "" {
		categories, err := api.PlayerMicroservice.Usecases.GetPlayerUsecase.Get(r.Context(), categoryId)
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

func (api *APIServer) addPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var request usecase.CreatePlayerRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	player, err := api.PlayerMicroservice.Usecases.CreatePlayerUsecase.CreatePlayer(r.Context(), &request)
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
	w.Header().Add("Content-Type", "application/json")
}

func (api *APIServer) deletePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if id := r.PathValue("id"); id != "" {
		err := api.PlayerMicroservice.Usecases.DeletePlayerUsecase.Delete(r.Context(), id)
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
