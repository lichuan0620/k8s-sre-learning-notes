# 常用 Metrics

## 前言

刚接触 Prometheus 时，我被公司 Prometheus 服务器里的 metrics 列表震惊了：这里面没有上千也有大几百个 metrics，难道我要全部理解？！最后事实证明，并不需要。我在负责监控相关功能一年后可能仅使用过其中的一半，而常用的则更少。这篇笔记将总结运维 K8s 集群时最常用的一些 Metrics，解释它们的含义，并讨论它们的最佳用法。

### 阅读前置条件

掌握 PromQL 基本语法；熟悉 Container 基本常识；了解 K8s 中 Node, Namespace, 和 Pod 资源。对 PromQL 有疑问的同学请移步 PromQL 语法相关的笔记，本文不做讲解。

### 实验前置条件

文中涉及的所有指标均来自 [node-exporter](https://github.com/prometheus/node_exporter), [cAdvisor](https://github.com/google/cadvisor) (内嵌在 [kubelet](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/) 里)，或 [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics)。

文中查询语句默认 K8s 为推行 [Kubernetes Metrics Overhaul](https://github.com/kubernetes/enhancements/issues/1206) (即 v1.14.0) 之前的版本。如有同学碰到查询失败的，可能是因为 metrics 的 labels 名称有变。例如 cAdvisor 指标中 Container name label 改版前为 `container_name`，改版后为 `container`。如果碰到，使用 `label_replace` function 修改即可（考验 PromQL 水平的时候到了！）。

## Container 资源

K8s 中大部分 Container metrics 由 cAdvisor 提供，特点是由 `container_` 开头。

### CPU 用量

Container CPU 主要关心 `container_cpu_usage_seconds_total`，即 Container CPU 用量总和。该指标的 `_total` 的后缀告诉我们它的属性是 counter，即累计值，所以使用时往往配合 `rate` 或 `irate` function 来查看实时值。常用查询方式：

```
# 基本款，irate 加上一分钟这种较短的 range vector selector，适用于监控图
irate(container_cpu_usage_seconds_total[1m])

# 报警优化款，使用 rate 加上较大的 range vector selector；由于 CPU 数据波动剧烈，这样的
# 写法能有效避免一个短暂的波动导致报警规则被触发/解除
rate(container_cpu_usage_seconds_total[5m])

# 注意每个 Pod 除了每个 Container 会输出一条 time series，Pod 本身也会输出一条代表 Pod
# 总用量的 time series，特点是 container_name=""。因此有两种方法查询 Pod 总和。
# 使用 container_name="" 查询条件
irate(container_cpu_usage_seconds_total{container_name=""}[1m])
# 使用 sum() function
sum(
	irate(container_cpu_usage_seconds_total{container_name!=""}[1m])
) by (namespace, pod_name)
```

### 内存用量

和内存用量相关的指标有三个：

### 配额使用率
