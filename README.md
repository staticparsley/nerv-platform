# NERV Platform

NERV Platform is an opinionated internal developer platform for my homelab Kubernetes cluster, NERV.

The goal is to provide a simple interface for understanding and interacting with the platform without replacing the tools that already operate it, such as Argo CD, Prometheus, Grafana, and Headlamp.

Rather than exposing raw Kubernetes resources, NERV Platform provides a developer-facing abstraction over the underlying infrastructure. Developers should be able to understand application and platform health without needing to know the details of Kubernetes resource schemas or use `kubectl` directly.

This project is also a hands-on platform engineering project for exploring Go, Kubernetes APIs, GitOps, observability, and internal developer tooling.

## Current Status

NERV Platform is in early development and currently runs as a service inside the NERV Kubernetes cluster.

The API communicates directly with the Kubernetes API using `client-go`. It authenticates in-cluster using a dedicated Kubernetes ServiceAccount with read-only RBAC and exposes a small set of developer-oriented APIs for cluster, namespace, and deployment health.

The service is containerized, built through CI, deployed through Argo CD, and available internally over HTTPS.

## Endpoints

### `GET /healthz`

Returns the health of the NERV Platform API.

```json
{
  "status": "healthy"
}
```

This endpoint is also used by the Kubernetes readiness and liveness probes.

### `GET /api/cluster`

Returns basic information about the NERV Kubernetes cluster.

```json
{
  "status": "healthy",
  "nodes": 1,
  "pods": 44
}
```

Node and pod counts are retrieved directly from the Kubernetes API.

### `GET /api/namespaces`

Returns namespaces and their current Kubernetes phase.

```json
[
  {
    "name": "argocd",
    "status": "Active"
  },
  {
    "name": "monitoring",
    "status": "Active"
  }
]
```

### `GET /api/namespaces/{namespace}/deployments`

Returns deployments within a namespace along with a simplified health status.

```json
[
  {
    "name": "argocd-server",
    "ready": 1,
    "desired": 1,
    "status": "healthy"
  }
]
```

NERV Platform derives deployment health from Kubernetes state. A deployment with fewer ready replicas than desired replicas is reported as `degraded`.

This is an early example of the platform's intended abstraction layer: Kubernetes exposes replica state, while NERV Platform translates that state into a simpler signal for developers.

## Architecture

```text
Developer / Future UI
        |
        | HTTPS
        v
ingress-nginx
        |
        v
NERV Platform API (Go)
        |
        | ServiceAccount
        | read-only RBAC
        v
Kubernetes API
        |
        +--> Argo CD API      [planned]
        |
        +--> Prometheus API   [planned]
```

NERV Platform is intended to sit above the existing platform tooling rather than replace it.

It is not intended to become a general-purpose Kubernetes dashboard. Tools such as Headlamp already provide direct visibility into Kubernetes resources. NERV Platform instead aims to expose opinionated workflows and abstractions that are useful to application developers.

## Kubernetes Access

When running inside Kubernetes, NERV Platform uses its own ServiceAccount and Kubernetes-native in-cluster authentication.

Its current RBAC permissions are read-only and limited to the resources required by the API:

- Nodes
- Pods
- Namespaces
- Deployments

NERV Platform does not require `cluster-admin` and does not have access to Kubernetes Secrets.

For local development, the service falls back to the developer's local kubeconfig at:

```text
~/.kube/config
```

This allows the same application to run locally during development and inside Kubernetes without embedding cluster credentials in the container image.

## Delivery

Application source is maintained in Forgejo and mirrored to GitHub.

```text
Forgejo
    |
    | repository mirror
    v
GitHub
    |
    | GitHub Actions
    v
GHCR
    |
    | container image
    v
Kubernetes
```

GitHub Actions runs the Go test suite, builds the Linux container image, and publishes it to GitHub Container Registry.

The Kubernetes manifests live in the NERV infrastructure repository and are reconciled through Argo CD using the same GitOps model as the rest of the cluster.

## Internal Access

Inside the NERV network, the API is available at:

```text
https://nerv-platform.nerv.local
```

Traffic enters through ingress-nginx and is routed to the NERV Platform Kubernetes Service. TLS is managed by cert-manager.

NERV currently uses an internal certificate authority, so clients that do not trust the NERV CA may require local certificate trust configuration.

## Tech Stack

- Go
- `net/http`
- Kubernetes `client-go`
- Docker
- Kubernetes
- GitHub Actions
- GitHub Container Registry
- Argo CD / GitOps
- ingress-nginx
- cert-manager
- Go CLI

Planned integrations include Prometheus, the Argo CD API, and a web frontend.

## Development

Run the service locally:

```bash
go run .
```

The API listens on port `8080`.

```bash
curl localhost:8080/healthz
curl localhost:8080/api/cluster
curl localhost:8080/api/namespaces
curl localhost:8080/api/namespaces/argocd/deployments
```

Run the test suite:

```bash
go test ./...
```

When running locally, NERV Platform uses the current kubeconfig. When deployed inside Kubernetes, it automatically uses its Pod ServiceAccount credentials.

## Roadmap

### v0.1

- [x] Go HTTP service
- [x] Health endpoint
- [x] Kubernetes API connectivity with `client-go`
- [x] Cluster node and pod counts
- [x] Namespace discovery
- [x] Deployment discovery by namespace
- [x] Developer-facing deployment health
- [x] Tests
- [x] Container image
- [x] CI image build and publication
- [x] Kubernetes deployment
- [x] ServiceAccount and minimal read-only RBAC
- [x] In-cluster authentication
- [x] GitOps deployment through Argo CD
- [x] Kubernetes Service
- [x] Internal HTTPS ingress
- [x] Developer CLI

### Future

- Application-centric health views
- Deployment and rollout visibility
- Prometheus metrics integration
- Argo CD integration
- Kubernetes warning/event timeline
- Failure explanations and actionable diagnostics
- Self-service platform workflows
- Web frontend
- Authentication and authorization

## Project Philosophy

NERV Platform is built around actual operational needs discovered while running the NERV homelab.

Features are added when they solve real platform friction rather than to reproduce Kubernetes functionality that already exists elsewhere. The platform should provide opinionated abstractions, useful defaults, and clear health signals while leaving lower-level operational tools available when deeper investigation is required.

## CLI

NERV Platform includes a small developer CLI for interacting with the platform API.

The CLI provides a simpler interface for common platform operations without requiring direct interaction with Kubernetes.

Build the CLI:

```bash
go build -o nerv ./cmd/nerv
```

### Cluster Status

View basic cluster health and resource counts:

```bash
./nerv cluster
```

Example:

```text
NERV Cluster

Status: healthy
Nodes:  1
Pods:   44
```

### Namespaces

List namespaces and their current status:

```bash
./nerv namespaces
```

### Deployments

View deployment health for a namespace:

```bash
./nerv deployments <namespace>
```

Example:

```bash
./nerv deployments nerv-platform
```

```text
Deployments in nerv-platform

nerv-platform                1/1  healthy
```

The CLI communicates with the NERV Platform API rather than directly with the Kubernetes API.

```text
Developer
    |
    v
NERV CLI
    |
    v
NERV Platform API
    |
    v
Kubernetes API
```




> **Note:** The current CLI is intended for the NERV homelab environment and skips TLS certificate verification when connecting to the platform API. A production implementation would validate the server certificate against a trusted CA.
