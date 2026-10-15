#include "bridge_util.h"
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

extern "C" int gorocksdb_verify_checksum(rocksdb_t* db, char** errptr) {
  if (db == nullptr || db->rep == nullptr) {
    return gorocksdb::SaveError(errptr, "Invalid argument: RocksDB handle is null");
  }
  return gorocksdb::Guard(
      errptr, "unknown C++ exception during VerifyChecksum",
      [&] { return gorocksdb::SaveStatus(errptr, db->rep->VerifyChecksum()); });
}
