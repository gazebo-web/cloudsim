package simulator

type IController interface {}

type Controller struct {
	Service IService
}
