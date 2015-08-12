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
}

func (cfopt *ColumnFamilyOptions) finalize() {
	var ccfopt *C.ColumnFamilyOptions_t = &cfopt.cfopt
	C.DeleteColumnFamilyOptionsT(ccfopt, toCBool(false))
}

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
	C.ColumnFamilyOptions_set_compression_per_level(ccfopt, &cints[0], C.uint64ToSizeT(C.uint64_t(len(cints))))
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

func (cdbopt *C.DBOptions_t) toDBOptions() (dbopt *DBOptions) {
	dbopt = &DBOptions{dbopt: *cdbopt}	
	runtime.SetFinalizer(dbopt, finalize)
	return
}

type Options struct {
	DBOptions
	ColumnFamilyOptions
	opt C.Options_t
}

func (opt *Options) finalize() {
	var copt *C.Options_t = &opt.opt
	C.DeleteOptionsT(copt, toCBool(false))
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
}

func (wopt *WriteOptions) finalize() {
	var cwopt *C.WriteOptions_t = &wopt.wopt
	C.DeleteWriteOptionsT(cwopt, toCBool(false))
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
}

func (ropt *ReadOptions) finalize() {
	var cropt *C.ReadOptions_t = &ropt.ropt
	C.DeleteReadOptionsT(cropt, toCBool(false))
}

func NewReadOptions() *ReadOptions {
	ropt := &ReadOptions{ropt: C.NewReadOptionsTDefault()}
	runtime.SetFinalizer(ropt, finalize)
	return ropt
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

