package dependency_container

import (
	"github.com/omatheuscaetano/planus-api/internal/handler"
	"github.com/omatheuscaetano/planus-api/internal/service"
)

type UserHandlers struct {
	User *handler.UserHandler
}

type UserInstances struct {
	Handlers *UserHandlers
}

func InstantiateUser() *UserInstances {
	return &UserInstances{
		Handlers: &UserHandlers{
			User: handler.NewUserHandler(service.NewUserService()),
		},
	}
}
