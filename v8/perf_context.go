package gorocksdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"
import "unsafe"

type PerfContext struct {
	c *C.rocksdb_perfcontext_t
}

func NewPerfContext() *PerfContext {
	return &PerfContext{c: C.rocksdb_perfcontext_create()}
}

func (c *PerfContext) Reset() {
	C.rocksdb_perfcontext_reset(c.c)
}

func (c *PerfContext) Report(excludeZeroCounters bool) string {
	cExcludeZeroCounters := boolToChar(excludeZeroCounters)
	cReport := C.rocksdb_perfcontext_report(c.c, cExcludeZeroCounters)
	defer C.rocksdb_free(unsafe.Pointer(cReport))
	return C.GoString(cReport)
}

func (c *PerfContext) Metric(metric PerfMetric) uint64 {
	return uint64(C.rocksdb_perfcontext_metric(c.c, C.int(metric)))
}

func (c *PerfContext) Destroy() {
	C.rocksdb_perfcontext_destroy(c.c)
	c.c = nil
}

func SetPerfLevel(level PerfLevel) {
	C.rocksdb_set_perf_level(C.int(level))
}

type PerfLevel uint

const (
	PerfLevelUninitialized            = PerfLevel(0)
	PerfLevelDisable                  = PerfLevel(1)
	PerfLevelEnableCount              = PerfLevel(2)
	PerfLevelEnableTimeExceptForMutex = PerfLevel(3)
	PerfLevelEnableTime               = PerfLevel(4)
	PerfLevelOutOfBounds              = PerfLevel(5)
)

type PerfMetric int

const (
	UserKeyComparisonCount            = PerfMetric(0)
	BlockCacheHitCount                = PerfMetric(1)
	BlockReadCount                    = PerfMetric(2)
	BlockReadByte                     = PerfMetric(3)
	BlockReadTime                     = PerfMetric(4)
	BlockChecksumTime                 = PerfMetric(5)
	BlockDecompressTime               = PerfMetric(6)
	GetReadBytes                      = PerfMetric(7)
	MultiGetReadBytes                 = PerfMetric(8)
	IterReadBytes                     = PerfMetric(9)
	InternalKeySkippedCount           = PerfMetric(10)
	InternalDeleteSkippedCount        = PerfMetric(11)
	InternalRecentSkippedCount        = PerfMetric(12)
	InternalMergeCount                = PerfMetric(13)
	GetSnapshotTime                   = PerfMetric(14)
	GetFromMemtableTime               = PerfMetric(15)
	GetFromMemtableCount              = PerfMetric(16)
	GetPostProcessTime                = PerfMetric(17)
	GetFromOutputFilesTime            = PerfMetric(18)
	SeekOnMemtableTime                = PerfMetric(19)
	SeekOnMemtableCount               = PerfMetric(20)
	NextOnMemtableCount               = PerfMetric(21)
	PrevOnMemtableCount               = PerfMetric(22)
	SeekChildSeekTime                 = PerfMetric(23)
	SeekChildSeekCount                = PerfMetric(24)
	SeekMinHeapTime                   = PerfMetric(25)
	SeekMaxHeapTime                   = PerfMetric(26)
	SeekInternalSeekTime              = PerfMetric(27)
	FindNextUserEntryTime             = PerfMetric(28)
	WriteWalTime                      = PerfMetric(29)
	WriteMemtableTime                 = PerfMetric(30)
	WriteDelayTime                    = PerfMetric(31)
	WritePreAndPostProcessTime        = PerfMetric(32)
	DbMutexLockNanos                  = PerfMetric(33)
	DbConditionWaitNanos              = PerfMetric(34)
	MergeOperatorTimeNanos            = PerfMetric(35)
	ReadIndexBlockNanos               = PerfMetric(36)
	ReadFilterBlockNanos              = PerfMetric(37)
	NewTableBlockIterNanos            = PerfMetric(38)
	NewTableIteratorNanos             = PerfMetric(39)
	BlockSeekNanos                    = PerfMetric(40)
	FindTableNanos                    = PerfMetric(41)
	BloomMemtableHitCount             = PerfMetric(42)
	BloomMemtableMissCount            = PerfMetric(43)
	BloomSstHitCount                  = PerfMetric(44)
	BloomSstMissCount                 = PerfMetric(45)
	KeyLockWaitTime                   = PerfMetric(46)
	KeyLockWaitCount                  = PerfMetric(47)
	EnvNewSequentialFileNanos         = PerfMetric(48)
	EnvNewRandomAccessFileNanos       = PerfMetric(49)
	EnvNewWritableFileNanos           = PerfMetric(50)
	EnvReuseWritableFileNanos         = PerfMetric(51)
	EnvNewRandomRwFileNanos           = PerfMetric(52)
	EnvNewDirectoryNanos              = PerfMetric(53)
	EnvFileExistsNanos                = PerfMetric(54)
	EnvGetChildrenNanos               = PerfMetric(55)
	EnvGetChildrenFileAttributesNanos = PerfMetric(56)
	EnvDeleteFileNanos                = PerfMetric(57)
	EnvCreateDirNanos                 = PerfMetric(58)
	EnvCreateDirIfMissingNanos        = PerfMetric(59)
	EnvDeleteDirNanos                 = PerfMetric(60)
	EnvGetFileSizeNanos               = PerfMetric(61)
	EnvGetFileModificationTimeNanos   = PerfMetric(62)
	EnvRenameFileNanos                = PerfMetric(63)
	EnvLinkFileNanos                  = PerfMetric(64)
	EnvLockFileNanos                  = PerfMetric(65)
	EnvUnlockFileNanos                = PerfMetric(66)
	EnvNewLoggerNanos                 = PerfMetric(67)
	NumberAsyncSeek                   = PerfMetric(68)
	BlobCacheHitCount                 = PerfMetric(69)
	BlobReadCount                     = PerfMetric(70)
	BlobReadByte                      = PerfMetric(71)
	BlobReadTime                      = PerfMetric(72)
	BlobChecksumTime                  = PerfMetric(73)
	BlobDecompressTime                = PerfMetric(74)
	InternalRangeDelReseekCount       = PerfMetric(75)
	BlockReadCpuTime                  = PerfMetric(76)
	InternalMergePointLookupCount     = PerfMetric(77)
	DataBlockReadByte                 = PerfMetric(78)
	IndexBlockReadByte                = PerfMetric(79)
	FilterBlockReadByte               = PerfMetric(80)
	CompressionDictBlockReadByte      = PerfMetric(81)
	MetadataBlockReadByte             = PerfMetric(82)
)
