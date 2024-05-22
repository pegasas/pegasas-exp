# pegasas-exp

docker run -d -p 9092:9092 -p 9093:9093 --name kafka \
-e "KAFKA_ENABLE_KRAFT=yes" \
-e "KAFKA_CFG_PROCESS_ROLES=broker,controller" \
-e "KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER" \
-e "KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093" \
-e "KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT" \
-e "KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://127.0.0.1:9092" \
-e "KAFKA_BROKER_ID=1" \
-e "KAFKA_KRAFT_CLUSTER_ID=iZWRiSqjZAlYwlKEqHFQWI" \
-e "KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@127.0.0.1:9093" \
-e "ALLOW_PLAINTEXT_LISTENER=yes" \
-v $PWD/kafka:/bitnami/kafka:rw \
-v $PWD/minio/config:/root/.minio \
-d bitnami/kafka:3.7.0

docker pull minio/minio

docker run -d -p 9000:9000 -p 9001:9001 --name minio \
-e "MINIO_ACCESS_KEY=minioadmin" \
-e "MINIO_SECRET_KEY=minioadmin" \
-v $PWD/minio/data:/data \
-v $PWD/minio/config:/root/.minio \
minio/minio \
server /data --console-address ":9001" --address ":9000"

docker run -d -p 2181:2181 -v $PWD/zookeeper:/data --name zookeeper -d zookeeper:3.8.0

docker run --name postgresql -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -v $PWD/postgresql:/var/lib/postgresql/data -p 5432:5432 -d postgres:9.6.2

进入容器：docker exec -it ec2143b01d5a bash
使用 root 登录：su root
连接数据库：psql -U root
修改 root 用户密码：Alter user root with password 'root';

psql guide: https://www.sjkjc.com/postgresql/show-tables/

SELECT datname FROM pg_database;
CREATE DATABASE dolphinscheduler;

退出数据库连接：\q


docker run -d --name dolphinscheduler-tools     -e DATABASE="postgresql"     -e SPRING_DATASOURCE_URL="jdbc:postgresql://localhost:5432/dolphinscheduler"     -e SPRING_DATASOURCE_USERNAME="root"     -e SPRING_DATASOURCE_PASSWORD="root"     -e SPRING_JACKSON_TIME_ZONE="UTC"     --net host     apache/dolphinscheduler-tools:"latest" tools/bin/upgrade-schema.sh

https://hub.docker.com/_/mysql/tags?page=1&name=5.7

docker run -d -p 3306:3306 --privileged=true -v $PWD/mysql/logs:/logs -v $PWD/mysql/data:/mysql_data -e MYSQL_ROOT_PASSWORD=123456 --name mysql mysql:5.7.44 --character-set-server=utf8mb4 --collation-server=utf8mb4_general_ci


<root level="ERROR">
    <appender-ref ref="STDOUT" />
</root>

datasource:
  driver-class-name: org.postgresql.Driver
  url: jdbc:postgresql://127.0.0.1:5432/dolphinscheduler
  username: postgres
  password: root

默认账号密码：

admin / dolphinscheduler123

./mvnw -B clean package \
          -Dmaven.test.skip \
          -Dmaven.javadoc.skip \
          -Dspotless.skip=true \
          -Dmaven.checkstyle.skip \
          -Dmaven.deploy.skip \
          -Ddocker.push.skip=true \
          -Pdocker,release -Ddocker.tag=ci \
          -pl org.apache.dolphinscheduler:dolphinscheduler-alert-server \
          -pl dolphinscheduler-tools \
          -pl dolphinscheduler-api \
          -pl dolphinscheduler-master \
          -pl dolphinscheduler-worker -am
    
./mvnw -B clean package \
          -Dmaven.test.skip \
          -Dmaven.javadoc.skip \
          -Dspotless.skip=true \
          -Dmaven.checkstyle.skip \
          -Dmaven.deploy.skip \
          -Ddocker.push.skip=true \
          -Pdocker,release -Ddocker.tag=ci \
          -pl dolphinscheduler-api \
          -pl dolphinscheduler-master \
          -pl dolphinscheduler-worker -am
          
./mvnw validate -P helm-doc -pl :dolphinscheduler

https://www.cnblogs.com/zhuochong/p/10064350.html

docker save -o dolphinscheduler-tools.tar apache/dolphinscheduler-tools:latest
docker save -o dolphinscheduler-master.tar apache/dolphinscheduler-master:latest
docker save -o dolphinscheduler-worker.tar apache/dolphinscheduler-worker:latest
docker save -o dolphinscheduler-api.tar apache/dolphinscheduler-api:latest
docker save -o dolphinscheduler-alert-server.tar apache/dolphinscheduler-alert-server:latest

docker load -i nginx.tar
其中-i和<表示从文件输入。会成功导入镜像及相关元数据，包括tag信息

docker restart 容器ID或容器名 ：不管容器是否启动，直接重启容器
1. 首先 docker ps 查看正在运行的容器信息，显示2分钟前启动运行
2. docker restart 59ec 重启容器
3. 再次 docker ps 查看容器信息 显示 2秒前启动运行

sudo adduser myfreax

sudo deluser username

sudo deluser --remove-home username

script/env/install_env.sh
-deployUser=${deployUser:-"dolphinscheduler"}
+deployUser=${deployUser:-"huangjunyao"}

# The user to deploy DolphinScheduler for all machine we config above. For now user must create by yourself before running `install.sh`
# script. The user needs to have sudo privileges and permissions to operate hdfs. If hdfs is enabled than the root directory needs
# to be created by this user
deployUser=${deployUser:-"huangjunyao"}

./bin/spark-submit \
  --class org.apache.spark.examples.SparkPi \
  --master local[2] \
  ./examples/jars/spark-examples_2.12-3.5.0.jar \
  100

export SPARK_HOME=/opt/apache/spark-3.5.0-bin-hadoop3

org.apache.spark.examples.SparkPi

lsof -i:10016

./sbin/start-master.sh

./sbin/stop-master.sh

cat /opt/apache/spark-3.5.0-bin-hadoop3/logs/spark-huangjunyao-org.apache.spark.deploy.master.Master-1-huangjunyao-ThinkPad-P14s-Gen-3.out

Starting Spark master at spark://huangjunyao-ThinkPad-P14s-Gen-3:7077

./sbin/start-worker.sh spark://huangjunyao-ThinkPad-P14s-Gen-3:7077

./sbin/stop-workers.sh spark://huangjunyao-ThinkPad-P14s-Gen-3:7077

http://localhost:8080/

./sbin/start-thriftserver.sh \
--hiveconf hive.server2.thrift.port=10000 \
--hiveconf hive.server2.thrift.bind.host=localhost \
--master spark://huangjunyao-ThinkPad-P14s-Gen-3:7077

./sbin/stop-thriftserver.sh

./bin/spark-sql

start_hive_metastore() {

        if [ ! -f ${HIVE_HOME}/formated ];then
                schematool -initSchema -dbType mysql --verbose >  ${HIVE_HOME}/formated
        fi

        $HIVE_HOME/bin/hive --service metastore

}

start_hive_hiveserver2() {

        $HIVE_HOME/bin/hive --service hiveserver2
}

----- flink sql gateway ----------

./bin/start-cluster.sh
./bin/stop-cluster.sh
./bin/sql-gateway.sh start -Dsql-gateway.endpoint.rest.address=localhost
./bin/sql-gateway.sh stop-all
https://nightlies.apache.org/flink/flink-docs-master/docs/dev/table/sql-gateway/overview/

com.pegasas.exp.Main

------ docker -------

docker tag ca1b6b825289 registry.cn-hangzhou.aliyuncs.com/xxxxxxx:v1.0

docker tag 7bca7ec275cd billingsi.azurecr.io/dolphinscheduler/dolphinscheduler-tools:20200420

kubectl create secret docker-registry billingsi \
    --namespace default \
    --docker-server=billingsi.azurecr.io \
    --docker-username=billingsi \
    --docker-password=
