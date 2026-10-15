#ifndef GOROCKSDB_BRIDGE_UTIL_H
#define GOROCKSDB_BRIDGE_UTIL_H

// Shared helpers for the C++ bridges (*.cc). C++ only: do not include this
// from gorocksdb.h or any file that cgo compiles as C.

#include <cstdlib>
#include <cstring>
#include <exception>

#include "rocksdb/status.h"

namespace gorocksdb {

// Stores a malloc-allocated copy of message in *errptr, freeing any existing
// message first, matching RocksDB's C API convention. Always returns 1 so
// callers can report failure even when the copy cannot be allocated, in which
// case *errptr is left null.
inline int SaveError(char** errptr, const char* message) noexcept {
  if (errptr == nullptr) {
    return 1;
  }
  std::free(*errptr);
  *errptr = nullptr;
  if (message == nullptr) {
    return 1;
  }
  const std::size_t size = std::strlen(message) + 1;
  char* copy = static_cast<char*>(std::malloc(size));
  if (copy != nullptr) {
    std::memcpy(copy, message, size);
    *errptr = copy;
  }
  return 1;
}

// Returns 0 for an OK status; otherwise saves its message and returns 1.
// May throw (Status::ToString allocates); call it inside Guard.
inline int SaveStatus(char** errptr, const ROCKSDB_NAMESPACE::Status& status) {
  return status.ok() ? 0 : SaveError(errptr, status.ToString().c_str());
}

// Runs f and returns its result. Converts any C++ exception into an error
// message and returns 1, so exceptions never cross the C/cgo boundary.
// unknown_error is used for exceptions not derived from std::exception.
template <typename F>
int Guard(char** errptr, const char* unknown_error, F f) noexcept {
  try {
    return f();
  } catch (const std::exception& e) {
    return SaveError(errptr, e.what());
  } catch (...) {
    return SaveError(errptr, unknown_error);
  }
}

}  // namespace gorocksdb

#endif  // GOROCKSDB_BRIDGE_UTIL_H
