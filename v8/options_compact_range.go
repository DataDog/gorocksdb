package gorocksdb

// #include "rocksdb/c.h"
import "C"

// BottommostLevelCompaction for level based compaction, we can configure if we want to skip/force
// bottommost level compaction.
type BottommostLevelCompaction byte

const (
	// KSkip skip bottommost level compaction
	KSkip BottommostLevelCompaction = 0
	// KIfHaveCompactionFilter only compact bottommost level if there is a compaction filter
	// This is the default option
	KIfHaveCompactionFilter BottommostLevelCompaction = 1
	// KForce always compact bottommost level
	KForce BottommostLevelCompaction = 2
	// KForceOptimized always compact bottommost level but in bottommost level avoid
	// double-compacting files created in the same compaction
	KForceOptimized BottommostLevelCompaction = 3
)

// CompactRangeOptions represent all of the available options for compact range.
type CompactRangeOptions struct {
	c *C.rocksdb_compactoptions_t
}

// NewCompactRangeOptions creates new compact range options.
func NewCompactRangeOptions() *CompactRangeOptions {
	return &CompactRangeOptions{
		c: C.rocksdb_compactoptions_create(),
	}
}

// Destroy deallocates the CompactionOptions object.
func (opts *CompactRangeOptions) Destroy() {
	C.rocksdb_compactoptions_destroy(opts.c)
	opts.c = nil
}

// SetBottommostLevelCompaction sets bottommost level compaction.
//
// Default: KIfHaveCompactionFilter
func (opts *CompactRangeOptions) SetBottommostLevelCompaction(value BottommostLevelCompaction) {
	C.rocksdb_compactoptions_set_bottommost_level_compaction(opts.c, C.uchar(value))
}
