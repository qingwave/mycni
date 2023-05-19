# MyCNI

A simple CNI plugin for kubernetes, similar with Flannel host-gw.

## Components
`mycni`: CNI plugin for containers, create Linux Bridge, veth pair and assign IP for Pod.
`mycnid`: daemon service on each host, watch Nodes and set routes, iptables for each host.

Network architecture:
![mycni-network-arch](./doc/k8s-mycni-arch.png)

## Quick Start

Deploy MyCNI into your kubernetes cluster

> !!! It will deploy a DaemonSet application in all nodes, please run this in your dev cluster or kind/minikube

```bash
kubectl apply -f https://raw.githubusercontent.com/qingwave/mycni/main/deploy/mycni.yaml
```

## Develop

It's more easier to use [kind](https://kind.sigs.k8s.io/) cluster for test. Create a kind cluster
```bash
make kind-cluster
```

Build image
```bash
make docker-build
```

Load image into kind cluster [optional]
```bash
make kind-image-load
```

Deploy CNI
```bash
make deploy
```
