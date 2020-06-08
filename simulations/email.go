package simulations

import (
	"errors"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/cloudsim/globals"
	"bytes"
	"encoding/json"
)

// SendEmail sends an email to a specific recipient. If the recipient is nil,
// then the default recipient defined in the IGN_FLAGS_EMAIL_TO env var will be
// used.
func SendEmail(recipient *[]string, sender *string, subject string, templateFilename string,
	templateData interface{}) *ign.ErrMsg {
	if recipient == nil {
		recipient = &globals.DefaultEmailRecipients
	}
	if sender == nil {
		sender = &globals.DefaultEmailSender
	}
	// If the sender or recipient are not defined, then don't send the email
	if (recipient != nil && len(*recipient) == 0) || (sender != nil && *sender == "") {
		return nil
	}

	// Prepare the template
	content, err := ign.ParseHTMLTemplate(templateFilename, templateData)
	if err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	// Send the email
	for _, r := range *recipient {
		err = ign.SendEmail(*sender, r, subject, content)
		if err != nil {
			return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
		}
	}

	return nil
}

// SendSimulationSummaryEmail sends a summary email to the user that created the simulation
func SendSimulationSummaryEmail(dep *SimulationDeployment, summary AggregatedSubTSimulationValues) *ign.ErrMsg {

	if dep.SummaryProcessed {
		return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, errors.New("summary has already been processed"))
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
	templateFilename := "simulations/email-templates/simulation_summary.html"
	templateData := struct {
		Name    string
		Circuit string
		SimName string
		GroupID string
		Summary string
	}{
		Name:    *user.Name,
		SimName: *dep.Name,
		Circuit: *dep.ExtraSelector,
		GroupID: *dep.GroupID,
		Summary: marshaledSummary.String(),
	}

	// Set the list of recipients
	recipients := make([]string, 0)
	// Include default recipients
	recipients = append(recipients, globals.DefaultEmailRecipients...)
	// Add the user email to the list
	if user.Email != nil {
		recipients = append(recipients, *user.Email)
	}
	// Add the organization email to the list if it differs from the user's
	if dep.Owner != nil {
		org, em := globals.UserAccessor.GetOrganization(*dep.Owner)
		if em != nil {
			return em
		}
		if org.Email != nil && (user.Email == nil || *org.Email != *user.Email) {
			recipients = append(recipients, *org.Email)
		}
	}

	if em := SendEmail(&recipients, nil, "Simulation summary", templateFilename, templateData); em != nil {
		return em
	}

	return nil
}
