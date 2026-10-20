package gorocksdb

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestColumnFamilyOpen(t *testing.T) {
	dir, err := ioutil.TempDir("", "gorocksdb-TestColumnFamilyOpen")
	ensure.Nil(t, err)

	givenNames := []string{"default", "guide"}
	opts := NewDefaultOptions()
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetCreateIfMissing(true)
	db, cfh, err := OpenDbColumnFamilies(opts, dir, givenNames, []*Options{opts, opts})
	ensure.Nil(t, err)
	defer db.Close()
	ensure.DeepEqual(t, len(cfh), 2)
	cfh[0].Destroy()
	cfh[1].Destroy()

	actualNames, err := ListColumnFamilies(opts, dir)
	ensure.Nil(t, err)
	ensure.SameElements(t, actualNames, givenNames)
}

func TestColumnFamilyCreateDrop(t *testing.T) {
	dir, err := ioutil.TempDir("", "gorocksdb-TestColumnFamilyCreate")
	ensure.Nil(t, err)

	opts := NewDefaultOptions()
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetCreateIfMissing(true)
	db, err := OpenDb(opts, dir)
	ensure.Nil(t, err)
	defer db.Close()
	cf, err := db.CreateColumnFamily(opts, "guide")
	ensure.Nil(t, err)
	defer cf.Destroy()

	actualNames, err := ListColumnFamilies(opts, dir)
	ensure.Nil(t, err)
	ensure.SameElements(t, actualNames, []string{"default", "guide"})

	ensure.Nil(t, db.DropColumnFamily(cf))

	actualNames, err = ListColumnFamilies(opts, dir)
	ensure.Nil(t, err)
	ensure.SameElements(t, actualNames, []string{"default"})
}

func TestColumnFamilyBatchPutGet(t *testing.T) {
	dir, err := ioutil.TempDir("", "gorocksdb-TestColumnFamilyPutGet")
	ensure.Nil(t, err)

	givenNames := []string{"default", "guide"}
	opts := NewDefaultOptions()
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetCreateIfMissing(true)
	db, cfh, err := OpenDbColumnFamilies(opts, dir, givenNames, []*Options{opts, opts})
	ensure.Nil(t, err)
	defer db.Close()
	ensure.DeepEqual(t, len(cfh), 2)
	defer cfh[0].Destroy()
	defer cfh[1].Destroy()

	wo := NewDefaultWriteOptions()
	defer wo.Destroy()
	ro := NewDefaultReadOptions()
	defer ro.Destroy()

	givenKey0 := []byte("hello0")
	givenVal0 := []byte("world0")
	givenKey1 := []byte("hello1")
	givenVal1 := []byte("world1")

	b0 := NewWriteBatch()
	defer b0.Destroy()
	b0.PutCF(cfh[0], givenKey0, givenVal0)
	ensure.Nil(t, db.Write(wo, b0))
	actualVal0, err := db.GetCF(ro, cfh[0], givenKey0)
	defer actualVal0.Free()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, actualVal0.Data(), givenVal0)

	b1 := NewWriteBatch()
	defer b1.Destroy()
	b1.PutCF(cfh[1], givenKey1, givenVal1)
	ensure.Nil(t, db.Write(wo, b1))
	actualVal1, err := db.GetCF(ro, cfh[1], givenKey1)
	defer actualVal1.Free()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, actualVal1.Data(), givenVal1)

	actualVal, err := db.GetCF(ro, cfh[0], givenKey1)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, actualVal.Size(), 0)
	actualVal, err = db.GetCF(ro, cfh[1], givenKey0)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, actualVal.Size(), 0)
}

func TestColumnFamilyPutGetDelete(t *testing.T) {
	dir, err := ioutil.TempDir("", "gorocksdb-TestColumnFamilyPutGet")
	ensure.Nil(t, err)

	givenNames := []string{"default", "guide"}
	opts := NewDefaultOptions()
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetCreateIfMissing(true)
	db, cfh, err := OpenDbColumnFamilies(opts, dir, givenNames, []*Options{opts, opts})
	ensure.Nil(t, err)
	defer db.Close()
	ensure.DeepEqual(t, len(cfh), 2)
	defer cfh[0].Destroy()
	defer cfh[1].Destroy()

	wo := NewDefaultWriteOptions()
	defer wo.Destroy()
	ro := NewDefaultReadOptions()
	defer ro.Destroy()

	givenKey0 := []byte("hello0")
	givenVal0 := []byte("world0")
	givenKey1 := []byte("hello1")
	givenVal1 := []byte("world1")

	ensure.Nil(t, db.PutCF(wo, cfh[0], givenKey0, givenVal0))
	actualVal0, err := db.GetCF(ro, cfh[0], givenKey0)
	defer actualVal0.Free()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, actualVal0.Data(), givenVal0)

	ensure.Nil(t, db.PutCF(wo, cfh[1], givenKey1, givenVal1))
	actualVal1, err := db.GetCF(ro, cfh[1], givenKey1)
	defer actualVal1.Free()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, actualVal1.Data(), givenVal1)

	actualVal, err := db.GetCF(ro, cfh[0], givenKey1)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, actualVal.Size(), 0)
	actualVal, err = db.GetCF(ro, cfh[1], givenKey0)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, actualVal.Size(), 0)

	ensure.Nil(t, db.DeleteCF(wo, cfh[0], givenKey0))
	actualVal, err = db.GetCF(ro, cfh[0], givenKey0)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, actualVal.Size(), 0)
}

func newTestDBCF(t *testing.T, name string) (db *DB, cfh []*ColumnFamilyHandle, cleanup func()) {
	dir, err := ioutil.TempDir("", "gorocksdb-TestColumnFamilyPutGet")
	ensure.Nil(t, err)

	givenNames := []string{"default", "guide"}
	opts := NewDefaultOptions()
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetCreateIfMissing(true)
	db, cfh, err = OpenDbColumnFamilies(opts, dir, givenNames, []*Options{opts, opts})
	ensure.Nil(t, err)
	cleanup = func() {
		for _, cf := range cfh {
			cf.Destroy()
		}
		db.Close()
	}
	return db, cfh, cleanup
}

func TestColumnFamilyMultiGet(t *testing.T) {
	db, cfh, cleanup := newTestDBCF(t, "TestDBMultiGet")
	defer cleanup()

	var (
		givenKey1 = []byte("hello1")
		givenKey2 = []byte("hello2")
		givenKey3 = []byte("hello3")
		givenVal1 = []byte("world1")
		givenVal2 = []byte("world2")
		givenVal3 = []byte("world3")
		wo        = NewDefaultWriteOptions()
		ro        = NewDefaultReadOptions()
	)

	// create
	ensure.Nil(t, db.PutCF(wo, cfh[0], givenKey1, givenVal1))
	ensure.Nil(t, db.PutCF(wo, cfh[1], givenKey2, givenVal2))
	ensure.Nil(t, db.PutCF(wo, cfh[1], givenKey3, givenVal3))

	// column family 0 only has givenKey1
	values, err := db.MultiGetCF(ro, cfh[0], []byte("noexist"), givenKey1, givenKey2, givenKey3)
	defer values.Destroy()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(values), 4)

	ensure.DeepEqual(t, values[0].Data(), []byte(nil))
	ensure.DeepEqual(t, values[1].Data(), givenVal1)
	ensure.DeepEqual(t, values[2].Data(), []byte(nil))
	ensure.DeepEqual(t, values[3].Data(), []byte(nil))

	// column family 1 only has givenKey2 and givenKey3
	values, err = db.MultiGetCF(ro, cfh[1], []byte("noexist"), givenKey1, givenKey2, givenKey3)
	defer values.Destroy()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(values), 4)

	ensure.DeepEqual(t, values[0].Data(), []byte(nil))
	ensure.DeepEqual(t, values[1].Data(), []byte(nil))
	ensure.DeepEqual(t, values[2].Data(), givenVal2)
	ensure.DeepEqual(t, values[3].Data(), givenVal3)

	// getting them all from the right CF should return them all
	values, err = db.MultiGetCFMultiCF(ro,
		ColumnFamilyHandles{cfh[0], cfh[1], cfh[1]},
		[][]byte{givenKey1, givenKey2, givenKey3},
	)
	defer values.Destroy()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(values), 3)

	ensure.DeepEqual(t, values[0].Data(), givenVal1)
	ensure.DeepEqual(t, values[1].Data(), givenVal2)
	ensure.DeepEqual(t, values[2].Data(), givenVal3)
}

func newTestDBCFPinned(t *testing.T, noBlockCache bool) (db *DB, cfh []*ColumnFamilyHandle, dir string, cleanup func()) {
	dir, err := ioutil.TempDir("", "gorocksdb-TestColumnFamilyGetPinned")
	ensure.Nil(t, err)

	bbto := NewDefaultBlockBasedTableOptions()
	bbto.SetNoBlockCache(noBlockCache)
	opts := NewDefaultOptions()
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetCreateIfMissing(true)
	opts.SetCompression(NoCompression)
	opts.SetBlockBasedTableFactory(bbto)
	db, cfh, err = OpenDbColumnFamilies(opts, dir, []string{"default", "guide"}, []*Options{opts, opts})
	ensure.Nil(t, err)
	cleanup = func() {
		for _, cf := range cfh {
			cf.Destroy()
		}
		db.Close()
		opts.Destroy()
		bbto.Destroy()
		os.RemoveAll(dir)
	}
	return db, cfh, dir, cleanup
}

func flushCF(t *testing.T, db *DB, cf *ColumnFamilyHandle) {
	fo := NewDefaultFlushOptions()
	defer fo.Destroy()
	fo.SetWait(true)
	ensure.Nil(t, db.FlushCF(cf, fo))
}

// ensureGetPinnedCFMatchesGetCF checks that GetPinnedCF returns the same value and the same
// found / not-found outcome as GetCF.
func ensureGetPinnedCFMatchesGetCF(t *testing.T, db *DB, ro *ReadOptions, cf *ColumnFamilyHandle, key []byte) {
	want, err := db.GetCF(ro, cf, key)
	ensure.Nil(t, err)
	defer want.Free()

	got, err := db.GetPinnedCF(ro, cf, key)
	ensure.Nil(t, err)
	defer got.Destroy()

	ensure.DeepEqual(t, got.Exists(), want.Exists())
	ensure.True(t, bytes.Equal(got.Data(), want.Data()))
	ensure.DeepEqual(t, len(got.Data()), want.Size())
}

func TestColumnFamilyGetPinned(t *testing.T) {
	givenVals := map[string][]byte{
		"small": []byte("world"),
		"empty": {},
		// Larger than a block, and large enough to be heap allocated on every Get.
		"large": bytes.Repeat([]byte("0123456789abcdef"), 2<<20),
		"":      []byte("value-for-empty-key"),
	}

	for _, tc := range []struct {
		name         string
		noBlockCache bool
		flush        bool
	}{
		// Memtable hits are always copied into the PinnableSlice.
		{name: "memtable", noBlockCache: true},
		// Without a block cache, SST values are copied into the PinnableSlice.
		{name: "sst-no-block-cache", noBlockCache: true, flush: true},
		// With a block cache, SST values are pinned to the cached block.
		{name: "sst-block-cache", noBlockCache: false, flush: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db, cfh, _, cleanup := newTestDBCFPinned(t, tc.noBlockCache)
			defer cleanup()

			wo := NewDefaultWriteOptions()
			defer wo.Destroy()
			ro := NewDefaultReadOptions()
			defer ro.Destroy()

			for k, v := range givenVals {
				ensure.Nil(t, db.PutCF(wo, cfh[1], []byte(k), v))
			}
			if tc.flush {
				flushCF(t, db, cfh[1])
			}

			for k := range givenVals {
				ensureGetPinnedCFMatchesGetCF(t, db, ro, cfh[1], []byte(k))
				// The keys only exist in cfh[1].
				ensureGetPinnedCFMatchesGetCF(t, db, ro, cfh[0], []byte(k))
			}
			ensureGetPinnedCFMatchesGetCF(t, db, ro, cfh[1], []byte("noexist"))

			// An empty value is found, it is not a missing key.
			v, err := db.GetPinnedCF(ro, cfh[1], []byte("empty"))
			ensure.Nil(t, err)
			ensure.True(t, v.Exists())
			ensure.DeepEqual(t, len(v.Data()), 0)
			v.Destroy()

			v, err = db.GetPinnedCF(ro, cfh[1], []byte("noexist"))
			ensure.Nil(t, err)
			ensure.False(t, v.Exists())
			ensure.True(t, v.Data() == nil)
			v.Destroy()

			// A deleted key is missing.
			ensure.Nil(t, db.DeleteCF(wo, cfh[1], []byte("small")))
			v, err = db.GetPinnedCF(ro, cfh[1], []byte("small"))
			ensure.Nil(t, err)
			ensure.False(t, v.Exists())
			v.Destroy()
		})
	}
}

func TestColumnFamilyGetPinnedOutlivesOverwrite(t *testing.T) {
	for _, noBlockCache := range []bool{true, false} {
		db, cfh, _, cleanup := newTestDBCFPinned(t, noBlockCache)
		defer cleanup()

		wo := NewDefaultWriteOptions()
		defer wo.Destroy()
		ro := NewDefaultReadOptions()
		defer ro.Destroy()

		givenKey := []byte("hello")
		ensure.Nil(t, db.PutCF(wo, cfh[1], givenKey, []byte("world1")))
		flushCF(t, db, cfh[1])

		v, err := db.GetPinnedCF(ro, cfh[1], givenKey)
		ensure.Nil(t, err)
		defer v.Destroy()

		// A value stays valid and unchanged while it is held, regardless of later writes and flushes.
		ensure.Nil(t, db.PutCF(wo, cfh[1], givenKey, []byte("world2")))
		flushCF(t, db, cfh[1])
		ensure.DeepEqual(t, v.Data(), []byte("world1"))
		ensureGetPinnedCFMatchesGetCF(t, db, ro, cfh[1], givenKey)
	}
}

func TestColumnFamilyGetPinnedError(t *testing.T) {
	db, cfh, dir, cleanup := newTestDBCFPinned(t, true)
	defer cleanup()

	wo := NewDefaultWriteOptions()
	defer wo.Destroy()
	ro := NewDefaultReadOptions()
	defer ro.Destroy()

	givenKey := []byte("hello")
	ensure.Nil(t, db.PutCF(wo, cfh[1], givenKey, bytes.Repeat([]byte("a"), 1<<20)))
	flushCF(t, db, cfh[1])

	// Corrupt the value's data block. Without a block cache every Get re-reads the block from disk,
	// so the checksum mismatch must surface as an error, not as a missing key or bad data.
	ssts, err := filepath.Glob(filepath.Join(dir, "*.sst"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(ssts), 1)
	f, err := os.OpenFile(ssts[0], os.O_RDWR, 0)
	ensure.Nil(t, err)
	_, err = f.WriteAt([]byte("b"), 1<<19)
	ensure.Nil(t, err)
	ensure.Nil(t, f.Close())

	_, wantErr := db.GetCF(ro, cfh[1], givenKey)
	ensure.NotNil(t, wantErr)

	v, err := db.GetPinnedCF(ro, cfh[1], givenKey)
	ensure.True(t, v == nil)
	ensure.NotNil(t, err)
	ensure.DeepEqual(t, err.Error(), wantErr.Error())
}
