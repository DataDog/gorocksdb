package gorocksdb

import (
	"errors"
	"strings"
)

// The C API flattens rocksdb::Status into a single char* holding
// Status::ToString(), losing the status code. Callers that need to tell apart
// "the operation was refused" from "the operation failed" therefore have to
// recover the code from the message prefix, which is what the helpers below do.
//
// The prefixes come from Status::ToString in util/status.cc and are part of
// RocksDB's observable output, but they are still a string contract rather than
// an API one. Keep the matching here, next to the binding, so that it is in one
// place if the C API ever grows a status-code-returning variant.
const incompletePrefix = "Result incomplete: "

// ErrIncomplete reports that RocksDB returned Status::Incomplete, meaning the
// operation was declined without being attempted rather than failed partway.
// Errors matching it wrap the original message and can be tested with
// errors.Is.
//
// The case that matters in practice is a write against WriteOptions with
// SetNoSlowdown(true) while the DB is stalling: RocksDB returns
// Incomplete("Write stall") instead of delaying the write. Such a write did not
// happen and may be retried later.
var ErrIncomplete = errors.New("rocksdb: result incomplete")

// IsIncomplete reports whether err came back from RocksDB as
// Status::Incomplete. See ErrIncomplete.
func IsIncomplete(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrIncomplete) {
		return true
	}
	// Errors from the C API are plain errors.New(C.GoString(errptr)) values, so
	// for those the prefix is all there is to go on.
	return strings.HasPrefix(err.Error(), incompletePrefix)
}
