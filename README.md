# Costco

This is a small project to build a private container registry to be run on a kubernetes cluster.

## Development

This repo uses skaffold and minikube to do local development.  To get started:

1. Setup minikube (https://minikube.sigs.k8s.io/docs/start/)
2. Install skaffold (https://skaffold.dev/docs/install/#standalone-binary)

To build the required images run: `skaffold build`

To run a development environment run: `skaffold run`

Skaffold also includes ability to watch the source code and trigger rebuilds
automatically.  This can be done with: `skaffold dev`

NOTE: The example setup in `/k8s` uses hard-coded secrets (access keys, etc).  These
are purely for development purposes and should be changed for a production setup.

## Roadmap

- [ ] Push containers
- [ ] Pull containers
- [ ] Read through cache
- [ ] Auth

## Design

This container registry is backed by an S3 like storage system, Minio for example, and uses a Postgres database to store metadata.

## Server

To see the endpoints the server must implement in order to be used see [this](https://distribution.github.io/distribution/spec/api/) link.