package dependency_container

type Instances struct {
	App *AppInstances
}

func InstantiateAll() *Instances {
	return &Instances{
		App: InstantiateApp(),
	}
}
