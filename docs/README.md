<div align="center">
  <img src="../assets/logo.png" width="200" alt="Ignition Robotics" />
  <h1>Ignition Robotics</h1>
</div>

## Table of contents

- What is Ignition Cloudsim?
- Components
    - Simulator
    - Platform
        - Machines
        - Storage
        - Orchestrator
        - Store
        - Secrets
    - Application services
        - Users
        - Simulations
- Starting a new application

## What is Ignition Cloudsim?

Ignition Cloudsim is an API that allows launching, running and processing simulations in the cloud. It currently has
support for AWS, but support for other providers can be implemented.

## Components

The cloudsim API consists of multiple sets of components that are called by applications. Every application needs at
least one set of components to run simulations in the cloud.

These components are an abstraction of a third-party service that will be consumed by cloudsim, and they were created in
order to let the application decide which cloud provider to use.

### Simulator

The simulator component is the most important component on Ignition Cloudsim. It's the one that is in charge of Starting
and Stopping simulations. To be able to perform these operations, the Simulator component is helped by a set of
Platforms and Application Services. In the following sections we'll describe how to create these components.

### Platform

A Platform is a meta-component that represents a certain region in a specific cloud provider where to launch
simulations. If you want to use `AWS` to launch simulations in `us-east-1` on a `EKS` cluster, you'll create a Platform
that represents that specific configuration. If you also need to launch on `us-east-2`, another Platform should be
created.

#### Machines

The Machines component is in charge of requesting instances to cloud providers in where to launch simulations. An EC2
implementation of this component can be found in the `ignitionrobotics.com/web/cloudsim/pkg/cloud/aws/ec2` package.

#### Storage

The Storage component is in charge of providing an API to upload simulation logs. An AWS S3 implementation of this
component can be found in the `ignitionrobotics.com/web/cloudsim/pkg/cloud/aws/s3` package.

### Application services

Ignition Cloudsim requires that the different applications implement a set of interfaces. These interfaces will allow
cloudsim to treat every application's simulations equally, but letting the developers of these applications to add
specific business logic.

## Starting a new application

