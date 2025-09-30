package dependency_container

import (
	"github.com/omatheuscaetano/planus-api/internal/handler"
	"github.com/omatheuscaetano/planus-api/internal/repository"
	"github.com/omatheuscaetano/planus-api/internal/service"
)

type UserServices struct {
	User *service.UserService
}

type UserHandlers struct {
	User *handler.UserHandler
}

type UserInstances struct {
	Services *UserServices
	Handlers *UserHandlers
}

func NewUserMongoRepository(c *InstancesConfig) repository.UserRepository {
	return repository.NewUserMongoRepository(c.Mongo)
}

func NewUserService(_ *InstancesConfig, repo repository.UserRepository) *service.UserService {
	return service.NewUserService(repo)
}

func NewUserHandler(_ *InstancesConfig, svc *service.UserService) *handler.UserHandler {
	return handler.NewUserHandler(svc)
}

func InstantiateUser(c *InstancesConfig) *UserInstances {
	services := &UserServices{
		User: NewUserService(c, NewUserMongoRepository(c)),
	}

	handlers := &UserHandlers{
		User: NewUserHandler(c, services.User),
	}

	return &UserInstances{
		Services: services,
		Handlers: handlers,
	}
}
