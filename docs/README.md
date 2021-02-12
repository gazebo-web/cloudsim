<div align="center">
  <img src="../assets/logo.png" width="200" alt="Ignition Robotics" />
  <h1>Ignition Robotics</h1>
</div>

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
and Stopping simulations. To be able to perform these operations, the Simulator component is helped by a Platform, and a
set of Services.

## Starting a new application

Ignition Cloudsim can be consumed by different applications