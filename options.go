// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

package rocksdb

/*
#include "options.h"
*/
import "C"

import (
	"runtime"
)

// DB contents are stored in a set of blocks, each of which holds a
// sequence of key,value pairs.  Each block may be compressed before
// being stored in a file.  The following enum describes which
// compression method (if any) is used to compress a block.
const (
	// NOTE: do not change the values of existing entries, as these are
	// part of the persistent format on disk.
	NoCompression int = iota
	SnappyCompression
	ZlibCompression
	BZip2Compression
	LZ4Compression
	LZ4HCCompression
)

type ColumnFamilyOptions struct {
	cfopt C.ColumnFamilyOptions_t
	// true if cfopt is deleted
	closed bool
}

func (cfopt *ColumnFamilyOptions) finalize() {
	if !cfopt.closed {
		cfopt.closed = true
		var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
		C.DeleteColumnFamilyOptionsT(ccfopt, toCBool(false))
	}
}

// Close the @cfopt
func (cfopt *ColumnFamilyOptions) Close() {
	runtime.SetFinalizer(cfopt, nil)
	cfopt.finalize()
}

// Create default ColumnFamilyOptions
func NewColumnFamilyOptions() *ColumnFamilyOptions {
	cfopt := &ColumnFamilyOptions{cfopt: C.NewColumnFamilyOptionsTDefault()}
	runtime.SetFinalizer(cfopt, finalize)
	return cfopt
}

// Compress blocks using the specified compression algorithm.  This
// parameter can be changed dynamically.
//
// Default: kSnappyCompression, which gives lightweight but fast
// compression.
//
// Typical speeds of kSnappyCompression on an Intel(R) Core(TM)2 2.4GHz:
//    ~200-500MB/s compression
//    ~400-800MB/s decompression
// Note that these speeds are significantly faster than most
// persistent storage speeds, and therefore it is typically never
// worth switching to kNoCompression.  Even if the input data is
// incompressible, the kSnappyCompression implementation will
// efficiently detect that and will switch to uncompressed mode.
func (cfopt *ColumnFamilyOptions) Compression() int {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	return int(C.ColumnFamilyOptions_get_compression(ccfopt))
}

func (cfopt *ColumnFamilyOptions) SetCompression(val int) {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	C.ColumnFamilyOptions_set_compression(ccfopt, C.int(val))
}

// -------------------
// Parameters that affect performance

// Amount of data to build up in memory (backed by an unsorted log
// on disk) before converting to a sorted on-disk file.
//
// Larger values increase performance, especially during bulk loads.
// Up to max_write_buffer_number write buffers may be held in memory
// at the same time,
// so you may wish to adjust this parameter to control memory usage.
// Also, a larger write buffer will result in a longer recovery time
// the next time the database is opened.
//
// Note that write_buffer_size is enforced per column family.
// See db_write_buffer_size for sharing memory across column families.
//
// Default: 4MB
//
// Dynamically changeable through SetOptions() API
func (cfopt *ColumnFamilyOptions) WriteBufferSize() uint64 {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	return uint64(C.ColumnFamilyOptions_get_write_buffer_size(ccfopt))
}

func (cfopt *ColumnFamilyOptions) SetWriteBufferSize(val uint64) {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	C.ColumnFamilyOptions_set_write_buffer_size(ccfopt, C.size_t(val))
}

// different options for compression algorithms
func (cfopt *ColumnFamilyOptions) SetCompressionOptions(wBits int, level int, strategy int) {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	C.ColumnFamilyOptions_set_compression_options(ccfopt, C.int(wBits), C.int(level), C.int(strategy))
}

// Different levels can have different compression policies. There
// are cases where most lower levels would like to use quick compression
// algorithms while the higher levels (which have more data) use
// compression algorithms that have better compression but could
// be slower. This array, if non-empty, should have an entry for
// each level of the database; these override the value specified in
// the previous field 'compression'.
//
// NOTICE if level_compaction_dynamic_level_bytes=true,
// compression_per_level[0] still determines L0, but other elements
// of the array are based on base level (the level L0 files are merged
// to), and may not match the level users see from info log for metadata.
// If L0 files are merged to level-n, then, for i>0, compression_per_level[i]
// determines compaction type for level n+i-1.
// For example, if we have three 5 levels, and we determine to merge L0
// data to L4 (which means L1..L3 will be empty), then the new files go to
// L4 uses compression type compression_per_level[1].
// If now L0 is merged to L2. Data goes to L2 will be compressed
// according to compression_per_level[1], L3 using compression_per_level[2]
// and L4 using compression_per_level[3]. Compaction for each level can
// change when data grows.
func (cfopt *ColumnFamilyOptions) SetCompressionPerLevel(levelValues []int) {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	cints := newCIntArrayFromArray(&levelValues)
	C.ColumnFamilyOptions_set_compression_per_level(ccfopt, &cints[0], C.size_t(len(cints)))
}

// This is a factory that provides MemTableRep objects.
// Default: a factory that provides a skip-list-based implementation of
// MemTableRep.
func (cfopt *ColumnFamilyOptions) SetMemtableFactory(mtf *MemTableRepFactory) {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	C.ColumnFamilyOptions_set_memtable_factory(ccfopt, &mtf.mtf)
}

// This is a factory that provides TableFactory objects.
// Default: a block-based table factory that provides a default
// implementation of TableBuilder and TableReader with default
// BlockBasedTableOptions.
func (cfopt *ColumnFamilyOptions) SetTableFactory(tbf *TableFactory) {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	C.ColumnFamilyOptions_set_table_factory(ccfopt, &tbf.tbf)
}

// REQUIRES: The client must provide a merge operator if Merge operation
// needs to be accessed. Calling Merge on a DB without a merge operator
// would result in Status::NotSupported. The client must ensure that the
// merge operator supplied here has the same name and *exactly* the same
// semantics as the merge operator provided to previous open calls on
// the same DB. The only exception is reserved for upgrade, where a DB
// previously without a merge operator is introduced to Merge operation
// for the first time. It's necessary to specify a merge operator when
// openning the DB in this case.
// Default: nullptr
func (cfopt *ColumnFamilyOptions) SetMergeOperator(mop *MergeOperator) {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	C.ColumnFamilyOptions_set_merge_operator(ccfopt, &mop.mop)
}

// If non-nullptr, use the specified function to determine the
// prefixes for keys.  These prefixes will be placed in the filter.
// Depending on the workload, this can reduce the number of read-IOP
// cost for scans when a prefix is passed via ReadOptions to
// db.NewIterator().  For prefix filtering to work properly,
// "prefix_extractor" and "comparator" must be such that the following
// properties hold:
//
// 1) key.starts_with(prefix(key))
// 2) Compare(prefix(key), key) <= 0.
// 3) If Compare(k1, k2) <= 0, then Compare(prefix(k1), prefix(k2)) <= 0
// 4) prefix(prefix(key)) == prefix(key)
//
// Default: nullptr
func (cfopt *ColumnFamilyOptions) SetPrefixExtractor(stf *SharedSliceTransform) {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	var cstf *C.PConstSliceTransform_t = nil
	if nil != stf {
		cstf = &stf.stf
	}
	C.ColumnFamilyOptions_set_prefix_extractor(ccfopt, cstf)
}

// -------------------
// Parameters that affect behavior

// Comparator used to define the order of keys in the table.
// Default: a comparator that uses lexicographic byte-wise ordering
//
// REQUIRES: The client must ensure that the comparator supplied
// here has the same name and orders keys *exactly* the same as the
// comparator provided to previous open calls on the same DB.
func (cfopt *ColumnFamilyOptions) SetComparator(cmp *Comparator) {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	if nil == cmp {
		cmp = NewDefaultComparator()
	}
	C.ColumnFamilyOptions_set_comparator(ccfopt, &cmp.cmp)
}

// A single CompactionFilter instance to call into during compaction.
// Allows an application to modify/delete a key-value during background
// compaction.
//
// If the client requires a new compaction filter to be used for different
// compaction runs, it can specify compaction_filter_factory instead of this
// option.  The client should specify only one of the two.
// compaction_filter takes precedence over compaction_filter_factory if
// client specifies both.
//
// If multithreaded compaction is being used, the supplied CompactionFilter
// instance may be used from different threads concurrently and so should be
// thread-safe.
//
// Default: nullptr
// The @cpf needs to be cleaned manually by SetCompactionFilter(nil)
func (cfopt *ColumnFamilyOptions) SetCompactionFilter(cpf *CompactionFilter) {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	if nil == cpf {
		cpf = NewDefaultCompactionFilter()
	}
	C.ColumnFamilyOptions_set_compaction_filter(ccfopt, &cpf.cpf)
}

// This is a factory that provides compaction filter objects which allow
// an application to modify/delete a key-value during background compaction.
//
// A new filter will be created on each compaction run.  If multithreaded
// compaction is being used, each created CompactionFilter will only be used
// from a single thread and so does not need to be thread-safe.
//
// Default: a factory that doesn't provide any object
func (cfopt *ColumnFamilyOptions) SetCompactionFilterFactory(cff *CompactionFilterFactory) {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	if nil == cff {
		cff = NewDefaultCompactionFilterFactory()
	}
	C.ColumnFamilyOptions_set_compaction_filter_factory(ccfopt, &cff.cff)
}

// Version TWO of the compaction_filter_factory
// It supports rolling compaction
//
// Default: a factory that doesn't provide any object
func (cfopt *ColumnFamilyOptions) SetCompactionFilterFactoryV2(cff *CompactionFilterFactoryV2) {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	if nil == cff {
		cff = NewDefaultCompactionFilterFactoryV2()
	}
	C.ColumnFamilyOptions_set_compaction_filter_factory_v2(ccfopt, &cff.cff)
}

type DBOptions struct {
	dbopt C.DBOptions_t
}

func (dbopt *DBOptions) finalize() {
	var cdbopt *C.DBOptions_t = &dbopt.dbopt
	C.DeleteDBOptionsT(cdbopt, toCBool(false))
}

func NewDBOptions() *DBOptions {
	dbopt := &DBOptions{dbopt: C.NewDBOptionsTDefault()}
	runtime.SetFinalizer(dbopt, finalize)
	return dbopt
}

// If true, the database will be created if it is missing.
// Default: false
func (dbopt *DBOptions) CreateIfMissing() bool {
	var cdbopt *C.DBOptions_t = &dbopt.dbopt
	return C.DBOptions_get_create_if_missing(cdbopt).toBool()
}

func (dbopt *DBOptions) SetCreateIfMissing(val bool) {
	var cdbopt *C.DBOptions_t = &dbopt.dbopt
	C.DBOptions_set_create_if_missing(cdbopt, toCBool(val))
}

// If true, an error is raised if the database already exists.
// Default: false
func (dbopt *DBOptions) ErrorIfExists() bool {
	var cdbopt *C.DBOptions_t = &dbopt.dbopt
	return C.DBOptions_get_error_if_exists(cdbopt).toBool()
}

func (dbopt *DBOptions) SetErrorIfExists(val bool) {
	var cdbopt *C.DBOptions_t = &dbopt.dbopt
	C.DBOptions_set_error_if_exists(cdbopt, toCBool(val))
}

// Use the specified object to interact with the environment,
// e.g. to read/write files, schedule background work, etc.
// Default: Env::Default()
func (dbopt *DBOptions) Env() (env *Env) {
	var cdbopt *C.DBOptions_t = &dbopt.dbopt
	cenv := C.DBOptions_get_env(cdbopt)
	// The wrapped Env is not deleted by garbage collector
	return cenv.toEnv(false)
}

func (dbopt *DBOptions) SetEnv(env *Env) {
	var cdbopt *C.DBOptions_t = &dbopt.dbopt
	C.DBOptions_set_env(cdbopt, &env.env)
}

// Allow the OS to mmap file for reading sst tables. Default: false
func (dbopt *DBOptions) SetAllowMmapReads(val bool) {
	var cdbopt *C.DBOptions_t = &dbopt.dbopt
	C.DBOptions_set_allow_mmap_reads(cdbopt, toCBool(val))
}

func (cdbopt *C.DBOptions_t) toDBOptions() (dbopt *DBOptions) {
	dbopt = &DBOptions{dbopt: *cdbopt}	
	runtime.SetFinalizer(dbopt, finalize)
	return
}

type Options struct {
	DBOptions
	ColumnFamilyOptions
	opt C.Options_t
	// True if he CompactionFilter is closed
	closed bool
}

// Release the @opt
func (opt *Options) finalize() {
	if !opt.closed {
		opt.closed = true
		var copt *C.Options_t = &opt.opt
		C.DeleteOptionsT(copt, toCBool(false))
	}
}

// Close the @opt
func (opt *Options) Close() {
	runtime.SetFinalizer(opt, nil)
	opt.finalize()
}

func NewOptions() *Options {
	opt := &Options{opt: C.NewOptionsTDefault()}
	C.OptionsTStaticCastToDBOptionsT(&opt.opt, &opt.DBOptions.dbopt)
	C.OptionsTStaticCastToColumnFamilyOptionsT(&opt.opt, &opt.ColumnFamilyOptions.cfopt)
	runtime.SetFinalizer(opt, finalize)
	return opt
}

func (copt *C.Options_t) toOptions() (opt *Options) {
	opt = &Options{opt: *copt}	
	opt.DBOptions.dbopt.rep = opt.opt.rep
	opt.ColumnFamilyOptions.cfopt.rep = opt.opt.rep
	runtime.SetFinalizer(opt, finalize)
	return opt
}

type WriteOptions struct {
	wopt C.WriteOptions_t
	// true if wopt is deleted
	closed bool
}

func (wopt *WriteOptions) finalize() {
	if !wopt.closed {
		wopt.closed = true
		var cwopt *C.WriteOptions_t = &wopt.wopt
		C.DeleteWriteOptionsT(cwopt, toCBool(false))
	}
}

// Close the @wopt
func (wopt *WriteOptions) Close() {
	runtime.SetFinalizer(wopt, nil)
	wopt.finalize()
}

func NewWriteOptions() *WriteOptions {
	wopt := &WriteOptions{wopt: C.NewWriteOptionsTDefault()}
	runtime.SetFinalizer(wopt, finalize)
	return wopt
}

func (wopt *WriteOptions) Sync() bool {
	var cwopt *C.WriteOptions_t = &wopt.wopt
	return C.WriteOptions_get_sync(cwopt).toBool()
}

func (wopt *WriteOptions) SetSync(val bool) {
	var cwopt *C.WriteOptions_t = &wopt.wopt
	C.WriteOptions_set_sync(cwopt, toCBool(val))
}

type ReadOptions struct {
	ropt C.ReadOptions_t
	// Keep snp from garbage collected
	snp *Snapshot
	// Keep slc from garbage collected
	slc *cSlice
	// true if ropt is deleted
	closed bool
}

func (ropt *ReadOptions) finalize() {
	if !ropt.closed {
		ropt.closed = true
		ropt.SetSnapshot(nil)
		ropt.SetIterateUpperBound(nil)
		var cropt *C.ReadOptions_t = &ropt.ropt
		C.DeleteReadOptionsT(cropt, toCBool(false))
	}
}

// Close the @ropt
func (ropt *ReadOptions) Close() {
	runtime.SetFinalizer(ropt, nil)
	ropt.finalize()
}

func NewReadOptions() *ReadOptions {
	ropt := &ReadOptions{ropt: C.NewReadOptionsTDefault()}
	runtime.SetFinalizer(ropt, finalize)
	return ropt
}

// If "snapshot" is non-nullptr, read as of the supplied snapshot
// (which must belong to the DB that is being read and which must
// not have been released).  If "snapshot" is nullptr, use an impliicit
// snapshot of the state at the beginning of this read operation.
// Default: nullptr
func (ropt *ReadOptions) SetSnapshot(snp *Snapshot) {
	ropt.snp = snp
	var (
		cropt *C.ReadOptions_t = &ropt.ropt
		csnp *C.Snapshot_t
	)

	if nil != snp {
		csnp = &snp.snp
	}

	C.ReadOptions_set_snapshot(cropt, csnp)
}

// "iterate_upper_bound" defines the extent upto which the forward iterator
// can returns entries. Once the bound is reached, Valid() will be false.
// "iterate_upper_bound" is exclusive ie the bound value is
// not a valid entry.  If iterator_extractor is not null, the Seek target
// and iterator_upper_bound need to have the same prefix.
// This is because ordering is not guaranteed outside of prefix domain.
// There is no lower bound on the iterator. If needed, that can be easily
// implemented
//
// Default: nullptr
func (ropt *ReadOptions) SetIterateUpperBound(upperbound []byte) {
	var cropt *C.ReadOptions_t = &ropt.ropt
	slc := newSliceFromBytes(upperbound)
	if nil != ropt.slc {
		ropt.slc.del()
	}
	ropt.slc = slc

	C.ReadOptions_set_iterate_upper_bound(cropt, &slc.slc)
}


type FlushOptions struct {
	fopt C.FlushOptions_t
}

func (fopt *FlushOptions) finalize() {
	var cfopt *C.FlushOptions_t = &fopt.fopt
	C.DeleteFlushOptionsT(cfopt, toCBool(false))
}

func NewFlushOptions() *FlushOptions {
	return &FlushOptions{fopt: C.NewFlushOptionsTDefault()}
}

type CompactionOptions struct {
	copt C.CompactionOptions_t
}

func (copt *CompactionOptions) finalize() {
	var ccopt *C.CompactionOptions_t = &copt.copt
	C.DeleteCompactionOptionsT(ccopt, toCBool(false))
}

func NewCompactionOptions() *CompactionOptions {
	copt := &CompactionOptions{copt: C.NewCompactionOptionsTDefault()}
	runtime.SetFinalizer(copt, finalize)
	return copt
}

