package factory

// Call contains a Factory and all parameters required to call Factory.New.
// It is used together with CallFactories to get the result of calling multiple factories in a single call.
type Call struct {
	Factory      Factory
	Config       *Config
	Dependencies Dependencies
	Out          interface{}
}

// Calls is a slice of Call.
type Calls []Call

// CallFactories receives a slice of Call and calls the Factory.New method of each Factory.
// If a factory return an error, the function returns the error immediately.
func CallFactories(factoryCalls Calls) error {
	for _, call := range factoryCalls {
		if err := call.Factory.New(call.Config, call.Dependencies, call.Out); err != nil {
			return err
		}
	}

	return nil
}
