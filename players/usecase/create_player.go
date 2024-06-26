package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/log"
	"github.com/Neniel/gotennis/lib/util"
)

type CreatePlayerRequest struct {
	GovernmentID string           `json:"government_id"`
	FirstName    string           `json:"first_name"`
	MiddleName   string           `json:"middle_name"`
	LastName     string           `json:"last_name"`
	Category     *entity.Category `json:"category,omitempty"`
	Birthdate    *time.Time       `json:"birthdate,omitempty"`
	PhoneNumber  string           `json:"phone_number"`
	Email        string           `json:"email"`
	Alias        *string          `json:"alias,omitempty"`
}

func (r *CreatePlayerRequest) Validate() error {
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

	if r.Alias != nil && *r.Alias == "" {
		r.Alias = nil
	}

	return nil
}

type CreatePlayer interface {
	Do(ctx context.Context, request *CreatePlayerRequest) (*entity.Player, error)
}

type internalCreatePlayer struct {
	ValidateGovernmentID ValidateGovernmentID
	ValidateEmail        ValidateEmail
	ValidateAlias        ValidateAlias
}

type createPlayer struct {
	*internalCreatePlayer
	DBWriter database.DBWriter
}

func NewCreatePlayer(dbWriter database.DBWriter, dbReader database.DBReader) CreatePlayer {
	return &createPlayer{
		DBWriter: dbWriter,
		internalCreatePlayer: &internalCreatePlayer{
			ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReader),
			ValidateEmail:        NewValidateEmailUsecase(dbReader),
			ValidateAlias:        NewValidateAliasUsecase(dbReader),
		},
	}
}

func (uc *createPlayer) Do(ctx context.Context, request *CreatePlayerRequest) (*entity.Player, error) {
	if err := request.Validate(); err != nil {
		log.Logger.Error(fmt.Errorf("couldn't create player. Error when validating request: %w", err).Error())
		return nil, err
	}

	isAvailableGovernmentID, err := uc.internalCreatePlayer.ValidateGovernmentID.IsAvailable(ctx, request.GovernmentID)
	if err != nil {
		log.Logger.Error(fmt.Errorf("couldn't create player. Error when validating government ID: %w", err).Error())
		return nil, err
	}

	if !isAvailableGovernmentID {
		log.Logger.Error(fmt.Errorf("couldn't create player. There is another player registered with the provided government ID").Error())
		return nil, fmt.Errorf("couldn't create player. There is another player registered with the provided government ID")
	}

	isAvailableEmail, err := uc.internalCreatePlayer.ValidateEmail.IsAvailable(ctx, request.Email)
	if err != nil {
		log.Logger.Error(fmt.Errorf("couldn't create player. Error when validating email: %w", err).Error())
		return nil, err
	}

	if !isAvailableEmail {
		log.Logger.Error(fmt.Errorf("couldn't create player. There is another player registered with the provided email").Error())
		return nil, fmt.Errorf("couldn't create player. There is another player registered with the provided email")
	}

	isAvailableAlias, err := uc.internalCreatePlayer.ValidateAlias.IsAvailable(ctx, request.Alias)
	if err != nil {
		log.Logger.Error(fmt.Errorf("couldn't create player. Error when validating alias: %w", err).Error())
		return nil, err
	}

	if !isAvailableAlias {
		log.Logger.Error(fmt.Errorf("couldn't create player. There is another player registered with the provided alias").Error())
		return nil, fmt.Errorf("couldn't create player. There is another player registered with the provided alias")
	}

	newPlayer := entity.NewPlayer(
		request.GovernmentID,
		request.FirstName,
		request.MiddleName,
		request.LastName,
		request.Birthdate,
		request.PhoneNumber,
		request.Email,
		request.Alias,
	)

	//uc.SystemDBWriter.AddUser()
	player, err := uc.DBWriter.AddPlayer(ctx, newPlayer)
	if err != nil {
		log.Logger.Error(fmt.Errorf("couldn't create player. Error when attempting add user to the database request: %w", err).Error())
		return nil, err
	}
	return player, nil
}
