package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"path/filepath"
)

// UploadLogs is a job in charge of uploading simulation logs.
var UploadLogs = &actions.Job{
	Name:       "upload-simulation-logs",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    uploadLogs,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// uploadLogsScript is used to set up parameters for the copy script.
type uploadLogsScript struct {
	Target   string
	Filename string
	Bucket   string
}

// uploadLogs is the execute function of the UploadLogs job.
func uploadLogs(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	// If logs aren't enabled, continue with the rest of the jobs
	if !s.Platform().Store().Ignition().LogsCopyEnabled() {
		return s, nil
	}

	// Get the simulation
	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	// Get the robot list
	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	// Get the namespace
	ns := s.Platform().Store().Orchestrator().Namespace()

	// Get the bucket to save logs into
	logsBucket := s.Platform().Store().Ignition().LogsBucket()

	// Attempt to upload robot logs
	for i, robot := range robots {
		robotID := subtapp.GetRobotID(i)
		name := subtapp.GetPodNameCommsBridgeCopy(s.GroupID, robotID)
		res, err := s.Platform().Orchestrator().Pods().Get(name, ns)
		if err != nil {
			continue
		}

		filename := subtapp.GetCommsBridgeLogsFilename(s.GroupID, robot.GetName())
		bucket := filepath.Join(logsBucket, subtapp.GetSimulationLogKey(s.GroupID, *sim.GetOwner()))

		scriptParams := uploadLogsScript{
			Target:   s.Platform().Store().Ignition().SidecarContainerLogsPath(),
			Filename: filename,
			Bucket:   s.Platform().Storage().PrepareAddress(bucket, filename),
		}

		exec := s.Platform().Orchestrator().Pods().Exec(res)
		containerName := subtapp.GetContainerNameCommsBridgeCopy()

		err = uploadSingleLogs(exec, containerName, "simulations/scripts/copy_to_s3.sh", scriptParams)
		if err != nil {
			return nil, err
		}
	}

	// Get gazebo copy pod name
	name := subtapp.GetPodNameGazeboServerCopy(s.GroupID)

	res, err := s.Platform().Orchestrator().Pods().Get(name, ns)
	if err != nil {
		return nil, err
	}

	filename := subtapp.GetGazeboLogsFilename(s.GroupID)
	bucket := filepath.Join(logsBucket, subtapp.GetSimulationLogKey(s.GroupID, *sim.GetOwner()))

	scriptParams := uploadLogsScript{
		Target:   s.Platform().Store().Ignition().SidecarContainerLogsPath(),
		Filename: filename,
		Bucket:   s.Platform().Storage().PrepareAddress(bucket, filename),
	}

	exec := s.Platform().Orchestrator().Pods().Exec(res)
	containerName := subtapp.GetContainerNameGazeboServerCopy()
	err = uploadSingleLogs(exec, containerName, "simulations/scripts/copy_to_s3.sh", scriptParams)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// uploadSingleLogs is a helper function in charge of running a certain script in scriptFilepath with the given scriptParams.
// It will run this script inside the containerName using the pods.Executor implementation passed as an argument.
// It will return an error if parsing the script template or executing the script returns an error.
func uploadSingleLogs(exec pods.Executor, containerName string, scriptFilepath string, scriptParams uploadLogsScript) error {
	script, err := ign.ParseTemplate(scriptFilepath, scriptParams)
	if err != nil {
		return err
	}

	err = exec.Script(containerName, script)
	if err != nil {
		return err
	}

	return nil
}
