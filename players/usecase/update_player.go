package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/entity"
	"github.com/Neniel/gotennis/util"
)

type UpdatePlayerRequest struct {
	ID           string           `json:"id"`
	GovernmentID string           `json:"government_id"`
	FirstName    string           `json:"first_name"`
	MiddleName   string           `json:"middle_name"`
	LastName     string           `json:"last_name"`
	Category     *entity.Category `json:"category"`
	Birthdate    *time.Time       `json:"birthdate"`
	PhoneNumber  string           `json:"phone_number"`
	Email        string           `json:"email"`
	Alias        string           `json:"alias"`
}

func (r *UpdatePlayerRequest) Validate(id string) error {
	if id == "" {
		return errors.New("id is required in request")
	}

	if r.ID == "" {
		return errors.New("id is required in request")
	}

	if r.ID != id {
		return errors.New("id in request and url must be equal")
	}

	if r.GovernmentID == "" {
		return util.ErrPlayerGovernmentIDIsEmpty
	}

	if r.Email == "" {
		return util.ErrPlayerEmailIsEmpty
	}

	if r.FirstName == "" {
		return util.ErrPlayerFirstNameIsEmpty
	}

	if r.LastName == "" {
		return util.ErrPlayerLastNameIsEmpty
	}

	if r.Birthdate != nil && r.Birthdate.IsZero() {
		return util.ErrPlayerBirthdateIsEmpty
	}

	if r.Birthdate != nil && r.Birthdate.After(time.Now().UTC()) {
		return util.ErrPlayerBirthdateIsFutureDate
	}

	return nil
}

type UpdatePlayerUsecase interface {
	UpdatePlayer(ctx context.Context, id string, request *UpdatePlayerRequest) (*entity.Player, error)
}

type internalUpdatePlayerUsecases struct {
	ValidateGovernmentID ValidateGovernmentIDUsecase
	ValidateEmail        ValidateEmailUsecase
	ValidateAlias        ValidateAliasUsecase
}

type updatePlayerUsecase struct {
	*internalUpdatePlayerUsecases
	DBWriter database.DBWriter
	DBReader database.DBReader
}

func NewUpdatePlayerUsecase(dbWriter database.DBWriter, dbReader database.DBReader) UpdatePlayerUsecase {
	return &updatePlayerUsecase{
		DBWriter: dbWriter,
		DBReader: dbReader,
		internalUpdatePlayerUsecases: &internalUpdatePlayerUsecases{
			ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReader),
			ValidateEmail:        NewValidateEmailUsecase(dbReader),
			ValidateAlias:        NewValidateAliasUsecase(dbReader),
		},
	}
}

func (uc *updatePlayerUsecase) UpdatePlayer(ctx context.Context, id string, request *UpdatePlayerRequest) (*entity.Player, error) {
	if err := request.Validate(id); err != nil {
		log.Println(fmt.Errorf("couldn't create player. Error when validating request: %w", err))
		return nil, err
	}

	isAvailableGovernmentID, err := uc.internalUpdatePlayerUsecases.ValidateGovernmentID.IsAvailable(ctx, request.GovernmentID)
	if err != nil {
		log.Println(fmt.Errorf("couldn't update player. Error when validating government ID: %w", err))
		return nil, err
	}

	if !isAvailableGovernmentID {
		log.Println(fmt.Errorf("couldn't update player. There is another player registered with the provided government ID"))
		return nil, fmt.Errorf("couldn't update player. There is another player registered with the provided government ID")
	}

	isAvailableEmail, err := uc.internalUpdatePlayerUsecases.ValidateEmail.IsAvailable(ctx, request.Email)
	if err != nil {
		log.Println(fmt.Errorf("couldn't update player. Error when validating email: %w", err))
		return nil, err
	}

	if !isAvailableEmail {
		log.Println(fmt.Errorf("couldn't update player. There is another player registered with the provided email"))
		return nil, fmt.Errorf("couldn't update player. There is another player registered with the provided email")
	}

	isAvailableAlias, err := uc.internalUpdatePlayerUsecases.ValidateAlias.IsAvailable(ctx, request.Alias)
	if err != nil {
		log.Println(fmt.Errorf("couldn't update player. Error when validating alias: %w", err))
		return nil, err
	}

	if !isAvailableAlias {
		log.Println(fmt.Errorf("couldn't update player. There is another player registered with the provided alias"))
		return nil, fmt.Errorf("couldn't update player. There is another player registered with the provided alias")
	}

	player, err := uc.DBReader.GetPlayer(ctx, id)
	if err != nil {
		return nil, err
	}

	player.GovernmentID = request.GovernmentID
	player.Email = request.Email
	player.Alias = request.Alias
	player.FirstName = request.FirstName
	player.MiddleName = request.MiddleName
	player.LastName = request.LastName
	player.Birthdate = request.Birthdate
	player.Category = request.Category
	player.PhoneNumber = request.PhoneNumber

	updatedPlayer, err := uc.DBWriter.UpdatePlayer(ctx, player)
	if err != nil {
		log.Println(fmt.Errorf("couldn't update player. Error when attempting add user to the database request: %w", err))
		return nil, err
	}
	return updatedPlayer, nil
}
