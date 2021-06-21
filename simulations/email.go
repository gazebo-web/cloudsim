package simulations

import (
	"bytes"
	"encoding/json"
	"errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/globals"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/email"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// SendSimulationSummaryEmail sends a summary email to the user that created the simulation
// `summary` contains summary information for the run. It should be passed for all types of simulation.
// `runData` contains specific simulation run information. It should only be passed for single sims.
func SendSimulationSummaryEmail(e email.Sender, dep *SimulationDeployment, summary AggregatedSubTSimulationValues,
	runData *string) *ign.ErrMsg {

	// Do not send emails for simulations that have already been processed
	if dep.Processed {
		return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, errors.New("simulation has already been processed"))
	}

	var marshaledSummary bytes.Buffer
	b, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}
	marshaledSummary.Write(b)

	// Get the org and user that started the simulation
	user, em := globals.UserAccessor.GetUserFromUsername(*dep.Creator)
	if em != nil {
		return em
	}

	// Set template info
	templateFilename := "simulations/email-templates/simulation_summary.gohtml"
	templateData := struct {
		Name    string
		Circuit string
		SimName string
		GroupID string
		Summary string
		RunData *string
	}{
		Name:    *user.Name,
		SimName: *dep.Name,
		Circuit: *dep.ExtraSelector,
		GroupID: *dep.GroupID,
		Summary: marshaledSummary.String(),
		RunData: runData,
	}

	// Set the list of recipients
	recipients := make([]string, 0)
	// Include default recipients
	recipients = append(recipients, globals.DefaultEmailRecipients...)
	// Add the user email to the list
	if user.Email != nil {
		recipients = append(recipients, *user.Email)
	}
	// Add the organization email to the list if the organization exists and it differs from the user's
	if dep.Owner != nil {
		org, em := globals.UserAccessor.GetOrganization(*dep.Owner)
		if em == nil && org.Email != nil && (user.Email == nil || *org.Email != *user.Email) {
			recipients = append(recipients, *org.Email)
		}
	}

	err = e.Send(recipients, globals.DefaultEmailSender, "Simulation summary", templateFilename, templateData)
	if err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	return nil
}
