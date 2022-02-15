# kubectl-pod-node-matrix

WORK IN PROGRESS!!

This plugin shows pod x node matrix with suitable colors to mitigate troubleshooting effort.

[![GoDoc](https://godoc.org/github.com/ardaguclu/kubectl-pod-node-matrix?status.svg)](https://godoc.org/github.com/ardaguclu/kubectl-pod-node-matrix)
[![Go Report Card](https://goreportcard.com/badge/github.com/ardaguclu/kubectl-pod-node-matrix)](https://goreportcard.com/report/github.com/ardaguclu/kubectl-pod-node-matrix)

## Details

Troubleshooting in Kubernetes takes some time and sorting out the root cause is sometimes overwhelming. 
Take an example of a couple of pods are not in running state, but the root cause is node has insufficient
disk space. To reduce the amount of time being spent to this troubleshooting, `pod-node-matrix` might provide a
place for "first look at". 

`pod-node-matrix` shows pods status x node in table view. This plugin can clearly indicate that if there is a general node problem,
or suggest that node has no problem and instead deployment, service, etc. of this pod can be checked.

Thanks to that, assuring that one part is working or not will narrow down the places will be examined.

## Installation

Use [krew](https://sigs.k8s.io/krew) plugin manager to install,

```shell script
kubectl krew install pod-node-matrix
kubectl pod-node-matrix --help
```

Or manually clone the repo and run;

```shell script
go run cmd/kubectl-pod-node-matrix.go
```

## Usage

```shell script
# Show pods statuses in default namespace per each node
kubectl pod-node-matrix

# Show pods statuses in given namespace per each node
kubectl pod-node-matrix -n ${NAMESPACE}

# Show pods statuses in all namespaces per each node
kubectl pod-node-matrix -A

# Shows pod statuses in verbose format. 
# It is suggested only there are few nodes.
kubectl pod-node-matrix --verbose
```

## Example

```shell script
$ go run cmd/kubectl-pod-node-matrix.go -n openshift-machine-api
+----------------------------------------------+---+---+---+---+---+
|                     PODS                     | 0 | 1 | 2 | 3 | 4 |
+----------------------------------------------+---+---+---+---+---+
| cluster-autoscaler-operator-7ffd7cb56d-nn8z6 |   | ✓ |   |   |   |
| cluster-baremetal-operator-bb769f679-h7k24   |   | ✓ |   |   |   |
| machine-api-controllers-f7c5dc6df-8v9nn      |   |   | x |   |   |
| machine-api-operator-7cdc47bb76-xddpt        |   | ✓ |   |   |   |
| metal3-f6cb6d88c-dl4h9                       | ✓ |   |   |   |   |
| metal3-image-cache-27v8f                     |   |   | x |   |   |
| metal3-image-cache-kx5xm                     |   | ✓ |   |   |   |
| metal3-image-cache-vtw26                     | ✓ |   |   |   |   |
| metal3-image-customization-c9bf787c8-vhr8b   |   |   | x |   |   |
+----------------------------------------------+---+---+---+---+---+
+---+----------+
| 0 | master-0 |
| 1 | master-1 |
| 2 | master-2 |
| 3 | worker-0 |
| 4 | worker-1 |
+---+----------+
```
