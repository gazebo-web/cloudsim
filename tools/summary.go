package tools

import "fmt"

func GenerateSummaryFilename(groupID string) string {
	return fmt.Sprintf("%s-summary.json", groupID)
}
