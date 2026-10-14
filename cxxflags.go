package gorocksdb

// RocksDB 6 headers require C++11 for the C++ bridge in verify_checksum.cc.

// #cgo CXXFLAGS: -std=c++11
import "C"
