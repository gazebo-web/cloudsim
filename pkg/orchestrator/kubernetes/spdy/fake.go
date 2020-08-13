package spdy

import (
	"k8s.io/client-go/tools/remotecommand"
	"net/url"
)

// Fake is a Initializer implementation.
type Fake struct {
	Calls int
	Stdin []byte
}

// Stream mocks the remotecommand.Executor Stream method.
func (f Fake) Stream(options remotecommand.StreamOptions) error {
	f.Calls++
	if options.Stdout != nil {
		_, err := options.Stdout.Write([]byte("stdout-test"))
		if err != nil {
			return err
		}
	}
	if options.Stderr != nil {
		_, err := options.Stdout.Write([]byte("stderr-test"))
		if err != nil {
			return err
		}
	}
	if options.Stdin != nil {
		_, err := options.Stdin.Read(f.Stdin)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewSPDYExecutor initializes a new remotecommand.Executor using Fake.
func (f Fake) NewSPDYExecutor(method string, url *url.URL) (remotecommand.Executor, error) {
	return &Fake{}, nil
}

// NewSPDYFakeInitializer initializes a new Fake.
func NewSPDYFakeInitializer() *Fake {
	return &Fake{}
}
