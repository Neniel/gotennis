package main

import (
	"context"

	"github.com/Neniel/gotennis/customers/usecase"
	"github.com/Neniel/gotennis/lib/app"
)

func main() {
	app := app.NewApp(context.Background())

	ms := &CustomerMicroservice{
		App: app,
		Usecases: &Usecases{
			CreateTenant: usecase.NewCreateTenant(app),
			ListTenants:  usecase.NewListTenants(app),
			GetTenant:    usecase.NewGetTenant(app),
			//UpdateCustomer: usecase.NewUpdateCustomer(app),
			DeleteTenant: usecase.NewDeleteTenant(app),
		},
	}

	ms.NewAPIServer().Run()
}
