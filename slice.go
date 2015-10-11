// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Slice is a simple structure containing a pointer into some external
// storage and a size.  The user of a Slice must ensure that the slice
// is not used after the corresponding external storage has been
// deallocated.
//
// Multiple threads can invoke const methods on a Slice without
// external synchronization, but if any of the threads may call a
// non-const method, all threads accessing the same Slice must use
// external synchronization.

package rocksdb

/*
#include <stdlib.h>
#include "slice.h"
*/
import "C"

import (
	"unsafe"
)

// Go wrap C slice
type cSlice struct {
	slc C.Slice_t
	// Allocated by cgo C.CString
	cptr *C.char
}

// Go wrap C slice array
type cSlicePtrAry []*cSlice

// Go wrap C slice to go bytes
func (slc *cSlice) goBytes(del bool) (val []byte) {
	var (
		cslc *C.Slice_t = &slc.slc
		cstr *C.char = C.SliceData(cslc)
		sz C.size_t = C.SliceSize(cslc)
	)
	if del {
		defer C.DeleteSliceT(cslc, toCBool(false))
	}

	if unsafe.Pointer(cstr) != nil && sz > 0 {
		val = C.GoBytes(unsafe.Pointer(cstr), C.int(sz));
	}
	return
}

// Bytes to Go wrap C slice. The caller is responsible for delete the returned slc
func newSliceFromBytes(bytes []byte) (slc *cSlice) {
	if  nil == bytes {
		// Create a cSlice with a NULL c slice
		slc = &cSlice{slc: C.NewSliceT(unsafe.Pointer(nil)), cptr: nil}
	} else {
		cptr := C.CString(string(bytes))
		slc = &cSlice{slc: C.NewSliceTRawArgs(cptr, C.size_t(len(bytes))), cptr: cptr}
	}
	return
}

// Delete Go wrap C slice
func (slc *cSlice) del()  {
	C.DeleteSliceT(&slc.slc, toCBool(false))
	C.free(unsafe.Pointer(slc.cptr))
}

// C slice to go bytes
// Delete cslc if delEle is true
func (cslc *C.Slice_t) cToBytes(delEle bool) (val []byte) {
	slc := cSlice{slc: *cslc}
	val = slc.goBytes(delEle)
	return
}

// C slice array to go bytes array
// Delete ccslc if delAry is true
// Delete ccslc[i] if delEle is true
func newBytesFromCSliceArray(ccslc *C.Slice_t, sz uint, delAry bool, delEle bool) (strs [][]byte) {
	if delAry {
		defer C.DeleteSliceTArray(ccslc)
	}

	strs = make([][]byte, sz)
	for i := uint(0); i < sz; i++ {
		cslc := cSlice{slc: (*[arrayDimenMax]C.Slice_t)(unsafe.Pointer(ccslc))[i]}
		strs[i] = cslc.goBytes(delEle)
	}
	return
}

// The caller is responsible for delete the returned slcs
func newSlicesFromBytesArray(bytess [][]byte) (slcs []*cSlice) {
	slcs = make([]*cSlice, len(bytess))
	for i, bytes := range bytess {
		slcs[i] = newSliceFromBytes(bytes)
	}
	return
}

// Delete Go wrap C slice array
func (slcs *cSlicePtrAry) del() {
	for _, slc := range *slcs {
		slc.del()
	}
}

// Go wrap C slice array to C slice array
func (slcs *cSlicePtrAry) toCArray() (cslcs []C.Slice_t) {
	cslcs = make([]C.Slice_t, len(*slcs))
	for i, slc := range *slcs {
		cslcs[i] = slc.slc
	}
	return
}

// C slice vector to go bytes array
func (slcv *C.SliceVector_t) toBytesArray() (slcs [][]byte) {
	sz := C.SliceVectorSize(slcv)
	if sz > 0 {
		slcs = make([][]byte, sz)
		for i, _ := range slcs {
			slc := C.SliceVectorIndex(slcv, C.size_t(i))
			slcs[i] = slc.cToBytes(false)
		}
	}
	return
}

// C slice deque to go bytes array
func (slcdq *C.SliceDeque_t) toBytesArray() (slcs [][]byte) {
	sz := C.SliceDequeSize(slcdq)
	if sz > 0 {
		slcs = make([][]byte, sz)
		for i, _ := range slcs {
			slc := C.SliceDequeIndex(slcdq, C.size_t(i))
			slcs[i] = slc.cToBytes(false)
		}
	}
	return
}
