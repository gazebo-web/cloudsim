# Async design

- User makes request to launch a new simulation.
    - This creates a SimDeployment record with the request and its configuration.
Status set to "Pending", and added to a queue of pending jobs.
    - This handler won't block and will return immediately. Processing will happen
  asynchronously.

- There is a pool of workers (eg. 10) in charge of launching simulations (and their nodes). Workers are go routines. Each worker in the poll will:
    - Get a record from the queue of pending jobs. Update its status to "LaunchingNodes".
    - If needed, ask the EC2 service to launch nodes.
    - Wait (block) until Nodes are ready to be used.
    - Update status to "LaunchingPods".
    - Launch the K8 Pods/Services or Deployments. And configure security.
    - Wait (block) until the pods are ready and Gazebo is Running.
    - Update the SimDeployment status to "Running".


- User makes request to know the status of a simulation.
    - This will return the last know status set in DB record.
    - In case of error during a simulation launch, the SimDeployment record's ErrorStatus field will be set to "InitializationFailed", and the user will need to request the launch of another simulation (eg. start the process from scratch). If the error happened during Termination then the error status will be "TerminationFailed".
    - After setting the ErrorStatus, the error handling process will take care of it.
    - In SimulationDeployment records, there are 2 different status fields to keep the context of which stage the simulation was when it failed.


- Error handling (eg. to rollback a failed launch or to complete an unfinished termination):
    - The error handling process tries to rollback failed launches, or continue failed terminations. The underlying idea is to release resources automatically without interrupting the flow of the simulations server.
    - There will be one async worker to deal with error handling.
    - Rollbacking a failed launch:
        - Depending on the last valid recorded status , it will start undoing the steps.
    - Completing unfinished Terminations:
        - Will redo the termination steps starting from the last valid recorded status.  
    - NOTE: in case of a 2nd error during the errorHandling process, the SimDeployment record will be marked to be manually reviewed by an administrator, by setting its error status field to 'AdminReview'.



- User makes request to stop a simulation.
    - Finds the corresponding SimDeployment record, and updates its status to "ToBeTerminated".
    - This handler will return immediately. Processing will happen asynchronously.


- There is another pool of workers (eg. 10) in charge of shutting down simulations. Each worker will:
    - Get a record from the queue of "ToBeTerminated" SimDeployments.
    - Update its status to "DeletingPods".
    - Asks Kubernetes to stop pods/deployments associated to that Simulation (GroupID)
    - Wait until pods/services are stopped.
    - Update status to "DeletingNodes".
    - Ask K8 to remove the nodes from the cluster.
    - Update status to "TerminatingInstances".
    - Use the Ec2 service to terminate the instances (k8 nodes).
    - Update SimDeployment status to "Terminated".


