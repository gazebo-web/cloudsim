<div align="center">
  <img src="./assets/logo.png" width="200" alt="Ignition Robotics" />
  <h1>Ignition Robotics</h1>
  <p>Ignition Web Cloudsim is a web server that allows launching, running and processing simulations in the cloud. It currently has support for AWS, but support for other providers can be implemented.</p>
</div>

# Development and Code style

See https://github.com/golang/go/wiki/CodeReviewComments

The `main()` and the package `init()` functions are in `application.go`. If you are starting to read this code, probably you
will want to start from there.


# Install


## Dependencies

Install base dependencies

### System packages

```
sudo apt-get update
sudo apt-get install tar lsb-release gnupg pkg-config build-essential curl git mercurial
```

### Protobuf
```
# Sanity check
if [ -n "`which protoc`" ]; then   
    echo -e "\\e[33mWarning: protoc is already installed in this system. Proceed with caution.\\e[0m"; 
fi

curl -OL https://github.com/google/protobuf/releases/download/v3.12.3/protoc-3.12.3-linux-x86_64.zip
unzip protoc-3.12.3-linux-x86_64.zip -d protoc3
mv ./protoc3/*/ /usr/local/bin/
chown root:root /usr/local/bin/protoc
chown -R root:root /usr/local/include/google

# Install protoc-gen-go
go get -u google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

### Ign-transport

And then
```
echo "deb http://packages.osrfoundation.org/gazebo/ubuntu-stable $(lsb_release -cs) main" > /etc/apt/sources.list.d/gazebo-stable.list
echo "deb http://packages.osrfoundation.org/gazebo/ubuntu-nightly $(lsb_release -cs) main" > /etc/apt/sources.list.d/gazebo-nightly.list
apt-key adv --keyserver keyserver.ubuntu.com --recv-keys D2486D2DD83DB69272AFE98867170598AF249743
apt-get update
apt-get install -y libignition-transport7-dev
```


## Install Go

Go version 1.15 or above (NOTE: we are currently using 1.15.2)
-    Follow instructions from: https://golang.org/doc/install#install
- `curl -O https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.15.2.linux-amd64.tar.gz`


1. Make a workspace (if needed), for example:

```
mkdir -p ~/go_ws
```

1. Download server code into new directories in the workspace:

```
hg clone https://gitlab.com/ignitionrobotics/web/cloudsim ~/go_ws/src/gitlab.com/ignitionrobotics/web/cloudsim
```

1. Set necessary environment variable (needs to be set every time the environment is built)

```
export GOPATH=~/go_ws
```


## Install `dep` tool (we are currently using v0.4.1) to manage Go dependencies (in vendor/ folder)

Create a bin directory

```
mkdir ~/go_ws/bin
```

Move to the workspace's root

```
cd ~/go_ws
```

Install dep tool

```
export DEP_RELEASE_TAG=v0.4.1
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```


Download application dependencies (vendor/)

```
cd ~/go_ws/src/gitlab.com/ignitionrobotics/web/cloudsim
```

Download dependencies into `vendor` folder:

```
~/go_ws/bin/dep ensure
```

IMPORTANT NOTE: You should not use `go get` to download dependencies (instead use `dep ensure`).
Use `go get` only when you need to modify the source code of any dependency.
Alternatively, use `virtualgo` (see "Tips for local development" section below).



## Compile the protobuf files and then build the application

```
cd ~/go_ws/src/gitlab.com/ignitionrobotics/web/cloudsim/ign-transport/proto/
protoc --proto_path=. --go_out=. ignition/msgs/*.proto
```

Once proto files are generated, run:
```
cd ~/go_ws/src/gitlab.com/ignitionrobotics/web/cloudsim
go install
```


## Install mysql:

NOTE: Install a version greater than `v5.7`. In the servers, we are currently using MySQL v5.7.21


```
sudo apt-get install mysql-server
```

The installer will ask you to create a root password for mysql.


Then create the database and a user in mysql. Replace `'newuser'` with your username and `'password'` with your new password:

Login to the database server: `mysql -u root -p`

```
CREATE DATABASE cloudsim;
```

Also create a separate database to use with tests:

```
CREATE DATABASE cloudsim_test;
```

```
CREATE USER 'newuser'@'localhost' IDENTIFIED BY 'password';
```

```
GRANT ALL PRIVILEGES ON cloudsim.* TO 'newuser'@'localhost';
```

```
GRANT ALL PRIVILEGES ON cloudsim_test.* TO 'newuser'@'localhost';
```

```
FLUSH PRIVILEGES;
```

```
exit
```

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
- Label the nodes to have them ready to be used by Cloudsim:
  - Note, you can see node labels by: `kubectl get nodes --show-labels`
  - Disable master: `kubectl label nodes kube-master cloudsim_free_node=false`
  - Enable node 1: `kubectl label nodes kube-node-1 cloudsim_free_node=true` and then `kubectl label nodes kube-node-1 cloudsim_groupid=`
  - Enable node 2: `kubectl label nodes kube-node-2 cloudsim_free_node=true` and then `kubectl label nodes kube-node-2 cloudsim_groupid=`

# Environment Variables

Create an `.env` file in the root of web-cloudsim folder. They will be automatically loaded
each time the server runs.
Remember to add it to `.hgignore`.


Add the following content:

```
export KUBEADM_JOIN=this value can be ignored when using a local kubernetes

# disable max simulations per team limit
export SIMSVC_SIMULTANEOUS_SIMS_PER_TEAM=0
# Timeout to wait for newly created Pods
export SIMSVC_POD_READY_TIMEOUT_SECONDS=15

# ign-transport environment variables
export IGN_VERBOSE=1
export IGN_PARTITION=foo
# This IGN_IP is needed if your machine has multiple IPs and you want ign_transport to know which one to identify with (eg. the one from Weave)
#export IGN_IP=10.34.0.0

# SUBT specifics
# Grace period in secods to wait for the Gazebo container to finish
export SUBT_GZSERVER_TERMINATE_GRACE_PERIOD_SECONDS=30
# Disable requirement for robot images to be inside a specific ECR repo
export SUBT_DISABLE_ROBOT_IMAGE_ECR_CHECK=false
# Disable simulation score generation 
export IGN_DISABLE_SCORE_GENERATION=false

# backup gz logs to S3 ?
export AWS_GZ_LOGS_ENABLED=false
export AWS_GZ_LOGS_BUCKET=web-cloudsim-keys

export AWS_INSTANCE_NAME_PREFIX=your-prefix

# possible values: ec2, minikube
export IGN_CLOUDSIM_NODES_MGR_IMPL=minikube
export IGN_CLOUDSIM_CONNECT_TO_CLOUD=true

export IGN_CLOUDSIM_SYSTEM_ADMIN=the-sysadmin-username

# Fuel environment used by certain competitions to show and access assets
export IGN_FUEL_URL=https://fuel.ignitionrobotis.org/1.0

# DB for Users (this extra DB is needed to have access to Users data)
# The User Data lives at the ign-fuelserver's DB
export IGN_USER_DB_USERNAME=root
export IGN_USER_DB_PASSWORD=xxxxxxxx
export IGN_USER_DB_ADDRESS=localhost:3306
export IGN_USER_DB_NAME=fuel
export IGN_USER_DB_MAX_OPEN_CONNS=66

# Also add env var for the cloudsim DB
export IGN_DB_USERNAME=root
export IGN_DB_PASSWORD=xxxxxxxx
export IGN_DB_ADDRESS=localhost:3306
export IGN_DB_NAME=cloudsim
export IGN_DB_MAX_OPEN_CONNS=66

# Email
export IGN_DISABLE_SUMMARY_EMAILS=false
export IGN_DEFAULT_EMAIL_RECIPIENT=
export IGN_DEFAULT_EMAIL_SENDER=cloudsim@osrfoundation.org

# Logging
export IGN_LOGGER_LOG_STDOUT=true
# Verbosity - 4 debug, 3 info, 2 warning, 1 error, 0 critical
export IGN_LOGGER_VERBOSITY=4
export IGN_DB_LOG=false

# DO NOT log to Rollbar
export IGN_ROLLBAR_TOKEN=
export IGN_ROLLBAR_ENV=
export IGN_ROLLBAR_ROOT=

# EC2 Machines
export IGN_EC2_AVAILABILITY_ZONES=us-east-1a
export IGN_EC2_MACHINES_LIMIT=-1
# All of the following variables should be updated to match the corresponding ids in your AWS account
export IGN_EC2_AMI=ami-08861f7e7b409ed0c
export IGN_EC2_SECURITY_GROUPS=sg-0c5c791266694a3ca
export IGN_EC2_SUBNETS=subnet-0e632d68a9032ab9d

# Configure AWS access
export AWS_ACCOUNT=your account id (eg. used by AWS ECR)
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your access key id
export AWS_SECRET_ACCESS_KEY=your secret access key

# Auth0
export AUTH0_RSA256_PUBLIC_KEY=xxxxxxxxxxxxxxxxxxx
```

# User registration

An authorized user is required to perform operations. Users are planned to be 
stored in a common database shared by all services. For now, registered 
users are currently being stored in the `fuel` database setup by 
`ign-fuelserver`. 

The following steps create a user in the `fuel` database. 
This step requires a running instance of 
[ign-fuelserver](https://bitbucket.org/ignitionrobotics/ign-fuelserver) 
and [web-app](https://bitbucket.org/ignitionrobotics/web-app). If the user is
 required for local development, then these instances should be local as well.

1. (Local development) Start `ign-fuelserver`.
1. (Local development) Start `web-app`.
1. Access the `ign-fuelserver` UI.
1. Login with an Auth0 account.
1. Insert your username to create a user and assign it the logged in Auth0 
account.

# Example routes to submit and shutdown simulations

Submit a new simulation:
```
curl -k -X POST --url http://localhost:8001/1.0/simulations -F name=testSim --header 'authorization: Bearer <token>'
```

Delete an existing simulation:
```
curl -k -X DELETE --url http://localhost:8001/1.0/simulations/{simulation-groupID} --header 'authorization: Bearer <token>'
```

Tip: an easy way to test if the web-cloudsim can connect to a kubernetes cluster is to run:
```
curl http://localhost:8001/1.0/k8/countpods --header 'authorization: Bearer <token>'
```

Note: To obtain `<token>` for a user, login with the user's Auth0 account on 
your target `ign-fuelserver` instance. The JWT bearer token can be found in the 
`token` value of the Local Storage.

# Test

1. Create a Test JWT token (this is needed for tests to pass -- `go test`)

    TL;DR: Just copy and paste the following env vars in your system (`.env`)

        # Test RSA256 Private key WITHOUT the -----BEGIN RSA PRIVATE KEY----- and -----END RSA PRIVATE KEY-----
        # It is used by token-generator to generate the Test JWT Token
        export TOKEN_GENERATOR_PRIVATE_RSA256_KEY=MIICWwIBAAKBgQDdlatRjRjogo3WojgGHFHYLugdUWAY9iR3fy4arWNA1KoS8kVw33cJibXr8bvwUAUparCwlvdbH6dvEOfou0/gCFQsHUfQrSDv+MuSUMAe8jzKE4qW+jK+xQU9a03GUnKHkkle+Q0pX/g6jXZ7r1/xAK5Do2kQ+X5xK9cipRgEKwIDAQABAoGAD+onAtVye4ic7VR7V50DF9bOnwRwNXrARcDhq9LWNRrRGElESYYTQ6EbatXS3MCyjjX2eMhu/aF5YhXBwkppwxg+EOmXeh+MzL7Zh284OuPbkglAaGhV9bb6/5CpuGb1esyPbYW+Ty2PC0GSZfIXkXs76jXAu9TOBvD0ybc2YlkCQQDywg2R/7t3Q2OE2+yo382CLJdrlSLVROWKwb4tb2PjhY4XAwV8d1vy0RenxTB+K5Mu57uVSTHtrMK0GAtFr833AkEA6avx20OHo61Yela/4k5kQDtjEf1N0LfI+BcWZtxsS3jDM3i1Hp0KSu5rsCPb8acJo5RO26gGVrfAsDcIXKC+bQJAZZ2XIpsitLyPpuiMOvBbzPavd4gY6Z8KWrfYzJoI/Q9FuBo6rKwl4BFoToD7WIUS+hpkagwWiz+6zLoX1dbOZwJACmH5fSSjAkLRi54PKJ8TFUeOP15h9sQzydI8zJU+upvDEKZsZc/UhT/SySDOxQ4G/523Y0sz/OZtSWcol/UMgQJALesy++GdvoIDLfJX5GBQpuFgFenRiRDabxrE9MNUZ2aPFaFp+DyAe+b4nDwuJaW2LURbr8AEZga7oQj0uYxcYw==

        # JWT Token generated by the token-generator program using the above Test RSA keys
        # This token does not expire.
        export IGN_TEST_JWT=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXItaWRlbnRpdHkifQ.iV59-kBkZ86XKKsph8fxEeyxDiswY1zvPGi4977cHbbDEkMA3Y3t_zzmwU4JEmjbTeToQZ_qFNJGGNufK2guLy0SAicwjDmv-3dHDfJUH5x1vfi1fZFnmX_b8215BNbCBZU0T2a9DEFypxAQCQyiAQDE9gS8anFLHHlbcWdJdGw

        # A Test RSA256 Public key, without the -----BEGIN CERTIFICATE----- and -----END CERTIFICATE-----.
        # It is used to override the AUTH0_RSA256_PUBLIC_KEY when tests are run.
        export TEST_RSA256_PUBLIC_KEY=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDdlatRjRjogo3WojgGHFHYLugdUWAY9iR3fy4arWNA1KoS8kVw33cJibXr8bvwUAUparCwlvdbH6dvEOfou0/gCFQsHUfQrSDv+MuSUMAe8jzKE4qW+jK+xQU9a03GUnKHkkle+Q0pX/g6jXZ7r1/xAK5Do2kQ+X5xK9cipRgEKwIDAQAB

    In summary, in order to make `go test` work with JWT you will need to set the following env vars:

    * `TOKEN_GENERATOR_PRIVATE_RSA256_KEY`
    * `TEST_RSA256_PUBLIC_KEY`
    * `IGN_TEST_JWT`


# Run the backend server

First, make sure to set the `AUTH0_RSA256_PUBLIC_KEY` environment variable with the Auth0 RSA256 public key. This env var will be used by the backend to decode and validate any received Auth0 JWT tokens.
Note: You can get this key from: <https://osrf.auth0.com/.well-known/jwks.json> (or from your own auth0 user). Open that url in the browser and copy the value of the `x5c` field.

```
$GOPATH/bin/web-cloudsim
```


## Cloudsim optional Env vars (in the .env file)

- export IGN_CLOUDSIM_SSL_PORT=(default 4431)
- export IGN_CLOUDSIM_HTTP_PORT=(default 8001)

Also, if using Kubernetes with EC2:

- export KUBEADM_JOIN=kubeadm join IP:PORT --token ...
- export KUBE_MASTER_IP=the ip of the master node.
- export IGN_TRANSPORT_IP_INTERFACE=`ethwe` (if using Weave Network in the kubernetes cluster)
- export AWS_GZ_LOGS_BUCKET=a bucket name to upload gazebo logfiles.

Note: To backup GZ logs into S3 it will also require a `kubernetes secret` called `aws-secrets` in the cluster.
That secret should at least contain the following `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` entries.
You can create the secret by running the script `create-secrets.sh` included with the web-cloudsim code.
These are the contents of that script:

```
#/bin/bash

# This script is used to create a Kubernetes secret with name 'aws-secrets' based on existing
# environment variables: AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY.
# This Secret will be used by the cloudsim server.

kubectl create secret generic aws-secrets --from-literal=aws-access-key-id=${AWS_ACCESS_KEY_ID} --from-literal=aws-secret-access-key=${AWS_SECRET_ACCESS_KEY}
```


Note2: to find the value for the KUBEADM_JOIN env variable, from the kubernetes master, run:
`sudo kubeadm token create --ttl 0 --print-join-command` and copy the results.
Read the k8-readme.md for more info.



# Kubernates
The cloudsim server will try to connect to a Kubernates master. To do this, it will look
for the following config file: `<home>/.kube/config`.


For more details and tips about kubernetes, see the `k8-readme.md` file.



# Linter

1. Get the linter (gometalinter)

    ```
    cd ~/go_ws
    ```

    ```
    curl -L https://git.io/vp6lP | sh
    ```

1. Run the linter

    ```
    ./bin/gometalinter $(go list gitlab.com/ignitionrobotics/web/cloudsim/...)
    ```

Note you can create this bash script:

```
#!/bin/bash
curl -L https://git.io/vp6lP | sh -s -- -b $GOPATH/bin
$GOPATH/bin/gometalinter $(go list gitlab.com/ignitionrobotics/web/cloudsim/...)
```

# Troubleshooting tips

- When running in Elasticbeanstalk (EBS from now) you can login to the ebs host with `eb ssh`, and from there
you can see the web-cloudsim server logs by running `docker logs -f <container>`.
- You can tweak and restart the cloudsim server docker container running in the EBS instance by opening shell to the container
`docker exec -ti container-name /bin/bash` , update the code, rebuild the code inside the container by doing `go install`, and finally
exiting the container and restarting it `docker restart container-name`. That will relaunch the cloudsim server process.
- You can see the logs of a kubernetes pod by `kubectl logs pod-name` or `kubectl logs -f pod-name`. If the pod has multiple
containers then you can `kubectl logs pod-name container-name`.
- You can get a shell to a running pod container by `kubectl exec -it pod-name -- /bin/bash`. If the pod has multiple
containers then you can `kubectl exec -it pod-name -c container-name -- /bin/bash` to open a shell into a specific container.
- You can see status and errors of a pod by `kubectl describe pod pod-name`.
- If you need to manually update a Kubernetes Node's `labels`:
  - To remove a label: `kubectl label nodes nodeName labelName-`
  - To update a label: `kubectl label --overwrite nodes nodeName key=value`


# Development

## Mysql tips
- Backup a DB: `mysqldump -P <port> -h <dbserver> -u <user> -p <clouddbname> > mybackup.sql`
- Tip: Backup only a set of tables: `mysqldump -P <port> -h <dbserver> -u <user> -p <clouddbname> <t1 t2 t3> > mybackup.sql`
- Restore a DB: `mysql -u <user> -p <clouddbname> < mybackup.sql`


## Debugging inside a docker container

If you ever need to debug the application as if it were running in AWS or the pipelines, you need to do it from inside its docker containter.
To do that:

Most ideas taken from here:
Mysql and Docker https://docs.docker.com/samples/library/mysql/#-via-docker-stack-deploy-or-docker-compose

1. First create the docker image for the web-cloudsim server. `docker build web-cloudsim` . Write down its image ID.

1. Then run a dockerized mysql database. `docker run --name my-mysql -e MYSQL_ROOT_PASSWORD=<desired-root-pwd> -d mysql:5.7.21`
This will create a mysql docker container with an empty mysql in it.

1. Then you need to connect to that mysql container and run some commands: `docker exec -it my-mysql bash`. From inside the container, connect to mysql using the client (eg. `mysql -u root -p`) and create databases cloudsim and cloudsim_test. eg: `create database cloudsim_test;`.

1. Run the web-cloudsim docker container and link it to the database. `docker run --name web-cloudsim --rm --link my-mysql:mysql -ti <web-cloudsim-image-id> /bin/bash`. This will open a new terminal inside the server container.

1. Then from inside the server container you will need to set the Env Var that points to the linked docker mysql. eg. `export IGN_DB_ADDRESS="172.17.0.2:3306"`

After that you can run the server (from inside the container) invoking the `web-cloudsim` command, or run tests by doing `go test`.


