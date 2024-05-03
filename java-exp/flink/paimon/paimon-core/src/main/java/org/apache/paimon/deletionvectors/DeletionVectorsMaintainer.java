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

package org.apache.paimon.deletionvectors;

import org.apache.paimon.annotation.VisibleForTesting;
import org.apache.paimon.data.BinaryRow;
import org.apache.paimon.index.IndexFileHandler;
import org.apache.paimon.index.IndexFileMeta;

import javax.annotation.Nullable;

import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;

import static org.apache.paimon.deletionvectors.DeletionVectorsIndexFile.DELETION_VECTORS_INDEX;

/** Maintainer of deletionVectors index. */
public class DeletionVectorsMaintainer {

    private final IndexFileHandler indexFileHandler;
    private final Map<String, DeletionVector> deletionVectors;
    private boolean modified;

    private DeletionVectorsMaintainer(
            IndexFileHandler fileHandler,
            @Nullable Long snapshotId,
            BinaryRow partition,
            int bucket) {
        this.indexFileHandler = fileHandler;
        IndexFileMeta indexFile =
                snapshotId == null
                        ? null
                        : fileHandler
                                .scan(snapshotId, DELETION_VECTORS_INDEX, partition, bucket)
                                .orElse(null);
        this.deletionVectors =
                indexFile == null
                        ? new HashMap<>()
                        : new HashMap<>(indexFileHandler.readAllDeletionVectors(indexFile));
        this.modified = false;
    }

    /**
     * Notifies a new deletion which marks the specified row position as deleted with the given file
     * name.
     *
     * @param fileName The name of the file where the deletion occurred.
     * @param position The row position within the file that has been deleted.
     */
    public void notifyNewDeletion(String fileName, long position) {
        DeletionVector deletionVector =
                deletionVectors.computeIfAbsent(fileName, k -> new BitmapDeletionVector());
        if (deletionVector.checkedDelete(position)) {
            modified = true;
        }
    }

    /**
     * Removes the specified file's deletion vector, this method is typically used for remove before
     * files' deletion vector in compaction.
     *
     * @param fileName The name of the file whose deletion vector should be removed.
     */
    public void removeDeletionVectorOf(String fileName) {
        if (deletionVectors.containsKey(fileName)) {
            deletionVectors.remove(fileName);
            modified = true;
        }
    }

    /**
     * Prepares to commit: write new deletion vectors index file if any modifications have been
     * made.
     *
     * @return A list containing the metadata of the deletion vectors index file, or an empty list
     *     if no changes need to be committed.
     */
    public List<IndexFileMeta> prepareCommit() {
        if (modified) {
            IndexFileMeta entry = indexFileHandler.writeDeletionVectorsIndex(deletionVectors);
            modified = false;
            return Collections.singletonList(entry);
        }
        return Collections.emptyList();
    }

    /**
     * Retrieves the deletion vector associated with the specified file name.
     *
     * @param fileName The name of the file for which the deletion vector is requested.
     * @return An {@code Optional} containing the deletion vector if it exists, or an empty {@code
     *     Optional} if not.
     */
    public Optional<DeletionVector> deletionVectorOf(String fileName) {
        return Optional.ofNullable(deletionVectors.get(fileName));
    }

    @VisibleForTesting
    public Map<String, DeletionVector> deletionVectors() {
        return deletionVectors;
    }

    /** Factory to restore {@link DeletionVectorsMaintainer}. */
    public static class Factory {

        private final IndexFileHandler handler;

        public Factory(IndexFileHandler handler) {
            this.handler = handler;
        }

        public DeletionVectorsMaintainer createOrRestore(
                @Nullable Long snapshotId, BinaryRow partition, int bucket) {
            return new DeletionVectorsMaintainer(handler, snapshotId, partition, bucket);
        }
    }
}
