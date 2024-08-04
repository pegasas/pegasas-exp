# YARN

yarn container 默认不支持对cpu进行资源隔离，一些计算密集型任务甚至可能占满NM节点的cpu资源，从而影响到其他任务的执行效率。

默认情况下，NodeManager 使用 DefaultContainerExecutor 以 NodeManager 启动者的身份来执行启动Container等操作，安全性低且没有任何CPU资源隔离机制。

要达到这种目的，必须要使用 LinuxContainerExecutor，从而以应用提交者的身份创建文件，运行/销毁 Container。允许用户在启动Container后直接将CPU份额和进程ID写入cgroup路径的方式实现CPU资源隔离。

Cgroup 是linux kernel的一个功能，可以资源进行隔离，Yarn中现在支持对cpu/mem/io三种资源进行隔离。

https://blog.csdn.net/lsshlsw/article/details/81365050