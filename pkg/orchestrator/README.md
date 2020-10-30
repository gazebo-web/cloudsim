# Orchestrator
The `Orchestrator` package represents a component in charge of managing a single cluster. In order to do that, this 
component exposes a set of interfaces that are used as subcomponents for managing different resources in a cluster.

# Cluster
The `Cluster` interface grants you access to the different interfaces available to manage resources. 
It's the main entry point when using the orchestrator component.

### Kubernetes
Three functions have been added to initialize a Kubernetes cluster. A default implementation, a custom implementation 
and a fake one for testing purposes.

**Default implementation**: With the default implementation you avoid creating all the components that will be used by 
the orchestrator component to manage the cluster.

```go
config, err := kubernetes.GetConfig()
if err != nil {
    panic(err)
}

client, err := kubernetes.NewAPI(config)
if err != nil {
    panic(err)
}

initializer := spdy.NewSPDYInitializer(config)

logger := ign.NewLoggerNoRollbar("Kubernetes", ign.VerbosityDebug)

ks := kubernetes.NewDefaultKubernetes(client, initializer, logger)
```

These steps are in the `InitializeKubernetes`.

<hr>

**Custom implementation**: When you need to initialize the cluster with different components, `NewCustomKubernetes` 
allows you to manually pass the different subcomponents.
```go
config := Config{
    Nodes:           nil,
    Pods:            nil,
    Ingresses:       nil,
    IngressRules:    nil,
    Services:        nil,
    NetworkPolicies: nil,
}

ks := kubernetes.NewCustomKubernetes(config)
```

# Nodes
The `Nodes` interface allows you to manage nodes. Nodes are a representation of a machine that's running inside the 
cluster. As it is right now, this interface exposes a method to wait until 
nodes have joined the cluster and are ready to be used.

```go
config, err := kubernetes.GetConfig()
if err != nil {
    panic(err)
}

client, err := kubernetes.NewAPI(config)
if err != nil {
    panic(err)
}

logger := ign.NewLoggerNoRollbar("Kubernetes", ign.VerbosityDebug)

ns := NewNodes(client, logger)
```

# Pods
The `Pods` interface allows you to manage pods. Pods are a representation of a Unit of Work running inside 
different Nodes. 
- Pods can be created by passing a `CreatePodInput` instance to the `Create` method.
- Scripts and commands can be executed inside a Pod by using the `Executor`.
- Files and logs can be read from a running Pod by using the `Reader`.
- Just like the `Nodes` interface, the `Pods` interface exposes a method to wait for a certain condition to happen.
- Pods can be deleted by using the `Delete` method.

```mermaid
graph TD;
  subgraph node1[Node]
    pod1[Pod];
    pod2[Pod];
  end
  subgraph node2[Node]
    pod3[Pod];
    pod4[Pod];
  end
```

# Services
The `Services` interface allows you to manage services. A Service is an abstraction that defines a policy on how to 
access `Pods`. 
- Services can be created by passing a `CreateServiceInput` instance to the `Create` method.
- Services can be removed by using the `Delete` method.

```mermaid
graph TD;
  client[Client]-->service[Service];
  subgraph cluster[Cluster]
  pod1[Pod]-->service;
  service-->pod2[Pod];
  service-->pod3[Pod];
  end
```

# Ingresses
The `Ingresses` interface allows you to manage ingresses. An Ingress exposes endpoints that allow external access 
to services within the cluster.

Two different implementations were create for Cloudsim, one using the default Kubernetes Ingress implementation, and 
another one using a third-party application called Gloo.

From Gloo's documentation:
> Gloo is a feature-rich, Kubernetes-native ingress controller, and next-generation API gateway. 
> Gloo is exceptional in its function-level routing; its support for legacy apps, microservices and serverless; 
> its discovery capabilities; its numerous features; and its tight integration with leading open-source projects.

```mermaid
graph LR;
  client[Client]-->ingress[Ingress];
  ingress-->service[Service];
  subgraph cluster[Cluster]
  ingress;
  service-->pod1[Pod];
  service-->pod2[Pod];
  end
```

A table has been added below to help you understand how these two implementations compare to each other by the 
implementation-specific resource they consume:

| Cloudsim | Kubernetes | Gloo |
| ----- | ----- | ----- |
| Ingresses | Ingress | Virtual Services |
| IngressRules | IngressRule | Virtual Host |
| Rule | HTTPIngressRuleValue | Route |
| Path | HTTPIngressPath | RouteAction + Matcher |

# Ingress Rules
The `IngressRules` interface allows you to manage rules from a certain `Ingress`.

# Rules
The `Rule` interface groups a set of methods to upsert and remove paths from a certain rule.

# Paths
`Path` is data structure that groups all the fields needed to link a certain address (or route) to a specific service.

**Path data structure**

| Field         | Kubernetes                                 | Gloo          |
| ------------- | ------------------------------------------ | ------------- |
| UID           | HTTPIngressPath.Backend.ServiceName        | Route name    |
| Address       | HTTPIngressPath.Path                       | Regex matcher |
| Endpoint.Name | HTTPIngressPath.Backend.ServiceName        | Upstream name |
| Endpoint.Port | HTTPIngressPath.Backend.ServicePort.IntVal | -             |