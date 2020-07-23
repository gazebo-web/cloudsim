package tracks

import (
	"encoding/json"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/structs"
)

// CreateTrackInput is an input for creating a new track.
type CreateTrackInput struct {
	Name        string `json:"name" validate:"required,gt=10"`
	Image       string `json:"image" validate:"required"`
	BridgeImage string `json:"bridge_image" validate:"required"`
	// Topic used to track general stats of the simulation (runtime, sim runtime, etc.)
	StatsTopic string `json:"stats_topic" validate:"required"`
	// Topic used to track when the simulation officially starts and ends
	WarmupTopic string `json:"warmup_topic" validate:"required"`
	// Maximum number of allowed "simulation seconds" for each world. 0 means unlimited.
	MaxSimSeconds int `json:"max_sim_seconds" validate:"required"`
	// Public makes a track available for launching directly.
	// Tracks that are not public can only be launched as part of a Circuit.
	Public bool `json:"public" validate:"required"`
}

// Value converts the current struct into a map[string]interface{}
func (c CreateTrackInput) Value() (interface{}, error) {
	return CreateTrackFromInput(c), nil
}

// NewCreateTrackInput creates a new CreateTrackInput from the given []byte.
func NewCreateTrackInput(body []byte) (domain.DTO, error) {
	var createTrackInput CreateTrackInput
	err := json.Unmarshal(body, &createTrackInput)
	if err != nil {
		return nil, err
	}
	return &createTrackInput, nil
}

// UpdateTrackInput is an input for updating an existent track.
type UpdateTrackInput struct {
	Name        *string `json:"name,omitempty" validate:"gt=10"`
	Image       *string `json:"image,omitempty"`
	BridgeImage *string `json:"bridge_image,omitempty"`
	// Topic used to track general stats of the simulation (runtime, sim runtime, etc.)
	StatsTopic *string `json:"stats_topic,omitempty"`
	// Topic used to track when the simulation officially starts and ends
	WarmupTopic *string `json:"warmup_topic,omitempty"`
	// Maximum number of allowed "simulation seconds" for each world. 0 means unlimited.
	MaxSimSeconds *int `json:"max_sim_seconds,omitempty"`
	// Public makes a track available for launching directly.
	// Tracks that are not public can only be launched as part of a Circuit.
	Public *bool `json:"public,omitempty"`
}

// Value converts the current struct into a map[string]interface{}
func (u UpdateTrackInput) Value() (interface{}, error) {
	return structs.ToMap(&u)
}

// NewUpdateTrackInput initializes a new UpdateTrackInput from the given []byte.
func NewUpdateTrackInput(body []byte) (domain.DTO, error) {
	var updateTrackInput UpdateTrackInput
	err := json.Unmarshal(body, &updateTrackInput)
	if err != nil {
		return nil, err
	}
	return &updateTrackInput, nil
}
