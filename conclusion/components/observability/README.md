# Monitoring
# [Observability](http://pigpi.cn/17779.html)

![observability.png](./observability.png)

- 日志（Logging）：日志的职责是记录离散事件，通过这些记录事后分析出程序的行为，譬如曾经调用过什么方法，曾经操作过哪些数据，等等。输出日志的确很容易，但收集和分析日志却可能会很复杂，面对成千上万的集群节点，面对迅速滚动的事件信息，面对数以 TB 计算的文本，传输与归集都并不简单。对大多数程序员来说，分析日志也许就是最常遇见也最有实践可行性的“大数据系统”了。
- 追踪（Tracing）：单体系统时代追踪的范畴基本只局限于栈追踪（Stack Tracing），调试程序时，在 IDE 打个断点，看到的 Call Stack 视图上的内容便是追踪；编写代码时，处理异常调用了 Exception::printStackTrace()方法，它输出的堆栈信息也是追踪。微服务时代，追踪就不只局限于调用栈了，一个外部请求需要内部若干服务的联动响应，这时候完整的调用轨迹将跨越多个服务，同时包括服务间的网络传输信息与各个服务内部的调用堆栈信息，因此，分布式系统中的追踪在国内常被称为“全链路追踪”（后文就直接称“链路追踪”了），许多资料中也称它为“分布式追踪”（Distributed Tracing）。追踪的主要目的是排查故障，如分析调用链的哪一部分、哪个方法出现错误或阻塞，输入输出是否符合预期，等等。
- 度量（Metrics）：度量是指对系统中某一类信息的统计聚合。譬如，JVM的内存大小、各分代的用量、峰值的线程数、垃圾收集的吞吐量、频率等等。度量的主要目的是监控（Monitoring）和预警（Alert），如某些度量指标达到风险阈值时触发事件，以便自动处理或者提醒管理员介入。

所谓“可观测性”，就是在任何相关方采取任何类型的行动时，都会产生包含了数字化特征的数据，这些数据都可以称之为“可观测数据”，如日志、痕迹、API 调用、停留时间、下载和文件传输等。应用可观测性用一种高度统筹和整合的方式将所有可观测的特征数据进行反馈，创造出一个决策循环，从而提高组织决策的有效性。

主流开源技术
鉴于开源社区主流技术方案还是围绕三大支柱来做的，所以本章还是按照日志、追踪、度量三方面来描述。

监控（Monitoring）
监控领域的主要技术栈包括：Prometheus，Cortex，ZABBIX，Grafana，sysdig，Prometheus已经成为事实上的监控标准，主流监控技术方案都会围绕Prometheus来做。

Prometheus提供了通用的数据模型和快捷数据采集、存储和查询接口，主要包括以下组件：

Prometheus server：服务核心组件，采用pull方式收集监控数据，通过http协议传输。并存储时间序列数据。Prometheus server 由三个部分组成：Retrival，Storage，PromQL

Retrieval：负责在活跃的target主机上抓取监控指标数据。

Storage：存储，主要是把采集到的数据存储到磁盘中，默认为15天。

PromQL：是Prometheus提供的查询语言模块。

client Library：客户端库，目的在于为那些期望原生提供数据展示功能的应用程序提供便捷的开发途径。

Exporters：指标暴露器，负责收集应用程序或服务的性能指标数据，并通过HTTP接口供Prometheus Server获取。社区提供了丰富的Exporter选型，几乎可以收集一切想要收集的数据，还可以自己编写Exporter。

以Prometheus为中心辐射了很多周边生态组件，比如：

Service Discovery，用于自动发现target主机
Alertmanager，独立的告警模块，支持邮件、钉钉、企业微信等主流方式
Pushgateway，支持以push方式暂存数据
Grafana，将数据可视化和分析
Prometheus本身也有局限性，比如：

不支持集群，当然社区也提供了方案，比如Cortex、thanos
不适合存储事件及日志等，它更多地展示的是趋势性的监控，而非精准数据；
不长存数据，Prometheus认为度量数据无需长期存储，所有如果有需求，需要自己提供长存储方案。
日志（Logging）
日志领域的主流技术方案除却经典Elastic Stack外，新兴的如Grafana Loki、Graylog。日志的技术方案链路较长，主要分为日志采集、日志缓冲、日志清洗、日志存储、日志可视化。

日志采集

常见的日志采集工具包括Logstash、Filebeat、Fluentd/Fluent-Bit。

Logstash由于对cpu、内存、io占用太大已经基本不使用在日志采集上；

Filebeat非常轻量稳定，主要采集目标是文件；

Fluentd/Fluent-Bit对传统日志采集的流程进一步的细分，能够将散乱的日志来源统一解析过滤并存储，是容器环境下日志采集的首选。

此外Loki使用Promtail采集日志，还有Flume、Loggie等等。

日志缓冲

Kafka等消息队列都可以用作日志缓冲，除非日志数据量过于庞大，否则一般会省去日志缓冲。

日志清洗

我所知的专门用于日志清洗的组件就是logstash，其他如新兴的日志方案Loki、Graylog等都将日志过滤、清洗功能集成在组件内部。

日志存储

首当其冲的是使用范围最广的Elasticsearch，基于倒排索引的特性使得ES查询速度快、QPS较高。

但是由于大规模建设成本问题，如今人们也在探索新的存储方式，比如clickhouse，吞吐量更高，查询性能更快。

而Loki直接使用本地存储或者对象存储。

日志可视化

可视化上主要有Elastic Stack中的Kibana，Grafana Stack中的Grafana，而Graylog自带可视化模块Graylog-web。

追踪（Tracing）
追踪领域最早可以追溯到Google的Dapper系统，比起日志与度量，追踪这个领域的产品竞争要相对激烈得多。一方面，目前还没有像日志、度量那样出现具有明显统治力的产品，仍处于群雄混战的状态。另一方面，几乎市面上所有的追踪系统都是以 Dapper 的论文为原型发展出来的，基本上都算是同门师兄弟，功能上并没有太本质的差距，却又受制于实现细节，彼此互斥，很难搭配工作。

一个完整的分布式追踪系统应该由数据收集、数据存储和数据展示三个子系统构成，而狭义上讲的追踪则就只是特指链路追踪数据的收集部分。譬如Spring Cloud Sleuth就属于狭义的追踪系统，通常会搭配 Zipkin 作为数据展示，搭配 Elasticsearch 作为数据存储来组合使用。除此之外，下面谈到的大多属于完整的追踪系统。

追踪的数据通常由业务应用产生，主流的收集方式包括：

（1）日志，Dapper系统就是通过从日志中筛选出追踪类型数据进行汇聚。

（2）应用探针，以Java应用为例，通常采用Java Agent字节码增强技术注入应用，对业务零侵入。

（3）服务网格下的SideCar边车代理，技术无关性，通过网络流量采集数据，缺点就是无法深入到方法级别，比如Envoy。

而在存储上大部分组件也都以ES为主，强如SkyWalking支持ES、H2、Mysql、TIDB、Sharding sphere等方式。

追踪领域目前比较成熟、出现较早的的一些组件包括：Twitter开源的Zipkin、Apache的SkyWalking、PinPoint、大众点评开源的CAT。新兴的比如CNCF jaeger、Grafana Tempo。

主流落地方案
ELK/EFK + GPE + SkyWalking
这套方案就是奔着建设三套系统的思路去的，也是最经典的可观测性方案：

日志方面Elasticsearch做存储，FileBeat或者Fluentd做采集，Kibana做展现，必要的话接入Logstash进行日志清洗，Kafka做日志缓冲
监控方面同样Elasticsearch做存储，Prometheus结合各种sdk、exporter、agent做指标采集和接入，Grafana做展示
追踪方面自由度大一些，我个人习惯使用SkyWalking
这套方案能力均衡，如果没有什么特殊需求，选万金油组合总没错的。

OpenTelemetry + Loki + Tempo + Prometheus + Grafana
这套方案算是新兴的轻量化方案，主要围绕Grafana Stack构建。

日志使用Loki，日志的收集可以通过loki4j收集应用打印日志，也可以通过Promtail收集其他类型日志
监控依然基于Prometheus和Grafana的无缝集成
追踪方面使用OpenTelemetry采集器进行数据采集，发送到Tempo存储并集成Grafana进行展示
ElasticSearch + Kibana + Beats + Elastic APM
这套方案主要围绕Elastic Stack构建，数据采集基于Beats家族组件

日志还是基于EFK
监控指标的采集利用Beats家族的Metricbeat组件
追踪使用Elastic APM
统一由Kibana进行数据展示

