# FLEET-observer

The FLEET observer is a tool that observes a Kubernetes cluster of Dolittle resources, and stores data over time - that
can be exported to the FLEET model. It is highly specialised for the current Dolittle structure of Kubernetes resources
and highly experimental - so this is probably not what you're looking for.

## FLEET domain-model
The FLEET domain-model produced by the FLEET observer is defined by the following entities and relationships:
```mermaid
  classDiagram
    class Customer {
      id: guid
      name: string
    }

    class Artifact {
      id: guid
    }
    class ArtifactVersion {
      name: string
      released: datetime
    }

    class Application {
      id: guid
      name: string
    }
    class Environment {
      name: string
    }

    class RuntimeVersion {
      major: number
      minor: number
      patch: number
      prerelease: string
      released: datetime
    }

    class Node {
      hostname: string
      image: string
      type: string
    }

    class Deployment {
      id: number
      name: string
      created: datetime
    }

    class ArtifactConfiguration {
      contentHash: string
    }
    class RuntimeConfiguration {
      contentHash: string
    }

    class DeploymentInstance {
      id: string
      started: datetime
      stopped: datetime
    }

    class Event {
      count: number
      firstTime: datetime
      lastTime: datetime
      platform: boolean
    }

    Customer <-- Artifact : developedBy
    Artifact <-- ArtifactVersion : versionOf
    ArtifactVersion <-- Deployment : usesArtifact

    Customer <-- Application : ownedBy
    Application <-- Environment : environmentOf
    Environment <-- Deployment : deployedIn

    RuntimeVersion <-- Deployment : usesRuntime

    RuntimeConfiguration <-- DeploymentInstance : usesRuntimeConfiguration
    ArtifactConfiguration <-- DeploymentInstance : usesArtifactConfiguration
    Deployment <-- DeploymentInstance : instanceOf
    Node <-- DeploymentInstance : scheduledOn
    DeploymentInstance <-- Event : happenedTo
```

## Architecture
The main usage of the FLEET observer is the `observe` command. In this mode, the FLEET observer lists and watches all known resources in the Kubernetes API, transforms them into the FLEET domain-model entities, and persists the entities and links to either a MongoDB or Neo4j database.

```mermaid
  graph LR;
    api[Kubernetes API server];
    observer(FLEET observer);
    db[(MongoDB \n or \n Neo4j)];

    api <-- watches --> observer;
    observer -- writes --> db;
```

Internally, there are multiple _observers_ that are responsible for listing and watching native Kubernetes _resources_. Whenever a change is detected (and at a regular sync-interval), thee resources are transformed into FLEET _entities_, and persisted to a _storage_ implementation. These transformations are pure functions, meaning that transforming the same resource and overwriting the resulting entities will not change the previous result. This means that (as long as the resources are not deleted in Kubernetes), the FLEET observer is stateless and should produce the same results every time it is run.

```mermaid
  graph TD;
    client[Kubernetes client];

    o_nodes[Node observer];
    o_namespaces[Namespace observer];
    o_replicasets[ReplicaSet observer];
    o_pods[Pod observer];
    o_events[Event observer];

    client --> o_nodes;
    client --> o_namespaces;
    client --> o_replicasets;
    client --> o_pods;
    client --> o_events;

    e_nodes[Nodes];
    e_customers[Customers];
    e_applications[Applications];
    e_environments[Environments]; 
    e_artifacts[Artifacts];
    e_artifact_versions[ArtifactVersions];
    e_runtime_versions[RuntimeVersions];
    e_deployments[Deployments];
    e_artifact_configurations[ArtifactConfigurations];
    e_runtime_configurations[RuntimeConfigurations];
    e_deployment_instances[DeploymentInstances];
    e_events[Events];

    o_nodes --> e_nodes;
    o_namespaces --> e_customers;
    o_namespaces --> e_applications;
    o_replicasets --> e_environments;
    o_replicasets --> e_artifacts;
    o_replicasets --> e_artifact_versions;
    o_replicasets --> e_runtime_versions;
    o_replicasets --> e_deployments;
    o_pods --> e_artifact_configurations;
    o_pods --> e_runtime_configurations;
    o_pods --> e_deployment_instances;
    o_pods --> e_events;
    o_events --> e_events;

    storage[Storage];
    e_nodes --> storage;
    e_customers --> storage;
    e_applications --> storage;
    e_environments --> storage;
    e_artifacts --> storage;
    e_artifact_versions --> storage;
    e_runtime_versions --> storage;
    e_deployments --> storage;
    e_artifact_configurations --> storage;
    e_runtime_configurations --> storage;
    e_deployment_instances --> storage;
    e_events --> storage;

    mongo[MongoDB];
    neo4j[Neo4j];
    storage --> mongo;
    storage --> neo4j;
```

## Deployment
The FLEET observer is designed to be deployed in Kubernetes as a `Deployment`, using the [dolittle/fleet-observer](https://hub.docker.com/r/dolittle/fleet-observer) Docker image. It should be configured to persist data to either a MongoDB or a Neo4j database, and it needs to run with a `ServiceAccount` that has permissions to `get`, `list`, `watch` the following resources:
 - Nodes
 - Namespaces
 - ReplicaSets
 - Pods
 - Events

## Usage

The FLEET observer currently requires `Go 1.18`, and you can run it from source using `go run . <command>` from the root
of the repository.

### Command: Observe
````shell
Starts the observer

Usage:
  fleet-observer observe [flags]

Flags:
      --cleanup.interval string           The interval to run cleanup jobs (default "1m")
  -h, --help                              help for observe
      --kubernetes.sync-interval string   The Kubernetes informer sync interval (default "1m")

Global Flags:
      --config strings                     A configuration file to load, can be specified multiple times.
      --logger.format string               The logging format to use, 'json' or 'console'. (default "console")
      --logger.level string                The logging minimum log level to output. (default "info")
      --mongodb.connection-string string   The connection string to MongoDB (default "mongodb://localhost:27017/observer")
      --neo4j.connection-string string     The connection string string to Neo4j. If not set, MongoDB will be used as storage
      --neo4j.password string              The password to use for authenticating with Neo4j. If not set, authentication will not be performed.
      --neo4j.username string              The username to use for authenticating with Neo4j. (default "neo4j")
````

### Command: Export
````shell
$ go run . export -h
Exports the stored data in the database as NDJSON

Usage:
  fleet-observer export [flags]

Flags:
  -h, --help            help for export
      --output string   The output file to export to (default "./export.ndjson")

Global Flags:
      --config strings                     A configuration file to load, can be specified multiple times.
      --logger.format string               The logging format to use, 'json' or 'console'. (default "console")
      --logger.level string                The logging minimum log level to output. (default "info")
      --mongodb.connection-string string   The connection string to MongoDB (default "mongodb://localhost:27017/observer")
      --neo4j.connection-string string     The connection string string to Neo4j. If not set, MongoDB will be used as storage
      --neo4j.password string              The password to use for authenticating with Neo4j. If not set, authentication will not be performed.
      --neo4j.username string              The username to use for authenticating with Neo4j. (default "neo4j")
````

### Command: Drop
> Note: This command only works with MongoDB at the moment
````shell
$ go run . drop -h
Drops the stored data in the database

Usage:
  fleet-observer drop [flags]

Flags:
  -h, --help   help for drop

Global Flags:
      --config strings                     A configuration file to load, can be specified multiple times.
      --logger.format string               The logging format to use, 'json' or 'console'. (default "console")
      --logger.level string                The logging minimum log level to output. (default "info")
      --mongodb.connection-string string   The connection string to MongoDB (default "mongodb://localhost:27017/observer")
      --neo4j.connection-string string     The connection string string to Neo4j. If not set, MongoDB will be used as storage
      --neo4j.password string              The password to use for authenticating with Neo4j. If not set, authentication will not be performed.
      --neo4j.username string              The username to use for authenticating with Neo4j. (default "neo4j")
````
