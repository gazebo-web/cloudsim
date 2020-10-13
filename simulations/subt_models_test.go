package simulations

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubTCreateSimulation_robotImageBelongsToECROwner(t *testing.T) {
	test := func(valid bool, owner string, robotImages []string) {
		cs := &SubTCreateSimulation{
			CreateSimulation: CreateSimulation{
				Owner: owner,
			},
			RobotImage: robotImages,
		}
		assert.Equal(t, valid, cs.robotImagesBelongToECROwner())
	}

	subTRedTeam := "OSRF SubT RedTeam"

	// If all images are stored in the team's ECR repo this should succeed
	test(true, subTRedTeam, []string{
		"200670743174.dkr.ecr.us-east-1.amazonaws.com/osrf_subt_redteam:subt_seed",
		"200670743174.dkr.ecr.us-east-1.amazonaws.com/osrf_subt_redteam:robotika_unittest",
		"200670743174.dkr.ecr.us-east-2.amazonaws.com/osrf_subt_redteam:virtual_stix_arl",
	})

	// A non-ECR repo should always succeed
	test(true, subTRedTeam, []string{
		"https://hub.docker.com/r/osrf/subt-virtual-testbed",
	})

	// A mix between images stored in the team's ECR repo and in non-ECR repos should succeed
	test(true, subTRedTeam, []string{
		"200670743174.dkr.ecr.us-east-1.amazonaws.com/osrf_subt_redteam:subt_seed",
		"200670743174.dkr.ecr.us-east-1.amazonaws.com/osrf_subt_redteam:robotika_unittest",
		"https://hub.docker.com/r/osrf/subt-virtual-testbed",
	})

	// If an image is stored in an ECR repo from another team then it should fail
	test(false, subTRedTeam, []string{
		"200670743174.dkr.ecr.us-east-1.amazonaws.com/osrf_subt_redteam:robotika_unittest",
		"200670743174.dkr.ecr.us-east-1.amazonaws.com/subt_sim:cloudsim_sim_latest",
		"200670743174.dkr.ecr.us-east-2.amazonaws.com/osrf_subt_redteam:virtual_stix_arl",
	})
}
