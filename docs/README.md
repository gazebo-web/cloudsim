<div align="center">
  <img src="../assets/logo.png" width="200" alt="Ignition Robotics" />
  <h1>Ignition Robotics</h1>
</div>

## Table of contents

- [What is Ignition Cloudsim?](#what-is-ignition-cloudsim)
- [Components](#components)
    - [Simulator](#simulator)
    - [Platform](#platform)
        - [Machines](#machines)
        - [Storage](#storage)
        - [Orchestrator](#orchestrator)
            - [Nodes](#nodes)
            - [Pods](#pods)
            - [Ingresses](#ingresses)
            - [Ingress Rules](#ingress-rules)
            - [Services](#services)
            - [Network Policies](#network-policies)
        - [Store](#store)
        - [Secrets](#secrets)
    - [Application services](#application-services)
        - [Users](#users)
        - [Simulations](#simulations)
- [Configuring a Platform](#configuring-a-platform)

## What is Ignition Cloudsim?

Ignition Cloudsim is an API that allows launching, running and processing simulations in the cloud. It currently has
support for AWS and Kubernetes, but support for other providers can be implemented.

## Components

The cloudsim API consists of multiple sets of components that are called by applications. Every application needs at
least one set of components to run simulations in the cloud.

These components are an abstraction of a third-party service that will be consumed by cloudsim, and they were created in
order to let the application decide which cloud provider to use.

### Simulator

The Simulator component is the most important component on Ignition Cloudsim. It's the one that is in charge of Starting
and Stopping simulations. To be able to perform these operations, the Simulator component is helped by a set of
Platforms and Application Services. In the following sections we'll describe how to create these components.

In order to start and stop simulations, the Simulator component uses actions. An action is a list of jobs that describe
how a simulation should be launched or terminated, step by step.

These jobs run on the application-side, but they usually perform requests against the Cloudsim API. Generic jobs that
reach the Clodusim API can be found in the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs` package.

### Platform

A Platform is a meta-component that represents a certain region in a specific cloud provider where to launch
simulations. If you want to use `AWS` to launch simulations in `us-east-1` on an `EKS` cluster, you'll create a Platform
that represents that specific configuration. If you also need to launch on `us-east-2`, another Platform should be
created.

In the following sections we'll describe each component that a Platform is composed with.

---

#### Machines

The Machines component is in charge of requesting instances to cloud providers in where to launch simulations. It also
provides methods to terminate instances, and count the amount of instances running on a specific Platform.

An AWS EC2 implementation of this component can be found in the `ignitionrobotics.com/web/cloudsim/pkg/cloud/aws/ec2`
package.

---

#### Storage

The Storage component is in charge of providing an API to upload simulation logs. These logs are useful because it
allows to run the result of a certain simulation locally.

An AWS S3 implementation of this component can be found in the `ignitionrobotics.com/web/cloudsim/pkg/cloud/aws/s3`
package.

---

#### Orchestrator

The Orchestrator component is also a meta-component that includes a set of different sub-components to interact with
different resources inside a cluster.

The next sections will include a brief description of the managers available for these resources.

##### Nodes

The Nodes sub-component is used to wait for recent nodes that have been created with the Machines component to join the
cluster.

An implementation of the Nodes sub-component using Kubernetes can be found in
the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/nodes` package.

##### Pods

The Pods sub-component is used to operate over a set of pods, it allows Cloudsim to create and destroy pods on a certain
cluster.

Cloudsim launches at least 3 pods for each simulation. One pod running an Ignition Gazebo Server, and two pods (Comms
bridge and Field computer) running robot code. This code is usually provided by users that consume applications that are
using cloudsim, and they're delivered in a container image.

An implementation of the Pods sub-component using Kubernetes can be found in
the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/pods` package.

##### Ingresses

The Ingresses sub-component is in charge of managing ingresses. These ingresses are used to route traffic from users to
simulations.

An implementation of the Ingresses sub-component using Kubernetes can be found in
the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/ingresses` package. Or if you prefer using the
Gloo Ingress Controller, an implementation with Gloo can be found
here: `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/gloo`.

##### Ingress Rules

The Ingress Rules sub-component relies on the Ingresses component, and it's in charge of managing rules for a certain
Ingress. These rules describe how a certain endpoint should route into a specific simulation.

An implementation of the Ingress Rules sub-component using Kubernetes can be found in
the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/ingresses/rules` package. Or if you prefer
using the Gloo Ingress Controller, an implementation with Gloo can be found
here: `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/gloo`.

##### Services

The Services sub-component is used to manage Services.

An implementation of the Pods sub-component using Kubernetes can be found in
the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/services` package.

##### Network Policies

The Network Policies sub-component is used to manage Network Policies. Network policies provider configuration to avoid
robot pods to communicate with other robots.

An implementation of the Network Policies sub-component using Kubernetes can be found in
the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/network` package.

---

#### Store

The Store component is used to provide configuration that is requested by different jobs.

---

#### Secrets

The Secrets component allows application to save secret data.

An implementation using Kubernetes can be found in the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets` package.

---

### Application services

Ignition Cloudsim requires that the different applications consuming the API implement a set of interfaces.

These interfaces will allow Ignition Cloudsim to treat every application's simulations equally, but letting the
developers of these applications to add specific business logic.

In the following sections you'll find a brief description of the interfaces that your application should implement.

#### Users

The users interface represents a service that manages a set of users. These users are usually the ones that consume the
application, and therefore, the Ignition Cloudsim API.

#### Simulations

The simulations interface represents a service that manages a set of simulations. These simulations are a representation
of a simulation running on a specific Platform. A simulation includes specific configuration to launch, like the image
that needs to be used for robots.

## Configuring a Platform

In the following tutorial we'll configure a platform that can be used by your simulator. Our current config will use:

- Orchestrator: Kubernetes
- Machines: EC2
- Storage: S3
- Secrets: Kubernetes
- Ingress controller: Kubernetes
- Store: Environment variables

### Setting up a new Ignition Logger

```go
logger := ign.NewLoggerNoRollbar("Application", ign.VerbosityDebug)
```

### Initializing a Kubernetes orchestrator

```go
orchestrator := kubernetes.InitializeKubernetes(logger)
```

### Initializing AWS session

Using `github.com/aws/aws-sdk-go/aws/session`

```go
session, err := session.NewSession()
```

### Starting Machines using EC2

```go
ec2api := ec2.NewAPI(session)
machines := ec2.NewMachines(ec2api, logger)
```

### Starting Storage using S3

```go
s3api := s3.NewAPI(session)
storage := s3.NewStorage(s3api, logger)
```

### Initializing config store

```go
configStore := env.NewStore()
```

### Initializing secrets

```go
config, err := kubernetes.GetConfig()
if err != nil {
	panic(err)
}

clientset, err := kubernetes.NewAPI(config)
if err != nil {
    panic(err)
}

secrets := secrets.NewKubernetesSecrets(clientset.CoreV1())
```

### Initializing Platform

```go
// Components
c := platform.Components{
    Machines: machines,
    Storage:  storage,
    Cluster:  orchestrator,
    Store:    configStore,
    Secrets:  secrets,
}

// Platform
p := platform.NewPlatform(c)

```
