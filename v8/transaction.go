package gorocksdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import (
	"errors"
	"unsafe"
)

// Transaction is used with TransactionDB for transaction support.
type Transaction struct {
	c *C.rocksdb_transaction_t
}

// NewNativeTransaction creates a Transaction object.
func NewNativeTransaction(c *C.rocksdb_transaction_t) *Transaction {
	return &Transaction{c}
}

// Commit commits the transaction to the database.
func (transaction *Transaction) Commit() error {
	var (
		cErr *C.char
	)
	C.rocksdb_transaction_commit(transaction.c, &cErr)
	if cErr != nil {
		defer C.rocksdb_free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// Rollback performs a rollback on the transaction.
func (transaction *Transaction) Rollback() error {
	var (
		cErr *C.char
	)
	C.rocksdb_transaction_rollback(transaction.c, &cErr)

	if cErr != nil {
		defer C.rocksdb_free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// Get returns the data associated with the key from the database given this transaction.
func (transaction *Transaction) Get(opts *ReadOptions, key []byte) (*Slice, error) {
	var (
		cErr    *C.char
		cValLen C.size_t
		cKey    = byteToChar(key)
	)
	cValue := C.rocksdb_transaction_get(
		transaction.c, opts.c, cKey, C.size_t(len(key)), &cValLen, &cErr,
	)
	if cErr != nil {
		defer C.rocksdb_free(unsafe.Pointer(cErr))
		return nil, errors.New(C.GoString(cErr))
	}
	return NewSlice(cValue, cValLen), nil
}

// GetCF returns the data associated with the key in a given column family from the database given this transaction.
func (transaction *Transaction) GetCF(opts *ReadOptions, cf *ColumnFamilyHandle, key []byte) (*Slice, error) {
	var (
		cErr    *C.char
		cValLen C.size_t
		cKey    = byteToChar(key)
	)
	cValue := C.rocksdb_transaction_get_cf(
		transaction.c, opts.c, cf.c, cKey, C.size_t(len(key)), &cValLen, &cErr,
	)
	if cErr != nil {
		defer C.rocksdb_free(unsafe.Pointer(cErr))
		return nil, errors.New(C.GoString(cErr))
	}
	return NewSlice(cValue, cValLen), nil
}

// GetForUpdate queries the data associated with the key and puts an exclusive lock on the key from the database given this transaction.
func (transaction *Transaction) GetForUpdate(opts *ReadOptions, key []byte) (*Slice, error) {
	var (
		cErr    *C.char
		cValLen C.size_t
		cKey    = byteToChar(key)
	)
	cValue := C.rocksdb_transaction_get_for_update(
		transaction.c, opts.c, cKey, C.size_t(len(key)), &cValLen, C.uchar(byte(1)) /*exclusive*/, &cErr,
	)
	if cErr != nil {
		defer C.rocksdb_free(unsafe.Pointer(cErr))
		return nil, errors.New(C.GoString(cErr))
	}
	return NewSlice(cValue, cValLen), nil
}

// GetForUpdateCF queries the data associated with the key in a given column family
// and puts an exclusive lock on the key from the database given this transaction.
func (transaction *Transaction) GetForUpdateCF(opts *ReadOptions, cf *ColumnFamilyHandle, key []byte) (*Slice, error) {
	var (
		cErr    *C.char
		cValLen C.size_t
		cKey    = byteToChar(key)
	)
	cValue := C.rocksdb_transaction_get_for_update_cf(
		transaction.c, opts.c, cf.c, cKey, C.size_t(len(key)), &cValLen, C.uchar(byte(1)) /*exclusive*/, &cErr,
	)
	if cErr != nil {
		defer C.rocksdb_free(unsafe.Pointer(cErr))
		return nil, errors.New(C.GoString(cErr))
	}
	return NewSlice(cValue, cValLen), nil
}

// GetPinnedForUpdateCF queries the data associated with the key in a given column family
// and puts an exclusive lock on the key from the database given this transaction.
// It uses a pinnable slice to improve performance by avoiding a memcpy.
func (transaction *Transaction) GetPinnedForUpdateCF(opts *ReadOptions, cf *ColumnFamilyHandle, key []byte) (*PinnableSliceHandle, error) {
	var (
		cErr *C.char
		cKey = byteToChar(key)
	)

	cHandle := C.rocksdb_transaction_get_pinned_for_update_cf(
		transaction.c, opts.c, cf.c, cKey, C.size_t(len(key)), C.uchar(byte(1)) /*exclusive*/, &cErr,
	)
	if cErr != nil {
		defer C.rocksdb_free(unsafe.Pointer(cErr))
		return nil, errors.New(C.GoString(cErr))
	}
	return NewNativePinnableSliceHandle(cHandle), nil
}

// Put writes data associated with a key to the transaction.
func (transaction *Transaction) Put(key, value []byte) error {
	var (
		cErr   *C.char
		cKey   = byteToChar(key)
		cValue = byteToChar(value)
	)
	C.rocksdb_transaction_put(
		transaction.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr,
	)
	if cErr != nil {
		defer C.rocksdb_free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// PutCF writes data associated with a key in a given family to the transaction.
func (transaction *Transaction) PutCF(cf *ColumnFamilyHandle, key, value []byte) error {
	var (
		cErr   *C.char
		cKey   = byteToChar(key)
		cValue = byteToChar(value)
	)
	C.rocksdb_transaction_put_cf(
		transaction.c, cf.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr,
	)
	if cErr != nil {
		defer C.rocksdb_free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// Delete removes the data associated with the key from the transaction.
func (transaction *Transaction) Delete(key []byte) error {
	var (
		cErr *C.char
		cKey = byteToChar(key)
	)
	C.rocksdb_transaction_delete(transaction.c, cKey, C.size_t(len(key)), &cErr)
	if cErr != nil {
		defer C.rocksdb_free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// DeleteCF removes the data in a given column family associated with the key from the transaction.
func (transaction *Transaction) DeleteCF(cf *ColumnFamilyHandle, key []byte) error {
	var (
		cErr *C.char
		cKey = byteToChar(key)
	)
	C.rocksdb_transaction_delete_cf(transaction.c, cf.c, cKey, C.size_t(len(key)), &cErr)
	if cErr != nil {
		defer C.rocksdb_free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// NewIterator returns an Iterator over the database that uses the
// ReadOptions given.
func (transaction *Transaction) NewIterator(opts *ReadOptions) *Iterator {
	return NewNativeIterator(unsafe.Pointer(C.rocksdb_transaction_create_iterator(transaction.c, opts.c)))
}

// NewIteratorCF returns an Iterator over the column family that uses the
// ReadOptions given.
func (transaction *Transaction) NewIteratorCF(opts *ReadOptions, cf *ColumnFamilyHandle) *Iterator {
	return NewNativeIterator(unsafe.Pointer(C.rocksdb_transaction_create_iterator_cf(transaction.c, opts.c, cf.c)))
}

// Destroy deallocates the transaction object.
func (transaction *Transaction) Destroy() {
	C.rocksdb_transaction_destroy(transaction.c)
	transaction.c = nil
}