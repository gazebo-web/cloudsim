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
- [Starting a new application](#starting-a-new-application)

## What is Ignition Cloudsim?

Ignition Cloudsim is an API that allows launching, running and processing simulations in the cloud. It currently has
support for AWS, but support for other providers can be implemented.

## Components

The cloudsim API consists of multiple sets of components that are called by applications. Every application needs at
least one set of components to run simulations in the cloud.

These components are an abstraction of a third-party service that will be consumed by cloudsim, and they were created in
order to let the application decide which cloud provider to use.

### Simulator

The Simulator component is the most important component on Ignition Cloudsim. It's the one that is in charge of Starting
and Stopping simulations. To be able to perform these operations, the Simulator component is helped by a set of
Platforms and Application Services. In the following sections we'll describe how to create these components.

In order to start and stop simulations, the Simulator component uses actions. An action is a set of jobs that describes
how a simulation should be launched, step by step.

These jobs run on the application, but they usually perform request against the Cloudsim API. Generic jobs can be found
in the
`gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs` package.

### Platform

A Platform is a meta-component that represents a certain region in a specific cloud provider where to launch
simulations. If you want to use `AWS` to launch simulations in `us-east-1` on an `EKS` cluster, you'll create a Platform
that represents that specific configuration. If you also need to launch on `us-east-2`, another Platform should be
created.

---

#### Machines

The Machines component is in charge of requesting instances to cloud providers in where to launch simulations. An EC2
implementation of this component can be found in the `ignitionrobotics.com/web/cloudsim/pkg/cloud/aws/ec2` package.

---

#### Storage

The Storage component is in charge of providing an API to upload simulation logs. An AWS S3 implementation of this
component can be found in the `ignitionrobotics.com/web/cloudsim/pkg/cloud/aws/s3` package.

---

#### Orchestrator

The Orchestrator component is also a meta-component that includes a set of different sub-components to interact with
different resources inside a cluster.

##### Nodes

The Nodes sub-component is used to wait for recent nodes that have been created with the Machines component to join the
cluster. An implementation of the Nodes sub-component using Kubernetes can be found in the `` package.

##### Pods

The Pods sub-component is used to operate over a set of pods, it allows Cloudsim to create and destroy pods on a certain
cluster. An implementation of the Pods sub-component using Kubernetes can be found in
the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/pods` package.

##### Ingresses

An implementation of the Ingresses sub-component using Kubernetes can be found in
the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/ingresses` package. Or if you prefer using the
Gloo Ingress Controller, an implementation with Gloo can be found
here: `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/gloo`.

##### Ingress Rules

An implementation of the Ingress Rules sub-component using Kubernetes can be found in
the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/ingresses/rules` package. Or if you prefer
using the Gloo Ingress Controller, an implementation with Gloo can be found
here: `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/gloo`.

##### Services

An implementation of the Pods sub-component using Kubernetes can be found in
the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/services` package.

##### Network Policies

An implementation of the Network Policies sub-component using Kubernetes can be found in
the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/network` package.

---

#### Store

The Store component is used to provide configuration. Stores are used usually used in

---

#### Secrets

The Secrets component allows application to save secret data. An implementation using Kubernetes can be found in
the `gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets` package.

---

### Application services

Ignition Cloudsim requires that the different applications implement a set of interfaces. These interfaces will allow
cloudsim to treat every application's simulations equally, but letting the developers of these applications to add
specific business logic.

#### Users

#### Simulations

## Starting a new application

