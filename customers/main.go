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
			CreateCustomer: usecase.NewCreateCustomer(app),
			ListCustomers:  usecase.NewListCustomers(app),
			GetCustomer:    usecase.NewGetCustomer(app),
			UpdateCustomer: usecase.NewUpdateCustomer(app),
			DeleteCustomer: usecase.NewDeleteCustomer(app),
		},
	}

	ms.NewAPIServer().Run()
}
