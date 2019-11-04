# Kubeconfig generator for Magnum and Keystone

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
Add more info here

### Installation of Keystone Server side component
Add more info here