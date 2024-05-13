package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/entity"
	"github.com/Neniel/gotennis/util"
)

type CreatePlayerRequest struct {
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

	return nil
}

type CreatePlayerUsecase interface {
	CreatePlayer(ctx context.Context, request *CreatePlayerRequest) (*entity.Player, error)
}

type internalCreatePlayerUsecases struct {
	ValidateGovernmentID ValidateGovernmentIDUsecase
	ValidateEmail        ValidateEmailUsecase
	ValidateAlias        ValidateAliasUsecase
}

type createPlayerUsecase struct {
	*internalCreatePlayerUsecases
	DBWriter database.DBWriter
}

func NewCreatePlayerUsecase(dbWriter database.DBWriter, dbReader database.DBReader) CreatePlayerUsecase {
	return &createPlayerUsecase{
		DBWriter: dbWriter,
		internalCreatePlayerUsecases: &internalCreatePlayerUsecases{
			ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReader),
			ValidateEmail:        NewValidateEmailUsecase(dbReader),
			ValidateAlias:        NewValidateAliasUsecase(dbReader),
		},
	}
}

func (uc *createPlayerUsecase) CreatePlayer(ctx context.Context, request *CreatePlayerRequest) (*entity.Player, error) {
	if err := request.Validate(); err != nil {
		log.Println(fmt.Errorf("couldn't create player. Error when validating request: %w", err))
		return nil, err
	}

	isAvailableGovernmentID, err := uc.internalCreatePlayerUsecases.ValidateGovernmentID.IsAvailable(ctx, request.GovernmentID)
	if err != nil {
		log.Println(fmt.Errorf("couldn't create player. Error when validating government ID: %w", err))
		return nil, err
	}

	if !isAvailableGovernmentID {
		log.Println(fmt.Errorf("couldn't create player. There is another player registered with the provided government ID"))
		return nil, fmt.Errorf("couldn't create player. There is another player registered with the provided government ID")
	}

	isAvailableEmail, err := uc.internalCreatePlayerUsecases.ValidateEmail.IsAvailable(ctx, request.Email)
	if err != nil {
		log.Println(fmt.Errorf("couldn't create player. Error when validating email: %w", err))
		return nil, err
	}

	if !isAvailableEmail {
		log.Println(fmt.Errorf("couldn't create player. There is another player registered with the provided email"))
		return nil, fmt.Errorf("couldn't create player. There is another player registered with the provided email")
	}

	isAvailableAlias, err := uc.internalCreatePlayerUsecases.ValidateAlias.IsAvailable(ctx, request.Alias)
	if err != nil {
		log.Println(fmt.Errorf("couldn't create player. Error when validating alias: %w", err))
		return nil, err
	}

	if !isAvailableAlias {
		log.Println(fmt.Errorf("couldn't create player. There is another player registered with the provided alias"))
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

	player, err := uc.DBWriter.AddPlayer(ctx, newPlayer)
	if err != nil {
		log.Println(fmt.Errorf("couldn't create player. Error when attempting add user to the database request: %w", err))
		return nil, err
	}
	return player, nil
}
