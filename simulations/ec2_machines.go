package simulations

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/caarlos0/env"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/ign-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"sync"
	"time"
)

// This file defines functions to launch and terminate AWS EC2 instances for Cloudsim.
// It also keeps track of created instances, and updates the MachineInstance
// records in the DB based on their status.

const (
	nodeLabelKeyGroupID          = "cloudsim_groupid"
	nodeLabelKeyCloudsimNodeType = "cloudsim_node_type"
	nodeLabelKeySubTRobotName    = "robot_name"
)

// MaxAWSRetries holds how many retries will be done against AWS. It is a var
// to allow tests to change this value.
var MaxAWSRetries = 8

type awsConfig struct {
	NamePrefix string `env:"AWS_INSTANCE_NAME_PREFIX,required"`
	// ShouldTerminateInstances is used to define if to 'Stop' or 'Terminate' EC2 instances
	// when deleting the nodes. Default value is to Terminate the instances.
	// You can change this value with env var `EC2_NODE_MGR_TERMINATE_INSTANCES`.
	ShouldTerminateInstances bool   `env:"EC2_NODE_MGR_TERMINATE_INSTANCES" envDefault:"true"`
	IamInstanceProfile       string `env:"AWS_IAM_INSTANCE_PROFILE_ARN" envDefault:"arn:aws:iam::200670743174:instance-profile/aws-eks-role-cloudsim-worker"`
}

type ec2Config struct {
	// ClusterName contains the name of the cluster EC2 instances will join.
	ClusterName string `env:"AWS_CLUSTER_NAME,required"`
	// Subnets is a slice of AWS subnet IDs where to launch simulations (Example: subnet-1270518251)
	Subnets []string `env:"IGN_EC2_SUBNETS,required" envSeparator:","`
	// AvailabilityZones is a slice of AWS availability zones where to launch simulations. (Example: us-east-1a)
	AvailabilityZones                  []string `env:"IGN_EC2_AVAILABILITY_ZONES,required" envSeparator:","`
	OnDemandCapacityReservationEnabled bool     `env:"IGN_EC2_ODCR_ENABLED" envDefault:"false"`
	// AvailableEC2Machines is the maximum number of machines that Cloudsim can have running at a single time.
	AvailableEC2Machines int `env:"IGN_EC2_MACHINES_LIMIT" envDefault:"-1"`
}

// Ec2Client is an implementation of NodeManager interface. It is the client to use
// when creating AWS EC2 instances for k8 cluster nodes.
type Ec2Client struct {
	awsCfg awsConfig
	ec2Cfg ec2Config
	// ec2 clients are safe to use concurrently.
	ec2Svc ec2iface.EC2API
	// Mutex to ensure AWS resource availability checks are not invalidated by other workers
	lockRunInstances sync.Mutex
	// A reference to the kubernetes client
	clientset kubernetes.Interface
	platforms map[string]PlatformType
	// availabilityZoneIndex holds the value of the latest availability zone index used to launch a simulation in AWS.
	availabilityZoneIndex int
}

// PlatformType is used to tailor an instance that is being created.
type PlatformType interface {
	getPlatformName() string
	// setupEC2InstanceSpecifics is invoked by the EC2 NodeManager to describe the needed EC2 instance details.
	setupEC2InstanceSpecifics(ctx context.Context, s *Ec2Client, tx *gorm.DB, dep *SimulationDeployment,
		template *ec2.RunInstancesInput) ([]*ec2.RunInstancesInput, error)
}

// NewEC2Client creates a new client to interact with EC2 machines.
func NewEC2Client(ctx context.Context, kcli kubernetes.Interface, ec2Svc ec2iface.EC2API) (*Ec2Client, error) {
	logger(ctx).Info("Creating ec2 Nodes Manager")
	ec := Ec2Client{}
	ec.platforms = map[string]PlatformType{}

	// Read configuration from environment
	ec.awsCfg = awsConfig{}
	if err := env.Parse(&ec.awsCfg); err != nil {
		return nil, err
	}

	ec.ec2Cfg = ec2Config{}
	if err := env.Parse(&ec.ec2Cfg); err != nil {
		return nil, err
	}

	avalabilityZoneSize := len(ec.ec2Cfg.AvailabilityZones)
	if len(ec.ec2Cfg.Subnets) != avalabilityZoneSize {
		return nil, errors.New("Subnet and AZ list length mismatch")
	}

	ec.clientset = kcli

	ec.ec2Svc = ec2Svc
	return &ec, nil
}

// Stop stops this EC2 client
func (s *Ec2Client) Stop() {
	// nothing to do at the moment
}

// RegisterPlatform registers a new Platform type.
func (s *Ec2Client) RegisterPlatform(ctx context.Context, p PlatformType) {
	logger(ctx).Info(fmt.Sprintf("EC2 Nodes Manager - Registered new platform [%s]", p.getPlatformName()))
	s.platforms[p.getPlatformName()] = p
}

// CloudMachinesList returns a paginated list with all cloud machines.
// @public
func (s *Ec2Client) CloudMachinesList(ctx context.Context, p *ign.PaginationRequest, tx *gorm.DB,
	byStatus *MachineStatus, invertStatus bool, groupID *string, application *string) (*MachineInstances, *ign.PaginationResult, *ign.ErrMsg) {

	// Create the DB query
	var machines MachineInstances

	q := tx.Model(&MachineInstance{})
	if byStatus != nil {
		if invertStatus {
			q = q.Where("last_known_status != ?", byStatus.ToStringPtr())
		} else {
			q = q.Where("last_known_status = ?", byStatus.ToStringPtr())
		}
	}

	if application != nil {
		q = q.Where("application = ?", *application)
	}

	if groupID != nil && len(strings.TrimSpace(*groupID)) > 0 {
		// Replace * with the SQL equivalient
		pattern := strings.Replace(*groupID, "*", "%", -1)
		// Replace ? with the SQL equivalient
		pattern = strings.Replace(pattern, "?", "_", -1)
		q = q.Where("group_id LIKE ?", pattern)
	}

	// Return the newest machines first
	q = q.Order("created_at desc, id", true)

	pagination, err := ign.PaginateQuery(q, &machines, *p)
	if err != nil {
		return nil, nil, ign.NewErrorMessageWithBase(ign.ErrorInvalidPaginationRequest, err)
	}
	if !pagination.PageFound {
		return nil, nil, ign.NewErrorMessage(ign.ErrorPaginationPageNotFound)
	}

	return &machines, pagination, nil
}

// buildUserDataString returns the UserData string to be used when creating a new EC2
// instance.
// @param extraLabels is an array of labels to set on the node. Each label has the form label=value.
// @return the userData in base64
func (s *Ec2Client) buildUserDataString(groupID string, extraLabels ...string) (base64Data, userData string) {

	const constLaunchEc2UserData = `#!/bin/bash
	set -x
	exec > >(tee /var/log/user-data.log|logger -t user-data ) 2>&1
	echo BEGIN
	date '+%Y-%m-%d %H:%M:%S'
	`

	// Include a label containing the Group ID
	extraLabels = append(extraLabels, getNodeLabelForGroupID(groupID))

	// NOTE: this nodeLabels trick helps setting labels to the new Node at creation time.
	nodeLabels := `cat > /etc/systemd/system/kubelet.service.d/20-labels-taints.conf <<EOF
[Service]
Environment="KUBELET_EXTRA_ARGS=--node-labels=` + strings.Join(extraLabels, ",") + `"
EOF
`

	userData = constLaunchEc2UserData + nodeLabels + s.buildClusterJoinCommand()
	base64Data = base64.StdEncoding.EncodeToString([]byte(userData))
	return
}

// buildClusterJoinCommand prepares the join command used to add an EC2 instance to the associated cluster.
func (s *Ec2Client) buildClusterJoinCommand() string {
	// This command runs the EKS cluster's join script. It requires that AWS env vars are configured.
	// The script can be found here: https://github.com/awslabs/amazon-eks-ami/blob/master/files/bootstrap.sh
	command := `
		set -o xtrace
		/etc/eks/bootstrap.sh %s %s
	`
	arguments := []string{
		// Allow the node to contain unlimited pods
		"--use-max-pods false",
	}

	return fmt.Sprintf(command, s.ec2Cfg.ClusterName, strings.Join(arguments, " "))
}

// buildClusterTag returns an EC2 tag required by clusters to mark worker nodes.
func (s *Ec2Client) buildClusterTag() *ec2.Tag {
	// Prepare the key
	key := fmt.Sprintf("kubernetes.io/cluster/%s", s.ec2Cfg.ClusterName)

	return &ec2.Tag{
		Key:   &key,
		Value: sptr("owned"),
	}
}

// setupInstanceSpecifics finds the platform handler and ask it to describe the needed instance details.
// To do this, it will pass a point to the 'RunInstancesInput' as argument, expecting the
// specific platform handler to update it with the specific details.
func (s *Ec2Client) setupInstanceSpecifics(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment,
	input *ec2.RunInstancesInput) ([]*ec2.RunInstancesInput, error) {

	// Invoke the specific handler (eg. SubT) to describe the needed instances.
	return s.platforms[*dep.Platform].setupEC2InstanceSpecifics(ctx, s, tx, dep, input)
}

// checkNodeAvailability checks that there are enough available EC2 instances
// to launch a specific group of instances. `inputs` should contain a list of
// all the instances required for a single simulation.
func (s *Ec2Client) checkNodeAvailability(ctx context.Context, simDep *SimulationDeployment,
	inputs []*ec2.RunInstancesInput) (bool, *ign.ErrMsg) {
	// Sanity check
	if len(inputs) == 0 {
		logger(ctx).Warning("checkNodeAvailability - Attempted to check availability for 0 instances.\n")
		return false, ign.NewErrorMessage(ign.ErrorUnexpected)
	}

	// Count how many ec2 instances are running at the given time and
	// compare it to the limit set by cloudsim
	// and the amount of new instances it's trying to launch.
	requestedInstances := len(inputs)

	if s.ec2Cfg.AvailableEC2Machines >= 0 {
		// Having 0 machines available stops launching new machines.
		if s.ec2Cfg.AvailableEC2Machines == 0 {
			return false, ign.NewErrorMessage(ign.ErrorLaunchingCloudInstanceNotEnoughResources)
		}

		reservedInstances := s.countReservedEC2Machines(ctx)

		logger(ctx).Debug(fmt.Sprintf("checkNodeAvailability - [%d] Requested instances | [%d] Reserved instances | [%d] Available instances.", requestedInstances, reservedInstances, s.ec2Cfg.AvailableEC2Machines))

		// Check if the number of required machines is greater than the current available amount of machines.
		if requestedInstances > s.ec2Cfg.AvailableEC2Machines-reservedInstances {
			return false, ign.NewErrorMessage(ign.ErrorLaunchingCloudInstanceNotEnoughResources)
		}
	} else if s.ec2Cfg.OnDemandCapacityReservationEnabled {
		for _, in := range inputs {
			availableInstances := s.getOnDemandCapacityReservation(ctx, *in.InstanceType, *in.Placement.AvailabilityZone)

			if availableInstances < int64(requestedInstances) {
				return false, ign.NewErrorMessage(ign.ErrorLaunchingCloudInstanceNotEnoughResources)
			}
		}
	}

	// Get and prepare the template and total number of EC2 instances for each
	// instance type. The first template found of each instance type is used to
	// check for resource availability.
	instanceTypes := map[string]int{}
	instanceTemplates := map[string]*ec2.RunInstancesInput{}
	for _, input := range inputs {
		instanceTypes[*input.InstanceType]++

		// Get the template for this instance type if it hasn't been defined
		if _, ok := instanceTemplates[*input.InstanceType]; !ok {
			template, err := cloneRunInstancesInput(input)
			if err != nil {
				return false, ign.NewErrorMessage(ign.ErrorUnexpected)
			}
			template.SetDryRun(true)
			instanceTemplates[*input.InstanceType] = template
		}
	}

	// Check the availability of each instance type
	for instanceType, count := range instanceTypes {
		template := instanceTemplates[instanceType]
		template.SetMinCount(int64(count))
		template.SetMaxCount(int64(count))
		// Check availability
		_, err := s.ec2Svc.RunInstances(template)
		// DryRun AWS requests always return an error. This error indicates
		// whether or not the request was successful.
		awsErr := err.(awserr.Error)
		if awsErr.Code() != AWSErrCodeDryRunOperation {
			logger(ctx).Info(fmt.Sprintf(
				"checkNodeAvailability - Not enough %s instances available (%d requested) for simulation [%s]: %s\n",
				*template.InstanceType,
				*template.MinCount,
				*simDep.GroupID,
				awsErr.Message(),
			))
			return false, ign.NewErrorMessageWithBase(ign.ErrorLaunchingCloudInstanceNotEnoughResources, awsErr.OrigErr())
		}
	}

	return true, nil
}

// runInstanceCall requests a single new EC2 instance to AWS.
func (s *Ec2Client) runInstanceCall(ctx context.Context, input *ec2.RunInstancesInput) (runResult *ec2.Reservation,
	err error) {
	for try := 1; try <= MaxAWSRetries; try++ {
		// First do a DryRun to check permissions etc.
		input.SetDryRun(true)
		_, err = s.ec2Svc.RunInstances(input)
		// DryRun AWS requests always return an error
		// This error indicates whether or not the request was successful
		awsErr, ok := err.(awserr.Error)
		if ok && AWSErrorIsRetryable(awsErr) {
			// If the error is non-fatal retry the launch
			logger(ctx).Info(fmt.Sprintf("launchNodes - %s retryable error: %s\n", awsErr.Code(), awsErr.Message()))
			if try != MaxAWSRetries {
				Sleep(time.Second * time.Duration(try))
			}
			continue
		} else if ok && awsErr.Code() == AWSErrCodeDryRunOperation {
			// If the error code is `DryRunOperation` it means we have the necessary
			// permissions to do the operation
			input.SetDryRun(false)
			// Perform the actual request to create the instances
			runResult, err = s.ec2Svc.RunInstances(input)
			// Do not log this message as a warning if the error is due to insufficient AWS capacity
			if awsErr, ok = err.(awserr.Error); ok && awsErr.Code() == AWSErrCodeInsufficientInstanceCapacity {
				logger(ctx).Info(fmt.Sprintf("launchNodes - error launching: %s\n", err))
			} else if err != nil {
				logger(ctx).Warning(fmt.Sprintf("launchNodes - error launching: %s\n", err))
			}
			return
		} else {
			// Return the RunInstance error otherwise
			return
		}
	}

	return
}

// setSourceDestCheck sets the Source/Dest. check of an EC2 instance to a specific value.
// The source/dest. check ensures that in instance is the source or destination of any traffic it sends or receives.
// There are some cases where this check prevents things from working. Examples include a group of clients behind a NAT
// or encapsulation of network traffic in an overlay network/
func (s *Ec2Client) setSourceDestCheck(instanceID *string, value *bool) error {
	sourceDestCheckInput := &ec2.ModifyInstanceAttributeInput{
		InstanceId:      instanceID,
		SourceDestCheck: &ec2.AttributeBooleanValue{Value: value},
	}
	_, err := s.ec2Svc.ModifyInstanceAttribute(sourceDestCheckInput)

	return err
}

// launchInstances starts a group of EC2 instances. The ids and machine info
// of new instances are returned to be accessed later during the node launch
// process.
func (s *Ec2Client) launchInstances(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment,
	instanceInputs []*ec2.RunInstancesInput) (instanceIds []string, machines []*MachineInstance, err error) {
	instanceIds = make([]string, 0)
	machines = make([]*MachineInstance, 0)

	for _, input := range instanceInputs {
		var runResult *ec2.Reservation
		runResult, err = s.runInstanceCall(ctx, input)
		// If the instance could not be started, stop launching new
		// instances and terminate launched instances
		if err != nil {
			return
		}

		// If we are here the AWS call was OK and we have the launched instances
		// info in 'runResult'.
		for _, ins := range runResult.Instances {
			// Get the created Instance ID(s).
			iID := *ins.InstanceId
			instanceIds = append(instanceIds, iID)

			// Disable Source/Dest. checks to allow cross subnet traffic.
			if err = s.setSourceDestCheck(&iID, aws.Bool(false)); err != nil {
				return
			}

			// And create a DB record to track the machine instance in case of errors later
			machine := MachineInstance{
				InstanceID:      &iID,
				LastKnownStatus: macInitializing.ToStringPtr(),
				GroupID:         dep.GroupID,
				Application:     dep.Application,
			}
			err = tx.Create(&machine).Error
			if err != nil {
				return
			}
			machines = append(machines, &machine)
		}
	}

	return
}

// terminateInstances terminates a group of EC2 instances.
func (s *Ec2Client) terminateInstances(ctx context.Context, machines []*MachineInstance) {
	terminateIds := make([]*string, 0)
	for _, machine := range machines {
		terminateIds = append(terminateIds, machine.InstanceID)
	}
	_, err := s.ec2Svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: terminateIds,
	})
	if err != nil {
		errorMsg := err.Error()
		if awsErr, ok := err.(awserr.Error); ok && awsErr != nil {
			errorMsg = fmt.Sprintf("%s - %s", awsErr.Code(), awsErr.Message())
		}
		logger(ctx).Error(fmt.Sprintf(
			"launchNodes - error while attempting to shutdown EC2 instances: %s\n",
			errorMsg,
		))
	}
}

// launchNodes will try to launch new EC2 instances, and register them with the k8 master.
// @public
func (s *Ec2Client) launchNodes(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) (*string, *ign.ErrMsg) {
	// TODO need to design a rollback mechanism to shutdown created machines/nodes/etc
	// if something fails during the creation (even the db)

	tstart := time.Now()

	groupID := *dep.GroupID

	// This will be the return value if everything is ok
	nodeSelectorGroupID := getNodeLabelForGroupID(groupID)
	instanceName := s.getInstanceNameFor(groupID, subtTypeGazebo)
	userData, _ := s.buildUserDataString(groupID, labelAndValue(nodeLabelKeyCloudsimNodeType, subtTypeGazebo))

	ignlog := logger(ctx)

	s.availabilityZoneIndex = (s.availabilityZoneIndex + 1) % len(s.ec2Cfg.AvailabilityZones)

	// Set the initial details of the instances that will be created.
	// Set the default values for the instances.
	// Note: specific Applications/Platforms can override these.
	instanceTemplate := &ec2.RunInstancesInput{
		DryRun: aws.Bool(true),
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			// This IAM role is assigned to the EC2 instance so it can join the EKS cluster,
			// write logs to AWS CloudWatch and access ECR.
			// It can be configured using an env var: AWS_IAM_INSTANCE_PROFILE_ARN.
			// The default value is: "arn:aws:iam::200670743174:instance-profile/aws-eks-role-cloudsim-worker"
			Arn: aws.String(s.awsCfg.IamInstanceProfile),
		},
		// IMPORTANT: the 'KeyName' is the name of the ssh key to use to remotely access this instance.
		KeyName:          aws.String("ignitionFuel"),
		MaxCount:         aws.Int64(1),
		MinCount:         aws.Int64(1),
		SecurityGroupIds: aws.StringSlice([]string{"sg-0c5c791266694a3ca"}),
		SubnetId:         aws.String(s.ec2Cfg.Subnets[s.availabilityZoneIndex]),
		Placement: &ec2.Placement{
			AvailabilityZone: aws.String(s.ec2Cfg.AvailabilityZones[s.availabilityZoneIndex]),
		},
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{Key: aws.String("Name"), Value: aws.String(instanceName)},
					{Key: aws.String("CloudsimGroupID"), Value: aws.String(groupID)},
					{Key: aws.String("project"), Value: aws.String("cloudsim")},
					{Key: dep.Platform, Value: aws.String("True")},
					s.buildClusterTag(),
				},
			},
		},
		// UserData is an initialization script run after the instance is launched.
		// We use it to make the Node 'join' the k8 cluster.
		UserData: aws.String(userData),
	}

	// Delegate the creation of specific "instanceInputs" to the chosen "platform"
	// and to add specific details if needed.
	var instanceInputs []*ec2.RunInstancesInput
	var err error
	if instanceInputs, err = s.setupInstanceSpecifics(ctx, tx, dep, instanceTemplate); err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorLaunchingCloudInstance, err)
	}

	// Now we try (and retry) to launch the EC2 instances needed for the simulation
	var instanceIds []string
	var machines []*MachineInstance
	for try := 1; try <= MaxAWSRetries; try++ {
		// Lock to ensure other worker threads don't claim EC2 instances after
		// it's been determined there's enough resources available to launch all
		// instances for this simulation
		s.lockRunInstances.Lock()

		// Check there's enough instances available. If not, wait some time and retry
		if ok, em := s.checkNodeAvailability(ctx, dep, instanceInputs); ok {
			// There are enough instances, launch instances and stop retrying
			instanceIds, machines, err = s.launchInstances(ctx, tx, dep, instanceInputs)
			break
		} else {
			// Used for error handling if the retry limit is exceeded
			err = em.BaseError
			// Do not retry if the error is fatal or this is the last try
			if em.ErrCode != ign.ErrorLaunchingCloudInstanceNotEnoughResources || try >= MaxAWSRetries {
				break
			}
			// Let other workers attempt to get instances
			s.lockRunInstances.Unlock()
			// Wait before retrying
			Sleep(time.Minute)

			ignlog.Debug(fmt.Sprintf(
				"launchNodes - not enough instances to start simulation for groupid [%s]. retrying: %s",
				*dep.GroupID,
				em.Msg,
			))

			// Retry claiming instances
			continue
		}
	}
	s.lockRunInstances.Unlock()

	// TODO It makes more sense for this to be handled inside the errorHandler in the future.
	// Check that everything was setup properly. If not, terminate launched instances.
	// There were some cases where the previous block succeeded but AWS was unable
	// to grant instances. A sanity check for this is made in order to handle this
	// situation.

	// Count how many machines were requested
	var requestedMachines int
	for _, input := range instanceInputs {
		requestedMachines += int(*input.MinCount)
	}

	// Check if there are no machines available or the number of instances created does not match the amount
	// of requested machines.
	invalidInstanceCount := len(instanceIds) != requestedMachines
	if err != nil || invalidInstanceCount {
		timeTrack(ctx, tstart, "launchNodes - launchInstances ended with error")

		// Terminate launched EC2 instances
		if len(instanceIds) > 0 {
			s.terminateInstances(ctx, machines)
		}

		// Set error type
		var errType int64
		awsErr, ok := err.(awserr.Error)
		retry := (ok && AWSErrorIsRetryable(awsErr)) ||
			(err != nil && err.Error() == ign.NewErrorMessage(ign.ErrorLaunchingCloudInstanceNotEnoughResources).Msg) ||
			invalidInstanceCount
		if retry {
			errType = ign.ErrorLaunchingCloudInstanceNotEnoughResources
		} else {
			errType = ign.ErrorLaunchingCloudInstance
		}

		return nil, ign.NewErrorMessageWithBase(errType, err)
	}

	timeTrack(ctx, tstart, "launchNodes - launchInstances")
	ignlog.Info(fmt.Sprintf("Launch instance requests succeeded. Instance Ids: %v", instanceIds))

	// Use a waiter function to Block until the instances are initialized
	describeInstanceStatusInput := &ec2.DescribeInstanceStatusInput{
		InstanceIds: aws.StringSlice(instanceIds),
	}
	ignlog.Info(fmt.Sprintf("About to WaitUntilInstanceStatusOk. Instance Ids: %v", instanceIds))
	if err := s.ec2Svc.WaitUntilInstanceStatusOk(describeInstanceStatusInput); err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorLaunchingCloudInstance, err)
	}
	timeTrack(ctx, tstart, "launchNodes - WaitUntilInstanceStatusOk")

	// Mark the Machines in the DB as 'Running'.
	for _, m := range machines {
		if em := m.updateMachineStatus(tx, macRunning); em != nil {
			return nil, em
		}
	}

	ignlog.Info(fmt.Sprintf("Instances are now running: %v. Cloudsim GroupID: %s\n", instanceIds, groupID))
	return &nodeSelectorGroupID, nil
}

func (s *Ec2Client) getInstanceNameFor(groupID, suffix string) string {
	return fmt.Sprintf("%s-node-group-%s-%s", s.awsCfg.NamePrefix, groupID, suffix)
}

// appendTags adds new tags to a RunInstancesInput.
func appendTags(input *ec2.RunInstancesInput, tags ...*ec2.Tag) {
	// Tag with SubT (hack: we are assuming the base Tags structure)
	input.TagSpecifications[0].Tags = append(input.TagSpecifications[0].Tags, tags...)
}

// replaceTag replaces the specified tag values. If a tag is not found no changes are performed.
func replaceTag(input *ec2.RunInstancesInput, tags ...*ec2.Tag) {
	for _, tag := range input.TagSpecifications[0].Tags {
		for _, newTag := range tags {
			if *tag.Key == *newTag.Key {
				*tag.Value = *newTag.Value
				break
			}
		}
	}
}

func replaceInstanceNameTag(input *ec2.RunInstancesInput, name string) {
	// hack: we are assuming the internals of the Tags structure
	nameTag := &ec2.Tag{Key: aws.String("Name"), Value: aws.String(name)}
	input.TagSpecifications[0].Tags[0] = nameTag
}

func getNodeLabelForGroupID(groupID string) string {
	return labelAndValue(nodeLabelKeyGroupID, groupID)
}

func labelAndValue(key, value string) string {
	return key + "=" + value
}

// deleteK8Nodes deletes the kubernetes nodes used to run a GroupID.
// It is expected that if the labeled Nodes cannot be found, then this function should return
// an ErrorLabeledNodeNotFound.
// @public
func (s *Ec2Client) deleteK8Nodes(ctx context.Context, tx *gorm.DB, groupID string) (interface{}, *ign.ErrMsg) {

	// Find and Delete all k8 Nodes associated to the GroupID.
	nodeLabel := getNodeLabelForGroupID(groupID)
	nodesInterface := s.clientset.CoreV1().Nodes()
	nodes, err := nodesInterface.List(metav1.ListOptions{LabelSelector: nodeLabel})
	if err != nil {
		return nil, NewErrorMessageWithBase(ErrorLabeledNodeNotFound, err)
	}
	for _, n := range nodes.Items {
		// Delete the node
		err = nodesInterface.Delete(n.Name, &metav1.DeleteOptions{})
		if err != nil {
			// There was an error deleting the Node. We log the error and continue,
			// as we want to Stop/Terminate the ec2 instance.
			em := ign.NewErrorMessageWithBase(ign.ErrorK8Delete, err)
			logger(ctx).Error("Error while invoking k8 Delete Node. Make sure a sysadmin deletes the Node manually", em)
		}
	}
	return nodes, nil
}

// deleteHosts is a helper function that sends a request AWS to terminate
// all the EC2 instances associated to a given GroupID.
// It also updates the MachineInstance DB records with the status of the terminated instances.
// @public
func (s *Ec2Client) deleteHosts(ctx context.Context, tx *gorm.DB,
	dep *SimulationDeployment) (interface{}, *ign.ErrMsg) {

	// Get the EC2 instance Ids for the given groupID.
	var machines MachineInstances
	if err := tx.Model(&MachineInstance{}).Where("group_id = ?", *dep.GroupID).Find(&machines).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorSimGroupNotFound, err)
	}

	var instIds []string
	// Update each machine instance DB record, to mark the termination process as 'started'.
	for _, m := range machines {
		instIds = append(instIds, *m.InstanceID)
		if em := m.updateMachineStatus(tx, macTerminating); em != nil {
			return nil, em
		}
	}

	instanceIds := aws.StringSlice(instIds)
	var result interface{}
	var err error

	// Terminate or Stop EC2 instances ?
	// Note: if the Simulation was marked with error then we Stop the instance, to
	// allow further investigation.
	stopOnEnd := dep.StopOnEnd != nil && *dep.StopOnEnd
	if s.awsCfg.ShouldTerminateInstances && !stopOnEnd {
		input := &ec2.TerminateInstancesInput{
			DryRun:      aws.Bool(true),
			InstanceIds: instanceIds,
		}
		result, err = s.ec2Svc.TerminateInstances(input)
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == AWSErrCodeDryRunOperation {
			// If the error code is `DryRunOperation` it means we have the necessary
			// permissions to do the operation
			input.DryRun = aws.Bool(false)
			result, err = s.ec2Svc.TerminateInstances(input)
		}
	} else {
		input := &ec2.StopInstancesInput{
			DryRun:      aws.Bool(true),
			InstanceIds: instanceIds,
		}
		result, err = s.ec2Svc.StopInstances(input)
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == AWSErrCodeDryRunOperation {
			input.DryRun = aws.Bool(false)
			result, err = s.ec2Svc.StopInstances(input)
		}
	}
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorStoppingCloudInstance, err)
	}

	// TODO there should be a watcher that marks machine DB records as "terminated"
	// once the EC2 termination succeeds. At the time being, we just mark it as terminated.
	for _, m := range machines {
		if em := m.updateMachineStatus(tx, macTerminated); em != nil {
			return nil, em
		}
	}

	logger(ctx).Info(fmt.Sprintf("Instances terminated: %s. Result: %v", instIds, result))
	return machines, nil
}

// countReservedEC2Machines returns the number of reserved EC2 machines that match the cloudsim-simulation-worker tag.
func (s *Ec2Client) countReservedEC2Machines(ctx context.Context) (instances int) {
	describeInstancesInput := &ec2.DescribeInstancesInput{
		MaxResults: aws.Int64(1000),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:cloudsim-simulation-worker"),
				Values: []*string{
					aws.String(s.awsCfg.NamePrefix),
				},
			},
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("pending"),
					aws.String("running"),
				},
			},
		},
	}

	describeInstancesOutput, err := s.ec2Svc.DescribeInstances(describeInstancesInput)
	if err != nil {
		logger(ctx).Warning("countReservedEC2Machines - There was an error getting the list of available machines")
		instances = 0
	} else {
		instances = len(describeInstancesOutput.Reservations)
	}
	return
}

// getZoneIDFromName receives a zone name like `us-east-1a` and returns the zone ID `use1-az1`.
func (s *Ec2Client) getZoneIDFromName(name string) (string, error) {
	out, err := s.ec2Svc.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{
		ZoneNames: aws.StringSlice([]string{name}),
	})

	if err != nil {
		return "", err
	}

	return *out.AvailabilityZones[0].ZoneId, nil
}

// getOnDemandCapacityReservation gets the amount of machines available to launch using On-Demand Capacity Reservation.
func (s *Ec2Client) getOnDemandCapacityReservation(ctx context.Context, instanceType string, zone string) int64 {
	zoneId, err := s.getZoneIDFromName(zone)
	if err != nil {
		return 0
	}

	input := &ec2.DescribeCapacityReservationsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-type"),
				Values: []*string{
					aws.String(instanceType),
				},
			},
			{
				Name: aws.String("availability-zone-id"),
				Values: []*string{
					aws.String(zoneId),
				},
			},
		},
	}

	output, err := s.ec2Svc.DescribeCapacityReservations(input)
	if err != nil {
		return 0
	}

	var reservation int64
	for _, r := range output.CapacityReservations {
		reservation += *r.AvailableInstanceCount
	}
	return reservation
}
