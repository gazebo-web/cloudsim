package application

// RegisterApplications adds a given application to the platform.
// Returns the list of applications
func RegisterApplications(applications map[string]*IApplication, app *IApplication) map[string]*IApplication {
	if app == nil {
		panic("Invalid application")
	}
	name := (*app).Name()
	applications[name] = app
	return applications
}
