package dependency_container

import (
	mongo "github.com/omatheuscaetano/planus-api/pkg/db"
)

type InstancesConfig struct {
	Mongo *mongo.MongoDB
}

type Instances struct {
	Auth *AuthInstances
	App  *AppInstances
	User *UserInstances
}

func InstantiateAll(c *InstancesConfig) *Instances {
	return &Instances{
		Auth: InstantiateAuth(c),
		App: InstantiateApp(),
		User: InstantiateUser(c),
	}
}
