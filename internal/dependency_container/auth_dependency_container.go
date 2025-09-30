package dependency_container

import (
	"github.com/omatheuscaetano/planus-api/internal/handler"
	"github.com/omatheuscaetano/planus-api/internal/service"
)

type AuthHandlers struct {
	Auth *handler.AuthHandler
}

type AuthInstances struct {
	Handlers *AuthHandlers
}

func NewAuthService(_ *InstancesConfig, svc *service.UserService) *service.AuthService {
	return service.NewAuthService(svc)
}

func NewAuthHandler(_ *InstancesConfig, svc *service.AuthService) *handler.AuthHandler {
	return handler.NewAuthHandler(svc)
}


func InstantiateAuth(c *InstancesConfig) *AuthInstances {
	return &AuthInstances{
		Handlers: &AuthHandlers{
			Auth: NewAuthHandler(c, NewAuthService(c, NewUserService(c, NewUserMongoRepository(c)))),
		},
	}
}
