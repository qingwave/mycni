# MyCNI

A simple CNI plugin for kubernetes, similar with Flannel host-gw.

## Components
`mycni`: CNI plugin for containers, create Linux Bridge, veth pair and assign IP for Pod.
`mycnid`: daemon service on each host, watch Nodes and set routes, iptables for each host.

Network architecture:
![mycni-network-arch](./doc/k8s-mycni-arch.png)

## Build
`make`

## Deploy
`make deploy` or `kubectl apply -f deploy/mycni.yaml`

## Debugs
`apk update` may failed in Dockerfile:
```shell
ERROR: https://dl-cdn.alpinelinux.org/alpine/v3.17/main: temporary error (try again later)
WARNING: Ignoring https://dl-cdn.alpinelinux.org/alpine/v3.17/main: No such file or directory
```
Please see[here](https://github.com/gliderlabs/docker-alpine/issues/386#issuecomment-376523853)
