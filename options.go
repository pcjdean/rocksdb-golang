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

func (opt *Options) Finalize() {
	var copt *C.Options_t = unsafe.Pointer(&opt.opt)
	C.DeleteOptionsT(copt, false)
}

type DBOptions struct {
	dbopt C.DBOptions_t
}

func (dbopt *DBOptions) Finalize() {
	var cdbopt *C.DBOptions_t = unsafe.Pointer(&dbopt.dbopt)
	C.DeleteDBOptionsT(cdbopt, false)
}

type WriteOptions struct {
	wopt C.WriteOptions_t
}

func (wopt *WriteOptions) Finalize() {
	var cwopt *C.WriteOptions_t = unsafe.Pointer(&wopt.wopt)
	C.DeleteWriteOptionsT(cwopt, false)
}

