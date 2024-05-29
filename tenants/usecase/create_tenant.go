package usecase

import (
	"context"
	"fmt"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/log"
	"github.com/Neniel/gotennis/lib/util"
)

type CreateTenant interface {
	Do(ctx context.Context, request *CreateTenantRequest) (*entity.Tenant, error)
}

type createTenant struct {
	DBWriter database.DBWriter
}

func NewCreateTenant(app app.IApp) CreateTenant {
	systemMongoDBClient := app.GetSystemMongoDBClient()
	return &createTenant{
		DBWriter: database.NewDatabaseWriter(systemMongoDBClient.MongoDBClient, systemMongoDBClient.DatabaseName),
	}
}

type CreateTenantRequest struct {
	Name string `json:"name"`
}

func (r *CreateTenantRequest) Validate() error {
	if r.Name == "" {
		return util.ErrCategoryNameIsEmpty
	}

	return nil
}

func (uc *createTenant) Do(ctx context.Context, request *CreateTenantRequest) (*entity.Tenant, error) {
	if err := request.Validate(); err != nil {
		log.Logger.Info(fmt.Errorf("could not create customer: %w", err).Error())
		return nil, err
	}

	customer := &entity.Tenant{
		Name: request.Name,
	}

	newCustomer, err := uc.DBWriter.AddTenant(ctx, customer)
	if err != nil {
		log.Logger.Info(fmt.Errorf("could not create customer: %w", err).Error())
		return nil, err
	}
	return newCustomer, nil
}
