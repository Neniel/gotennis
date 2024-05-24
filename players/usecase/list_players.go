package usecase

import (
	"context"
	"log"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
)

type ListPlayersUsecase interface {
	List(ctx context.Context) ([]entity.Player, error)
}

type listPlayersUsecase struct {
	DBReader database.DBReader
}

func NewListPlayersUsecase(dbReader database.DBReader) ListPlayersUsecase {
	return &listPlayersUsecase{
		DBReader: dbReader,
	}
}

func (uc *listPlayersUsecase) List(ctx context.Context) ([]entity.Player, error) {
	players, err := uc.DBReader.GetPlayers(ctx)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return players, nil
}
