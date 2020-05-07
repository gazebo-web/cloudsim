package robots

type IService interface {}

type Service struct {}

func NewService() IService {
	var s IService
	s = &Service{}
	return s
}