

version: "3.6"
services:
  kafka1:
    container_name: kafka1
    image: 'bitnami/kafka:3.3.1'
    user: root
    ports:
      - '19092:9092'
      - '19093:9093'
    environment:
      # 允许使用Kraft
      - KAFKA_ENABLE_KRAFT=yes
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      # 定义kafka服务端socket监听端口（Docker内部的ip地址和端口）
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      # 定义安全协议
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      #定义外网访问地址（宿主机ip地址和端口）
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://11.21.13.15:19092
      - KAFKA_BROKER_ID=1
      - KAFKA_KRAFT_CLUSTER_ID=iZWRiSqjZAlYwlKEqHFQWI
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@172.23.0.11:9093,2@172.23.0.12:9093,3@172.23.0.13:9093
      - ALLOW_PLAINTEXT_LISTENER=yes
      # 设置broker最大内存，和初始内存
      - KAFKA_HEAP_OPTS=-Xmx512M -Xms256M
    volumes:
      - /opt/volume/kafka/broker01:/bitnami/kafka:rw
    networks:
      netkafka:
        ipv4_address: 172.23.0.11
  kafka2:
    container_name: kafka2
    image: 'bitnami/kafka:3.3.1'
    user: root
    ports:
      - '29092:9092'
      - '29093:9093'
    environment:
      - KAFKA_ENABLE_KRAFT=yes
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://11.21.13.15:29092  #修改宿主机ip
      - KAFKA_BROKER_ID=2
      - KAFKA_KRAFT_CLUSTER_ID=iZWRiSqjZAlYwlKEqHFQWI #哪一，三个节点保持一致
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@172.23.0.11:9093,2@172.23.0.12:9093,3@172.23.0.13:9093
      - ALLOW_PLAINTEXT_LISTENER=yes
      - KAFKA_HEAP_OPTS=-Xmx512M -Xms256M
    volumes:
      - /opt/volume/kafka/broker02:/bitnami/kafka:rw
    networks:
      netkafka:
        ipv4_address: 172.23.0.12
  kafka3:
    container_name: kafka3
    image: 'bitnami/kafka:3.3.1'
    user: root
    ports:
      - '39092:9092'
      - '39093:9093'
    environment:
      - KAFKA_ENABLE_KRAFT=yes
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://11.21.13.15:39092  #修改宿主机ip
      - KAFKA_BROKER_ID=3
      - KAFKA_KRAFT_CLUSTER_ID=iZWRiSqjZAlYwlKEqHFQWI
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@172.23.0.11:9093,2@172.23.0.12:9093,3@172.23.0.13:9093
      - ALLOW_PLAINTEXT_LISTENER=yes
      - KAFKA_HEAP_OPTS=-Xmx512M -Xms256M
    volumes:
      - /opt/volume/kafka/broker03:/bitnami/kafka:rw
    networks:
      netkafka:
        ipv4_address: 172.23.0.13
networks:
  name:
  netkafka:
    driver: bridge
    name: netkafka
    ipam:
      driver: default
      config:
        - subnet: 172.23.0.0/25
          gateway: 172.23.0.1
