# Kubeconfig generator for Magnum and Keystone

## Prerequisites
* Keystone server side auth component
* Keystone client side auth component
* Keystone policy ConfigMap

## Notes
This small tool does the following:
1. Fetches the cluster specific CA certificate stored in OpenStack. Only the creator of the cluster can fetch this at the moment. Will be used as CA to be able to validate the Kubernetes API server certificate.
2. Creates a custom made kubeconfig that will utilize the client side Keystone binary when authenticting against Kubernetes.

What i wanted to create was an automated way of getting the kubectl and configuration part configured as when running e.g. EKS when running `aws eks update-kubeconfig --name <cluster name> --region <region>`.

