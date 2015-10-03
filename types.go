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

type SequenceNumber uint64

const (
	// Max array dimension
	arrayDimenMax = 0xFFFFFFFF

	// The initial size of callbackInterfaces
	initialInterfacesSize = 100
)

var (
	// Map to keep all the Interfaces callbacks from garbage collected
	callbackInterfaces map[unsafe.Pointer]interface{}

	// Mutext to protect callbackInterfaces
	callbackInterfacesMutex sync.Mutex = sync.Mutex{}
)

// Interface to release C pointer
type finalizer interface {
	finalize()
}

// Called by go finalizer
func finalize(obj finalizer) {
	obj.finalize()
}

//export InterfacesRemoveReference
// Remove interface citf from the callbackInterfaces 
// to leave citf garbage collected
func InterfacesRemoveReference(citf unsafe.Pointer) {
	defer callbackInterfacesMutex.Unlock()
	callbackInterfacesMutex.Lock()
	if nil != callbackInterfaces {
		delete(callbackInterfaces, citf)
	} else {
		log.Println("InterfacesRemoveReference: callbackInterfaces is not created!")
	}
}

// Get interface itf from the callbackInterfaces
// with citf as the key
func InterfacesGet(citf unsafe.Pointer) (itf interface{}) {
	defer callbackInterfacesMutex.Unlock()
	callbackInterfacesMutex.Lock()
	if nil != callbackInterfaces {
		itf = callbackInterfaces[citf]
	} else {
		log.Println("InterfacesGet: callbackInterfaces is not created!")
	}
	return
}

// Add interface itf to the callbackInterfaces to keep itf alive
// Return the key of the Interfaces in map callbackInterfaces
func InterfacesAddReference(itf interface{}) (citf unsafe.Pointer) {
	defer callbackInterfacesMutex.Unlock()
	callbackInterfacesMutex.Lock()
	if nil == callbackInterfaces {
		callbackInterfaces = make(map[unsafe.Pointer]interface{}, initialInterfacesSize)
	}
	citf = unsafe.Pointer(&itf)
	callbackInterfaces[citf] = itf
	return
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
