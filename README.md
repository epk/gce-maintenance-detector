## GCE Maintenance Detector

This is a simple program that detects maintenance events on Google Compute Engine (GCE) instances (for use with Google Kubernetes Engine).

It uses the Google Cloud Metadata API to subscribe to maintenance events and logs them to stdout.

### Usage

```
KO_DOCKER_REPO=$MY_DOCKER_REPO ko resolve -f infrastructure/daemonset.yml | kubectl apply --context $MY_GKE_CONTEXT --namespace $MY_NAMESPACE -f -
```
