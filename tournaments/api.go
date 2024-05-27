package main

import (
	"encoding/json"
	"errors"

	"net/http"
	"os"

	"github.com/Neniel/gotennis/tournaments/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/log"
	"github.com/Neniel/gotennis/lib/telemetry/grafana"
)

type Usecases struct {
	CreateTournament usecase.CreateTournament
	ListTournaments  usecase.ListTournaments
	GetTournament    usecase.GetTournament
	UpdateTournament usecase.UpdateTournament
	DeleteTournament usecase.DeleteTournament
}

type TournamentMicroservice struct {
	App      app.IApp
	Usecases *Usecases
}

type APIServer struct {
	TournamentMicroservice *TournamentMicroservice
}

func (ms *TournamentMicroservice) NewAPIServer() *APIServer {
	return &APIServer{
		TournamentMicroservice: &TournamentMicroservice{
			App:      ms.App,
			Usecases: ms.Usecases,
		},
	}

}

func (api *APIServer) Run() {
	log.Logger.Info("Starting API Server")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/ping", api.pingHandler)
	mux.HandleFunc("GET /api/tournaments", api.listTournaments)
	mux.HandleFunc("GET /api/tournaments/{id}", api.getTournament)
	mux.HandleFunc("POST /api/tournaments", api.addTournament)
	mux.HandleFunc("PUT /api/tournaments/{id}", api.updateTournament)
	mux.HandleFunc("DELETE /api/tournaments/{id}", api.deleteTournament)

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

func (api *APIServer) listTournaments(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	categories, err := api.TournamentMicroservice.Usecases.ListTournaments.Do(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		grafana.SendMetric("tournament.list", 1, 1, map[string]interface{}{
			"status_code": http.StatusInternalServerError,
		})
		return
	}

	err = json.NewEncoder(w).Encode(&categories)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		grafana.SendMetric("tournament.list", 1, 1, map[string]interface{}{
			"status_code": http.StatusInternalServerError,
		})
		return
	}

	grafana.SendMetric("tournament.list", 1, 1, map[string]interface{}{
		"status_code": http.StatusOK,
	})
}

func (api *APIServer) getTournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	if categoryId := r.PathValue("id"); categoryId != "" {
		categories, err := api.TournamentMicroservice.Usecases.GetTournament.Do(r.Context(), categoryId)
		if errors.Is(err, primitive.ErrInvalidHex) {
			w.WriteHeader(http.StatusBadRequest)
			grafana.SendMetric("tournament.get", 1, 1, map[string]interface{}{
				"status_code": http.StatusBadRequest,
			})
			return
		}

		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			grafana.SendMetric("tournament.get", 1, 1, map[string]interface{}{
				"status_code": http.NotFound,
			})
			return
		}

		err = json.NewEncoder(w).Encode(&categories)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			grafana.SendMetric("tournament.get", 1, 1, map[string]interface{}{
				"status_code": http.StatusInternalServerError,
			})
			return
		}

		grafana.SendMetric("tournament.get", 1, 1, map[string]interface{}{
			"status_code": http.StatusOK,
		})
	}
}

func (api *APIServer) addTournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	var request usecase.CreateTournamentRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		grafana.SendMetric("tournaments.add", 1, 1, map[string]interface{}{
			"status_code": http.StatusBadRequest,
		})
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	tournament, err := api.TournamentMicroservice.Usecases.CreateTournament.CreateTournament(r.Context(), &request)
	if err != nil {
		grafana.SendMetric("tournaments.add", 1, 1, map[string]interface{}{
			"status_code": http.StatusBadRequest,
		})
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&tournament)
	if err != nil {
		grafana.SendMetric("tournaments.add", 1, 1, map[string]interface{}{
			"status_code": http.StatusInternalServerError,
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	grafana.SendMetric("tournaments.add", 1, 1, map[string]interface{}{
		"status_code": http.StatusCreated,
	})
}
func (api *APIServer) updateTournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	if id := r.PathValue("id"); id != "" {
		var request usecase.UpdateTournamentRequest
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			grafana.SendMetric("tournaments.update", 1, 1, map[string]interface{}{
				"status_code": http.StatusBadRequest,
			})
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		category, err := api.TournamentMicroservice.Usecases.UpdateTournament.Do(r.Context(), id, &request)
		if err != nil {
			grafana.SendMetric("tournaments.update", 1, 1, map[string]interface{}{
				"status_code": http.StatusBadRequest,
			})
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&category)
		if err != nil {
			grafana.SendMetric("tournaments.update", 1, 1, map[string]interface{}{
				"status_code": http.StatusInternalServerError,
			})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
}

func (api *APIServer) deleteTournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	if id := r.PathValue("id"); id != "" {
		err := api.TournamentMicroservice.Usecases.DeleteTournament.Do(r.Context(), id)
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
