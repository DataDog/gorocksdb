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

## C++ bridge for `DB.VerifyChecksum`

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

Because of the bridge, building this package now needs a C++ compiler. The
required C++ standard is set in `cxxflags.go`.
