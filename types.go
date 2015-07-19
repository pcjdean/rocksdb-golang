// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

package rocksdb

type finalizer interface {
	finalize()
}

func finalize(obj *finalizer) {
	obj.finalize()
}

func newUint64ArrayFromCArray(cuints *C.uint64_t, sz uint) (vals []uint64) {
	defer C.Deleteuint64TArray(cuints)
	vals = make([]uint64, sz)
	for i := 0; i < sz; i++ {
		vals[i] = (*[sz]C.uint64_t)(unsafe.Pointer(cuints))[i]
	}
	return
}
