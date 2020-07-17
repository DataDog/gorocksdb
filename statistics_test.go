package gorocksdb

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

type stringMerge struct {
}

func (sm stringMerge) FullMerge(key, existingValue []byte, operands [][]byte) ([]byte, bool) {
	ret := []byte{}
	if existingValue != nil {
		ret = append(existingValue, ',')
		if len(operands) >= 1 {
			ret = append(ret, operands[0]...)
		}
	} else {
		if len(operands) >= 1 {
			ret = operands[0]
		}
	}

	for _, op := range operands[1:] {
		ret = append(ret, ',')
		ret = append(ret, op...)
	}
	return ret, true
}

func (sm stringMerge) Name() string {
	return "StringMerge"
}

func TestStatistics(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	dir, err := ioutil.TempDir("/tmp", "gorocksdb-TestStatistics")
	assert.NoError(err, "Could not create temp directory")

	opts := NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	mergeOpt := stringMerge{}
	opts.SetMergeOperator(mergeOpt)
	opts.EnableStatistics()
	opts.SetStatsLevel(StatsExceptDetailedTimers)

	cache := NewLRUCache(1000)
	bbto := NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(cache)
	bbto.SetFlushEveryKeyPolicy()
	opts.SetBlockBasedTableFactory(bbto)
	bbto.SetIndexType(KBinarySearchIndexType)

	db, err := OpenDb(opts, dir)
	assert.NoError(err, "Unable to open database")
	assert.NotNil(db, "Unable to open database")

	ro := NewDefaultReadOptions()
	ro.SetFillCache(true)

	wo := NewDefaultWriteOptions()
	fo := NewDefaultFlushOptions()

	assert.NoError(db.Merge(wo, []byte("a1"), []byte("x1")), "Unable to Merge a1/x1")
	assert.NoError(db.Merge(wo, []byte("b1"), []byte("y1")), "Unable to Merge b1/y1")
	assert.NoError(db.Merge(wo, []byte("c1"), []byte("z1")), "Unable to Merge c1/z1")
	assert.NoError(db.Flush(fo), "Unable to db.Flush()")

	assert.NoError(db.Merge(wo, []byte("a2"), []byte("x2")), "Unable to Merge a2/x2")
	assert.NoError(db.Merge(wo, []byte("b2"), []byte("y2")), "Unable to Merge b2/y2")
	assert.NoError(db.Merge(wo, []byte("c2"), []byte("z2")), "Unable to Merge c2/z2")
	assert.NoError(db.Flush(fo), "Unable to db.Flush()")

	assert.NoError(db.Merge(wo, []byte("a3"), []byte("x3")), "Unable to Merge a3/x3")
	assert.NoError(db.Merge(wo, []byte("b3"), []byte("y3")), "Unable to Merge b3/y3")
	assert.NoError(db.Merge(wo, []byte("c3"), []byte("z3")), "Unable to Merge c3/z3")
	assert.NoError(db.Flush(fo), "Unable to db.Flush()")

	it := db.NewIterator(ro)
	assert.NotNil(it, "Unable to create iterator")

	opts.SetTickerCount(BlockCacheMiss, 0)
	opts.SetTickerCount(BlockCacheHit, 0)

	it.Seek([]byte("b2"))
	assert.True(it.Valid(), "Failed iterator seek")
	assert.Equal("b2", string(it.Key().Data()), "Incorrect key found")
	assert.Equal("y2", string(it.Value().Data()), "Incorrect value found")
	assert.Equal(uint64(5), opts.GetTickerCount(BlockCacheMiss), "Incorrect block miss count")
	assert.Equal(uint64(0), opts.GetTickerCount(BlockCacheHit), "Incorrect block hit count")
	assert.Equal(uint64(0), opts.GetTickerCount(BlockCacheDataMiss), "Incorrect block data miss count")
	assert.Equal(uint64(0), opts.GetTickerCount(BlockCacheDataHit), "Incorrect block data hit count")

	it.Seek([]byte("a2"))
	assert.True(it.Valid(), "Failed iterator seek")
	assert.Equal("a2", string(it.Key().Data()), "Incorrect key found")
	assert.Equal("x2", string(it.Value().Data()), "Incorrect value found")
	assert.Equal(uint64(8), opts.GetTickerCount(BlockCacheMiss), "Incorrect block miss count")
	assert.Equal(uint64(2), opts.GetTickerCount(BlockCacheHit), "Incorrect block hit count")
	assert.Equal(uint64(0), opts.GetTickerCount(BlockCacheDataMiss), "Incorrect block data miss count")
	assert.Equal(uint64(0), opts.GetTickerCount(BlockCacheDataHit), "Incorrect block data hit count")

	it.Seek([]byte("c2"))
	assert.True(it.Valid(), "Failed iterator seek")
	assert.Equal("c2", string(it.Key().Data()), "Incorrect key found")
	assert.Equal("z2", string(it.Value().Data()), "Incorrect value found")
	assert.Equal(uint64(9), opts.GetTickerCount(BlockCacheMiss), "Incorrect block miss count")
	assert.Equal(uint64(3), opts.GetTickerCount(BlockCacheHit), "Incorrect block hit count")

	it.Seek([]byte("a1"))
	assert.True(it.Valid(), "Failed iterator seek")
	assert.Equal("a1", string(it.Key().Data()), "Incorrect key found")
	assert.Equal("x1", string(it.Value().Data()), "Incorrect value found")
	assert.Equal(uint64(9), opts.GetTickerCount(BlockCacheMiss), "Incorrect block miss count")
	assert.Equal(uint64(7), opts.GetTickerCount(BlockCacheHit), "Incorrect block hit count")

	it.Seek([]byte("a3"))
	assert.True(it.Valid(), "Failed iterator seek")
	assert.Equal("a3", string(it.Key().Data()), "Incorrect key found")
	assert.Equal("x3", string(it.Value().Data()), "Incorrect value found")
	assert.Equal(uint64(9), opts.GetTickerCount(BlockCacheMiss), "Incorrect block miss count")
	assert.Equal(uint64(11), opts.GetTickerCount(BlockCacheHit), "Incorrect block hit count")

	misses := opts.GetAndResetTickerCount(BlockCacheMiss)
	hits := opts.GetAndResetTickerCount(BlockCacheHit)

	assert.Equal(uint64(9), misses, "Incorrect block miss count")
	assert.Equal(uint64(11), hits, "Incorrect block hit count")

	assert.Equal(uint64(0), opts.GetTickerCount(BlockCacheMiss), "Incorrect block miss count")
	assert.Equal(uint64(0), opts.GetTickerCount(BlockCacheHit), "Incorrect block hit count")
}
