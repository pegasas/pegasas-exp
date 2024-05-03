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

package org.apache.paimon.utils;

import org.apache.paimon.data.Segments;
import org.apache.paimon.options.MemorySize;

import org.apache.paimon.shade.caffeine2.com.github.benmanes.caffeine.cache.Cache;
import org.apache.paimon.shade.caffeine2.com.github.benmanes.caffeine.cache.Caffeine;
import org.apache.paimon.shade.guava30.com.google.common.util.concurrent.MoreExecutors;

import java.util.function.Function;

/** Cache {@link Segments}. */
public class SegmentsCache<T> {

    private static final int OBJECT_MEMORY_SIZE = 1000;

    private final int pageSize;
    private final Cache<T, Segments> cache;

    public SegmentsCache(int pageSize, MemorySize maxMemorySize) {
        this.pageSize = pageSize;
        this.cache =
                Caffeine.newBuilder()
                        .weigher(this::weigh)
                        .maximumWeight(maxMemorySize.getBytes())
                        .executor(MoreExecutors.directExecutor())
                        .build();
    }

    public int pageSize() {
        return pageSize;
    }

    public Segments getSegments(T key, Function<T, Segments> viewFunction) {
        return cache.get(key, viewFunction);
    }

    private int weigh(T cacheKey, Segments segments) {
        return OBJECT_MEMORY_SIZE + segments.segments().size() * pageSize;
    }
}
