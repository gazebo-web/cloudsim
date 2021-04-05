package jobs

import (
	"bytes"
	"encoding/json"
	"github.com/jinzhu/gorm"
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// SendSummaryEmail is a job in charge of sending an email to participants with the simulation's statistics and score.
var SendSummaryEmail = &actions.Job{
	Name:       "send-summary-email",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    sendSummaryEmail,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// sendSummaryEmail is the execute function of the SendSummaryEmail job.
func sendSummaryEmail(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	// Get the simulation
	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	// Only send summary emails for single simulations.
	if !sim.IsKind(simulations.SimSingle) {
		return s, nil
	}

	// Get the user info
	user, em := s.Services().Users().GetUserFromUsername(sim.GetCreator())
	if em != nil {
		return nil, em.BaseError
	}

	// Generate list of recipients
	var recipients []string

	// Append default recipients
	recipients = append(recipients, s.Platform().Store().Ignition().DefaultRecipients()...)

	// Add user's email.
	if user.Email != nil {
		recipients = append(recipients, *user.Email)
	}

	// Get default sender
	sender := s.Platform().Store().Ignition().DefaultSender()

	// If there's an owner assigned, add the organization email
	if sim.GetOwner() != nil {
		org, em := s.Services().Users().GetOrganization(*sim.GetOwner())
		if em != nil {
			return nil, em.BaseError
		}

		if org.Email != nil {
			recipients = append(recipients, *org.Email)
		}
	}

	subtSim := sim.(subt.Simulation)

	var marshaledSummary bytes.Buffer
	b, err := json.MarshalIndent(s.Summary, "", "  ")
	if err != nil {
		return nil, err
	}

	_, err = marshaledSummary.Write(b)
	if err != nil {
		return nil, err
	}

	// Send the email
	err = s.Platform().EmailSender().Send(recipients, sender, "Simulation summary", "simulations/email-templates/simulation_summary.gohtml", map[string]interface{}{
		"Name":    *user.Name,
		"Circuit": subtSim.GetTrack(),
		"SimName": subtSim.GetName(),
		"GroupID": subtSim.GetGroupID(),
		"Summary": marshaledSummary,
		"RunData": s.Summary.RunData,
	})
	if err != nil {
		return nil, err
	}

	return s, nil
}
