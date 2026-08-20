# Kubeconfig generator for Magnum and Keystone

[![CI](https://github.com/mikejoh/kubeconfgen/actions/workflows/go.yml/badge.svg)](https://github.com/mikejoh/kubeconfgen/actions/workflows/go.yml)
[![Release](https://img.shields.io/github/v/release/mikejoh/kubeconfgen)](https://github.com/mikejoh/kubeconfgen/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/mikejoh/kubeconfgen)](https://goreportcard.com/report/github.com/mikejoh/kubeconfgen)

The reason i created this small tool was to have an automated way of configuring `kubectl`. Just like running e.g. `aws eks update-kubeconfig --name <cluster name> --region <region>` in AWS but in this case for a OpenStack created Magnum cluster.

### Notes
This small tool does the following:
1. Fetches the cluster specific CA certificate stored in OpenStack. Only the creator of the cluster can fetch this at the moment. Will be used as CA to be able to validate the Kubernetes API server certificate.
2. Creates a custom made kubeconfig that will utilize the client side Keystone binary when authenticting against Kubernetes.

## Prerequisites
* OpenStack Magnum (stable/stein)
* Magnum created k8s cluster of version >1.12
* Keystone server side auth component >1.16
* Keystone client side auth component >1.16
* Keystone policy ConfigMap with a `v2` auth policy

### Overview of k8s authn and authz through Keystone
`kubeconfgen` generates a `kubeconfig` for a Magnum-created Kubernetes cluster that authenticates through OpenStack Keystone, the same way `aws eks update-kubeconfig --name <cluster name> --region <region>` does for an EKS cluster on AWS. Given a cluster name and a set of Keystone credentials (as flags or the standard `OS_*` environment variables), it fetches the cluster's CA certificate from OpenStack Magnum and writes a `kubeconfig` that uses the `client-keystone-auth` exec plugin to authenticate against the Kubernetes API server, so `kubectl` can talk to the cluster without any further manual configuration.

### Installation of Keystone Server side component
Install with `go install`:
```
go install github.com/mikejoh/kubeconfgen/cmd/kubeconfgen@latest
```

Or download a prebuilt binary for Linux, macOS or Windows from the [GitHub Releases](https://github.com/mikejoh/kubeconfgen/releases/latest) page.