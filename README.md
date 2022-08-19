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

## 注意
`Dockerfile` 中的 `apk update` 执行，可能遇到如下错误:
```shell
ERROR: https://dl-cdn.alpinelinux.org/alpine/v3.17/main: temporary error (try again later)
WARNING: Ignoring https://dl-cdn.alpinelinux.org/alpine/v3.17/main: No such file or directory
```
解决方式可以[参考这里](https://github.com/gliderlabs/docker-alpine/issues/386#issuecomment-376523853)
