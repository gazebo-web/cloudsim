package pool

// Factory is a function responsible of initializing and returning a JobPool.
// Dev note: we created the PoolFactory abstraction to allow tests use
// synchronic pools.
type Factory func(poolSize int, jobF func(interface{})) (IPool, error)
