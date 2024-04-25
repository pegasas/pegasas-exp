package org.apache.flink.connector.sqlserver;

import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.table.api.bridge.java.StreamTableEnvironment;

public class Main {
    public static void main(String[] args) throws Exception {
        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
        StreamTableEnvironment tableEnv = StreamTableEnvironment.create(env);

        tableEnv.executeSql("CREATE TABLE XPay_Test (\n" +
                "  CauseId STRING\n" +
                ") WITH (\n" +
                "   'connector' = 'sqlserver',\n" +
                "   'url' = 'jdbc:sqlserver://xpay-sql-int-eus.database.windows.net:1433;encrypt=true;databaseName=xpay-sql-p1;" +
                "user=;password=',\n" +
                "   'table-name' = 'Donation.NPO'\n" +
                ")");
        tableEnv.executeSql("SELECT CauseId, count(1) FROM XPay_Test where CauseId = '840-352587426' group by CauseId limit 3").print();
    }
}
