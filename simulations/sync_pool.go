package simulations

import (
	"github.com/panjf2000/ants"
	"sync"
)

// SynchronicPool is used to schedule Jobs and make the caller wait (block)
// until they are executed.
// Dev note: it wraps an ants.job pool and uses a sync.WaitGroup to make the
// calling routine block until the job finishes.
type SynchronicPool struct {
	antsp *ants.PoolWithFunc
}

// syncPoolPair is an internal type of wrap several arguments in a struct
type syncPoolPair struct {
	wg       *sync.WaitGroup
	origArgs interface{}
}

// Serve submits a task to pool.
// Implements simulations.JobPool interface.
func (sp *SynchronicPool) Serve(args interface{}) error {

	// Creates a WaitGroup to make the routine block on it.
	var wg sync.WaitGroup
	wg.Add(1)
	pair := &syncPoolPair{
		wg:       &wg,
		origArgs: args,
	}
	// Delegate to the internal ants.pool
	if err := sp.antsp.Serve(pair); err != nil {
		return err
	}
	// This call will block until the ants job finishes.
	wg.Wait()
	return nil
}

// SynchronicPoolFactory is a Factory function that creates a new SynchronicPool using
// the given arguments.
func SynchronicPoolFactory(poolSize int, jobF func(interface{})) (JobPool, error) {

	jobWithMultipleArgs := func(payload interface{}) {
		// This is a wrapper on top of the original job function
		// that receives a WaitGroup and mark it as Done after running the job func.
		pair, ok := payload.(*syncPoolPair)
		if !ok {
			return
		}
		defer func() {
			// check for panic
			if p := recover(); p != nil {
				pair.wg.Done()
				panic(p) // re-throw panic
			}
		}()
		jobF(pair.origArgs)
		pair.wg.Done()
	}

	antspool, err := ants.NewPoolWithFunc(poolSize, jobWithMultipleArgs)
	if err != nil {
		return nil, err
	}

	syncPool := SynchronicPool{
		antsp: antspool,
	}
	return &syncPool, nil
}

/////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////
