// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

package rocksdb

/*
#include "types.h"
*/
import "C"

type SequenceNumber uint64

const arrayDimenMax = 0xFFFFFFFF

type finalizer interface {
	finalize()
}

func finalize(obj finalizer) {
	obj.finalize()
}

func newUint64ArrayFromCArray(cuints *[]C.uint64_t) (vals []uint64) {
	vals = make([]uint64, len(*cuints))
	for i, v := range *cuints {
		vals[i] = uint64(v)
	}
	return
}

func newCIntArrayFromArray(vals *[]int) (cints []C.int) {
	cints = make([]C.int, len(*vals))
	for i, v := range *vals {
		cints[i] = C.int(v)
	}
	return
}

func toCBool(b bool) (ret C.bool) {
	if b {
		ret = C.true
	} else {
		ret = C.false
	}
	return
}

func (b C.bool) toBool() (ret bool)  {
	if b == C.true {
		ret = true
	} else {
		ret = false
	}
	return
}
