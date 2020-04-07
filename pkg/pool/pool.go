package pool

import "github.com/panjf2000/ants"

// IPool is a pool of workers that can accept jobs to be executed.
// For more details see project "github.com/panjf2000/ants".
type IPool interface {
	Serve(args interface{}) error
}


// NewPool is the default implementation of the PoolFactory interface.
// It creates an ants.PoolWithFunc.
func NewPool(poolSize int, jobF func(interface{})) (IPool, error) {
	return ants.NewPoolWithFunc(poolSize, jobF)
}
