// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package rocksdb

/*
#include "columnfamilyhandle.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type ColumnFamilyHandle struct {
	cfh C.ColumnFamilyHandle_t
}

func (cfh *ColumnFamilyHandle) finalize() {
	var ccfh *C.ColumnFamilyHandle_t = &cfh.cfh
	C.DeleteColumnFamilyHandleT(ccfh, toCBool(false))
}

func (cfh *ColumnFamilyHandle) GetName() string {
	var ptr *C.ColumnFamilyHandle_t = &cfh.cfh
	rstr := cString{C.ColumnFamilyGetName(ptr)}
	return rstr.goString(true);
}
    
func (cfh *ColumnFamilyHandle) GetID() uint32 {
	var ptr *C.ColumnFamilyHandle_t = &cfh.cfh
	return uint32(C.ColumnFamilyGetID(ptr))
}

func (ccfh *C.ColumnFamilyHandle_t) toColumnFamilyHandle() (cfh *ColumnFamilyHandle) {
	cfh = &ColumnFamilyHandle{cfh: *ccfh}	
	runtime.SetFinalizer(cfh, finalize)
	return
}

func newColumnFamilyHandleArrayFromCArray(cfh *C.ColumnFamilyHandle_t, sz uint) (cfhs []*ColumnFamilyHandle) {
	defer C.DeleteColumnFamilyHandleTArray(cfh)
	cfhs = make([]*ColumnFamilyHandle, sz)
	for i := uint(0); i < sz; i++ {
		cfhs[i] = &ColumnFamilyHandle{cfh: (*[arrayDimenMax]C.ColumnFamilyHandle_t)(unsafe.Pointer(cfh))[i]}
		runtime.SetFinalizer(cfhs[i], finalize)
	}
	return
}

func newCArrayFromColumnFamilyHandleArray(cfhs ...*ColumnFamilyHandle) (ccfhs []C.ColumnFamilyHandle_t) {
	var cfhlen int
	if cfhs != nil {
		cfhlen = len(cfhs)
		ccfhs = make([]C.ColumnFamilyHandle_t, cfhlen)
		for i := 0; i < cfhlen; i++ {
			ccfhs[i] = cfhs[i].cfh
		}
	}
	return
}
