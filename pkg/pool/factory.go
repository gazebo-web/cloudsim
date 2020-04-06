package pool

type Factory func(poolSize int, jobF func(interface{})) (IPool, error)
