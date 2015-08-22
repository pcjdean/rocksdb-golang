// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

package rocksdb

/*
#include "status.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// Go Status
type Status struct {
	sta C.Status_t
}

// Release resources
func (stat *Status) finalize() {
	var cstat *C.Status_t = &stat.sta
	C.DeleteStatusT(cstat, toCBool(false))
}

// C Status to Go Status
func (csta *C.Status_t) toStatus() (sta *Status) {
	sta = &Status{sta: *csta}	
	runtime.SetFinalizer(sta, finalize)
	return
}

// Create a new 'DB closed' go status
func NewDBClosedStatus() *Status {
	csta := C.StatusDBClosedStatus()
	return csta.toStatus()
}

// C Status array to Go Status array
func newStatusArrayFromCArray(csta *C.Status_t, sz uint) (stas []*Status) {
	defer C.DeleteStatusTArray(csta)
	stas = make([]*Status, sz)
	for i, _ := range stas {
		stas[i] = &Status{sta: (*[arrayDimenMax]C.Status_t)(unsafe.Pointer(csta))[i]}	
		runtime.SetFinalizer(stas[i], finalize)
	}
	return
}

// Returns true iff the status indicates success.
func (stat *Status) Ok() bool {
	var cstat *C.Status_t = &stat.sta
	return C.StatusOk(cstat).toBool()
}

// Returns true iff the status indicates a NotFound error.
func (stat *Status) IsNotFound() bool {
	var cstat *C.Status_t = &stat.sta
	return C.StatusIsNotFound(cstat).toBool()
}

// Returns true iff the status indicates a Corruption error.
func (stat *Status) IsCorruption() bool {
	var cstat *C.Status_t = &stat.sta
	return C.StatusIsCorruption(cstat).toBool()
}

// Returns true iff the status indicates a NotSupported error.
func (stat *Status) IsNotSupported() bool {
	var cstat *C.Status_t = &stat.sta
	return C.StatusIsNotSupported(cstat).toBool()
}

// Returns true iff the status indicates an IOError.
func (stat *Status) IsInvalidArgument() bool {
	var cstat *C.Status_t = &stat.sta
	return C.StatusIsInvalidArgument(cstat).toBool()
}

// Returns true iff the status indicates an MergeInProgress.
func (stat *Status) IsMergeInProgress() bool {
	var cstat *C.Status_t = &stat.sta
	return C.StatusIsMergeInProgress(cstat).toBool()
}

// Returns true iff the status indicates Incomplete
func (stat *Status) IsIncomplete() bool {
	var cstat *C.Status_t = &stat.sta
	return C.StatusIsIncomplete(cstat).toBool()
}

// Returns true iff the status indicates Shutdown In progress
func (stat *Status) IsShutdownInProgress() bool {
	var cstat *C.Status_t = &stat.sta
	return C.StatusIsShutdownInProgress(cstat).toBool()
}

func (stat *Status) IsTimedOut() bool {
	var cstat *C.Status_t = &stat.sta
	return C.StatusIsTimedOut(cstat).toBool()
}

func (stat *Status) IsAborted() bool {
	var cstat *C.Status_t = &stat.sta
	return C.StatusIsAborted(cstat).toBool()
}

// Returns true iff the status indicates that a resource is Busy and
// temporarily could not be acquired.
func (stat *Status) IsBusy() bool {
	var cstat *C.Status_t = &stat.sta
	return C.StatusIsBusy(cstat).toBool()
}

// Return a string representation of this status suitable for printing.
// Returns the string "OK" for success.
func (stat *Status) String() string {
	var cstat *C.Status_t = &stat.sta
	str := cString{str: C.StatusToString(cstat)}
	return str.goString(false)
}
