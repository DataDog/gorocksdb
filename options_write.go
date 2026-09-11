package gorocksdb

// #include "rocksdb/c.h"
import "C"

// WriteOptions represent all of the available options when writing to a
// database.
type WriteOptions struct {
	c *C.rocksdb_writeoptions_t
}

// NewDefaultWriteOptions creates a default WriteOptions object.
func NewDefaultWriteOptions() *WriteOptions {
	return NewNativeWriteOptions(C.rocksdb_writeoptions_create())
}

// NewNativeWriteOptions creates a WriteOptions object.
func NewNativeWriteOptions(c *C.rocksdb_writeoptions_t) *WriteOptions {
	return &WriteOptions{c}
}

// SetSync sets the sync mode. If true, the write will be flushed
// from the operating system buffer cache before the write is considered complete.
// If this flag is true, writes will be slower.
// Default: false
func (opts *WriteOptions) SetSync(value bool) {
	C.rocksdb_writeoptions_set_sync(opts.c, boolToChar(value))
}

// DisableWAL sets whether WAL should be active or not.
// If true, writes will not first go to the write ahead log,
// and the write may got lost after a crash.
// Default: false
func (opts *WriteOptions) DisableWAL(value bool) {
	C.rocksdb_writeoptions_disable_WAL(opts.c, C.int(btoi(value)))
}

// SetNoSlowdown sets whether a write that would have to wait is rejected
// instead of blocking.
//
// If true and the write would need to stall or be delayed — because too many
// memtables are pending flush, or because too many L0 files have accumulated —
// the write returns Status::Incomplete("Write stall") immediately rather than
// sleeping until there is room. See IsIncomplete for recognizing that status,
// which the C API only exposes as an error string.
//
// This is useful for a caller that has something better to do than occupy its
// thread waiting, such as dropping the write entirely. Note that it changes
// only the behavior under back pressure; it is not a hint and does not make
// writes faster when the DB is keeping up.
//
// Writes with different no_slowdown values are not batched into the same write
// group, so mixing values across concurrent writers fragments group commit.
//
// Default: false
// https://github.com/facebook/rocksdb/wiki/Write-Stalls
func (opts *WriteOptions) SetNoSlowdown(value bool) {
	C.rocksdb_writeoptions_set_no_slowdown(opts.c, boolToChar(value))
}

// SetLowPri sets whether this write is deprioritized.
// If true, a write to a DB that is already experiencing back pressure is
// stalled or rejected (per SetNoSlowdown) at a lower threshold than a normal
// write, so that low-priority traffic yields to the rest.
// Default: false
func (opts *WriteOptions) SetLowPri(value bool) {
	C.rocksdb_writeoptions_set_low_pri(opts.c, boolToChar(value))
}

// Destroy deallocates the WriteOptions object.
func (opts *WriteOptions) Destroy() {
	C.rocksdb_writeoptions_destroy(opts.c)
	opts.c = nil
}
