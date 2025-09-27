package dependency_container

import "github.com/omatheuscaetano/planus-api/internal/handler"

type UserHandlers struct {
	User *handler.UserHandler
}

type UserInstances struct {
	Handlers *UserHandlers
}

func InstantiateUser() *UserInstances {
	return &UserInstances{
		Handlers: &UserHandlers{
			User: handler.NewUserHandler(),
		},
	}
}
