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

package org.apache.paimon.mergetree.compact;

import org.apache.paimon.KeyValue;

import javax.annotation.Nullable;

/**
 * A {@link MergeFunction} where key is primary key (unique) and value is the full record, only keep
 * the latest one.
 */
public class DeduplicateMergeFunction implements MergeFunction<KeyValue> {

    private KeyValue latestKv;

    @Override
    public void reset() {
        latestKv = null;
    }

    @Override
    public void add(KeyValue kv) {
        latestKv = kv;
    }

    @Override
    public KeyValue getResult() {
        return latestKv;
    }

    public static MergeFunctionFactory<KeyValue> factory() {
        return new Factory();
    }

    private static class Factory implements MergeFunctionFactory<KeyValue> {

        private static final long serialVersionUID = 1L;

        @Override
        public MergeFunction<KeyValue> create(@Nullable int[][] projection) {
            return new DeduplicateMergeFunction();
        }
    }
}
