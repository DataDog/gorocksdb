package gorocksdb

// #include "rocksdb/c.h"
// #include "gorocksdb.h"
import "C"

// MemoryAllocator wraps a memory allocator for rocksdb.
type MemoryAllocator struct {
	c *C.rocksdb_memory_allocator_t
}

// Destroy the allocator.
func (m *MemoryAllocator) Destroy() {
	C.rocksdb_memory_allocator_destroy(m.c)
}
