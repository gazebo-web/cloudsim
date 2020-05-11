package simulations

type SimulationCreateInput interface {
	Input() SimulationCreate
}

type SimulationCreate struct {
	Name  string `json:"name" validate:"required,min=3,alphanum" form:"name"`
	Owner string `json:"owner" form:"owner"`
	// The docker image(s) that will be used for the Field Computer(s)
	Image       []string `json:"image" form:"image"`
	Platform    string   `json:"platform" form:"platform"`
	Application string   `json:"application" form:"application"`
	Private     *bool    `json:"private" validate:"omitempty" form:"private"`
	// When shutting down simulations, stop EC2 instances instead of terminating them. Requires admin privileges.
	StopOnEnd *bool `json:"stop_on_end" validate:"omitempty" form:"stop_on_end"`
	// Extra: it is expected that this field will be set by the Application logic
	// and not From form values.
	Extra *string `form:"-"`
	// ExtraSelector: it is expected that this field will be set by the Application logic
	// and not from Form values.
	ExtraSelector *string `form:"-"`
	// TODO: This is a field specific to SubT. As such this is a temporary field
	//  that should be included in the same separate table where Extra and
	//  ExtraSelector should reside.
	// Contains the names of all robots in the simulation in a comma-separated list.
	Robots *string `form:"-"`
}

func (sc SimulationCreate) Input() SimulationCreate {
	return sc
}

type SimulationCreateOutput interface {
	Output() Simulation
}
