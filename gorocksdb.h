#ifndef GOROCKSDB_H
#define GOROCKSDB_H

#include <stdlib.h>
#include "rocksdb/c.h"

#ifdef __cplusplus
extern "C" {
#endif

// This API provides convenient C wrapper functions for rocksdb client.

/* Base */

extern void gorocksdb_destruct_handler(void* state);
/*
 * Verifies block checksums in all live SST files.
 *
 * Returns 0 on success and nonzero on failure. On failure, *errptr is set
 * to a malloc-allocated message when allocation succeeds. Any existing
 * *errptr is freed first, matching RocksDB's C API convention.
 */
extern int gorocksdb_verify_checksum(rocksdb_t* db, char** errptr);

/* CompactionFilter */

extern rocksdb_compactionfilter_t* gorocksdb_compactionfilter_create(uintptr_t idx);

/* Comparator */

extern rocksdb_comparator_t* gorocksdb_comparator_create(uintptr_t idx);

/* Filter Policy */

extern rocksdb_filterpolicy_t* gorocksdb_filterpolicy_create(uintptr_t idx);
extern void gorocksdb_filterpolicy_delete_filter(void* state, const char* v, size_t s);

/* Merge Operator */

extern rocksdb_mergeoperator_t* gorocksdb_mergeoperator_create(uintptr_t idx);
extern void gorocksdb_mergeoperator_delete_value(void* state, const char* v, size_t s);

/* Slice Transform */

extern rocksdb_slicetransform_t* gorocksdb_slicetransform_create(uintptr_t idx);

#ifdef __cplusplus
}  // extern "C"
#endif

#endif  /* GOROCKSDB_H */
