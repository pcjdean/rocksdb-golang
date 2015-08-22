// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// An iterator yields a sequence of key/value pairs from a source.
// The following class defines the interface.  Multiple implementations
// are provided by this library.  In particular, iterators are provided
// to access the contents of a Table or a DB.
//
// Multiple threads can invoke const methods on an Iterator without
// external synchronization, but if any of the threads may call a
// non-const method, all threads accessing the same Iterator must use
// external synchronization.

package rocksdb

/*
#include "iterator.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
	"sync"
)

// Go Iterator
type Iterator struct {
	it C.Iterator_t
	// Thread safe
	mutex sync.Mutex
	db *DB // make sure the iterator is deleted before the db
	// true if the underlying c object is deleted
	closed bool
}

// Release resources
func (it *Iterator) finalize() {
	if !it.closed {
		it.closed = true
		var cit *C.Iterator_t = &it.it
		C.DeleteIteratorT(cit, toCBool(false))
	}
}

// Close the Iterator
func (it *Iterator) Close() {
	runtime.SetFinalizer(it, nil)
	it.finalize()
}

// Iterator of C to go iterator
func (cit *C.Iterator_t) toIterator(db *DB) (it *Iterator) {
	it = &Iterator{it: *cit, mutex: sync.Mutex{}, db: db}	
	runtime.SetFinalizer(it, finalize)
	return
}

// Array of C iterators to array of go iterators
func newIteratorArrayFromCArray(cit *C.Iterator_t, sz uint, db *DB) (its []*Iterator) {
	defer C.DeleteIteratorTArray(cit)
	its = make([]*Iterator, sz)
	for i := uint(0); i < sz; i++ {
		its[i] = &Iterator{it: (*[arrayDimenMax]C.Iterator_t)(unsafe.Pointer(cit))[i], mutex: sync.Mutex{}, db: db}
		runtime.SetFinalizer(its[i], finalize)
	}
	return
}

// An iterator is either positioned at a key/value pair, or
// not valid.  This method returns true iff the iterator is valid.
func (it *Iterator) Valid() bool {
	defer it.mutex.Unlock()
	it.mutex.Lock()

	var cit *C.Iterator_t = &it.it
	return C.IteratorValid(cit).toBool()
}

// Position at the first key in the source.  The iterator is Valid()
// after this call iff the source is not empty.
func (it *Iterator) SeekToFirst() {
	defer it.mutex.Unlock()
	it.mutex.Lock()

	var cit *C.Iterator_t = &it.it
	C.IteratorSeekToFirst(cit)
}

// Position at the last key in the source.  The iterator is
// Valid() after this call iff the source is not empty.
func (it *Iterator) SeekToLast() {
	defer it.mutex.Unlock()
	it.mutex.Lock()

	var cit *C.Iterator_t = &it.it
	C.IteratorSeekToLast(cit)
}

// Position at the first key in the source that at or past target
// The iterator is Valid() after this call iff the source contains
// an entry that comes at or past target.
func (it *Iterator) Seek(key []byte) {
	defer it.mutex.Unlock()
	it.mutex.Lock()

	ckey := newSliceFromBytes(key)
	defer ckey.del()

	var cit *C.Iterator_t = &it.it
	C.IteratorSeek(cit, &ckey.slc)
}

// Moves to the next entry in the source.  After this call, Valid() is
// true iff the iterator was not positioned at the last entry in the source.
// REQUIRES: Valid()
func (it *Iterator) Next() {
	defer it.mutex.Unlock()
	it.mutex.Lock()

	var cit *C.Iterator_t = &it.it
	C.IteratorNext(cit)
}

// Moves to the previous entry in the source.  After this call, Valid() is
// true iff the iterator was not positioned at the first entry in source.
// REQUIRES: Valid()
func (it *Iterator) Prev() {
	defer it.mutex.Unlock()
	it.mutex.Lock()

	var cit *C.Iterator_t = &it.it
	C.IteratorPrev(cit)
}

// Return the key for the current entry.  The underlying storage for
// the returned slice is valid only until the next modification of
// the iterator.
// REQUIRES: Valid()
func (it *Iterator) Key() (key []byte){
	defer it.mutex.Unlock()
	it.mutex.Lock()

	var cit *C.Iterator_t = &it.it
	ckey := C.IteratorKey(cit)
	key = ckey.cToBytes()
	return
}

// Return the value for the current entry.  The underlying storage for
// the returned slice is valid only until the next modification of
// the iterator.
// REQUIRES: !AtEnd() && !AtStart()
func (it *Iterator) Value() (val []byte){
	defer it.mutex.Unlock()
	it.mutex.Lock()

	var cit *C.Iterator_t = &it.it
	cval := C.IteratorValue(cit)
	val = cval.cToBytes()
	return
}

// If an error has occurred, return it.  Else return an ok status.
// If non-blocking IO is requested and this operation cannot be
// satisfied without doing some IO, then this returns Status::Incomplete().
func (it *Iterator) Status() (val *Status){
	defer it.mutex.Unlock()
	it.mutex.Lock()

	var cit *C.Iterator_t = &it.it
	cval := C.IteratorStatus(cit)
	val = cval.toStatus()
	return
}
