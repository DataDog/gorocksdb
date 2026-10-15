package gorocksdb

// #include "gorocksdb.h"
import "C"

import (
	"errors"
	"strings"
	"unsafe"
)

// SSTFileReader reads a single SST file produced by RocksDB or SSTFileWriter,
// without opening a database. Its methods must not be called concurrently.
// Call Destroy when finished to release the file and native resources.
type SSTFileReader struct {
	c *C.gorocksdb_sstfilereader_t

	// Keep Go references to the options' dependencies alive.
	opts *Options
}

// NewSSTFileReader creates a reader using a copy of opts. Use options compatible
// with the SST's comparator, table format, and environment. The caller must not
// destroy opts or its dependencies until after the reader is destroyed.
func NewSSTFileReader(opts *Options) (*SSTFileReader, error) {
	if opts == nil || opts.c == nil {
		return nil, errors.New("rocksdb: SST reader requires options")
	}
	var cErr *C.char
	c := C.gorocksdb_sstfilereader_create(opts.c, &cErr)
	if c == nil {
		return nil, sstFileReaderError(cErr)
	}
	return &SSTFileReader{c: c, opts: opts}, nil
}

// Open opens the SST at path. A successful call releases any previously opened
// file; a failed call leaves the previous file, if any, open. If the file belongs
// to a running DB, the caller must prevent its deletion during use, for example
// with DB.DisableFileDeletions.
func (r *SSTFileReader) Open(path string) error {
	if r == nil || r.c == nil {
		return errors.New("rocksdb: SST reader is destroyed")
	}
	if strings.IndexByte(path, 0) >= 0 {
		return errors.New("rocksdb: SST path contains a NUL byte")
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	var cErr *C.char
	if C.gorocksdb_sstfilereader_open(r.c, cPath, &cErr) != 0 {
		return sstFileReaderError(cErr)
	}
	return nil
}

// VerifyChecksum verifies the block checksums in the opened SST. It reads the
// file and blocks until verification finishes. It does not verify a manifest's
// whole-file checksum, other SSTs, WALs, or referenced blob files. Open must have
// succeeded before calling VerifyChecksum.
func (r *SSTFileReader) VerifyChecksum() error {
	if r == nil || r.c == nil {
		return errors.New("rocksdb: SST reader is destroyed")
	}
	var cErr *C.char
	if C.gorocksdb_sstfilereader_verify_checksum(r.c, &cErr) != 0 {
		return sstFileReaderError(cErr)
	}
	return nil
}

// Destroy releases the reader and closes its file. It is safe to call more
// than once or on a nil receiver.
func (r *SSTFileReader) Destroy() {
	if r == nil {
		return
	}
	C.gorocksdb_sstfilereader_destroy(r.c)
	r.c = nil
	r.opts = nil
}

func sstFileReaderError(cErr *C.char) error {
	if cErr == nil {
		return errors.New("rocksdb: SST reader operation failed")
	}
	defer C.rocksdb_free(unsafe.Pointer(cErr))
	return errors.New(C.GoString(cErr))
}
