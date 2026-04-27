package gorocksdb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestPerfContext(t *testing.T) {
	db := newTestDB(t, "TestPerfContext", nil)
	defer db.Close()

	wo := NewDefaultWriteOptions()
	defer wo.Destroy()
	ro := NewDefaultReadOptions()
	defer ro.Destroy()

	SetPerfLevel(KEnableCount)
	defer SetPerfLevel(KDisable)

	pc := NewPerfContext()
	defer pc.Destroy()

	// write some keys
	for i := 0; i < 10; i++ {
		key := []byte(fmt.Sprintf("key-%02d", i))
		val := []byte(fmt.Sprintf("val-%02d", i))
		ensure.Nil(t, db.Put(wo, key, val))
	}

	// read them back
	for i := 0; i < 10; i++ {
		key := []byte(fmt.Sprintf("key-%02d", i))
		v, err := db.Get(ro, key)
		ensure.Nil(t, err)
		v.Free()
	}

	// iterate over all keys
	it := db.NewIterator(ro)
	defer it.Close()
	for it.SeekToFirst(); it.Valid(); it.Next() {
	}
	ensure.Nil(t, it.Err())

	// human-readable report should mention key comparison and memtable reads
	report := pc.Report(true)
	ensure.True(t, strings.Contains(report, "user_key_comparison_count"))
	ensure.True(t, strings.Contains(report, "get_from_memtable_count"))

	// reads and comparisons should have accumulated
	ensure.True(t, pc.Metric(0) > 0)  // rocksdb_user_key_comparison_count
	ensure.True(t, pc.Metric(7) > 0)  // rocksdb_get_read_bytes

	// reset zeroes all counters
	pc.Reset()
	ensure.DeepEqual(t, pc.Metric(16), uint64(0)) // rocksdb_get_from_memtable_count
	ensure.DeepEqual(t, pc.Metric(0), uint64(0))  // rocksdb_user_key_comparison_count
	ensure.DeepEqual(t, pc.Metric(7), uint64(0))  // rocksdb_get_read_bytes
}

func TestPerfContextTimingMetrics(t *testing.T) {
	db := newTestDB(t, "TestPerfContextTimingMetrics", nil)
	defer db.Close()

	wo := NewDefaultWriteOptions()
	defer wo.Destroy()
	ro := NewDefaultReadOptions()
	defer ro.Destroy()

	SetPerfLevel(KEnableTimeExceptForMutex)
	defer SetPerfLevel(KDisable)

	pc := NewPerfContext()
	defer pc.Destroy()

	for i := 0; i < 10; i++ {
		key := []byte(fmt.Sprintf("key-%02d", i))
		val := []byte(fmt.Sprintf("val-%02d", i))
		ensure.Nil(t, db.Put(wo, key, val))
	}

	for i := 0; i < 10; i++ {
		key := []byte(fmt.Sprintf("key-%02d", i))
		v, err := db.Get(ro, key)
		ensure.Nil(t, err)
		v.Free()
	}

	// timing metrics should be populated at this perf level
	ensure.True(t, pc.Metric(29) > 0) // rocksdb_write_wal_time
	ensure.True(t, pc.Metric(30) > 0) // rocksdb_write_memtable_time
	ensure.True(t, pc.Metric(15) > 0) // rocksdb_get_from_memtable_time

	// mutex timing metrics should remain zero at this perf level
	ensure.DeepEqual(t, pc.Metric(33), uint64(0)) // rocksdb_db_mutex_lock_nanos
	ensure.DeepEqual(t, pc.Metric(34), uint64(0)) // rocksdb_db_condition_wait_nanos

	// report should include timing field names
	report := pc.Report(true)
	ensure.True(t, strings.Contains(report, "write_wal_time"))
	ensure.True(t, strings.Contains(report, "get_from_memtable_time"))
}
