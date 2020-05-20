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

**Kubernetes version**: `1.14`

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
| cloudsim-server-with-weave | sg-047577a416acc18d7 |


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

## We are currently using these env variables values
- env `SUBT_GZSERVER_LOGS_VOLUME_MOUNT_PATH` value `/tmp/ign`
- env `SIMSVC_NODE_READY_TIMEOUT_SECONDS` value `600`
- env `SIMSVC_POD_READY_TIMEOUT_SECONDS` value `900`

## Troubleshooting
- We've noticed that after a Cloudsim server restart, if the server had running simulations,
then the simulations will be regenerated but the `ign-transport topics and connections`
(and thus `/stats` messages) will be lost. We suggest Shutting down and restart those simulations too.

# We are using the following AMIs:
- Kubernetes GPU (Worker) Nodes: `ami-0884e51dacccc6d23`, name `preyna-ubuntu-18_04-CUDA_10_1-nvidia-docker_2-kubernetes_1_14-v0.2.1`. Used with `g3.4xlarge` instances. Note: in SubT these AMI and g3 instance are used for both gzserver and field-computer nodes.
- Kubernetes Master Node: `ami-05cc5ecb0a82d6c3d`, name: `cloudsim-ubuntu-bionic-18.04-docker_2_18.09.6-kubernetes_1_14-master-v0.2`. Used with `t2.medium` instance.
- VPC ID: `vpc-12af6375`
- Subnet ID: `subnet-0e632d68a9032ab9d`
- Security Group ID: `sg-0c5c791266694a3ca` (name: `kubernetes`)
- IMPORTANT: All EC2 instances (master and nodes) have the IAM Role: `arn:aws:iam::200670743174:instance-profile/cloudsim-ec2-node`. This IAM Role has attached
policies to access ECR and CloudWatch Logs.
- Using `ignitionFuel.pem` as ssh key.

# Elastic Beanstalk configuration (ie. for the Cloudsim server)
- `t2.small` instances with AMI `ami-038a24558e3a38586`.
- VPC: `vpc-12af6375`
- Region: `us-east-1`
- We are enabling Logging to CloudWatch Logs.
- Service Role (Security): `aws-elasticbeanstalk-service-role`.
- Use this AIM Instance Profile: `aws-elasticbeanstalk-ec2-role-cloudsim-server`
  - Taken from here: https://stackoverflow.com/questions/21653176/grant-s3-access-to-elastic-beanstalk-instances
  - How to add ElasticBeantalk hooks to run before the docker build:
    - https://blog.eq8.eu/article/aws-elasticbeanstalk-hooks.html
    - https://stackoverflow.com/questions/32449701/how-can-i-set-a-per-instance-env-variable-for-an-elastic-beanstalk-docker-contai
- Enabled Security Groups:
  - `sg-08a042dfd1e4a105d` (name: `awseb-e-hhjidw3spq-stack-AWSEBSecurityGroup-1JJDHJPRN3LVZ`). Still not sure if this one is needed (it was enabled when I looked)
  - `sg-047577a416acc18d7` (name: `cloudsim-server-with-weave`). Note: this one
  is needed to support Weave Network between Kubernetes nodes and this EBS server.
  This security group is tightly coupled with the `Kubernetes security group`.
  - `sg-9d31e8e6`  (name: `rds-launch-wizard`). Still not sure if this one is needed (it was enabled when I looked)
- Upload the `kube/config` needed by kubectl to connect to the master to the following S3 bucket,
`s3://web-cloudsim-keys/$IGN_CLOUDSIM_K8_CONFIG_FILENAME`. The EBS deployment process will use the filename listed
in `IGN_CLOUDSIM_K8_CONFIG_FILENAME` to download the kubernetes config file.

#
# Cloudsim server (Elasticbeanstalk): how to switch to a different kubernetes master?
- Upload the `kube config` file needed by kubectl to connect to the new master to the,
`s3://web-cloudsim-keys/$IGN_CLOUDSIM_K8_CONFIG_FILENAME` S3 bucket.
- Update the cloudsim server's `IGN_CLOUDSIM_K8_CONFIG_FILENAME` env var with the kube config file name.
- Update `KUBE_MASTER_IP` env var with the new master's amazon internal IP address.
- Update `KUBEADM_JOIN` env var with the new "kubeadm join" command needed to join Nodes to the new master.


# Instructions to setup the Kubernetes' "ECR Credentials Updater" in AWS Cluster

Using ECR to host private docker images require extra configuration steps in the Cluster.
By default `docker pull` done by kubernetes won't work with ECR. To make that work
we need to create a kubernetes `Secret` with valid tokens to access ECR.
For this, we deploy a kubernetes cronjob in the cluster that will get and refresh the
tokens periodically and will update the Kubernetes Secret so any Pod can use it.

There is a `ecr-cred-updater.yaml` in the web-cloudsim's `aws-k8/ecr` folder.
That yaml file will deploy a cronjob, a serviceaccount, a role and a rolebinding.

Note: this is only needed once per cluster.

Required steps:

1. Make sure the Secret `aws-secrets` is already present in the cluster. If not, create it
using the `create-secrets.sh` script from web-cloudsim.

1. Also make sure the ec2 instance have the correct IAM Role set (eg. `arn:aws:iam::200670743174:instance-profile/cloudsim-ec2-node`).

1. Then run `kubectl create -f ecr-cred-updater.yaml`. The output should be like:

```
role.rbac.authorization.k8s.io/ecr-cred-updater created
serviceaccount/ecr-cred-updater created
rolebinding.rbac.authorization.k8s.io/ecr-cred-updater created
cronjob.batch/ecr-cred-updater created
```

The above command created a cronjob that will refresh the ECR credentials each
8hrs(Credentials expire after ~12 hours).

Important: after setting up the cronjob , we need to manually run it once (ie. the first time),
to create get the initial ECR credentials and initialize the Secret that will be used
by kubernetes. To do this, run:

`kubectl create job --from=cronjob/ecr-cred-updater ecr-cred-first-launch-manual`

To verify that the credentials were successfully created, you can run:

1. `kubectl get secrets --all-namespaces`. There should be `aws-secrets` and a `aws-registry`.

1. `kubectl describe sa`. The serviceAccount with name `default` should list
`Image pull secrets:  aws-registry`.

1. `kubectl get cronjobs` . You should see a cronjob named `ecr-cred-updater`.

## If you want to update the ecr-cred-updater schedule

You can change the frequency at which ECR credentials are updated by patching 
the cronjob schedule:

- Once per minute: `kubectl patch cronjob/ecr-cred-updater -p '{"spec":{"schedule": "* * * * *"}}'`
- Once per hour: `kubectl patch cronjob/ecr-cred-updater -p '{"spec":{"schedule": "0 * * * *"}}'`
- Once every 8 hours: `kubectl patch cronjob/ecr-cred-updater -p '{"spec":{"schedule": "0 */8 * * *"}}'`

Note that this will not run the cronjob. The job must be run manually if you 
wish to update the `aws-registry` secret.

## If you want to delete the ecr-cred-updater

You will need to undo all the steps from above. Specifically:

- `kubectl delete job ecr-cred-first-launch-manual`
- `kubectl delete cronjob ecr-cred-updater`
- `kubectl delete rolebinding.rbac.authorization.k8s.io ecr-cred-updater`
- `kubectl delete serviceaccount ecr-cred-updater`
- `kubectl delete role.rbac.authorization.k8s.io ecr-cred-updater`
- `kubectl patch serviceaccount default -p '{"imagePullSecrets":[]}'`

# Manually Join a Node to the cluster:
- From a node:
`sudo kubeadm join <ip>:6443 --token <token> --discovery-token-ca-cert-hash <hash>`
- Tip: if you don't know the data to use, you can run this line on the master node:
`sudo kubeadm token create --print-join-command`. It will create a new token and print the join command to use by new nodes.
Tokens will be valid for 24hr by default. To change the token validity use the `--ttl 0` argument.
Eg. `sudo kubeadm token create --ttl 0 --print-join-command` will create tokens that will not expire.
- Specify node label on node join (discussion): https://github.com/kubernetes/kubeadm/issues/202
- More info on joining nodes: https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-join/
- More info on creating tokens in the master node: https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-token/



# Creating everything from scratch
In case you need to re-do everything...

## Creating the AWS Security Group
- created security group for our current VPC:
```
aws ec2 create-security-group --group-name kubernetes --description "Kubernetes Security Group" --vpc-id vpc-12af6375
```
After creating the Security Group we got its ID, and use it to add inbound rules to it.

```
aws ec2 authorize-security-group-ingress --group-id sg-0c5c791266694a3ca --protocol tcp --port 22 --cidr 0.0.0.0/0
aws ec2 authorize-security-group-ingress --group-id sg-0c5c791266694a3ca --protocol tcp --port 80 --cidr 0.0.0.0/0
aws ec2 authorize-security-group-ingress --group-id sg-0c5c791266694a3ca --protocol tcp --port 6443 --cidr 0.0.0.0/0
aws ec2 authorize-security-group-ingress --group-id sg-0c5c791266694a3ca --source-group sg-0c5c791266694a3ca --protocol all
aws ec2 authorize-security-group-ingress --group-id sg-0c5c791266694a3ca --protocol udp --port 6783-6784 --cidr 172.30.0.0/32
aws ec2 authorize-security-group-ingress --group-id sg-0c5c791266694a3ca --protocol tcp --port 6783 --cidr 172.30.0.0/32
```
Where `172.30.0.0/32` is the CIDR that both the cloudsim EBS instance and the Kubernetes EC2 instances share.

## Security Group for the EBS instance
The EBS instance (elasticbeanstalk) that will run the Cloudsim server will need to access the Weave network.
For this we need to open ports (inbound):
- Weave ports: You must permit traffic to flow through `TCP 6783` and `UDP 6783/6784`, which are Weave’s control and data ports.
- Note: The Weave Net daemon listens on localhost (127.0.0.1) TCP port 6784 for commands from other Weave Net components. This port should not be opened to other hosts.


## Steps to create a new kubernetes Master Node (and save its AMI), with Weave support and Ign-gazebo.
- Created an EC2 `t2.medium` instance with 128Gb disk, and:
  - using AMI: `ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-20190212.1` (ami-0a313d6098716f372)
  - using with "Kubernetes Security Group" using base AMI "Ubuntu 18 + HVM".
  - Choose `ignitionFuel.pem` as ssh key.
- SSH into the instance.
- Installed `docker 18.09.6` following default instructions from here: https://docs.docker.com/install/linux/docker-ce/ubuntu/
- installed mercurial and git
- Run `swapoff -a` (before installing kubernetes)
- Run `sudo su -`
- Installed Kubernetes and Kubeadm `1.14.2` following this: https://kubernetes.io/docs/setup/independent/install-kubeadm/
- Then marked package with `apt-mark hold`, to avoid unexpected updates (`apt-mark hold kubelet kubeadm kubectl`).
- (At this point you may want to create an AMI, before launching the master. You can reuse the AMI for other masters or workers).
- We've created the AMI `cloudsim-ubuntu-bionic-18.04-docker_2_18.09.6-kubernetes_1_14-master-v0.2` (`ami-05cc5ecb0a82d6c3d`)


## Launching a new Master node
- After following instructions from above...OR after launching a new ec2 instance based on the Master Node AMI...ssh into the ec2 machine and:
- Launch the master and init the cluster with `sudo kubeadm init --ignore-preflight-errors=all --node-name master --token-ttl 0`. 
  - Follow the instructions from the command output, to configure `.kube/config`.
  - Also write down the `kubeadm join xxx.xx.x.xxx:6443 --token xxxxxxxxxxx --discovery-token-ca-cert-hash yyyyyyyyyyyyyyyyyyy`. You will use it to join nodes.
- Then `sudo sysctl net.bridge.bridge-nf-call-iptables=1` to pass bridged IPv4 traffic to iptables’ chains. This is a requirement for some CNI plugins to work.
- Apply the Weave pod network: `kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"`. We're using Weave Net, as it has support for multicast and attaching EC2 instances that are not part of the kubernetes cluster.
- Once the Pod network has been installed, you can confirm that it is working by checking that the CoreDNS pod is Running in the output of `kubectl get pods --all-namespaces`. And once the CoreDNS pod is up and running, you can continue by joining your nodes.
- Install latest version of NVIDIA Device Plugin: `kubectl create -f https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin/v1.12/nvidia-device-plugin.yml`.
It is installed as a DaemonSet. You can check it was installed by: `kubectl get ds --all-namespaces`.
This plugin helps the Master to detect and use GPU based Worker nodes.
  - Other details and ideas for troubleshooting, here: https://github.com/kubernetes/kops/tree/master/hooks/nvidia-device-plugin
  - Official doc: https://github.com/NVIDIA/k8s-device-plugin
  - Tip: in the future, to check that nodes are detected to have GPUs: `kubectl describe nodes|grep -E 'gpu:\s.*[1-9]'`
- Create the `aws-secrets` from the web-cloudsim's `create-secrets.sh` script. Note that this script will need several AWS_* env
variables to be exported in the invoking bash terminal. After using the script, remember to remove the exported AWS
variables from the system (if needed).
  - You can check it was created by running `kubectl get secrets`. You should see a `aws-secrets` Secret.
- Install the `ecr-cred-updater.yaml`, to install a cronjob and its associated service-account, role and role binding to read AWS ECR tokens and
enable the kubernetes cluster to access private images from ECR repositories. See section "Instructions to setup
the ECR Credentials Updater in AWS Cluster" on this doc.
- Install the default Network Policy to block "ingress" and "egress" of all Cloudsim Pods by default. To do this run: `kubectl create -f default_deny_cloudsim.yaml`
from the `aws-k8/network_policies` folder. This policy will block connections to/from any Pod and to/from the external world
unless allowed by other policy.

- PENDING here: install addons for monitoring and centralized logging


## Instructions to create the Worker Nodes AMI (with NVIDIA)
(Note: no need to follow these steps if you are using the AMIs that we've already created)

- Note: This was done to create the cloudsim K8 Node AMI (ie. The Worker Nodes with GPU).
- Using instance type `g3.4xlarge` (with 1 GPU Testla M60).
  - Used the Kubernetes Security Group (sg-0c5c791266694a3ca)
  - Used 128 GB disk (General Purpose SSD -- gp2)
- Using base community AMI: `ubuntu-18_04-CUDA_10_1-nvidia-docker_2 (ami-0891f5dcc59fc5285)`. It comes with Go 1.10.8, CUDA 10,
nvidia-docker 2, and Nvidia as docker runtime by default.
- Then followed instructions from `optimize GPU`: https://docs.amazonaws.cn/en_us/AWSEC2/latest/UserGuide/optimize_gpu.html
- Installed mercurial, nano, git.
- Run `swapoff -a` (before installing kubernetes)
- Installed Kubernetes and Kubeadm 1.14.1 following this: https://kubernetes.io/docs/setup/independent/install-kubeadm/
- Then marked package with `apt-mark hold`, to avoid unexpected updates (`apt-mark hold kubelet kubeadm kubectl`).
- Important: Make sure to update `/etc/docker/daemon.json` to make nvidia the docker default runtime. See: https://github.com/NVIDIA/k8s-device-plugin#preparing-your-gpu-nodes.
- Created the new AMI (latest, `cloudsim-ubuntu-18_04-CUDA_10_1-nvidia-docker_2-kubernetes_1_14.10-v0.2.2`). This AMI is expected to be used with `G3` EC2 instances.

#
# Some Tips for Kubernetes

- The Node's client config is at `sudo cat /etc/kubernetes/kubelet.conf`
- To remove a node from cluster
    - 0) (at master) `kubectl drain <node-name> --delete-local-data --force --ignore-daemonsets`
    - 1) (at node) `sudo kubeadm reset`.
    - 2) (at master) `kubectl delete node <node-name>`.
- Use node affinity to send pods to specific node instance groups: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
- shutdown idle pods: https://carlosbecker.com/posts/k8s-sandbox-costs/
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



#
# Local development
#

## Install Dind-Cluster to have a local Kubernetes with 1 master and 2 slave nodes

- Pre-requisite: install liblz4-tool: `apt-get install liblz4-tool`.
- Clone the repository: https://github.com/kubernetes-sigs/kubeadm-dind-cluster
- From the dind-cluster root folder , run `./build/build-local.sh` to build its docker images.
- After the images are succesfully built, do the following:
```
$ export DIND_IMAGE="mirantis/kubeadm-dind-cluster:local"
$ sudo CNI_PLUGIN="weave" ./dind-cluster.sh up
```
- Check the cluster is working by: `kubectl get nodes` . You should see 3 nodes.
- Now `Label` the nodes to have them ready to be used by Cloudsim:
  - Note, you can see current node labels by: `kubectl get nodes --show-labels`
  1. Disable master: `kubectl label nodes kube-master cloudsim_free_node=false`
  1. Enable node 1: `kubectl label nodes kube-node-1 cloudsim_free_node=true` and then `kubectl label nodes kube-node-1 cloudsim_groupid=`
  1. Enable node 2: `kubectl label nodes kube-node-2 cloudsim_free_node=true` and then `kubectl label nodes kube-node-2 cloudsim_groupid=`


## Deploying local development changes to elasticbeanstalk server

If you want to deploy a local development version to `staging` (just as an example), you can have these 2 scripts in your
local web-cloudsim root folder:

1. A modified version of `beanstalk_deploy.py` that includes this line:
```
 BUCKET_KEY = os.getenv('APPLICATION_NAME') + '/' + VERSION_LABEL + \
     '-preyna_local_builds.zip'
```
Note the "-preyna_local_builds.zip". Let's call this new file `local_beanstalk_deploy.py` just for this example (see script below).
1. Also a new `deploy-staging.sh` script with contents:
```
set -e
set -v
rm -f /tmp/artifact.zip
zip -r /tmp/artifact.zip * .ebextensions -x vendor/\* -x .env

export S3_BUCKET="web-cloudsim-deploy"
export APPLICATION_NAME="web-cloudsim"
export APPLICATION_ENVIRONMENT="web-cloudsim-staging"
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your key here
export AWS_SECRET_ACCESS_KEY= your secret here
python local_beanstalk_deploy.py
```

With the above scripts you just need to run `deploy-staging.sh` and your current branch will be deployed to the
staging server.


## Running Gazebo from its docker image
- In your local dev machine: `sudo IGN_VERBOSE=1 docker run -it -e DISPLAY -e QT_X11_NO_MITSHM=1 -e XAUTHORITY=$XAUTH -e IGN_PARTITION -v "$XAUTH:$XAUTH" -v "/tmp/.X11-unix:/tmp/.X11-unix" -v "/etc/localtime:/etc/localtime:ro" -v "/dev/input:/dev/input" --security-opt seccomp=unconfined nkoenig/ign-gazebo:nightly ign-gazebo-server -v 4`.
- Shorter gzserver version: `sudo IGN_VERBOSE=1 docker run -it nkoenig/ign-gazebo:nightly ign-gazebo-server -v 4`.


## Some useful Gazebo commands
- To see a topic in the command line: `IGN_VERBOSE=0 ign topic -e -t /world/default/stats`
- To resume a simulation: `IGN_VERBOSE=0 ign service -s /world/default/control --reqtype ignition.msgs.WorldControl --reptype ignition.msgs.Boolean --timeout 2000 --req 'pause:false'`. Use `pause:true` if you to pause the simulation.


# 
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


# Related to ign-transport multicast and the required Weave Net plugin

## To manually join a node (eg. ec2 machine) to the k8 cluster and participate in the Weave network:

1. To manually add an external ec2 machine to the Weave's Network:
    - "External" refers to adding a new EC2 machine to an existing Weave Network, where the Weave Network is running inside a Kubernetes
    cluster, and the new EC2 machine won't participate of that cluster. As a concrete example, think of the Web-Cloudsim backend server
    (it doesn't belong to the kubernetes cluster but it has read ign-transport messages sent within the cluster).
    - First, at the new  EC2, install Weave command (https://www.weave.works/docs/net/latest/install/installing-weave/) and then run `sudo weave launch`. Also `ip route` or `ifconfig` to get the main IP address of the EC2 instance (eg. 172.30.0.170).
1. Then, there are some alternatives to join the new EC2 node to the Weave Network:
    - From the k8 master node `kubectl exec -n kube-system <weave-net-xxxxx> -c weave -- /home/weave/weave --local connect 172.30.0.170` (this is done to manually add an extra node to the Weave Network from inside the cluster). With that done, the machines will define new routes among them.
    - As an alternative, you can install the Weave command in the k8 master and perform the `weave connect <New_host_IP>` from there.
    - The best approach (to me): from the New_host run `sudo weave launch --ipalloc-init observer 172.30.0.186`, where the IP is the k8 master's IP (or another existing Weave host).
1. Once the extra host was added to the Weave cluster, if you want a child docker container to also participate in the Weave cluster, then from the host run `sudo weave attach <container>`. This command will give the containera a Weave's IP. You can see that by doing `ip route` within the container.
    - Another way to have a docker container automatically join an existing Weave network is to launch it using: `docker $(weave config) run ...`.
1. (For local dev only) If you want to make the software running in the "Host" (and not in a docker container) participate and use an IP from Weave, then after adding the Node to the Weave Network, run `sudo weave expose` in the host.
1. Note: To make the `ign-transport broadcast` work, it is important to set the `IGN_IP` environmnent variable to the current Weave's IP (usually the result of running `weave attach <container>`). Also to run it with `IGN_PARTITION=foo`.


## Automatic approach (using elasticbeanstalk):

- We've updated Cloudsim server's Dockerfile to detect the Weave's IP. This is done with the Entrypoint `docker-entrypoint.sh` script.
- We've also updated `.elasticbeanstalk` scripts to install the Weave command and automatically join the kubernetes/Weave network. For this, the EBS
instance expects the following environment variables to be set:
  - `export IGN_TRANSPORT_IP_INTERFACE=ethwe`. Setting the `IGN_TRANSPORT_IP_INTERFACE` env var with `ethwe` value will make the
  cloudsim server docker container to use the Weave IP for `IGN_IP` (needed by ign-transport process).
  - `export KUBE_MASTER_IP=<ip>`. The `KUBE_MASTER_IP` is used to automatically join the Weave Network at elasticbeanstalk startup (even
  before the cloudsim docker container is run).


## Weave Tips:

1. To debug Weave connections: `kubectl exec -n kube-system weave-net-xxxx -c weave -- /home/weave/weave --local status connections`
1. To remove a host from the Weave network: `weave forget IP` (eg. from the master).



#
# Connecting from Go code
- The client machine needs to have a `~/.kube/config` file. The cloudsim server will look for that file to connect.
- You can get it by copying the master's admin.conf file from the "Start the cluster" (ie. /etc/kubernetes/admin.conf) to the Go client machine.
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


#
# Connecting to an existing Kubernetes cluster from Staging and Production servers
- It is important to set `IGN_CLOUDSIM_K8_CONFIG_FILENAME` environment variable (in `staging` and `production` servers)
with the filename of a kubernetes config file. That file will be later used by the Golang cloudsim server to connect to
the kubernetes cluster.
- The EBS deployment process will use the filename listed in `IGN_CLOUDSIM_K8_CONFIG_FILENAME` to download a kubernetes config file from the following
S3 bucket, `s3://web-cloudsim-keys/$IGN_CLOUDSIM_K8_CONFIG_FILENAME`.
- When the server is deployed in EBS, that config file will be automatically downloaded and injected into the created docker container.
- For more information, see project files `.ebextensions/01_files.config` and `Dockerfile`.


# 
# Notes about copying Gazebo log files to S3 buckets
- We are using a side container approach, using a shared `volume` between the `gzserver` container and the `copy-to-s3` container.
The type of the shared volume is `EmptyDir`.
- The `copy-to-s3` container has a `preStop` lifecycle hook that will run the following command: `aws s3 cp <logfile> s3://<route>`.
- The GZ logs are copied into the S3 bucket identified by this `AWS_GZ_LOGS_BUCKET` env var, and under a folder named after the Team's name.
Eg. `s3://web-cloudsim-logs/my-cool-team/xxxxx.log`.
- Details about `copy-to-s3` container:
  - It is based on `infrastructureascode/aws-cli` dockerhub image.
  - Needs the `AWS_SECRET_ACCESS_KEY` and `AWS_ACCESS_KEY_ID` environment variables set.

Tip: You can use this docker image `preyna/tests:bigfile` as the `SUBT_GZSERVER_IMAGE` simulation image.
This docker image will start by creating a log file of 1 GB at `/tmp/ign/log/` and later will keep running forever in
a for loop, appending 1 extra MB each 2 seconds to the log file.
The file was created as an attempt to test big gz log files in simulations that are left there until they expire.


#
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

