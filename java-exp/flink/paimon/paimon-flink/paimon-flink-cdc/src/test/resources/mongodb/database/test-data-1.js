// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
//  -- this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
//  the License.  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

db.getCollection('t1').insertMany([
    {
        "_id": ObjectId("100000000000000000000101"),
        "name": "scooter",
        "description": "Small 2-wheel scooter",
        "weight": 3.14
    },
    {
        "_id": ObjectId("100000000000000000000102"),
        "name": "car battery",
        "description": "12V car battery",
        "weight": 8.1
    },
    {
        "_id": ObjectId("100000000000000000000103"),
        "name": "12-pack drill bits",
        "description": "12-pack of drill bits with sizes ranging from #40 to #3",
        "weight": 0.8
    }
]);
