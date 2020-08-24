package workers

import "fmt"

// Worker is a function that will run in a separated thread.
// The worker will be waiting for a given payload, and will perform an operation on it.
type Worker func(payload interface{})

// Launch is a worker that launches a simulation. It receives a LaunchInput as a payload.
// If the payload is wrong, it will panic.
func Launch(payload interface{}) {
	input, ok := payload.(LaunchInput)
	if !ok {
		panic("Wrong input")
	}
	fmt.Println(fmt.Sprintf("Launch simulation [%s]", input.GroupID))
}

// Terminate is a worker that terminates a simulation. It receives a TerminateInput as a payload.
// If the payload is wrong, it will panic.
func Terminate(payload interface{}) {
	input, ok := payload.(TerminateInput)
	if !ok {
		panic("Wrong input")
	}
	fmt.Println(fmt.Sprintf("Terminate simulation [%s]", input.GroupID))
}
