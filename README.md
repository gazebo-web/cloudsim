<div align="center">
  <img src="./assets/logo.png" width="200" alt="Gazebo" />
  <h1>Gazebo Cloudsim</h1>
  <p>Cloudsim is an open source library for creating robotic applications in the cloud. It allows to create applications that run Gazebo simulations in Kubernetes. It currently has support for AWS, but support for other providers can be implemented.</p>


[![GitHub release](https://img.shields.io/github/release/gazebo-web/cloudsim?include_prereleases=&sort=semver&color=blue)](https://github.com/gazebo-web/cloudsim/releases/)
[![License](https://img.shields.io/badge/license-Apache-blue)](#license)
[![issues - cloudsim](https://img.shields.io/github/issues/gazebo-web/cloudsim)](https://github.com/gazebo-web/cloudsim/issues)

[![Format](https://img.shields.io/github/actions/workflow/status/gazebo-web/cloudsim/format.yaml?branch=main&label=Format)](https://github.com/gazebo-web/cloudsim/actions?query=workflow:"Format")
[![Test](https://img.shields.io/github/actions/workflow/status/gazebo-web/cloudsim/test.yaml?branch=main&label=Test)](https://github.com/gazebo-web/cloudsim/actions?query=workflow:"Test")
</div>

## Getting started

The cloudsim library provides multiple components that are used by applications to run simulations in the cloud. These components usually are an abstraction of a third-party service that will be consumed by cloudsim, and they were created in
order to let the application decide which cloud provider to use.

We'll cover each component in more detail later in this document.

### Features
- A **Simulator** component that allows to schedule simulations in a certain **Platform**.
- A **Platform** component that groups other components in such a way that represents a certain region in a cloud provider or a custom setup used by applications.
- A **Machines** component in charge of interacting with a cloud provider and performing various operations with machine instances like _EC2_ instances.
- A **Storage** component that uploads artifacts produced by a Simulation to a cloud storage like _S3_.
- A **Cost Calculator** component that allows application to keep track of their expenses and the cost of running simulations in a certain cloud provider.
- An **Orchestrator** component that interacts with a Kubernetes cluster to launch, stop and restart simulation resources, including but not limited to the following sub-components:
    - Pods
    - Nodes
    - Secrets
    - Configurations
    - Network
    - Ingresses
    - Services
- A **Cycler** component that allows cycling over different regions in a cloud provider, implementing different strategies to enable multi-region support.
- An **Email Sender** that allows notifying users about their simulations.

## Installation
Using Go CLI
```shell
go get github.com/gazebo-web/cloudsim@4.0.0
```

## Contribute
There are many ways to contribute to Gazebo Cloudsim.
* Reviewing source code changes.
* Reporting bugs.
* Creating new issues to discuss potential new features.

## Documentation

Work in progress

## License

Released under [Apache](/LICENSE) by [@gazebo-web](https://github.com/gazebo-web).
