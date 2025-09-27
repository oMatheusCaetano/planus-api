package dependency_container

type Instances struct {
	App *AppInstances
	User *UserInstances
}

func InstantiateAll() *Instances {
	return &Instances{
		App: InstantiateApp(),
		User: InstantiateUser(),
	}
}
