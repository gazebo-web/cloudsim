package quals

type IRepository interface {
	GetByOwnerAndCircuit(owner, circuit string) (*Qualification, error)
}