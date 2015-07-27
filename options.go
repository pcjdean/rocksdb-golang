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

type Options struct {
	opt C.Options_t
}

func (opt *Options) finalize() {
	var copt *C.Options_t = &opt.opt
	C.DeleteOptionsT(copt, toCBool(false))
}

func (copt *C.Options_t) toOptions() (opt *Options) {
	opt = &Options{opt: *copt}	
	runtime.SetFinalizer(opt, finalize)
	return
}

type DBOptions struct {
	dbopt C.DBOptions_t
}

func (dbopt *DBOptions) finalize() {
	var cdbopt *C.DBOptions_t = &dbopt.dbopt
	C.DeleteDBOptionsT(cdbopt, toCBool(false))
}

func (cdbopt *C.DBOptions_t) toDBOptions() (dbopt *DBOptions) {
	dbopt = &DBOptions{dbopt: *cdbopt}	
	runtime.SetFinalizer(dbopt, finalize)
	return
}

type WriteOptions struct {
	wopt C.WriteOptions_t
}

func (wopt *WriteOptions) finalize() {
	var cwopt *C.WriteOptions_t = &wopt.wopt
	C.DeleteWriteOptionsT(cwopt, toCBool(false))
}

type ReadOptions struct {
	ropt C.ReadOptions_t
}

func (ropt *ReadOptions) finalize() {
	var cropt *C.ReadOptions_t = &ropt.ropt
	C.DeleteReadOptionsT(cropt, toCBool(false))
}

type FlushOptions struct {
	fopt C.FlushOptions_t
}

func (fopt *FlushOptions) finalize() {
	var cfopt *C.FlushOptions_t = &fopt.fopt
	C.DeleteFlushOptionsT(cfopt, toCBool(false))
}

type CompactionOptions struct {
	copt C.CompactionOptions_t
}

func (copt *CompactionOptions) finalize() {
	var ccopt *C.CompactionOptions_t = &copt.copt
	C.DeleteCompactionOptionsT(ccopt, toCBool(false))
}

