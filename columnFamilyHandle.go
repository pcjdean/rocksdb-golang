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

// Go wrap of C column family
type ColumnFamilyHandle struct {
	cfh C.ColumnFamilyHandle_t
	// The cfh is deleted
	closed bool
	db *DB // make sure the ColumnFamilyHandle is deleted before the db
}

// Release the column family
func (cfh *ColumnFamilyHandle) finalize() {
	if !cfh.closed {
		cfh.closed = true
		cfh.db.removeFromCfhmap(cfh)
		var ccfh *C.ColumnFamilyHandle_t = &cfh.cfh
		C.DeleteColumnFamilyHandleT(ccfh, toCBool(false))
	}
}

// Close the @cfh
func (cfh *ColumnFamilyHandle) Close() {
	runtime.SetFinalizer(cfh, nil)
	cfh.finalize()
}

// Return name of the column family
func (cfh *ColumnFamilyHandle) GetName() string {
	var ptr *C.ColumnFamilyHandle_t = &cfh.cfh
	rstr := cString{C.ColumnFamilyGetName(ptr)}
	return rstr.goString(true);
}
    
// Return ID of the column family
func (cfh *ColumnFamilyHandle) GetID() uint32 {
	var ptr *C.ColumnFamilyHandle_t = &cfh.cfh
	return uint32(C.ColumnFamilyGetID(ptr))
}

// C ColumnFamilyHandle_t to go ColumnFamilyHandle
func (ccfh *C.ColumnFamilyHandle_t) toColumnFamilyHandle(db *DB) (cfh *ColumnFamilyHandle) {
	cfh = &ColumnFamilyHandle{cfh: *ccfh, db: db}	
	db.addToCfhmap(cfh)
	runtime.SetFinalizer(cfh, finalize)
	return
}

// C array of ColumnFamilyHandle_t to go array of ColumnFamilyHandle
func newColumnFamilyHandleArrayFromCArray(db *DB, cfh *C.ColumnFamilyHandle_t, sz uint) (cfhs []*ColumnFamilyHandle) {
	defer C.DeleteColumnFamilyHandleTArray(cfh)
	cfhs = make([]*ColumnFamilyHandle, sz)
	for i := uint(0); i < sz; i++ {
		cfhs[i] = &ColumnFamilyHandle{cfh: (*[arrayDimenMax]C.ColumnFamilyHandle_t)(unsafe.Pointer(cfh))[i], db: db}
		db.addToCfhmap(cfhs[i])
		runtime.SetFinalizer(cfhs[i], finalize)
	}
	return
}

// Go array of ColumnFamilyHandle to C array of ColumnFamilyHandle_t
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
