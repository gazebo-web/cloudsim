package simulations

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

type Controller struct {

}
