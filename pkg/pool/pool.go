package pool

import "github.com/panjf2000/ants"

// Event registers a single pool event listener that will receive
// notifications any time a pool worker "finishes" its job (either with result or error).
type Event int

// Pool is a pool of workers that can accept jobs to be executed.
// For more details see project "github.com/panjf2000/ants".
type Pool interface {
	Serve(args interface{}) error
}

// NewPool is the default implementation of the PoolFactory interface.
// It creates an ants.PoolWithFunc.
func NewPool(poolSize int, jobF func(interface{})) (Pool, error) {
	return ants.NewPoolWithFunc(poolSize, jobF)
}
