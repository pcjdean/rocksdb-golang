// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

package rocksdb

/*
#include "options.h"
*/
import "C"

type Options struct {
	opt C.Options_t
}

func (opt *Options) finalize() {
	var copt *C.Options_t = unsafe.Pointer(&opt.opt)
	C.DeleteOptionsT(copt, false)
}

type DBOptions struct {
	dbopt C.DBOptions_t
}

func (dbopt *DBOptions) finalize() {
	var cdbopt *C.DBOptions_t = unsafe.Pointer(&dbopt.dbopt)
	C.DeleteDBOptionsT(cdbopt, false)
}

type WriteOptions struct {
	wopt C.WriteOptions_t
}

func (wopt *WriteOptions) finalize() {
	var cwopt *C.WriteOptions_t = unsafe.Pointer(&wopt.wopt)
	C.DeleteWriteOptionsT(cwopt, false)
}

type ReadOptions struct {
	ropt C.ReadOptions_t
}

func (ropt *ReadOptions) finalize() {
	var cropt *C.ReadOptions_t = unsafe.Pointer(&ropt.ropt)
	C.DeleteReadOptionsT(cropt, false)
}

