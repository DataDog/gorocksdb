# gorocksdb, a Go wrapper for RocksDB

[![Build Status](https://travis-ci.org/tecbot/gorocksdb.svg)](https://travis-ci.org/tecbot/gorocksdb) [![GoDoc](https://godoc.org/github.com/tecbot/gorocksdb?status.svg)](http://godoc.org/github.com/tecbot/gorocksdb)

## Install

You'll need to build [RocksDB](https://github.com/facebook/rocksdb) v6.16+ on your machine. Currently this needs to be 6.x, and as of March 2025 it needs to also be < 6.24.2 (I recommend 6.23.3). Later versions break from, e.g. https://github.com/facebook/rocksdb/commit/2a383f21f46fc11687a2fdd3e81616e93e980df4 .

After that, you can install gorocksdb using the following command:

    CGO_CFLAGS="-I/path/to/rocksdb/include" \
    CGO_LDFLAGS="-L/path/to/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" \
      go get github.com/tecbot/gorocksdb

Please note that this package might upgrade the required RocksDB version at any moment.
Vendoring is thus highly recommended if you require high stability.
