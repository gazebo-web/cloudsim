package application

// IApplication describes a set of methods for an Application.
type IApplication interface {
	Name() string
}

type Application struct {}

func (app *Application) Name() string {
	panic("Name should be implemented by the application")
}