---
title: "SQL Write"
weight: 2
type: docs
aliases:
- /flink/sql-write.html
---
<!--
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
-->

# SQL Write

## Syntax

```sql
INSERT { INTO | OVERWRITE } table_identifier [ part_spec ] [ column_list ] { value_expr | query };
```

For more information, please check the syntax document:

[Flink INSERT Statement](https://nightlies.apache.org/flink/flink-docs-stable/docs/dev/table/sql/insert/)

## INSERT INTO

Use `INSERT INTO` to apply records and changes to tables.

```sql
INSERT INTO my_table SELECT ...
```

`INSERT INTO` supports both batch and streaming mode. In Streaming mode, by default, it will also perform compaction, 
snapshot expiration, and even partition expiration in Flink Sink (if it is configured).

For multiple jobs to write the same table, you can refer to [dedicated compaction job]({{< ref "maintenance/dedicated-compaction#dedicated-compaction-job" >}}) for more info.

## Overwriting the Whole Table

For unpartitioned tables, Paimon supports overwriting the whole table.
(or for partitioned table which disables `dynamic-partition-overwrite` option).

Use `INSERT OVERWRITE` to overwrite the whole unpartitioned table.

```sql
INSERT OVERWRITE my_table SELECT ...
```

### Overwriting a Partition

For partitioned tables, Paimon supports overwriting a partition.

Use `INSERT OVERWRITE` to overwrite a partition.

```sql
INSERT OVERWRITE my_table PARTITION (key1 = value1, key2 = value2, ...) SELECT ...
```

### Dynamic Overwrite

Flink's default overwrite mode is dynamic partition overwrite (that means Paimon only deletes the partitions
appear in the overwritten data). You can configure `dynamic-partition-overwrite` to change it to static overwritten.

```sql
-- MyTable is a Partitioned Table

-- Dynamic overwrite
INSERT OVERWRITE my_table SELECT ...

-- Static overwrite (Overwrite whole table)
INSERT OVERWRITE my_table /*+ OPTIONS('dynamic-partition-overwrite' = 'false') */ SELECT ...
```

## Truncate tables

{{< tabs "truncate-tables-syntax" >}}

{{< tab "Flink 1.17-" >}}

You can use `INSERT OVERWRITE` to purge tables by inserting empty value.

```sql
INSERT OVERWRITE my_table /*+ OPTIONS('dynamic-partition-overwrite'='false') */ SELECT * FROM my_table WHERE false;
```

{{< /tab >}}

{{< tab "Flink 1.18+" >}}

```sql
TRUNCATE TABLE my_table;
```

{{< /tab >}}

{{< /tabs >}}

## Purging Partitions

Currently, Paimon supports two ways to purge partitions.

1. Like purging tables, you can use `INSERT OVERWRITE` to purge data of partitions by inserting empty value to them.

2. Method #1 does not support to drop multiple partitions. In case that you need to drop multiple partitions, you can submit the drop_partition job through `flink run`.

```sql
-- Syntax
INSERT OVERWRITE my_table /*+ OPTIONS('dynamic-partition-overwrite'='false') */ 
PARTITION (key1 = value1, key2 = value2, ...) SELECT selectSpec FROM my_table WHERE false;

-- The following SQL is an example:
-- table definition
CREATE TABLE my_table (
    k0 INT,
    k1 INT,
    v STRING
) PARTITIONED BY (k0, k1);

-- you can use
INSERT OVERWRITE my_table /*+ OPTIONS('dynamic-partition-overwrite'='false') */ 
PARTITION (k0 = 0) SELECT k1, v FROM my_table WHERE false;

-- or
INSERT OVERWRITE my_table /*+ OPTIONS('dynamic-partition-overwrite'='false') */ 
PARTITION (k0 = 0, k1 = 0) SELECT v FROM my_table WHERE false;
```

## Updating tables

{{< hint info >}}
Important table properties setting:
1. Only [primary key table]({{< ref "primary-key-table" >}}) supports this feature.
2. [MergeEngine]({{< ref "primary-key-table/merge-engine" >}}) needs to be [deduplicate]({{< ref "primary-key-table/merge-engine#deduplicate" >}})
   or [partial-update]({{< ref "primary-key-table/merge-engine#partial-update" >}}) to support this feature.
3. Do not support updating primary keys.
{{< /hint >}}

Currently, Paimon supports updating records by using `UPDATE` in Flink 1.17 and later versions. You can perform `UPDATE` in Flink's `batch` mode.

```sql
-- Syntax
UPDATE table_identifier SET column1 = value1, column2 = value2, ... WHERE condition;

-- The following SQL is an example:
-- table definition
CREATE TABLE my_table (
	a STRING,
	b INT,
	c INT,
	PRIMARY KEY (a) NOT ENFORCED
) WITH ( 
	'merge-engine' = 'deduplicate' 
);

-- you can use
UPDATE my_table SET b = 1, c = 2 WHERE a = 'myTable';
```

## Deleting from table

{{< tabs "delete-from-table" >}}

{{< tab "Flink 1.17+" >}}

{{< hint info >}}
Important table properties setting:
1. Only primary key tables support this feature.
2. If the table has primary keys, [MergeEngine]({{< ref "primary-key-table/merge-engine" >}}) needs to be [deduplicate]({{< ref "primary-key-table/merge-engine#deduplicate" >}}) to support this feature.
3. Do not support deleting from table in streaming mode.
{{< /hint >}}

```sql
-- Syntax
DELETE FROM table_identifier WHERE conditions;

-- The following SQL is an example:
-- table definition
CREATE TABLE my_table (
    id BIGINT NOT NULL,
    currency STRING,
    rate BIGINT,
    dt String,
    PRIMARY KEY (id, dt) NOT ENFORCED
) PARTITIONED BY (dt) WITH ( 
    'merge-engine' = 'deduplicate' 
);

-- you can use
DELETE FROM my_table WHERE currency = 'UNKNOWN';
```

{{< /tab >}}

{{< /tabs >}}
