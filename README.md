# kubectl-pod-node-matrix

WORK IN PROGRESS!!

This plugin shows pod x node matrix with suitable colors to mitigate troubleshooting effort.

[![GoDoc](https://godoc.org/github.com/ardaguclu/kubectl-pod-node-matrix?status.svg)](https://godoc.org/github.com/ardaguclu/kubectl-pod-node-matrix)
[![Go Report Card](https://goreportcard.com/badge/ardaguclu/kubectl-pod-node-matrix)](https://goreportcard.com/report/ardaguclu/kubectl-pod-node-matrix)

# Details

Troubleshooting in Kubernetes takes some time and sorting out the real cause is sometimes overwhelming. 
Take an example of a couple of pods are not in running state, but the actual cause is node has insufficient
disk space. To reduce the amount of time being spent to this troubleshooting, `pod-node-matrix` might provide a
place for "first look at". 

`pod-node-matrix` shows pods status x node in table view. This plugin can clearly indicate that if there is a general node problem,
or can strongly suggest that node has no problem and instead deployment, service, etc. of this pod have problem.

Thanks to that, assuring that one part is working or not will definitely narrow down the places should be 
checked.

## Installation

Use [krew](https://sigs.k8s.io/krew) plugin manager to install,

```shell script
kubectl krew install pod-node-matrix
kubectl pod-node-matrix --help
```

Or manually,

`kubectl-pod-node-matrix` can be installed via:

```shell
go get github.com/ardaguclu/kubectl-pod-node-matrix/cmd/kubectl-pod-node-matrix
```

## Usage

```shell script
# Just shows the pods in default namespace
kubectl-pod-node-matrix

# Shows the pods in given namespace
kubectl-pod-node-matrix -n ${NAMESPACE}

# Shows all pods in all namespaces
kubectl-pod-node-matrix -A
```
