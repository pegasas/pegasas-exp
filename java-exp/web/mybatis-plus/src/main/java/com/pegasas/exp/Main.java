package com.pegasas.exp;

import com.zaxxer.hikari.HikariDataSource;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.ResultSetMetaData;

public class Main {
    public static void main(String[] args) throws Exception {
        HikariDataSource dataSource = new HikariDataSource();
        dataSource.setJdbcUrl("jdbc:postgresql://127.0.0.1:5432/dolphinscheduler");
        dataSource.setDriverClassName("org.postgresql.Driver");
        dataSource.setUsername("root");
        dataSource.setPassword("root");
        dataSource.setIdleTimeout(60000);
        dataSource.setAutoCommit(true);
        dataSource.setMaximumPoolSize(5);
        dataSource.setMinimumIdle(1);
        dataSource.setMaxLifetime(60000 * 10);
        dataSource.setConnectionTestQuery("SELECT 1");
        Connection connection = dataSource.getConnection();
        PreparedStatement preparedStatement = connection.prepareStatement("SELECT * FROM pg_database");
        ResultSet resultSet = preparedStatement.executeQuery();
        while(resultSet.next()){
            ResultSetMetaData metaData = resultSet.getMetaData();
            int columnCount = metaData.getColumnCount();
            for (int i = 1; i <= columnCount; i++) {
                String name = metaData.getColumnLabel(i);
                String field = resultSet.getString(i);
                System.out.printf("%s:%s\t",name,field);
            }
            System.out.println();
        }
    }
}