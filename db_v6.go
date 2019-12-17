package gorocksdb

// #include "rocksdb/c.h"
import "C"
import(
	"errors"
	"unsafe"
)

// FlushCF triggers a manual flush for the column camily.
func (db *DB) FlushCF(cf *ColumnFamilyHandle, opts *FlushOptions) error {
	var cErr *C.char
	C.rocksdb_flush_cf(db.c, opts.c, cf.c, &cErr)
	if cErr != nil {
		defer C.rocksdb_free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}