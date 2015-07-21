// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package rocksdb

/*
#include "transaction_log.h"
*/
import "C"

type LogFile struct {
	logf C.LogFile_t
}

func (logf *LogFile) finalize() {
	var clogf *C.LogFile_t = unsafe.Pointer(&logf.logf)
	C.DeleteLogFileT(clogf, false)
}

func newLogFileArrayFromCArray(clogfs *C.LogFile_t, sz uint) (logfs []*LogFile) {
	defer C.DeleteLogFileTArray(clogfs)
	logfs = make([]*LogFile, sz)
	for i := 0; i < sz; i++ {
		logf := &LogFile{logf: (*[sz]C.LogFile_t)(unsafe.Pointer(clogfs))[i]}
		logfs[i] = logf
		runtime.SetFinalizer(logf, finalize)
	}
	return
}

type TransactionLogIterator struct {
	tranit C.TransactionLogIterator_t
}

func (tranit *TransactionLogIterator) finalize() {
	var ctranit *C.TransactionLogIterator_t = unsafe.Pointer(&tranit.tranit)
	C.DeleteTransactionLogIteratorT(ctranit, false)
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
	var ctranropt *C.TransactionLogIterator_ReadOptions_t = unsafe.Pointer(&tranropt.tranropt)
	C.DeleteTransactionLogIterator_ReadOptionsT(ctranropt, false)
}
