/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.paimon.flink.action;

import org.apache.paimon.CoreOptions;
import org.apache.paimon.flink.compact.UnawareBucketCompactionTopoBuilder;
import org.apache.paimon.flink.sink.CompactorSinkBuilder;
import org.apache.paimon.flink.source.CompactorSourceBuilder;
import org.apache.paimon.table.FileStoreTable;

import org.apache.flink.api.common.RuntimeExecutionMode;
import org.apache.flink.configuration.ExecutionOptions;
import org.apache.flink.configuration.ReadableConfig;
import org.apache.flink.streaming.api.datastream.DataStreamSource;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.table.data.RowData;

import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/** Table compact action for Flink. */
public class CompactAction extends TableActionBase {

    private List<Map<String, String>> partitions;

    public CompactAction(String warehouse, String database, String tableName) {
        this(warehouse, database, tableName, Collections.emptyMap(), Collections.emptyMap());
    }

    public CompactAction(
            String warehouse,
            String database,
            String tableName,
            Map<String, String> catalogConfig,
            Map<String, String> tableConf) {
        super(warehouse, database, tableName, catalogConfig);
        if (!(table instanceof FileStoreTable)) {
            throw new UnsupportedOperationException(
                    String.format(
                            "Only FileStoreTable supports compact action. The table type is '%s'.",
                            table.getClass().getName()));
        }
        HashMap<String, String> dynamicOptions = new HashMap<>(tableConf);
        dynamicOptions.put(CoreOptions.WRITE_ONLY.key(), "false");
        table = table.copy(dynamicOptions);
    }

    // ------------------------------------------------------------------------
    //  Java API
    // ------------------------------------------------------------------------

    public CompactAction withPartitions(List<Map<String, String>> partitions) {
        this.partitions = partitions;
        return this;
    }

    @Override
    public void build() {
        ReadableConfig conf = env.getConfiguration();
        boolean isStreaming =
                conf.get(ExecutionOptions.RUNTIME_MODE) == RuntimeExecutionMode.STREAMING;
        FileStoreTable fileStoreTable = (FileStoreTable) table;
        switch (fileStoreTable.bucketMode()) {
            case UNAWARE:
                {
                    buildForUnawareBucketCompaction(env, fileStoreTable, isStreaming);
                    break;
                }
            case FIXED:
            case DYNAMIC:
            default:
                {
                    buildForTraditionalCompaction(env, fileStoreTable, isStreaming);
                }
        }
    }

    private void buildForTraditionalCompaction(
            StreamExecutionEnvironment env, FileStoreTable table, boolean isStreaming) {
        CompactorSourceBuilder sourceBuilder =
                new CompactorSourceBuilder(identifier.getFullName(), table);
        CompactorSinkBuilder sinkBuilder = new CompactorSinkBuilder(table);

        sourceBuilder.withPartitions(partitions);
        DataStreamSource<RowData> source =
                sourceBuilder.withEnv(env).withContinuousMode(isStreaming).build();
        sinkBuilder.withInput(source).build();
    }

    private void buildForUnawareBucketCompaction(
            StreamExecutionEnvironment env, FileStoreTable table, boolean isStreaming) {
        UnawareBucketCompactionTopoBuilder unawareBucketCompactionTopoBuilder =
                new UnawareBucketCompactionTopoBuilder(env, identifier.getFullName(), table);

        unawareBucketCompactionTopoBuilder.withPartitions(partitions);
        unawareBucketCompactionTopoBuilder.withContinuousMode(isStreaming);
        unawareBucketCompactionTopoBuilder.build();
    }

    @Override
    public void run() throws Exception {
        build();
        execute("Compact job");
    }

    public List<Map<String, String>> getPartitions() {
        return partitions;
    }
}
