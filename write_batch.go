// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// WriteBatch holds a collection of updates to apply atomically to a DB.
//
// The updates are applied in the order in which they are added
// to the WriteBatch.  For example, the value of "key" will be "v3"
// after the following batch is written:
//
//    batch.Put("key", "v1");
//    batch.Delete("key");
//    batch.Put("key", "v2");
//    batch.Put("key", "v3");
//
// Multiple threads can invoke const methods on a WriteBatch without
// external synchronization, but if any of the threads may call a
// non-const method, all threads accessing the same WriteBatch must use
// external synchronization.

package rocksdb

/*
#include "write_batch.h"
*/
import "C"

type WriteBatch struct {
	wbt C.WriteBatch_t
}

func (wbt *WriteBatch) Finalize() {
	var cwbt *C.WriteBatch_t = unsafe.Pointer(&wbt.wbt)
	C.DeleteWriteBatchT(cwbt, false)
}
