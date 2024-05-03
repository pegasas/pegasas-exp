package org.apache.paimon.flink;

import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.table.api.bridge.java.StreamTableEnvironment;

public class Main {
    public static void main(String[] args) throws Exception {
        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
        StreamTableEnvironment tableEnv = StreamTableEnvironment.create(env);

        tableEnv.executeSql("CREATE CATALOG my_catalog WITH (\n" +
                "    'type'='paimon',\n" +
                "    'warehouse'='/home/huangjunyao/paimon'\n" +
                ")");
        tableEnv.executeSql("USE CATALOG my_catalog");
//        tableEnv.executeSql("USE CATALOG my_catalog;").print();
    }
}
