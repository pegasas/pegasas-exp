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

package org.apache.paimon.flink.source.table;

import org.apache.paimon.flink.source.FlinkTableSource;

import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.table.connector.source.LookupTableSource;
import org.apache.flink.table.connector.source.abilities.SupportsDynamicFiltering;
import org.apache.flink.table.connector.source.abilities.SupportsStatisticReport;
import org.apache.flink.table.connector.source.abilities.SupportsWatermarkPushDown;
import org.apache.flink.table.data.RowData;
import org.apache.flink.table.plan.stats.TableStats;

import java.util.List;

/** The {@link BaseTableSource} with lookup, watermark, statistic and dynamic filtering. */
public class RichTableSource extends BaseTableSource
        implements LookupTableSource,
                SupportsWatermarkPushDown,
                SupportsStatisticReport,
                SupportsDynamicFiltering {

    private final FlinkTableSource source;

    public RichTableSource(FlinkTableSource source) {
        super(source);
        this.source = source;
    }

    @Override
    public RichTableSource copy() {
        return new RichTableSource(source.copy());
    }

    @Override
    public LookupRuntimeProvider getLookupRuntimeProvider(LookupContext context) {
        return source.getLookupRuntimeProvider(context);
    }

    @Override
    public void applyWatermark(WatermarkStrategy<RowData> watermarkStrategy) {
        source.pushWatermark(watermarkStrategy);
    }

    @Override
    public TableStats reportStatistics() {
        return source.reportStatistics();
    }

    @Override
    public List<String> listAcceptedFilterFields() {
        return source.listAcceptedFilterFields();
    }

    @Override
    public void applyDynamicFiltering(List<String> candidateFilterFields) {
        source.applyDynamicFiltering(candidateFilterFields);
    }
}
