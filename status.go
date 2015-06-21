// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

package rocksdb

/*
#include "status.h"
*/
import "C"

type Status struct {
	sta C.Status_t
}

func (stat *Status) Finalize() {
	var cstat *C.Status_t = unsafe.Pointer(&stat.sta)
	C.DeleteStatusT(cstat, false)
}

// Returns true iff the status indicates success.
func (stat *Status) Ok() bool {
	var cstat *C.Status_t = unsafe.Pointer(&stat.sta)
	return StatusOk(cstat)
}

// Returns true iff the status indicates a NotFound error.
func (stat *Status) IsNotFound() bool {
	var cstat *C.Status_t = unsafe.Pointer(&stat.sta)
	return StatusIsNotFound(cstat)
}

// Returns true iff the status indicates a Corruption error.
func (stat *Status) IsCorruption() bool {
	var cstat *C.Status_t = unsafe.Pointer(&stat.sta)
	return StatusIsCorruption(cstat)
}

// Returns true iff the status indicates a NotSupported error.
func (stat *Status) IsNotSupported() bool {
	var cstat *C.Status_t = unsafe.Pointer(&stat.sta)
	return StatusIsNotSupported(cstat)
}

// Returns true iff the status indicates an IOError.
func (stat *Status) IsInvalidArgument() bool {
	var cstat *C.Status_t = unsafe.Pointer(&stat.sta)
	return StatusIsInvalidArgument(cstat)
}

// Returns true iff the status indicates an MergeInProgress.
func (stat *Status) IsMergeInProgress() bool {
	var cstat *C.Status_t = unsafe.Pointer(&stat.sta)
	return StatusIsMergeInProgress(cstat)
}

// Returns true iff the status indicates Incomplete
func (stat *Status) IsIncomplete() bool {
	var cstat *C.Status_t = unsafe.Pointer(&stat.sta)
	return StatusIsIncomplete(cstat)
}

// Returns true iff the status indicates Shutdown In progress
func (stat *Status) IsShutdownInProgress() bool {
	var cstat *C.Status_t = unsafe.Pointer(&stat.sta)
	return StatusIsShutdownInProgress(cstat)
}

func (stat *Status) IsTimedOut() bool {
	var cstat *C.Status_t = unsafe.Pointer(&stat.sta)
	return StatusIsTimedOut(cstat)
}

func (stat *Status) IsAborted() bool {
	var cstat *C.Status_t = unsafe.Pointer(&stat.sta)
	return StatusIsAborted(cstat)
}

// Returns true iff the status indicates that a resource is Busy and
// temporarily could not be acquired.
func (stat *Status) IsBusy() bool {
	var cstat *C.Status_t = unsafe.Pointer(&stat.sta)
	return StatusIsBusy(cstat)
}

// Return a string representation of this status suitable for printing.
// Returns the string "OK" for success.
func (stat *Status) ToString() String {
	var cstat *C.Status_t = unsafe.Pointer(&stat.sta)
	str := String{StatusToString(cstat)}
	return str.GoString(false)
}
