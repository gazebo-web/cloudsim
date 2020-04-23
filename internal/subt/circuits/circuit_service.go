package circuits

type IService interface {
	GetByName(name string) (*Circuit, error)
}

type Service struct {

}