package application

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// GetTagsInstanceBase returns the base tags to identify cloud instances.
func GetTagsInstanceBase(gid simulations.GroupID) []machines.Tag {
	return []machines.Tag{
		{
			Resource: "instance",
			Map: map[string]string{
				"cloudsim_groupid":     gid.String(),
				"CloudsimGroupID":      gid.String(),
				"project":              "cloudsim",
				"Cloudsim":             "True",
				"SubT":                 "True",
				"cloudsim-application": "SubT",
			},
		},
	}
}

// GetTagsInstanceSpecific returns the specific tags to identify a single cloud instance.
func GetTagsInstanceSpecific(prefix string, gid simulations.GroupID, suffix string, clusterName, nodeType string) []machines.Tag {
	name := fmt.Sprintf("%s-%s-%s", prefix, gid.String(), suffix)
	clusterKey := fmt.Sprintf("kubernetes.io/cluster/%s", clusterName)
	return []machines.Tag{
		{
			Resource: "instance",
			Map: map[string]string{
				"Name":                       name,
				"cloudsim_groupid":           gid.String(),
				"CloudsimGroupID":            gid.String(),
				"project":                    "cloudsim",
				"Cloudsim":                   "True",
				"SubT":                       "True",
				"cloudsim-application":       "SubT",
				"cloudsim_node_type":         nodeType,
				clusterKey:                   "owned",
			},
		},
	}
}
