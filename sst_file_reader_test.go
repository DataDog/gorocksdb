package gorocksdb

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestSST(t *testing.T, opts *Options, keys ...string) string {
	t.Helper()
	envOpts := NewDefaultEnvOptions()
	defer envOpts.Destroy()
	writer := NewSSTFileWriter(envOpts, opts)
	defer writer.Destroy()
	path := filepath.Join(t.TempDir(), "test.sst")
	if err := writer.Open(path); err != nil {
		t.Fatal(err)
	}
	for _, key := range keys {
		if err := writer.Add([]byte(key), []byte("value-"+key)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Finish(); err != nil {
		t.Fatal(err)
	}
	return path
}

func newTestSSTReader(t *testing.T, opts *Options) *SSTFileReader {
	t.Helper()
	reader, err := NewSSTFileReader(opts)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(reader.Destroy)
	return reader
}

func TestSSTFileReaderVerifyChecksum(t *testing.T) {
	for _, customComparator := range []bool{false, true} {
		name := "default"
		if customComparator {
			name = "custom-comparator"
		}
		t.Run(name, func(t *testing.T) {
			opts := NewDefaultOptions()
			t.Cleanup(opts.Destroy)
			keys := []string{"aaa", "bbb", "ccc"}
			if customComparator {
				opts.SetComparator(&bytesReverseComparator{})
				keys = []string{"ccc", "bbb", "aaa"}
			}
			path := writeTestSST(t, opts, keys...)
			reader := newTestSSTReader(t, opts)
			if err := reader.Open(path); err != nil {
				t.Fatal(err)
			}
			// Repeated verification must remain valid.
			for i := 0; i < 2; i++ {
				if err := reader.VerifyChecksum(); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestSSTFileReaderDBGeneratedFile(t *testing.T) {
	opts := NewDefaultOptions()
	t.Cleanup(opts.Destroy)
	opts.SetCreateIfMissing(true)
	dir := t.TempDir()
	db, err := OpenDb(opts, dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if db != nil {
			db.Close()
		}
	})
	wo := NewDefaultWriteOptions()
	t.Cleanup(wo.Destroy)
	fo := NewDefaultFlushOptions()
	t.Cleanup(fo.Destroy)
	fo.SetWait(true)
	if err := db.Put(wo, []byte("key"), []byte("value")); err != nil {
		t.Fatal(err)
	}
	if err := db.Flush(fo); err != nil {
		t.Fatal(err)
	}
	files := db.GetLiveFilesMetaData()
	// The reader must work without a live DB handle.
	db.Close()
	db = nil
	if len(files) == 0 {
		t.Fatal("flush produced no SST files")
	}
	reader := newTestSSTReader(t, opts)
	if err := reader.Open(filepath.Join(dir, files[0].Name)); err != nil {
		t.Fatal(err)
	}
	if err := reader.VerifyChecksum(); err != nil {
		t.Fatal(err)
	}
}

func TestSSTFileReaderDetectsCorruption(t *testing.T) {
	opts := NewDefaultOptions()
	t.Cleanup(opts.Destroy)
	path := writeTestSST(t, opts, "aaa", "bbb")
	// Change a data-block byte, leaving the footer and metadata intact.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	cleanPath := filepath.Join(filepath.Dir(path), "clean.sst")
	if err := os.WriteFile(cleanPath, data, 0600); err != nil {
		t.Fatal(err)
	}
	data[0] ^= 0xff
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	reader := newTestSSTReader(t, opts)
	if err := reader.Open(path); err != nil {
		t.Fatal(err)
	}
	if err := reader.VerifyChecksum(); err == nil || !strings.Contains(err.Error(), "Corruption") {
		t.Fatalf("expected checksum corruption, got %v", err)
	}
	// Verification is limited to the selected file, not other SSTs beside it.
	if err := reader.Open(cleanPath); err != nil {
		t.Fatal(err)
	}
	if err := reader.VerifyChecksum(); err != nil {
		t.Fatalf("corrupt neighbor affected a clean SST: %v", err)
	}
}

func TestSSTFileReaderOpenErrors(t *testing.T) {
	opts := NewDefaultOptions()
	t.Cleanup(opts.Destroy)
	validPath := writeTestSST(t, opts, "key")
	data, err := os.ReadFile(validPath)
	if err != nil {
		t.Fatal(err)
	}
	for name, contents := range map[string][]byte{
		"empty":     {},
		"invalid":   []byte("not an SST file"),
		"truncated": data[:len(data)/2],
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "invalid.sst")
			if err := os.WriteFile(path, contents, 0600); err != nil {
				t.Fatal(err)
			}
			reader := newTestSSTReader(t, opts)
			if err := reader.Open(path); err == nil {
				t.Fatal("expected Open to fail")
			}
			if err := reader.VerifyChecksum(); err == nil {
				t.Fatal("expected verification without an open file to fail")
			}
		})
	}
	reader := newTestSSTReader(t, opts)
	if err := reader.Open(filepath.Join(t.TempDir(), "missing.sst")); err == nil {
		t.Fatal("expected missing-file error")
	}
	if err := reader.Open(validPath + "\x00suffix"); err == nil {
		t.Fatal("expected NUL-path error")
	}
	// A failed Open must not prevent a later successful Open.
	if err := reader.Open(validPath); err != nil {
		t.Fatal(err)
	}
	if err := reader.Open(filepath.Join(t.TempDir(), "missing.sst")); err == nil {
		t.Fatal("expected missing-file error")
	}
	if err := reader.VerifyChecksum(); err != nil {
		t.Fatalf("failed Open lost the previous file: %v", err)
	}
	if err := reader.Open(writeTestSST(t, opts, "another")); err != nil {
		t.Fatal(err)
	}
	if err := reader.VerifyChecksum(); err != nil {
		t.Fatal(err)
	}
}

func TestSSTFileReaderLifecycle(t *testing.T) {
	if reader, err := NewSSTFileReader(nil); err == nil || reader != nil {
		t.Fatal("expected nil options to fail")
	}
	opts := NewDefaultOptions()
	defer opts.Destroy()
	reader := newTestSSTReader(t, opts)
	if err := reader.VerifyChecksum(); err == nil {
		t.Fatal("expected verification before Open to fail")
	}
	reader.Destroy()
	reader.Destroy()
	for _, r := range []*SSTFileReader{reader, nil} {
		if err := r.Open("unused.sst"); err == nil {
			t.Fatal("expected Open on destroyed/nil reader to fail")
		}
		if err := r.VerifyChecksum(); err == nil {
			t.Fatal("expected verification on destroyed/nil reader to fail")
		}
		r.Destroy()
	}
}
