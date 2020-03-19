package tools

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	"strings"
)

// EnvVarToSlice reads the contents of an environment variable and splits it into an array of strings using comma as
// the separator.
func EnvVarToSlice(envVar string) []string {
	s, _ := ign.ReadEnvVar(envVar)
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
