package simulations

// IController represents a group of methods to expose in the API Rest.
type IController interface {
	Start()
	LaunchHeld()
	Restart()
	Shutdown()
	GetAll()
	Get()
	GetDownloadableLogs()
	GetLiveLogs()
}

// Controller is an IController implementation.
type Controller struct {

}
