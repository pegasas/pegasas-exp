# Kubernetes

## Kubernetes总体架构

k8s-overall-architecture.png

https://www.cnblogs.com/caodan01/p/15102328.html

https://www.guoshaohe.com/cloud-computing/kubernetes-source-read/1249

### kubectl
kubectl 是kuberntes官方提供的命令行工具
kubectl 以交互式命令行的方式，来与kube-apiserver组件交互，通信协议是 HTTP/JSON
kubectl 发送HTTP请求，到kube-apiserver接收、处理并将结果反馈给kubectl。kubectl接收到响应并展示结果。

### client-go

client-go 是通过编程的方式与kube-apiserver进行交互，实现与kubectl相同的功能。
client-go 简单、易用，Kubernetes系统的其他组件与kube-apiserver通信的方式也基于client-go实现
如果需要对kubernetes做二次开发，可以使用 client-go
熟练使用并掌握 client-go 对 k8s 的二次开发有着关键作用。

## K8s autoscaler

https://kingjcy.github.io/post/cloud/paas/base/kubernetes/k8s-autoscaler/