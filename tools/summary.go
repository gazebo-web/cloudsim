package tools

import "fmt"

// GenerateSummaryFilename gets the summary filename for the given group id.
func GenerateSummaryFilename(groupID string) string {
	return fmt.Sprintf("%s-summary.json", groupID)
}
