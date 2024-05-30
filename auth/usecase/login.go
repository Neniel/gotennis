package usecase

import (
	"context"
	"fmt"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/log"
)

type Login interface {
	Do(ctx context.Context, request *LoginRequest) error
}

type login struct {
	App      app.IApp
	DBReader database.DBReader
}

func NewLogin(dbReader database.DBReader) Login {
	return &login{
		DBReader: dbReader,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	return nil
}

func (uc *login) Do(ctx context.Context, request *LoginRequest) error {
	if err := request.Validate(); err != nil {
		log.Logger.Info(fmt.Errorf("could not login: %w", err).Error())
		return err
	}

	err := uc.DBReader.Login(ctx, request.Username, request.Password)
	if err != nil {
		log.Logger.Info(fmt.Errorf("could not create category: %w", err).Error())
		return err
	}
	return nil
}
