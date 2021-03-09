# Kubernetes Label Operator

This operator adds a `padok.fr/pod-name` label to Pods with the
`padok.fr/add-pod-name-label=true` annotation. The label's value is the Pod's
name.

## Usage

To run the operator locally on your Kubernetes cluster:

```bash
make run
```

To build and release a Docker image for the operator:

```bash
make IMG=docker.io/busser/label-operator docker-build docker-push
```

To deploy the operator to your Kubernetes cluster:

```bash
make IMG=docker.io/busser/label-operator deploy
```

## Testing

To run unit tests:

```bash
make test
```
