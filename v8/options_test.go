package gorocksdb

import (
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	opts := NewDefaultOptions()
	defer opts.Destroy()

	// Test setting max bg jobs
	assert.Equal(t, 2, opts.GetMaxBackgroundJobs())
	opts.SetMaxBackgroundJobs(10)
	assert.Equal(t, 10, opts.GetMaxBackgroundJobs())

	// Test setting max bg compactions
	assert.Equal(t, uint32(1), opts.GetMaxSubcompactions())
	opts.SetMaxSubcompactions(9)
	assert.Equal(t, uint32(9), opts.GetMaxSubcompactions())
}

func TestReadOptionsIterateBounds(t *testing.T) {
	db := newTestDB(t, "TestReadOptionsIterateBounds", nil)
	defer db.Close()

	// Insert keys
	givenKeys := [][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
		[]byte("e"),
	}
	wo := NewDefaultWriteOptions()
	defer wo.Destroy()
	for _, k := range givenKeys {
		ensure.Nil(t, db.Put(wo, k, []byte("val")))
	}

	t.Run("UpperBound", func(t *testing.T) {
		ro := NewDefaultReadOptions()
		defer ro.Destroy()

		// Set upper bound to "d" - should only see a, b, c
		ro.SetIterateUpperBound([]byte("d"))

		iter := db.NewIterator(ro)
		defer iter.Close()

		var actualKeys [][]byte
		for iter.SeekToFirst(); iter.Valid(); iter.Next() {
			key := make([]byte, len(iter.Key().Data()))
			copy(key, iter.Key().Data())
			actualKeys = append(actualKeys, key)
		}
		ensure.Nil(t, iter.Err())
		ensure.DeepEqual(t, actualKeys, [][]byte{[]byte("a"), []byte("b"), []byte("c")})
	})

	t.Run("LowerBound", func(t *testing.T) {
		ro := NewDefaultReadOptions()
		defer ro.Destroy()

		// Set lower bound to "c" - should only see c, d, e
		ro.SetIterateLowerBound([]byte("c"))

		iter := db.NewIterator(ro)
		defer iter.Close()

		var actualKeys [][]byte
		for iter.SeekToFirst(); iter.Valid(); iter.Next() {
			key := make([]byte, len(iter.Key().Data()))
			copy(key, iter.Key().Data())
			actualKeys = append(actualKeys, key)
		}
		ensure.Nil(t, iter.Err())
		ensure.DeepEqual(t, actualKeys, [][]byte{[]byte("c"), []byte("d"), []byte("e")})
	})

	t.Run("BothBounds", func(t *testing.T) {
		ro := NewDefaultReadOptions()
		defer ro.Destroy()

		// Set lower bound to "b" and upper bound to "d" - should only see b, c
		ro.SetIterateLowerBound([]byte("b"))
		ro.SetIterateUpperBound([]byte("d"))

		iter := db.NewIterator(ro)
		defer iter.Close()

		var actualKeys [][]byte
		for iter.SeekToFirst(); iter.Valid(); iter.Next() {
			key := make([]byte, len(iter.Key().Data()))
			copy(key, iter.Key().Data())
			actualKeys = append(actualKeys, key)
		}
		ensure.Nil(t, iter.Err())
		ensure.DeepEqual(t, actualKeys, [][]byte{[]byte("b"), []byte("c")})
	})

	t.Run("MultipleSets", func(t *testing.T) {
		ro := NewDefaultReadOptions()
		defer ro.Destroy()

		// Set bounds multiple times to ensure memory is properly managed
		ro.SetIterateUpperBound([]byte("e"))
		ro.SetIterateLowerBound([]byte("a"))

		// Change them
		ro.SetIterateUpperBound([]byte("d"))
		ro.SetIterateLowerBound([]byte("b"))

		// Change them again
		ro.SetIterateUpperBound([]byte("c"))
		ro.SetIterateLowerBound([]byte("b"))

		iter := db.NewIterator(ro)
		defer iter.Close()

		var actualKeys [][]byte
		for iter.SeekToFirst(); iter.Valid(); iter.Next() {
			key := make([]byte, len(iter.Key().Data()))
			copy(key, iter.Key().Data())
			actualKeys = append(actualKeys, key)
		}
		ensure.Nil(t, iter.Err())
		ensure.DeepEqual(t, actualKeys, [][]byte{[]byte("b")})
	})
}
