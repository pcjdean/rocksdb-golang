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

type cSlice struct {
	slc C.Slice_t
	cptr *C.char
}

type cSlicePtrAry []*cSlice

// The caller is responsible for delete the returned slc
func newSliceFromBytes(bytes []byte) (slc *cSlice) {
	cptr := C.CString(string(bytes))
	slc = &cSlice{slc: C.NewSliceTRawArgs(cptr, C.uint64ToSizeT(C.uint64_t(len(bytes)))), cptr: cptr}
	return
}

func (slc *cSlice) del()  {
	C.DeleteSliceT(&slc.slc, toCBool(false))
	C.free(unsafe.Pointer(slc.cptr))
}

// The caller is responsible for delete the returned slcs
func newSlicesFromBytesArray(bytess [][]byte) (slcs []*cSlice) {
	slcs = make([]*cSlice, len(bytess))
	for i, bytes := range bytess {
		slcs[i] = newSliceFromBytes(bytes)
	}
	return
}

func (slcs *cSlicePtrAry) del() {
	for _, slc := range *slcs {
		slc.del()
	}
}

func (slcs *cSlicePtrAry) toCArray() (cslcs []C.Slice_t) {
	cslcs = make([]C.Slice_t, len(*slcs))
	for i, slc := range *slcs {
		cslcs[i] = slc.slc
	}
	return
}
