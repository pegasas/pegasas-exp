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

package org.apache.paimon.flink.action.cdc.schema;

import org.apache.paimon.catalog.Identifier;
import org.apache.paimon.schema.Schema;

import java.util.Collections;
import java.util.List;

/** Describe a table that is not merged. */
public class UnmergedJdbcTableInfo implements JdbcTableInfo {

    private final Identifier identifier;
    private final Schema schema;

    public UnmergedJdbcTableInfo(Identifier identifier, Schema schema) {
        this.identifier = identifier;
        this.schema = schema;
    }

    @Override
    public String location() {
        return identifier.getFullName();
    }

    @Override
    public List<Identifier> identifiers() {
        return Collections.singletonList(identifier);
    }

    @Override
    public String tableName() {
        return identifier.getObjectName();
    }

    @Override
    public String toPaimonTableName() {
        // the Paimon table name should be compound of origin database name and table name
        // together to avoid name conflict
        return identifier.getDatabaseName() + "_" + identifier.getObjectName();
    }

    @Override
    public Schema schema() {
        return schema;
    }
}
