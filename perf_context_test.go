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

	SetPerfLevel(PerfLevelEnableCount)
	defer SetPerfLevel(PerfLevelDisable)

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
	ensure.True(t, pc.Metric(UserKeyComparisonCount) > 0)
	ensure.True(t, pc.Metric(GetReadBytes) > 0)

	// reset zeroes all counters
	pc.Reset()
	ensure.DeepEqual(t, pc.Metric(GetFromMemtableCount), uint64(0))
	ensure.DeepEqual(t, pc.Metric(UserKeyComparisonCount), uint64(0))
	ensure.DeepEqual(t, pc.Metric(GetReadBytes), uint64(0))
}

func TestPerfContextTimingMetrics(t *testing.T) {
	db := newTestDB(t, "TestPerfContextTimingMetrics", nil)
	defer db.Close()

	wo := NewDefaultWriteOptions()
	defer wo.Destroy()
	ro := NewDefaultReadOptions()
	defer ro.Destroy()

	SetPerfLevel(PerfLevelEnableTimeExceptForMutex)
	defer SetPerfLevel(PerfLevelDisable)

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
	ensure.True(t, pc.Metric(WriteWalTime) > 0)
	ensure.True(t, pc.Metric(WriteMemtableTime) > 0)
	ensure.True(t, pc.Metric(GetFromMemtableTime) > 0)

	// mutex timing metrics should remain zero at this perf level
	ensure.DeepEqual(t, pc.Metric(DbMutexLockNanos), uint64(0))
	ensure.DeepEqual(t, pc.Metric(DbConditionWaitNanos), uint64(0))

	// report should include timing field names
	report := pc.Report(true)
	ensure.True(t, strings.Contains(report, "write_wal_time"))
	ensure.True(t, strings.Contains(report, "get_from_memtable_time"))
}
