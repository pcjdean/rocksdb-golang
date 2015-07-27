// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package rocksdb

/*
#include "transaction_log.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type LogFile struct {
	logf C.LogFile_t
}

func (logf *LogFile) finalize() {
	var clogf *C.LogFile_t = &logf.logf
	C.DeleteLogFileT(clogf, toCBool(false))
}

func newLogFileArrayFromCArray(clogfs *C.LogFile_t, sz uint) (logfs []*LogFile) {
	defer C.DeleteLogFileTArray(clogfs)
	logfs = make([]*LogFile, sz)
	for i := uint(0); i < sz; i++ {
		logf := &LogFile{logf: (*[arrayDimenMax]C.LogFile_t)(unsafe.Pointer(clogfs))[i]}
		logfs[i] = logf
		runtime.SetFinalizer(logf, finalize)
	}
	return
}

type TransactionLogIterator struct {
	tranit C.TransactionLogIterator_t
}

func (tranit *TransactionLogIterator) finalize() {
	var ctranit *C.TransactionLogIterator_t = &tranit.tranit
	C.DeleteTransactionLogIteratorT(ctranit, toCBool(false))
}

func (ctranit *C.TransactionLogIterator_t) toTransactionLogIterator() (tranit *TransactionLogIterator) {
	tranit = &TransactionLogIterator{tranit: *ctranit}
	runtime.SetFinalizer(tranit, finalize)
	return
}

type TransactionLogIteratorReadOptions struct {
	tranropt C.TransactionLogIterator_ReadOptions_t
}

func (tranropt *TransactionLogIteratorReadOptions) finalize() {
	var ctranropt *C.TransactionLogIterator_ReadOptions_t = &tranropt.tranropt
	C.DeleteTransactionLogIterator_ReadOptionsT(ctranropt, toCBool(false))
}
