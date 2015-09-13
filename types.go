// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

package rocksdb

/*
#include "types.h"
*/
import "C"

import (
	"unsafe"
	"log"
	"sync"
)

var (
	// Map to keep all the go callbacks from garbage collected
	callbackInterfaceSet map[*interface{}]bool

	// Mutext to protect callbackInterfaceSet
	callbackInterfaceSetMutex sync.Mutex = sync.Mutex{}
)

type SequenceNumber uint64

const (
	// Max array dimension
	arrayDimenMax = 0xFFFFFFFF

	// The initial size of callbackInterfaceSet
	initialInterfaceSetSize = 100
)

// Interface to release C pointer
type finalizer interface {
	finalize()
}

// Called by go finalizer
func finalize(obj finalizer) {
	obj.finalize()
}

//export InterfaceRemoveReference
// Remove interface citf from the callbackInterfaceSet 
// to leave citf garbage collected
func InterfaceRemoveReference(citf unsafe.Pointer) {
	itf := (*interface{})(citf)
	defer callbackInterfaceSetMutex.Unlock()
	callbackInterfaceSetMutex.Lock()
	if nil != callbackInterfaceSet {
		delete(callbackInterfaceSet, itf)
	} else {
		log.Println("InterfaceRemoveReference: callbackInterfaceSet is not created!")
	}
}

//export InterfaceAddReference
// Add interface citf to the callbackInterfaceSet to keep itf alive
func InterfaceAddReference(citf unsafe.Pointer) {
	itf := (*interface{})(citf)
	defer callbackInterfaceSetMutex.Unlock()
	callbackInterfaceSetMutex.Lock()
	if nil == callbackInterfaceSet {
		callbackInterfaceSet = make(map[*interface{}]bool, initialInterfaceSetSize)
	}
	callbackInterfaceSet[itf] = true
}

// Convert C int64 array to go int64 array
func newUint64ArrayFromCArray(cuints *[]C.uint64_t) (vals []uint64) {
	vals = make([]uint64, len(*cuints))
	for i, v := range *cuints {
		vals[i] = uint64(v)
	}
	return
}

// Convert go int array to C int array
func newCIntArrayFromArray(vals *[]int) (cints []C.int) {
	cints = make([]C.int, len(*vals))
	for i, v := range *vals {
		cints[i] = C.int(v)
	}
	return
}

// Convert go bool to C bool
func toCBool(b bool) (ret C.bool) {
	if b {
		ret = C.true
	} else {
		ret = C.false
	}
	return
}

// Convert C bool to go bool
func (b C.bool) toBool() (ret bool)  {
	if b == C.true {
		ret = true
	} else {
		ret = false
	}
	return
}
