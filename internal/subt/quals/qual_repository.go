package quals

type IRepository interface {
	GetByOwnerAndCircuit(owner, circuit string) (*Qualification, error)
}

type Repository struct {

}