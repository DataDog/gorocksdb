#include <memory>
#include <utility>

#include "bridge_util.h"
#include "gorocksdb.h"
#include "rocksdb/options.h"
#include "rocksdb/sst_file_reader.h"

using gorocksdb::Guard;
using gorocksdb::SaveError;
using gorocksdb::SaveStatus;

// rocksdb_options_t is opaque in rocksdb/c.h. This mirrors db/c.cc in
// RocksDB 6.x, including v6.20.3. It is not a public ABI contract: re-check
// this layout whenever upgrading RocksDB.
struct rocksdb_options_t {
  ROCKSDB_NAMESPACE::Options rep;
};

struct gorocksdb_sstfilereader_t {
  explicit gorocksdb_sstfilereader_t(const ROCKSDB_NAMESPACE::Options& opts)
      : options(opts) {}

  ROCKSDB_NAMESPACE::Options options;
  std::unique_ptr<ROCKSDB_NAMESPACE::SstFileReader> reader;
};

extern "C" gorocksdb_sstfilereader_t* gorocksdb_sstfilereader_create(
    const rocksdb_options_t* options, char** errptr) {
  if (options == nullptr) {
    SaveError(errptr, "Invalid argument: RocksDB options are null");
    return nullptr;
  }
  gorocksdb_sstfilereader_t* handle = nullptr;
  Guard(errptr, "unknown C++ exception creating SST reader", [&] {
    handle = new gorocksdb_sstfilereader_t(options->rep);
    return 0;
  });
  return handle;
}

extern "C" int gorocksdb_sstfilereader_open(
    gorocksdb_sstfilereader_t* handle, const char* path, char** errptr) {
  if (handle == nullptr || path == nullptr) {
    return SaveError(errptr, "Invalid argument: SST reader or path is null");
  }
  return Guard(errptr, "unknown C++ exception opening SST file", [&] {
    // Open a fresh reader so a failed Open cannot leave partial state behind.
    std::unique_ptr<ROCKSDB_NAMESPACE::SstFileReader> reader(
        new ROCKSDB_NAMESPACE::SstFileReader(handle->options));
    const int rc = SaveStatus(errptr, reader->Open(path));
    if (rc == 0) {
      handle->reader = std::move(reader);
    }
    return rc;
  });
}

extern "C" int gorocksdb_sstfilereader_verify_checksum(
    gorocksdb_sstfilereader_t* handle, char** errptr) {
  if (handle == nullptr || !handle->reader) {
    return SaveError(errptr, "Invalid argument: SST reader has no open file");
  }
  return Guard(errptr, "unknown C++ exception verifying SST checksums", [&] {
    return SaveStatus(errptr, handle->reader->VerifyChecksum());
  });
}

extern "C" void gorocksdb_sstfilereader_destroy(
    gorocksdb_sstfilereader_t* handle) {
  delete handle;
}
