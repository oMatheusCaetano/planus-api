package dependency_container

import (
	"github.com/omatheuscaetano/planus-api/internal/handler"
	"github.com/omatheuscaetano/planus-api/internal/repository"
	"github.com/omatheuscaetano/planus-api/internal/service"
)

type UserHandlers struct {
	User *handler.UserHandler
}

type UserInstances struct {
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


	return &UserInstances{
		Handlers: &UserHandlers{
			User: NewUserHandler(c, NewUserService(c, NewUserMongoRepository(c))),
		},
	}
}
