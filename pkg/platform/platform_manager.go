package platform

type IPlatformManager interface {
	StartSimulation()
	LaunchSimulation()
	StopSimulation()
	RestartSimulation()
}