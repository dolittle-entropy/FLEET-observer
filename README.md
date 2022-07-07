# FLEET-observer

The FLEET observer is a tool that observes a Kubernetes cluster of Dolittle resources, and stores data over time - that
can be exported to the FLEET model. It is highly specialised for the current Dolittle structure of Kubernetes resources
and highly experimental - so this is probably not what you're looking for.

## Usage

The FLEET observer currently requires `Go 1.18`, and you can run it from source using `go run . <command>` from the root
of the repository.

### Command: Observe
````shell
$ go run . observe -h
Starts the observer

Usage:
  fleet-observer observe [flags]

Flags:
  -h, --help                              help for observe
      --kubernetes.sync-interval string   The Kubernetes informer sync interval (default "1m")

Global Flags:
      --config strings                     A configuration file to load, can be specified multiple times
      --logger.format string               The logging format to use, 'json' or 'console' (default "console")
      --logger.level string                The logging minimum log level to output (default "info")
      --mongodb.connection-string string   The connection string to MongoDB (default "mongodb://localhost:27017/observer")
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
      --config strings                     A configuration file to load, can be specified multiple times
      --logger.format string               The logging format to use, 'json' or 'console' (default "console")
      --logger.level string                The logging minimum log level to output (default "info")
      --mongodb.connection-string string   The connection string to MongoDB (default "mongodb://localhost:27017/observer")
````

### Command: Drop
````shell
$ go run . observe -h
Drops the stored data in the database

Usage:
  fleet-observer drop [flags]

Flags:
  -h, --help   help for drop

Global Flags:
      --config strings                     A configuration file to load, can be specified multiple times
      --logger.format string               The logging format to use, 'json' or 'console' (default "console")
      --logger.level string                The logging minimum log level to output (default "info")
      --mongodb.connection-string string   The connection string to MongoDB (default "mongodb://localhost:27017/observer")
````
