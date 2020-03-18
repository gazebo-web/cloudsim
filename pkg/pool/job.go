package pool

// IJob is a pool of Jobs that can accept jobs to be executed.
// For more details see project "github.com/panjf2000/ants".
type IJob interface {
	Serve(args interface{}) error
}