# gorocksdb, a Go wrapper for RocksDB

[![Build Status](https://travis-ci.org/tecbot/gorocksdb.svg)](https://travis-ci.org/tecbot/gorocksdb) [![GoDoc](https://godoc.org/github.com/tecbot/gorocksdb?status.svg)](http://godoc.org/github.com/tecbot/gorocksdb)

## Install

You'll need to build [RocksDB](https://github.com/facebook/rocksdb) v6.16+ on your machine.

After that, you can install gorocksdb using the following command:

    CGO_CFLAGS="-I/path/to/rocksdb/include" \
    CGO_CXXFLAGS="-I/path/to/rocksdb/include" \
    CGO_LDFLAGS="-L/path/to/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" \
      go get github.com/tecbot/gorocksdb

Please note that this package might upgrade the required RocksDB version at any moment.
Vendoring is thus highly recommended if you require high stability.

## Verify a single SST file

`SSTFileReader` opens an SST produced by RocksDB or `SSTFileWriter` without
opening a database:

```go
opts := gorocksdb.NewDefaultOptions()
defer opts.Destroy()

reader, err := gorocksdb.NewSSTFileReader(opts)
if err != nil {
    return err
}
defer reader.Destroy()

if err := reader.Open("/path/to/file.sst"); err != nil {
    return err
}
return reader.VerifyChecksum()
```

Use options compatible with the file's comparator, table format, and environment.
Keep the options and their dependencies undestroyed until after `reader.Destroy()`.
The reader's methods must not be called concurrently. `Open` can replace the
current file; a failed `Open` leaves the previous file open.

`VerifyChecksum` reads the opened SST and checks its block checksums. It blocks
until finished. It does not check manifest whole-file checksums, other SSTs,
WALs, or referenced blob files. For files owned by a running DB, prevent file
deletion while using the reader (for example, with `DB.DisableFileDeletions` and
a matching `DB.EnableFileDeletions` after destroying the reader).

This API is available only in the root package, not `v8`.

## C++ bridges for checksum verification

Most of this package uses RocksDB's public C API (`rocksdb/c.h`). The RocksDB
6.x versions supported by the root package do not expose `DB::VerifyChecksum`
in that API. So `DB.VerifyChecksum` goes through a small C++ bridge in
`verify_checksum.cc`. This method is not available in the `v8` package.

The bridge copies the private layout of RocksDB's `rocksdb_t` wrapper to reach
the underlying C++ `DB`. That layout is not a public contract, so re-check it
whenever you upgrade RocksDB.

Newer RocksDB releases add a native C function, `rocksdb_verify_checksum`. It
first appears in tagged release v11.8.0. We are pinned to an older RocksDB, so
we cannot use it yet. After upgrading to a release that has it, delete the
bridge and call `rocksdb_verify_checksum` directly.

`SSTFileReader` also uses a small C++ bridge, in `sst_file_reader.cc`, because
RocksDB 6.x has no public C API for it. This bridge copies the private layout of
`rocksdb_options_t` to obtain the C++ options. Re-check that layout as well when
upgrading RocksDB. Unlike the DB bridge, it owns its reader handle and does not
need access to a private DB handle.

Because of these bridges, building this package needs a C++ compiler and the
non-lite RocksDB library. The required C++ standard is C+11 or newer.
