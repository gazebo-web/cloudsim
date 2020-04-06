package pool

// IPool is a pool of workers that can accept jobs to be executed.
// For more details see project "github.com/panjf2000/ants".
type IPool interface {
	Serve(args interface{}) error
}
