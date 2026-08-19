# NERV Platform

NERV Platform is an opinionated internal developer platform for my homelab Kubernetes cluster, NERV.

The goal is to provide a simple interface for understanding and interacting with the platform without replacing the tools that already operate it, such as Argo CD, Prometheus, Grafana, and Headlamp.

This project is also a hands-on platform engineering project for exploring Go, Kubernetes APIs, GitOps, observability, and internal developer tooling.

## Current Status

NERV Platform is currently in early development.

The initial API communicates directly with the Kubernetes API using `client-go` and exposes basic cluster information.

### Endpoints

#### `GET /healthz`

Returns the health of the NERV Platform API.

```json
{
  "status": "healthy"
}
```

#### `GET /api/cluster`

Returns basic information about the NERV Kubernetes cluster.

```json
{
  "status": "healthy",
  "nodes": 1,
  "pods": 43
}
```

The node and pod counts are retrieved directly from the Kubernetes API.

## Architecture

```text
Browser / Future UI
        |
        v
NERV Platform API (Go)
        |
        +--> Kubernetes API
        |
        +--> Argo CD API      [planned]
        |
        +--> Prometheus API   [planned]
```

NERV Platform is intended to sit above the existing platform tooling rather than replace it.

## Tech Stack

- Go
- `net/http`
- Kubernetes `client-go`
- Docker
- Kubernetes
- Argo CD / GitOps

Planned integrations include Prometheus, the Argo CD API, and a React frontend.

## Development

NERV Platform currently uses the developer's local kubeconfig to authenticate with Kubernetes.

```text
Developer
    |
    v
kubeconfig
    |
    v
Kubernetes API
```

Run the service:

```bash
go run .
```

The API listens on port `8080`.

```bash
curl localhost:8080/healthz
curl localhost:8080/api/cluster
```

## Roadmap

### v0.1

- [x] Go HTTP service
- [x] Health endpoint
- [x] Kubernetes API connectivity with `client-go`
- [x] Cluster node count
- [x] Cluster pod count
- [ ] Tests
- [ ] Container image
- [ ] Kubernetes deployment
- [ ] ServiceAccount and minimal RBAC
- [ ] In-cluster authentication
- [ ] GitOps deployment through Argo CD

### Future

- Application and deployment health
- Namespace filtering
- Prometheus metrics integration
- Argo CD integration
- Kubernetes warning/event timeline
- Web frontend
- Authentication and authorization

## Project Philosophy

NERV Platform is built around actual operational needs discovered while running the NERV homelab.

Features are added when they solve real platform friction rather than simply because a technology is interesting.
