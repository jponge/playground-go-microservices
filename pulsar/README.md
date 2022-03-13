# Event-driven services with Apache Pulsar

## Notes

- Uses `viper` and `cobra` for CLI and configuration
- Uses the Apache Pulsar client
- Uses `net/http` with [Chi](github.com/go-chi/chi) as a modern *muxer*
- The gateway handles `/record` endpoints just like the other temperature store services, so the generator can be used
- The consumer CLI uses a _reader_ to replay all events
- Metrics are exposed on `/metrics` using Prometheus
  - A custom `gateway_ingested_updates` counter tracks the number of ingested updates
  - To do so we use a custom Chi middleware function (`gateway.TrackIngestionMetrics`) to avoid manual updates
  - The Pulsar client already uses Prometheus, so we get a bunch of metrics for free besides thoses from the Go runtime 

Notes on Pulsar:
- I did not use schemas
- There are many ways to consume events (with a consumer rather than a reader, by having the client post messages to a Go channel rather than manually iterating, etc)
- The Go client brings lots of dependencies
- I haven't checked how to unify Go logging across libraries (that client logging can be "configured" but it doesn't seem to work directly with `log` but another library)

## Usage

Start Pulsar from Docker:

    $ run-pulsar-in-docker.sh

Run the gateway (assuming Minikube as a Docker environment):

    $ go run main.go --pulsar.broker-address pulsar://$(minikube ip):6650

Run the consumer (same command):

    $ go run main.go --pulsar.broker-address pulsar://$(minikube ip):6650

Generate some workload with the [temperature generator](../temperature-generator):

    $ go run main.go run