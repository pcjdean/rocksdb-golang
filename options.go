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

func (dbopt *DBOptions) CreateIfMissing() bool {
	var cdbopt *C.DBOptions_t = &dbopt.dbopt
	return C.DBOptions_get_create_if_missing(cdbopt).toBool()
}

func (dbopt *DBOptions) SetCreateIfMissing(val bool) {
	var cdbopt *C.DBOptions_t = &dbopt.dbopt
	C.DBOptions_set_create_if_missing(cdbopt, toCBool(val))
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
	opt.DBOptions.dbopt.rep = opt.opt.rep
	opt.ColumnFamilyOptions.cfopt.rep = opt.opt.rep
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

