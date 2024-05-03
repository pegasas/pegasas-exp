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

package org.apache.paimon.table.sink;

import org.apache.paimon.data.BinaryRow;
import org.apache.paimon.data.InternalRow;
import org.apache.paimon.types.RowKind;

import static org.apache.paimon.utils.Preconditions.checkArgument;

/** A sink record contains key, row and partition, bucket information. */
public class SinkRecord {

    private final BinaryRow partition;

    private final int bucket;

    private final BinaryRow primaryKey;

    private final InternalRow row;

    public SinkRecord(BinaryRow partition, int bucket, BinaryRow primaryKey, InternalRow row) {
        checkArgument(partition.getRowKind() == RowKind.INSERT);
        checkArgument(primaryKey.getRowKind() == RowKind.INSERT);
        this.partition = partition;
        this.bucket = bucket;
        this.primaryKey = primaryKey;
        this.row = row;
    }

    public BinaryRow partition() {
        return partition;
    }

    public int bucket() {
        return bucket;
    }

    public BinaryRow primaryKey() {
        return primaryKey;
    }

    public InternalRow row() {
        return row;
    }
}
