package dependency_container

import (
	"github.com/omatheuscaetano/planus-api/internal/handler"
)

type AppHandlers struct {
	App *handler.AppHandler
}

type AppInstances struct {
	Handlers *AppHandlers
}

func InstantiateApp() *AppInstances {
	return &AppInstances{
		Handlers: &AppHandlers{
			App: handler.NewAppHandler(),
		},
	}
}
