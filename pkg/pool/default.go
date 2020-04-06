package pool

import "github.com/panjf2000/ants"

// DefaultFactory is the default implementation of the PoolFactory interface.
// It creates an ants.PoolWithFunc.
func DefaultFactory(poolSize int, jobF func(interface{})) (IPool, error) {
	return ants.NewPoolWithFunc(poolSize, jobF)
}
