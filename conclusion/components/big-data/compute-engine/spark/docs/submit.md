### Spark应用程序的提交

前面主要按照Spark基础组件、Spark内存管理、Spark存储体系、Spark的运算体系的顺序分析了Spark Core中涉及到的主要组件，这些组件是Spark Core的基石，
但是在分析Spark的运算体系这个部分的时候有两个涉及到的重要组件没有进行分析，它们是DAGScheduler、TaskScheduler，它们是Spark Core三大核心之一的
基于RDD的调度系统的核心组成部分。为了能够更加深入的理解它们，下面将从Spark应用程序提交到集群运行开始来对基于RDD的调度进行分析，当然这个过程中也将会
对DAGScheduler、TaskScheduler这两个组件进行深入的了解。

在将Spark应用程序提交到集群运行时，会使用spark-submit脚本，该脚本比较简单，分析该脚本知道其实最终调用的是spark-class脚本，传入的参数是
SparkSubmit及其他用户传入的参数。在spark-class中，首先会使用load-spark-env.sh加载spark的环境变量信息、定位spark jars文件等，然后调用
org.apache.spark.launcher.Main正式启动org.apache.spark.deploy.SparkSubmit的执行。

在submit()方法中，包含两个步骤：
  1. 基于集群管理器和部署模式为运行的主子类设置适当的类路径、系统属性和应用参数以准备运行环境;

  2. 使用运行环境去调用主子类的主方法;

在standalone模式下，有两种提交网关：
  1. 传统的Akka网关，使用org.apache.spark.deploy.Client作为封装;

  2. 从Spark 1.3引入的新版基于Rest的网关，这是从Spark 1.3起的默认网关，但是提交可能会失败。如果主节点终端不是一个REST服务器，可以考虑使用遗留的网关。

可以看到，其实最终调用的还是runMain()方法，在最后会调用JavaMainApplication方法，该方法是SparkApplication的实现，并且使用main()方法封装了
一个独立的Java类，在其中会最终调用mainMethod.invoke(null, args)方法，该方法执行完毕后会进入到用户编写的类的main()方法执行Spark应用程序。
至此就完成了Spark应用程序执行的提交。