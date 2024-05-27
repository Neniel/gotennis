package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Neniel/gotennis/lib/database"

	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/log"
	"github.com/Neniel/gotennis/lib/util"
)

type PartiallyUpdatePlayerRequest struct {
	ID           string           `json:"id"`
	GovernmentID string           `json:"government_id,omitempty"`
	FirstName    string           `json:"first_name,omitempty"`
	MiddleName   string           `json:"middle_name,omitempty"`
	LastName     string           `json:"last_name,omitempty"`
	Category     *entity.Category `json:"category,omitempty"`
	Birthdate    *time.Time       `json:"birthdate,omitempty"`
	PhoneNumber  string           `json:"phone_number,omitempty"`
	Email        string           `json:"email,omitempty"`
	Alias        *string          `json:"alias,omitempty"`
}

func (r *PartiallyUpdatePlayerRequest) Validate(id string) error {
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

	if r.Alias != nil && *r.Alias == "" {
		return util.ErrPlayerAliasIsEmpty
	}

	if r.Birthdate != nil && r.Birthdate.IsZero() {
		return util.ErrPlayerBirthdateIsEmpty
	}

	if r.Birthdate != nil && r.Birthdate.After(time.Now().UTC()) {
		return util.ErrPlayerBirthdateIsFutureDate
	}

	return nil
}

type PartialltUpdatePlayer interface {
	Do(ctx context.Context, id string, request *PartiallyUpdatePlayerRequest) (*entity.Player, error)
}

type internalPartiallyUpdatePlayer struct {
	ValidateGovernmentID ValidateGovernmentID
	ValidateEmail        ValidateEmail
	ValidateAlias        ValidateAlias
}

type partiallyUpdatePlayer struct {
	*internalPartiallyUpdatePlayer
	DBWriter database.DBWriter
	DBReader database.DBReader
}

func NewPartiallyUpdatePlayer(dbWriter database.DBWriter, dbReader database.DBReader) PartialltUpdatePlayer {
	return &partiallyUpdatePlayer{
		DBWriter: dbWriter,
		DBReader: dbReader,
		internalPartiallyUpdatePlayer: &internalPartiallyUpdatePlayer{
			ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReader),
			ValidateEmail:        NewValidateEmailUsecase(dbReader),
			ValidateAlias:        NewValidateAliasUsecase(dbReader),
		},
	}
}

func (uc *partiallyUpdatePlayer) Do(ctx context.Context, id string, request *PartiallyUpdatePlayerRequest) (*entity.Player, error) {
	if err := request.Validate(id); err != nil {
		log.Logger.Error(fmt.Errorf("couldn't create player. Error when validating request: %w", err).Error())
		return nil, err
	}

	isAvailableGovernmentID, err := uc.internalPartiallyUpdatePlayer.ValidateGovernmentID.IsAvailable(ctx, request.GovernmentID)
	if err != nil {
		log.Logger.Error(fmt.Errorf("couldn't update player. Error when validating government ID: %w", err).Error())
		return nil, err
	}

	if !isAvailableGovernmentID {
		log.Logger.Error(fmt.Errorf("couldn't update player. There is another player registered with the provided government ID").Error())
		return nil, fmt.Errorf("couldn't update player. There is another player registered with the provided government ID")
	}

	isAvailableEmail, err := uc.internalPartiallyUpdatePlayer.ValidateEmail.IsAvailable(ctx, request.Email)
	if err != nil {
		log.Logger.Error(fmt.Errorf("couldn't update player. Error when validating email: %w", err).Error())
		return nil, err
	}

	if !isAvailableEmail {
		log.Logger.Error(fmt.Errorf("couldn't update player. There is another player registered with the provided email").Error())
		return nil, fmt.Errorf("couldn't update player. There is another player registered with the provided email")
	}

	isAvailableAlias, err := uc.internalPartiallyUpdatePlayer.ValidateAlias.IsAvailable(ctx, request.Alias)
	if err != nil {
		log.Logger.Error(fmt.Errorf("couldn't update player. Error when validating alias: %w", err).Error())
		return nil, err
	}

	if !isAvailableAlias {
		log.Logger.Error(fmt.Errorf("couldn't update player. There is another player registered with the provided alias").Error())
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
		log.Logger.Error(fmt.Errorf("couldn't update player. Error when attempting add user to the database request: %w", err).Error())
		return nil, err
	}
	return updatedPlayer, nil
}
