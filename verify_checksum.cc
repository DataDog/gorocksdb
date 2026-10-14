#include <cstdlib>
#include <cstring>
#include <exception>
#include <string>

#include "gorocksdb.h"
#include "rocksdb/db.h"

using ROCKSDB_NAMESPACE::DB;

// rocksdb_t is opaque in rocksdb/c.h. The RocksDB versions this package
// supports do not expose DB::VerifyChecksum in their public C API, so this
// bridge mirrors the private wrapper in RocksDB's db/c.cc:
//
//   struct rocksdb_t { DB* rep; };
//
// This layout is not a public contract. Re-check it whenever upgrading RocksDB.
//
// Newer RocksDB releases add a native C function (first seen in tagged release
// v11.8.0):
//
//   void rocksdb_verify_checksum(rocksdb_t* db, char** errptr);
//
// We are pinned to an older RocksDB, so we cannot use it yet. After upgrading
// to a release that has it, delete this file and its gorocksdb.h declaration,
// and call C.rocksdb_verify_checksum directly from DB.VerifyChecksum.
struct rocksdb_t {
  DB* rep;
};

namespace {

// Stores a malloc-allocated copy of message in *errptr. Like RocksDB's
// SaveError helper, this frees any existing message before replacing it.
// It never throws. If allocation fails, *errptr is left null and the caller
// still reports failure through its return value.
void SaveError(char** errptr, const char* message) {
  if (errptr == nullptr) {
    return;
  }

  std::free(*errptr);
  *errptr = nullptr;

  if (message == nullptr) {
    return;
  }

  const std::size_t size = std::strlen(message) + 1;
  char* copy = static_cast<char*>(std::malloc(size));
  if (copy == nullptr) {
    return;
  }

  std::memcpy(copy, message, size);
  *errptr = copy;
}

}  // namespace

extern "C" int gorocksdb_verify_checksum(rocksdb_t* db, char** errptr) {
  if (db == nullptr || db->rep == nullptr) {
    SaveError(errptr, "Invalid argument: RocksDB handle is null");
    return 1;
  }

  // C++ exceptions must not cross the C/cgo boundary.
  try {
    const ROCKSDB_NAMESPACE::Status status = db->rep->VerifyChecksum();
    if (status.ok()) {
      return 0;
    }

    const std::string message = status.ToString();
    SaveError(errptr, message.c_str());
  } catch (const std::exception& e) {
    SaveError(errptr, e.what());
  } catch (...) {
    SaveError(errptr, "unknown C++ exception during VerifyChecksum");
  }

  return 1;
}
