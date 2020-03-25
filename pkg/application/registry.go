package application

func RegisterApplications(applications map[string]*IApplication, fn func() *IApplication) map[string]*IApplication {
	app := fn()
	name := (*app).Name()
	applications[name] = app
	return applications
}
