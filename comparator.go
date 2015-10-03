// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package rocksdb

/*
#include "comparator.h"
*/
import "C"

import (
	"unsafe"
	"runtime"
)

// A Comparator object provides a total order across slices that are
// used as keys in an sstable or a database.  A Comparator implementation
// must be thread-safe since rocksdb may invoke its methods concurrently
// from multiple threads.
type IComparator interface {

	// The name of the comparator.  Used to check for comparator
	// mismatches (i.e., a DB created with one comparator is
	// accessed using a different comparator.
	//
	// The client of this package should switch to a new name whenever
	// the comparator implementation changes in a way that will cause
	// the relative ordering of any two keys to change.
	//
	// Names starting with "rocksdb." are reserved and should not be used
	// by any clients of this package.
	Name() string

	// Three-way comparison.  Returns value:
	//   < 0 iff "a" < "b",
	//   == 0 iff "a" == "b",
	//   > 0 iff "a" > "b"
	Compare(a, b []byte) int

	// Advanced functions: these are used to reduce the space requirements
	// for internal data structures like index blocks.

	// If *start < limit, changes *start to a short string in [start,limit).
	// Simple comparator implementations may return with *start unchanged,
	// i.e., an implementation of this method that does nothing is correct.
	FindShortestSeparator(start, limit []byte) []byte

	// Changes *key to a short string >= *key.
	// Simple comparator implementations may return with *key unchanged,
	// i.e., an implementation of this method that does nothing is correct.
	FindShortSuccessor(key []byte) []byte
}

// Wrap functions for IComparator

//export IComparatorCompare
func IComparatorCompare(ccmp unsafe.Pointer, a, b *C.Slice_t) C.int {
	cmp := InterfacesGet(ccmp).(IComparator)
	return C.int(cmp.Compare(a.cToBytes(false), b.cToBytes(false)))
}

//export IComparatorName
func IComparatorName(ccmp unsafe.Pointer) *C.char {
	cmp := InterfacesGet(ccmp).(IComparator)
	return C.CString(cmp.Name())
}

//export IComparatorFindShortestSeparator
func IComparatorFindShortestSeparator(ccmp unsafe.Pointer, start *C.String_t, limit *C.Slice_t, sz *C.size_t) (val *C.char) {
	val = nil
	cmp := InterfacesGet(ccmp).(IComparator)
	sep := cmp.FindShortestSeparator(start.cToBytes(false), limit.cToBytes(false))
	if nil != sep {
		val = C.CString(string(sep))
		*sz =  C.uint64ToSizeT(C.uint64_t(len(sep)))
	}
	return 
}

//export IComparatorFindShortSuccessor
func IComparatorFindShortSuccessor(ccmp unsafe.Pointer, key *C.String_t, sz *C.size_t) (val *C.char) {
	val = nil
	cmp := InterfacesGet(ccmp).(IComparator)
	sep := cmp.FindShortSuccessor(key.cToBytes(false))
	if nil != sep {
		val = C.CString(string(sep))
		*sz =  C.uint64ToSizeT(C.uint64_t(len(sep)))
	}
	return 
}

// Wrap go Comparator
type Comparator struct {
	cmp C.Comparator_t
	// True if the Comparator is closed
	closed bool
}

// Release resources
func (cmp *Comparator) finalize() {
	if !cmp.closed {
		cmp.closed = true
		var ccmp *C.Comparator_t= &cmp.cmp
		C.DeleteComparatorT(ccmp, toCBool(false))
	}
}

// Close the @cmp
func (cmp *Comparator) Close() {
	runtime.SetFinalizer(cmp, nil)
	cmp.finalize()
}

// C Comparator to go Comparator
func (ccmp *C.Comparator_t) toComparator(del bool) (cmp *Comparator) {
	cmp = &Comparator{cmp: *ccmp}	
	if del {
		runtime.SetFinalizer(cmp, finalize)
	}
	return
}

// Return a new Comparator that uses IComparator
func NewComparator(itf IComparator) (cmp *Comparator) {
	var citf unsafe.Pointer = nil

	if nil != itf {
		citf = InterfacesAddReference(itf)
	}
	ccmp := C.NewComparator(citf)
	return ccmp.toComparator(true)
}

// Return a new default Comparator
func NewDefaultComparator() (cmp *Comparator) {
	cmp = &Comparator{cmp: C.Comparator_t{nil}}	
	return
}

// Return a builtin comparator that uses lexicographic byte-wise
// ordering.  The result remains the property of this module and
// must not be deleted.
func NewBytewiseComparator() (cmp *Comparator) {
	ccmp := C.GoBytewiseComparator()
	return ccmp.toComparator(false)
}

// Return a builtin comparator that uses reverse lexicographic byte-wise
// ordering.
func NewReverseBytewiseComparator() (cmp *Comparator) {
	ccmp := C.GoReverseBytewiseComparator()
	return ccmp.toComparator(false)
}
