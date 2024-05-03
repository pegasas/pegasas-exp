package org.apache.paimon.flink;

import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.table.api.bridge.java.StreamTableEnvironment;

public class Main {
    public static void main(String[] args) throws Exception {
        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
        StreamTableEnvironment tableEnv = StreamTableEnvironment.create(env);

        tableEnv.executeSql("CREATE CATALOG my_catalog WITH (\n" +
                "    'type'='paimon',\n" +
                "    'warehouse'='file:///home/huangjunyao/paimon'\n" +
                ")");
        tableEnv.executeSql("USE CATALOG my_catalog");
        tableEnv.executeSql("CREATE TABLE word_count (\n" +
                "    word STRING PRIMARY KEY NOT ENFORCED,\n" +
                "    cnt BIGINT\n" +
                ")");
        tableEnv.executeSql("CREATE TEMPORARY TABLE word_table (\n" +
                "    word STRING\n" +
                ") WITH (\n" +
                "    'connector' = 'datagen',\n" +
                "    'fields.word.length' = '1'\n" +
                ")");
        tableEnv.executeSql("SET 'execution.checkpointing.interval' = '10 s'");
        tableEnv.executeSql("INSERT INTO word_count SELECT word, COUNT(*) FROM word_table GROUP BY word");
//        tableEnv.executeSql("USE CATALOG my_catalog;").print();
    }
}
