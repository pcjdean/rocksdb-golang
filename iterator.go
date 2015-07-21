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

type Iterator struct {
	it C.Iterator_t
	db *DB // make sure the iterator is deleted before the db
}

func (it *Iterator) finalize() {
	var cit *C.Iterator_t = unsafe.Pointer(&it.it)
	C.DeleteIteratorT(cit, false)
}

func (cit *C.Iterator_t) toIterator(db *DB) (it *Iterator) {
	it = &Iterator{it: *cit, db: db}	
	runtime.SetFinalizer(it, finalize)
	return
}

func newIteratorArrayFromCArray(cit *C.Iterator_t, sz uint, db *DB) (its []*Iterator) {
	defer C.DeleteIteratorTArray(cit)
	its = make([]*Iterator, sz)
	for i := 0; i < sz; i++ {
		its[i] = &Iterator{it: (*[sz]C.Iterator_t)(unsafe.Pointer(cit))[i], db: db}
		runtime.SetFinalizer(its[i], finalize)
	}
	return
}
