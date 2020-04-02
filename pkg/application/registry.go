package application

// RegisterApplications adds a given application to the platform.
// Returns the list of applications
func RegisterApplication(applications *map[string]IApplication, app IApplication) {
	if app == nil || applications == nil {
		panic("Invalid application")
	}
	name := app.Name()
	(*applications)[name] = app
}
