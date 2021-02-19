<div align="center">
  <img src="./assets/logo.png" width="200" alt="Ignition Robotics" />
  <h1>Ignition Robotics</h1>
  <p>Setting up a new cloudsim environment</p>
</div>

# Amazon Web Services (AWS)

## Elastic Kubernetes Service (EKS)
We use Elastic Kubernetes Service to deploy our Kubernetes clusters. We will describe how to create a new EKS cluster in the following steps.

### Cluster
#### Configuration
**Name**: `web-cloudsim-[environment]`

**Kubernetes version**: `1.15`

**Cluster Service Role**: `aws-eks-role-cloudsim`

**Tags**:

| Key | Value |
| ------ | ------ |
| project | cloudsim |
| SubT | true |
| application | subt |
| enviroment | [environment] |
| platform | cloudsim |
| Name | cloudsim-eks-cluster-[environment] |

#### Networking

**VPC**: `vpc-12af6375 - 172.30.0.0/16`

**Subnets**:

| Name | Subnet ID |
| ------ | ------ |
| subnet-cloudsim-az-1a | subnet-0e632d68a9032ab9d |
| subnet-cloudsim-az-1b | subnet-03774a4f37672e4da |
| subnet-cloudsim-az-1c | subnet-00a9f3acf0ce3785a |
| subnet-cloudsim-az-1d | subnet-0614ac8a450d5d1d1 |
| subnet-cloudsim-az-1f | subnet-048c68c81c80a6636 |

**Additional security groups**:

| Name | ID |
| ------ | ------ |
| rds-launch-wizard | sg-9d31e8e6 |
| kubernetes | sg-0c5c791266694a3ca |
| cloudsim-server | sg-023c19380b48dcabb |


**Cluster endpoint access** set to `Public and private`

#### Logging
- **API server**: `Enabled`.
- **Audit**: `Enabled`.
- **Authenticator**: `Enabled`.
- **Controller manager**: `Enabled`.
- **Scheduler**: `Enabled`.


<hr />


### Workers - Elastic Compute Cloud (EC2)
After creating the cluster, we need to add a node group in order to have a place where the pods can live. Under the `Compute` tab in the EKS control panel, click the `Add Node Group` button.

#### Group configuration
**Name**: `web-cloudsim-[environment]-nodes`

**Node IAM Role**: `aws-eks-role-cloudsim-worker`

**Subnets**:

| Name | Subnet ID |
| ------ | ------ |
| subnet-cloudsim-az-1a | subnet-0e632d68a9032ab9d |
| subnet-cloudsim-az-1b | subnet-03774a4f37672e4da |
| subnet-cloudsim-az-1c | subnet-00a9f3acf0ce3785a |
| subnet-cloudsim-az-1d | subnet-0614ac8a450d5d1d1 |
| subnet-cloudsim-az-1f | subnet-048c68c81c80a6636 |

**Allow remote access to nodes** set to `Enabled`.

**SSH key pair**: `ignitionFuel`

**Allow remote access from**: `All`

#### Node compute configuration
**AMI type**: Amazon Linux 2 (AL2_x86_64)

**Instance type**: `t3.small`

**Disk size**: `20GiB`

#### Group size
**Minimum size**: `2` nodes

**Maximum size**: `2` nodes

**Desired size**: `2` nodes

#### Labels

| Key | Value | Description |
| ------ | ------ | ------ |
| gitlab | true | Allow the gitlab runner to run job pods on this node group |
| server | true | Allow the cloudsim deployment to launch pods on this node group |

### Calico

**TODO**
This section will be completed in another MR.

## Simple Storage Service (S3)
S3 buckets containing Gz log files should have the following `CORS configuration`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<CORSConfiguration xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
<CORSRule>
    <AllowedOrigin>*</AllowedOrigin>
    <AllowedMethod>GET</AllowedMethod>
    <MaxAgeSeconds>300</MaxAgeSeconds>
    <AllowedHeader>Authorization</AllowedHeader>
</CORSRule>
</CORSConfiguration>
```

This CORS configuration is needed to download the logs from the Portal web application.

# Subterranean Challenge

## Circuit rules
SubT circuits need the DB table `sub_t_circuit_rules` with data to fetch the rules
to launch SubT circuits. Make sure your server has that table loaded with rules.

# Future Work
- (Readiness Probe) Check if Gz is running: `ign topic -e -t /world/default/stats -n 1 | grep -c "sim_time"`. If the above return 1, then simulation is "Up".
- Nvidia: "Node ready" -> check for gpu count before launching the pods, to make sure the device pluing worked on the new node(s).
- install monitoring (read nvidia instructions)
- install logging system
- install k8dashboard in master and write instructions
- launching 2 pods using affinity (to different nodes)
- Add a systemadmin route to run a query and return "inconsistencies" between DB and AWS/K8. So the admin can then update DB records by hand.
- (future) investigate about EC2 Spot instances (maybe faster to launch, cheaper)
- (future) Investigate using autoscaling in AWS


## Backlog
- Consider making use of kubernetes 'namespaces' to separate applications types.
- Will need to use NodePort service to access a pod port from the outside world.
- Interesting read about Exposing pods: http://alesnosek.com/blog/2017/02/14/accessing-kubernetes-pods-from-outside-of-the-cluster/
- Networking
  - https://www.aquasec.com/wiki/display/containers/Kubernetes+Networking+101
  - https://github.com/ahmetb/kubernetes-network-policy-recipes
- Security in containers
  - https://kubernetes.io/blog/2018/07/18/11-ways-not-to-get-hacked/#8-run-containers-as-a-non-root-user
  - linux capabilities.
- Casbin related (authorization):
  - About Casbin and multithreading: https://casbin.org/docs/en/multi-threading
  - Casbin tutorial: https://zupzup.org/casbin-http-role-auth/
  - Casbin as a Service: https://github.com/casbin/casbin-server
  
# Tips
## AWS
- Running ssh commands remotely: https://docs.aws.amazon.com/systems-manager/latest/userguide/walkthrough-cli.html#walkthrough-cli-example-1
  - go-sdk version: https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/#SSM.SendCommand
- Sending user commands at Instance Launch time: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html
  - these scripts will be run as root user. So no need for using sudo in them.
  - console output in the instance will be at `/var/log/cloud-init-output.log`
- How to filter (get) ec2 instances from Go code: https://github.com/aws/aws-sdk-go/blob/master/example/service/ec2/filterInstances/filter_ec2_by_tag.go

<hr />

## We are currently using these env variables
- env `SUBT_GZSERVER_LOGS_VOLUME_MOUNT_PATH` with default value `/tmp/ign`
- env `SIMSVC_NODE_READY_TIMEOUT_SECONDS` with default value `300`
- env `SIMSVC_POD_READY_TIMEOUT_SECONDS` with default value `300`


## Troubleshooting
<!--
### TODO CHECK IF THIS IS STILL TRUE
- We've noticed that after a Cloudsim server restart, if the server had running simulations,
then the simulations will be regenerated but the `ign-transport topics and connections`
(and thus `/stats` messages) will be lost. We suggest Shutting down and restart those simulations too.
-->

# We are using the following AMIs:
- Kubernetes GPU (Worker) Nodes: `ami-08861f7e7b409ed0c`, name `cloudsim-worker-node-eks-gpu-optimized-1.0.0`. 
Used with `g3.4xlarge` instances. Note: in SubT these AMI and g3 instance are used for both gzserver and field-computer nodes.
- VPC ID: `vpc-12af6375`
- Subnets:
 
| Name | Subnet ID |
| ------ | ------ |
| subnet-cloudsim-az-1a | subnet-0e632d68a9032ab9d |
| subnet-cloudsim-az-1b | subnet-03774a4f37672e4da |
| subnet-cloudsim-az-1c | subnet-00a9f3acf0ce3785a |
| subnet-cloudsim-az-1d | subnet-0614ac8a450d5d1d1 |
| subnet-cloudsim-az-1f | subnet-048c68c81c80a6636 |

- Security Groups: 

| Name | ID |
| ------ | ------ |
| kubernetes | sg-0c5c791266694a3ca |

- IMPORTANT: All EC2 instances (master and nodes) have the IAM Role: 
`arn:aws:iam::200670743174:instance-profile/aws-eks-role-cloudsim-worker`. This IAM Role has attached policies to
 access ECR and CloudWatch Logs.
- Using `ignitionFuel.pem` as ssh key.

# Cloudsim server: how to switch to a different Kubernetes cluster?

On startup, Cloudsim gets or updates the configuration file used by `kubectl` to access its cluster. For this
 configuration to work, an EKS cluster must be up and running and the `AWS_CLUSTER_NAME` environment variable must be
 set to the EKS cluster name. Running or restarting the Cloudsim container will fetch the configuration if it
 does not exist or update it if it already exists. 

# Manually Join a Node to the cluster:

To join instances to an EKS cluster, instances must be running an EKS optimized AMI. These images contain a special
 `bootstrap.sh` script that adds the instance to a cluster. There is a custom EKS Cloudsim AMI prepared for nodes. Its
 id is specified in the "We are using the following AMIs" section. Always use the Cloudsim AMI, unless you explicitly
 want to create a node using base EKS AMI.
 
To make the instance join the cluster, run the following inside the instance

```
/etc/eks/bootstrap.sh <cluster_name>
```

Where `<cluster_name` is the name of the EKS cluster.


# TODO Needs rewrite
# Creating everything from scratch
In case you need to re-do everything...

## Creating the AWS Security Group
Created security group for our current VPC:

```
aws ec2 create-security-group --group-name kubernetes --description "Kubernetes Security Group" --vpc-id vpc-12af6375
```

After creating the Security Group we got its ID, and use it to add rules to it.

```
# Allow SSH connections from all sources
aws ec2 authorize-security-group-ingress --group-id sg-0c5c791266694a3ca --protocol tcp --port 22 --cidr 0.0.0.0/0
# Allow HTTP traffic from all sources
aws ec2 authorize-security-group-ingress --group-id sg-0c5c791266694a3ca --protocol tcp --port 80 --cidr 0.0.0.0/0
# Allow traffic from the secure kube api port from all sources 
aws ec2 authorize-security-group-ingress --group-id sg-0c5c791266694a3ca --protocol tcp --port 6443 --cidr 0.0.0.0/0
# Allow all traffic from instances with the `kubernetes` security group
aws ec2 authorize-security-group-ingress --group-id sg-0c5c791266694a3ca --source-group sg-0c5c791266694a3ca --protocol all
# Allow all traffic from instances with the `cloudsim-server` security group
aws ec2 authorize-security-group-ingress --group-id sg-0c5c791266694a3ca --source-group sg-023c19380b48dcabb --protocol all
# Allow all outbound traffic
aws ec2 authorize-security-group-egress --group-id sg-0c5c791266694a3ca --protocol all --cidr 0.0.0.0/0
```

## Security Group for EKS cluster instances
EKS cluster instances that will run Cloudsim servers will need to allow traffic from simulation instances.

```
aws ec2 create-security-group --group-name cloudsim-server --description "Cloudsim Server Security Group" --vpc-id vpc-12af6375
```

After creating the security group, use the ID from the previous command to open the following ports:

```
# Allow SSH connections from all sources
aws ec2 authorize-security-group-ingress --group-id sg-023c19380b48dcabb --protocol tcp --port 22 --cidr 0.0.0.0/0
# Allow all traffic from instances with the `kubernetes` security group
aws ec2 authorize-security-group-ingress --group-id sg-023c19380b48dcabb --source-group sg-0c5c791266694a3ca --protocol all
# Allow all outbound traffic
aws ec2 authorize-security-group-egress --group-id sg-023c19380b48dcabb --protocol all --cidr 0.0.0.0/0
```


## Instructions to create the Worker Nodes AMI (with NVIDIA)
(Note: no need to follow these steps if you are using the AMIs that we've already created)

- Note: This was done to create the cloudsim K8 Node AMI (ie. The Worker Nodes with GPU).
- Using instance type `g3.4xlarge` (with 1 GPU Tesla M60).
  - Used the Kubernetes Security Group (sg-0c5c791266694a3ca)
  - Used 128 GB disk (General Purpose SSD -- gp2)
- Using AWS EKS Amazon Linux 2 GPU-optimized AMI: [`ami-08dc081250e6c9d58`](https://console.aws.amazon.com/systems-manager/parameters/%252Faws%252Fservice%252Feks%252Foptimized-ami%252F1.14%252Famazon-linux-2-gpu%252Frecommended%252Fimage_id/description?region=us-east-1#). 
  It comes with Kubernetes 1.14, CUDA 10, nvidia-docker 2, and Nvidia as docker runtime by default.
- Added an X server configuration specific to G3 instance GPUs

```
cat > /etc/X11/xorg.conf.d/nvidia.conf <<-EOF
Section "Device"
    Identifier     "NVIDIA"
    Driver         "nvidia"
    BusID          "0:30:0"
EndSection
EOF
``` 

- Created a simple X server systemd service.

```
cat > /etc/systemd/system/xorg.service <<-EOF
[Unit]
Description=Xorg server

[Service]
Type=simple
SuccessExitStatus=0 1

ExecStart=/usr/bin/Xorg
ExecStop=/usr/bin/kill `pidof Xorg`

[Install]
WantedBy=default.target
EOF
``` 

- Enabled the systemd service.

```
sudo systemctl enable xorg
```

- Created the new AMI (latest, `cloudsim-worker-node-eks-gpu-optimized-1.0.0`). This AMI is expected to be used with `G3` EC2 instances.

#
# Some Tips for Kubernetes

- The `kubelet` config (node client) can be found in `/etc/kubernetes/kubelet.conf`.
- Use node affinity to send pods to specific node instance groups: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
- Shutdown idle pods: https://carlosbecker.com/posts/k8s-sandbox-costs/
- Long polling in Go: https://lucasroesler.com/2018/07/golang-long-polling-a-tale-of-server-timeouts/


## Kubernetes Security tips
- If you get an error when using kubectl remotely:  x509: certificate is valid for...
    - Fix: https://mk.imti.co/kubectl-remote-x509-valid/ and
https://stackoverflow.com/questions/46360361/invalid-x509-certificate-for-kubernetes-master?utm_medium=organic&utm_source=google_rich_qa&utm_campaign=google_rich_qa.
Then you need to re-copy the kunernetes/config (with cert data) to clients.
- K8 Security Best Practices: https://github.com/freach/kubernetes-security-best-practice/blob/master/README.md
- Use Kube-bench tool: https://github.com/aquasecurity/kube-bench
- Securing a cluster (official doc): https://kubernetes.io/docs/tasks/administer-cluster/securing-a-cluster/
- Trusting TLS in a Cluster https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/


## Kubernetes Network Policies
- recipes: https://github.com/ahmetb/kubernetes-network-policy-recipes
- Doc: https://kubernetes.io/docs/concepts/services-networking/network-policies/ 


# Local development

## Install Dind-Cluster to have a local Kubernetes with 1 master and 2 slave nodes

1. Install liblz4-tool
    ```
    sudo apt-get install liblz4-tool
    ```
1. Clone the kubeadm-dind-cluster repository
    ```
    git clone https://github.com/kubernetes-sigs/kubeadm-dind-cluster ~/kubeadm-dind-cluster
    ```
1. Build the kubeadm-dind-cluster
    ```
    cd ~/kubeadm-dind-cluster
    ./build/build-local.sh
    ```
1. After the images are succesfully built, do the following
    ```
    export DIND_IMAGE="mirantis/kubeadm-dind-cluster:local"
    sudo CNI_PLUGIN="weave" ./dind-cluster.sh up
    ```
1. Check the cluster is working by running the following, which should show
   3 nodes.
    ```
    kubectl get nodes
   ```
1. Now `Label` the nodes to have them ready to be used by Cloudsim:
  - Note, you can see current node labels by: `kubectl get nodes --show-labels`
  1. Disable master: `kubectl label nodes kube-master cloudsim_free_node=false`
  1. Enable node 1: `kubectl label nodes kube-node-1 cloudsim_free_node=true` and then `kubectl label nodes kube-node-1 cloudsim_groupid=`
  1. Enable node 2: `kubectl label nodes kube-node-2 cloudsim_free_node=true` and then `kubectl label nodes kube-node-2 cloudsim_groupid=`

## Running Gazebo from its docker image
- In your local dev machine: `sudo IGN_VERBOSE=1 docker run -it -e DISPLAY -e QT_X11_NO_MITSHM=1 -e XAUTHORITY=$XAUTH -e IGN_PARTITION -v "$XAUTH:$XAUTH" -v "/tmp/.X11-unix:/tmp/.X11-unix" -v "/etc/localtime:/etc/localtime:ro" -v "/dev/input:/dev/input" --security-opt seccomp=unconfined nkoenig/ign-gazebo:nightly ign-gazebo-server -v 4`.
- Shorter gzserver version: `sudo IGN_VERBOSE=1 docker run -it nkoenig/ign-gazebo:nightly ign-gazebo-server -v 4`.

## Some useful Gazebo commands
- To see a topic in the command line: `IGN_VERBOSE=0 ign topic -e -t /world/default/stats`
- To resume a simulation: `IGN_VERBOSE=0 ign service -s /world/default/control --reqtype ignition.msgs.WorldControl --reptype ignition.msgs.Boolean --timeout 2000 --req 'pause:false'`. Use `pause:true` if you to pause the simulation.

## Notes about local testing using ign_transport running inside Kubernetes (local machine using Docker-in-Docker):

I was able to successfully run ign-transport's examples `publisher_c` and `subscriber_c` inside a `dind-cluster` (local machine).
It seems that Weave's NPC (firewall) was blocking some packets. To make it work I had to:

1. (This is an optional step, only needed to support ign-transport locally) Install `dind-cluster` from scratch (remove previous installation first).
Install a modified version of dind-cluster. I've modified dind cluster's `image/wrapkubeadm` script to remove the `--masquerade-all` command line argument in function `dind::proxy-cidr-and-no-conntrack` , and also changed the "masqueradeAll" config setting to "false" in function `dind::frob-proxy-configmap`. Then, I rebuilt the dind docker image (`build/build-local.sh` script).
1. Then installed dind-cluster by:
```
$ export DIND_IMAGE="mirantis/kubeadm-dind-cluster:local"
$ sudo CNI_PLUGIN="weave" ./dind-cluster.sh up
```
1. And then I deleted the automatically installed Weave plugin from the cluster:
`$ kubectl delete -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')`
1. And later, I've reinstalled the Weave plugin but with an extra argument (disabling its NPC, which is its internal firewall)
`$ kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')&disable-npc=true"`
1. Once I had this cluster, I could run kubernetes deployments for `publisher_c` and `subscriber_c` by running commands:
  - `kubectl run testsub --image=transport_examples --env IGN_PARTITION=foo --env IGN_VERBOSE=1 --image-pull-policy=Never --command -- /osrf/ign-transport/example/build/subscriber_c`
  - `kubectl run testpub --image=transport_examples --env IGN_PARTITION=foo --env IGN_VERBOSE=1 --image-pull-policy=Never --command -- /osrf/ign-transport/example/build/publisher_c`

# Connecting from Go code
- The client machine needs to have a `~/.kube/config` file. The cloudsim server will look for that file to connect.
- This is obtained automatically when launching the Cloudsim container. It require the `AWS_CLUSTER_NAME` environment
  variable to be set to the EKS cluster name.
- Note: security configuration and authentication should be managed by the configuration file. The Go code should just 'use that config'.
- There should be an open port to connect to. By default we use the secure port 6443.
  - During development, to avoid Certificates Error, you can add `insecure-skip-tls-verify: true` to the config file in the client machine.
- Tips to configure authentication in the cluster: https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/


## Using Token based authentication
The `config` file will use certificate based authentication by default. You can switch to use Token (which won't expire)
by adding a bearer token in the `users > user` section of the config file. Use the `tiller` service account token.

Example of a `.kube/config` file using a token:
```
...

users:
- name: kubernetes-admin
  user:
    token: theBearerToken

...
```


You can see all secrets by doing `kubectl -n kube-system get secret`

Then to see a specific secret and get its token: `kubectl -n kube-system describe secret tiller-token-xxxxxx`


More info:
- https://docs.vmware.com/en/VMware-Cloud-PKS/services/com.vmware.cloudpks.using.doc/GUID-2C2ECEB3-7BFC-44F0-9870-591F52EBA4A9.html
- https://stackoverflow.com/questions/46664104/how-to-sign-in-kubernetes-dashboard
- https://github.com/kubernetes/dashboard/wiki/Access-control#bearer-token

# Notes about copying Gazebo log files to S3 buckets
- We are using an additional copy pod to upload simulation logs to S3. Copy pods share data with their target pods by
 using a Kubernetes shared `volume`.
- The GZ logs are copied into the S3 bucket identified by this `AWS_GZ_LOGS_BUCKET` env var, and under a folder named after the Team's name.
Eg. `s3://web-cloudsim-logs/my-cool-team/xxxxx.log`.

Tip: You can use this docker image `preyna/tests:bigfile` as the `SUBT_GZSERVER_IMAGE` simulation image.
This docker image will start by creating a log file of 1 GB at `/tmp/ign/log/` and later will keep running forever in
a for loop, appending 1 extra MB each 2 seconds to the log file.
The file was created as an attempt to test big gz log files in simulations that are left there until they expire.

# AWS Tips
- Running ssh commands remotely: https://docs.aws.amazon.com/systems-manager/latest/userguide/walkthrough-cli.html#walkthrough-cli-example-1
  - go-sdk version: https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/#SSM.SendCommand
- Sending user commands at Instance Launch time: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html
  - these scripts will be run as root user. So no need for using sudo in them.
  - console output in the instance will be at `/var/log/cloud-init-output.log`
- How to filter (get) ec2 instances from Go code: https://github.com/aws/aws-sdk-go/blob/master/example/service/ec2/filterInstances/filter_ec2_by_tag.go

# Important things to keep in mind
- Containers in a Pod share the same IPC namespace and they can also communicate with each other using standard inter-process communications like SystemV semaphores or POSIX shared memory.
  - https://linchpiner.github.io/k8s-multi-container-pods.html

# Next (priorities)
- (Readiness Probe) Check if Gz is running: `ign topic -e -t /world/default/stats -n 1 | grep -c "sim_time"`. If the above return 1, then simulation is "Up".
- Nvidia: "Node ready" -> check for gpu count before launching the pods, to make sure the device pluing worked on the new node(s).
- install monitoring (read nvidia instructions)
- install logging system
- install k8dashboard in master and write instructions
- launching 2 pods using affinity (to different nodes)
- Add a systemadmin route to run a query and return "inconsistencies" between DB and AWS/K8. So the admin can then update DB records by hand.
- (future) investigate about EC2 Spot instances (maybe faster to launch, cheaper)
- (future) Investigate using autoscaling in AWS

# PENDING:
- Consider making use of kubernetes 'namespaces' to separate applications types.
- Will need to use NodePort service to access a pod port from the outside world.
- Interesting read about Exposing pods: http://alesnosek.com/blog/2017/02/14/accessing-kubernetes-pods-from-outside-of-the-cluster/
- Networking
  - https://www.aquasec.com/wiki/display/containers/Kubernetes+Networking+101
  - https://github.com/ahmetb/kubernetes-network-policy-recipes
- Security in containers
  - https://kubernetes.io/blog/2018/07/18/11-ways-not-to-get-hacked/#8-run-containers-as-a-non-root-user
  - linux capabilities.
- Casbin related (authorization):
  - About Casbin and multithreading: https://casbin.org/docs/en/multi-threading
  - Casbin tutorial: https://zupzup.org/casbin-http-role-auth/
  - Casbin as a Service: https://github.com/casbin/casbin-server

