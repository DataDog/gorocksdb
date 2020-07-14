// +build !linux !static

package gorocksdb

// #cgo LDFLAGS: -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd -ldl
// #cgo CXXFLAGS: --std=c++11
import "C"
