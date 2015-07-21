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
#include "slice.h"
*/
import "C"

type cSlice struct {
	slc C.Slice_t
}

type cSlicePtrAry []*cSlice

func newSliceFromBytes(bytes []byte) (slc *cSlice) {
	slc = &cSlice{slc: C.NewSliceTRawArgs(unsafe.Pointer(&bytes[0]), len(bytes))}
	return
}

func (slc *cSlice) del()  {
	C.DeleteSliceT(unsafe.Pointers(&slc.slc), false)
}

func newSlicesFromBytesArray(bytess [][]byte) (slcs []*cSlice) {
	slcs = make([]*cSlice, len(bytess))
	for i, bytes := range bytess {
		slcs[i] = newSliceFromBytes(bytes)
	}
	return
}

func (slcs cSlicePtrAry) del() {
	for _, slc := range slcs {
		slc.del()
	}
}

func (slcs cSlicePtrAry) toCArray (cslcs []C.Slice_t) {
	cslcs = make([]C.Slice_t, len(bytess))
	for i, slc := range slcs {
		cslcs[i] = slc.slc
	}
	return
}
