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

import (
	"runtime"
	"sync"
)

type WriteBatch struct {
	wbt C.WriteBatch_t
	mutex sync.Mutex
	closed bool
}

func (wbt *WriteBatch) finalize() {
	if !wbt.closed {
		wbt.closed = true
		var cwbt *C.WriteBatch_t = &wbt.wbt
		C.DeleteWriteBatchT(cwbt, toCBool(false))
	}
}

func (wbt *WriteBatch) Close() {
	runtime.SetFinalizer(wbt, nil)
	wbt.finalize()
}

func NewWriteBatch() *WriteBatch {
	wbt:= &WriteBatch{wbt: C.NewWriteBatchTDefault(), mutex: sync.Mutex{}}
	runtime.SetFinalizer(wbt, finalize)
	return wbt
}

func NewWriteBatchFromBytes(bytes []byte) *WriteBatch {
	str := string(bytes)
	cstr := newCStringFromString(&str)
	defer cstr.del()
	wbt:= &WriteBatch{wbt: C.NewWriteBatchTRawArgs(&cstr.str), mutex: sync.Mutex{}}
	runtime.SetFinalizer(wbt, finalize)
	return wbt
}

// Store the mapping "key->value" in the database.
func (wbt *WriteBatch) Put(key []byte, val []byte, cfh ...*ColumnFamilyHandle) {
	if wbt.closed {
		return
	}

	ckey := newSliceFromBytes(key)
	defer ckey.del()
	cval := newSliceFromBytes(val)
	defer cval.del()

	var (
		cwbt *C.WriteBatch_t = &wbt.wbt
		ccfh *C.ColumnFamilyHandle_t
		cckey *C.Slice_t = &ckey.slc 
		ccval *C.Slice_t = &cval.slc 
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	defer wbt.mutex.Unlock()
	wbt.mutex.Lock()
	if ccfh != nil {
		C.WriteBatchPutWithColumnFamily(cwbt, ccfh, cckey, ccval)
	} else {
		C.WriteBatchPut(cwbt, cckey, ccval)
	}
	return
}

// Merge "value" with the existing value of "key" in the database.
// "key->merge(existing, value)"
func (wbt *WriteBatch) Merge(key []byte, val []byte, cfh ...*ColumnFamilyHandle) {
	if wbt.closed {
		return
	}

	ckey := newSliceFromBytes(key)
	defer ckey.del()
	cval := newSliceFromBytes(val)
	defer cval.del()

	var (
		cwbt *C.WriteBatch_t = &wbt.wbt
		ccfh *C.ColumnFamilyHandle_t
		cckey *C.Slice_t = &ckey.slc 
		ccval *C.Slice_t = &cval.slc 
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	defer wbt.mutex.Unlock()
	wbt.mutex.Lock()
	if ccfh != nil {
		C.WriteBatchMergeWithColumnFamily(cwbt, ccfh, cckey, ccval)
	} else {
		C.WriteBatchMerge(cwbt, cckey, ccval)
	}
	return
}

// If the database contains a mapping for "key", erase it.  Else do nothing.
func (wbt *WriteBatch) Delete(key []byte, cfh ...*ColumnFamilyHandle) {
	if wbt.closed {
		return
	}

	ckey := newSliceFromBytes(key)
	defer ckey.del()

	var (
		cwbt *C.WriteBatch_t = &wbt.wbt
		ccfh *C.ColumnFamilyHandle_t
		cckey *C.Slice_t = &ckey.slc 
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	defer wbt.mutex.Unlock()
	wbt.mutex.Lock()
	if ccfh != nil {
		C.WriteBatchDeleteWithColumnFamily(cwbt, ccfh, cckey)
	} else {
		C.WriteBatchDelete(cwbt, cckey)
	}
	return
}

// Clear all updates buffered in this batch.
func (wbt *WriteBatch) Clear() {
	if wbt.closed {
		return
	}

	var (
		cwbt *C.WriteBatch_t = &wbt.wbt
	)

	defer wbt.mutex.Unlock()
	wbt.mutex.Lock()

	C.WriteBatchClear(cwbt)
}

// Retrieve the serialized version of this batch.
func (wbt *WriteBatch) Data() []byte {
	if wbt.closed {
		return nil
	}

	var (
		cwbt *C.WriteBatch_t = &wbt.wbt
	)

	defer wbt.mutex.Unlock()
	wbt.mutex.Lock()

	cstr := C.WriteBatchData(cwbt)

	return cstr.cToBytes()
}

// Retrieve data size of the batch.
func (wbt *WriteBatch) GetDataSize() uint64 {
	if wbt.closed {
		return 0
	}

	var (
		cwbt *C.WriteBatch_t = &wbt.wbt
	)

	defer wbt.mutex.Unlock()
	wbt.mutex.Lock()

	sz := C.WriteBatchGetDataSize(cwbt)

	return uint64(sz)
}

// Returns the number of updates in the batch
func (wbt *WriteBatch) Count() uint64 {
	if wbt.closed {
		return 0
	}

	var (
		cwbt *C.WriteBatch_t = &wbt.wbt
	)

	defer wbt.mutex.Unlock()
	wbt.mutex.Lock()

	cnt := C.WriteBatchCount(cwbt)

	return uint64(cnt)
}
