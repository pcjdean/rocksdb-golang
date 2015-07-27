// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package rocksdb

/*
#include "metadata.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type ColumnFamilyMetaData struct {
	cfmd C.ColumnFamilyMetaData_t
}

func (cfmd *ColumnFamilyMetaData) finalize() {
	var ccfmd *C.ColumnFamilyMetaData_t = &cfmd.cfmd
	C.DeleteColumnFamilyMetaDataT(ccfmd, toCBool(false))
}

func (ccfmd *C.ColumnFamilyMetaData_t) toColumnFamilyMetaData() (cfmd *ColumnFamilyMetaData) {
	cfmd = &ColumnFamilyMetaData{cfmd: *ccfmd}	
	runtime.SetFinalizer(cfmd, finalize)
	return
}

type LiveFileMetaData struct {
	lfmd C.LiveFileMetaData_t
}

func (lfmd *LiveFileMetaData) finalize() {
	var clfmd *C.LiveFileMetaData_t = &lfmd.lfmd
	C.DeleteLiveFileMetaDataT(clfmd, toCBool(false))
}

func newLiveFileMetaDataArrayFromCArray(clfmd *C.LiveFileMetaData_t, sz uint) (lfmds []*LiveFileMetaData) {
	defer C.DeleteLiveFileMetaDataTArray(clfmd)
	lfmds = make([]*LiveFileMetaData, sz)
	for i := uint(0); i < sz; i++ {
		lfmds[i] = &LiveFileMetaData{lfmd: (*[arrayDimenMax]C.LiveFileMetaData_t)(unsafe.Pointer(clfmd))[i]}
		runtime.SetFinalizer(lfmds[i], finalize)
	}
	return
}
