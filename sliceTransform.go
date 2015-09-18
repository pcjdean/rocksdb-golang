// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Class for specifying user-defined functions which perform a
// transformation on a slice.  It is not required that every slice
// belong to the domain and/or range of a function.  Subclasses should
// define InDomain and InRange to determine which slices are in either
// of these sets respectively.

package rocksdb

/*
#include "sliceTransform.h"
*/
import "C"

import (
	"unsafe"
	"log"
	"sync"
)

const (
	// The initial size of callbackISliceTransform
	initialISliceTransformSize = 100
)

var (
	// Map to keep all the ISliceTransform callbacks from garbage collected
	callbackISliceTransform map[unsafe.Pointer]ISliceTransform

	// Mutext to protect callbackISliceTransform
	callbackISliceTransformMutex sync.Mutex = sync.Mutex{}
)

//export ISliceTransformRemoveReference
// Remove interface citf from the callbackISliceTransform 
// to leave citf garbage collected
func ISliceTransformRemoveReference(citf unsafe.Pointer) {
	defer callbackISliceTransformMutex.Unlock()
	callbackISliceTransformMutex.Lock()
	if nil != callbackISliceTransform {
		delete(callbackISliceTransform, citf)
	} else {
		log.Println("ISliceTransformRemoveReference: callbackISliceTransform is not created!")
	}
}

// Get interface itf from the callbackISliceTransform
// with citf as the key
func ISliceTransformGet(citf unsafe.Pointer) (itf ISliceTransform) {
	defer callbackISliceTransformMutex.Unlock()
	callbackISliceTransformMutex.Lock()
	if nil != callbackISliceTransform {
		itf = callbackISliceTransform[citf]
	} else {
		log.Println("ISliceTransformGet: callbackISliceTransform is not created!")
	}
	return
}

// Add interface itf to the callbackISliceTransform to keep itf alive
// Return the key of the ISliceTransform in map callbackISliceTransform
func ISliceTransformAddReference(itf ISliceTransform) (citf unsafe.Pointer) {
	defer callbackISliceTransformMutex.Unlock()
	callbackISliceTransformMutex.Lock()
	if nil == callbackISliceTransform {
		callbackISliceTransform = make(map[unsafe.Pointer]ISliceTransform, initialISliceTransformSize)
	}
	citf = unsafe.Pointer(&itf)
	callbackISliceTransform[citf] = itf
	return
}

type ISliceTransform interface {

	// Return the name of this transformation.
	Name() string

	// transform a src in domain to a dst in the range
	Transform(src []byte) []byte

	// determine whether this is a valid src upon the function applies
	InDomain(src []byte) bool

	// determine whether dst=Transform(src) for some src
	InRange(dst []byte) bool

	// Transform(s)=Transform(`prefix`) for any s with `prefix` as a prefix.
	//
	// This function is not used by RocksDB, but for users. If users pass
	// Options by string to RocksDB, they might not know what prefix extractor
	// they are using. This function is to help users can determine:
	//   if they want to iterate all keys prefixing `prefix`, whetherit is
	//   safe to use prefix bloom filter and seek to key `prefix`.
	// If this function returns true, this means a user can Seek() to a prefix
	// using the bloom filter. Otherwise, user needs to skip the bloom filter
	// by setting ReadOptions.total_order_seek = true.
	//
	// Here is an example: Suppose we implement a slice transform that returns
	// the first part of the string after spliting it using deimiter ",":
	// 1. SameResultWhenAppended("abc,") should return true. If aplying prefix
	//    bloom filter using it, all slices matching "abc:.*" will be extracted
	//    to "abc,", so any SST file or memtable containing any of those key
	//    will not be filtered out.
	// 2. SameResultWhenAppended("abc") should return false. A user will not
	//    guaranteed to see all the keys matching "abc.*" if a user seek to "abc"
	//    against a DB with the same setting. If one SST file only contains
	//    "abcd,e", the file can be filtered out and the key will be invisible.
	//
	// i.e., an implementation always returning false is safe.
	SameResultWhenAppended(prefix []byte) bool
}

// Wrap functions for ISliceTransform

//export ISliceTransformName
func ISliceTransformName(cflp unsafe.Pointer) *C.char {
	flp := ISliceTransformGet(cflp)
	return C.CString(flp.Name())
}

//export ISliceTransformTransform
func ISliceTransformTransform(cflp unsafe.Pointer, src *C.Slice_t) *C.char {
	flp := ISliceTransformGet(cflp)
	return C.CString(string(flp.Transform(src.cToBytes(false))))
}

//export ISliceTransformInDomain
func ISliceTransformInDomain(cflp unsafe.Pointer, src *C.Slice_t) *C.bool {
	flp := ISliceTransformGet(cflp)
	return toCBool(flp.InDomain(src.cToBytes(false)))
}

//export ISliceTransformInRange
func ISliceTransformInRange(cflp unsafe.Pointer, dst *C.Slice_t) *C.bool {
	flp := ISliceTransformGet(cflp)
	return toCBool(flp.InRange(dst.cToBytes(false)))
}

//export ISliceTransformSameResultWhenAppended
func ISliceTransformSameResultWhenAppended(cflp unsafe.Pointer, prefix *C.Slice_t) *C.bool {
	flp := ISliceTransformGet(cflp)
	return toCBool(flp.SameResultWhenAppended(prefix.cToBytes(false)))
}

// Wrap go SliceTransform
type SliceTransform struct {
	stf C.SliceTransform_t
}

// Release resources
func (stf *SliceTransform) finalize() {
	var cstf *C.SliceTransform_t= &stf.stf
	C.DeleteSliceTransformT(cstf, toCBool(false))
}

// C SliceTransform to go SliceTransform
func (cstf *C.SliceTransform_t) toSliceTransform() (stf *SliceTransform) {
	stf = &SliceTransform{stf: *cstf}	
	runtime.SetFinalizer(stf, finalize)
	return
}

// Return a new SliceTransform that uses ISliceTransform
func NewSliceTransform(itf ISliceTransform) (stf *SliceTransform) {
	var citf unsafe.Pointer = nil

	if nil != itf {
		citf =ISliceTransformAddReference(itf)
	}
	cstf := C.SliceTransformNewSliceTransform(citf)
	return cstf.toSliceTransform()
}
