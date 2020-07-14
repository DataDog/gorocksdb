#include <stdlib.h>
#include "rocksdb/c.h"

// This API provides convenient C wrapper functions for rocksdb client.

/* Base */

extern void gorocksdb_destruct_handler(void* state);

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

/* Statistics/Tickers */
#ifdef __cplusplus__
extern "C" {
#endif

uint64_t gorocksdb_get_ticker_count(rocksdb_options_t *options, uint32_t ticker);
uint64_t gorocksdb_get_and_reset_ticker_count(rocksdb_options_t *options, uint32_t ticker);
void gorocksdb_set_stats_level(rocksdb_options_t *options, uint8_t level);
uint8_t gorocksdb_get_stats_level(rocksdb_options_t *options);

#ifdef __cplusplus__
}
#endif