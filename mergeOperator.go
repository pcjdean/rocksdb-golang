// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package rocksdb

/*
#include "mergeOperator.h"
*/
import "C"

import (
	"unsafe"
	"runtime"
)

// The Merge Operator
//
// Essentially, a MergeOperator specifies the SEMANTICS of a merge, which only
// client knows. It could be numeric addition, list append, string
// concatenation, edit data structure, ... , anything.
// The library, on the other hand, is concerned with the exercise of this
// interface, at the right time (during get, iteration, compaction...)
//
// To use merge, the client needs to provide an object implementing one of
// the following interfaces:
//  a) AssociativeMergeOperator - for most simple semantics (always take
//    two values, and merge them into one value, which is then put back
//    into rocksdb); numeric addition and string concatenation are examples;
//
//  b) MergeOperator - the generic class for all the more abstract / complex
//    operations; one method (FullMerge) to merge a Put/Delete value with a
//    merge operand; and another method (PartialMerge) that merges multiple
//    operands together. this is especially useful if your key values have
//    complex structures but you would still like to support client-specific
//    incremental updates.
//
// AssociativeMergeOperator is simpler to implement. MergeOperator is simply
// more powerful.
//
// Refer to rocksdb-merge wiki for more details and example implementations.
//
type IMergeOperator interface {

	// The name of the MergeOperator. Used to check for MergeOperator
	// mismatches (i.e., a DB created with one MergeOperator is
	// accessed using a different MergeOperator)
	// TODO: the name is currently not stored persistently and thus
	//       no checking is enforced. Client is responsible for providing
	//       consistent MergeOperator between DB opens.
	Name() string

	// Gives the client a way to express the read -> modify -> write semantics
	// key:      (IN)    The key that's associated with this merge operation.
	//                   Client could multiplex the merge operator based on it
	//                   if the key space is partitioned and different subspaces
	//                   refer to different types of data which have different
	//                   merge operation semantics
	// existing: (IN)    null indicates that the key does not exist before this op
	// operand_list:(IN) the sequence of merge operations to apply, front() first.
	// new_value:(OUT)   Client is responsible for filling the merge result here
	// logger:   (IN)    Client could use this to log errors during merge.
	//
	// Return true on success.
	// All values passed in will be client-specific values. So if this method
	// returns false, it is because client specified bad data or there was
	// internal corruption. This will be treated as an error by the library.
	//
	// Also make use of the *logger for error messages.
	FullMerge(key []byte, exval []byte, opdlist [][]byte, logger *Logger) (suc bool, newval []byte)

	// This function performs merge(left_op, right_op)
	// when both the operands are themselves merge operation types
	// that you would have passed to a DB::Merge() call in the same order
	// (i.e.: DB::Merge(key,left_op), followed by DB::Merge(key,right_op)).
	//
	// PartialMerge should combine them into a single merge operation that is
	// saved into *new_value, and then it should return true.
	// *new_value should be constructed such that a call to
	// DB::Merge(key, *new_value) would yield the same result as a call
	// to DB::Merge(key, left_op) followed by DB::Merge(key, right_op).
	//
	// The default implementation of PartialMergeMulti will use this function
	// as a helper, for backward compatibility.  Any successor class of
	// MergeOperator should either implement PartialMerge or PartialMergeMulti,
	// although implementing PartialMergeMulti is suggested as it is in general
	// more effective to merge multiple operands at a time instead of two
	// operands at a time.
	//
	// If it is impossible or infeasible to combine the two operations,
	// leave new_value unchanged and return false. The library will
	// internally keep track of the operations, and apply them in the
	// correct order once a base-value (a Put/Delete/End-of-Database) is seen.
	//
	// TODO: Presently there is no way to differentiate between error/corruption
	// and simply "return false". For now, the client should simply return
	// false in any case it cannot perform partial-merge, regardless of reason.
	// If there is corruption in the data, handle it in the FullMerge() function,
	// and return false there.  The default implementation of PartialMerge will
	// always return false.
	PartialMerge(key []byte, leftopd []byte, rightopd []byte, logger *Logger) (suc bool, newval []byte)

	// This function performs merge when all the operands are themselves merge
	// operation types that you would have passed to a DB::Merge() call in the
	// same order (front() first)
	// (i.e. DB::Merge(key, operand_list[0]), followed by
	//  DB::Merge(key, operand_list[1]), ...)
	//
	// PartialMergeMulti should combine them into a single merge operation that is
	// saved into *new_value, and then it should return true.  *new_value should
	// be constructed such that a call to DB::Merge(key, *new_value) would yield
	// the same result as subquential individual calls to DB::Merge(key, operand)
	// for each operand in operand_list from front() to back().
	//
	// The PartialMergeMulti function will be called only when the list of
	// operands are long enough. The minimum amount of operands that will be
	// passed to the function are specified by the "min_partial_merge_operands"
	// option.
	//
	// In the default implementation, PartialMergeMulti will invoke PartialMerge
	// multiple times, where each time it only merges two operands.  Developers
	// should either implement PartialMergeMulti, or implement PartialMerge which
	// is served as the helper function of the default PartialMergeMulti.
	PartialMergeMulti(key []byte, opdlist [][]byte, logger *Logger) (suc bool, newval []byte)
}

// Wrap functions for IMergeOperator

//export IMergeOperatorFullMerge
func IMergeOperatorFullMerge(cmop unsafe.Pointer, key *C.Slice_t, exval *C.Slice_t, opdlist *C.StringDeque_t, cnewval *C.String_t, clogger *C.Logger_t) C.bool {
	mop := InterfacesGet(cmop).(IMergeOperator)
	logger := clogger.toLogger(false)
	suc, newval := mop.FullMerge(key.cToBytes(false), exval.cToBytes(false), opdlist.toBytesArray(), logger)
	if suc {
		cnewval.setBytes(newval)
	}
	return toCBool(suc)
}

//export IMergeOperatorPartialMerge
func IMergeOperatorPartialMerge(cmop unsafe.Pointer, key *C.Slice_t, leftopd *C.Slice_t, rightopd *C.Slice_t, cnewval *C.String_t, clogger *C.Logger_t) C.bool {
	mop := InterfacesGet(cmop).(IMergeOperator)
	logger := clogger.toLogger(false)
	suc, newval := mop.PartialMerge(key.cToBytes(false), leftopd.cToBytes(false), rightopd.cToBytes(false), logger)
	if suc {
		cnewval.setBytes(newval)
	}
	return toCBool(suc)
}

//export IMergeOperatorPartialMergeMulti
func IMergeOperatorPartialMergeMulti(cmop unsafe.Pointer, key *C.Slice_t, opdlist *C.SliceDeque_t, cnewval *C.String_t, clogger *C.Logger_t) C.bool {
	mop := InterfacesGet(cmop).(IMergeOperator)
	logger := clogger.toLogger(false)
	suc, newval := mop.PartialMergeMulti(key.cToBytes(false), opdlist.toBytesArray(), logger)
	if suc {
		cnewval.setBytes(newval)
	}
	return toCBool(suc)
}

//export IMergeOperatorName
func IMergeOperatorName(cmop unsafe.Pointer) *C.char {
	mop := InterfacesGet(cmop).(IMergeOperator)
	return C.CString(mop.Name())
}

type IAssociativeMergeOperator interface {
	// Inherit from IMergeOperator
	IMergeOperator

	// Gives the client a way to express the read -> modify -> write semantics
	// key:           (IN) The key that's associated with this merge operation.
	// existing_value:(IN) null indicates the key does not exist before this op
	// value:         (IN) the value to update/merge the existing_value with
	// new_value:    (OUT) Client is responsible for filling the merge result here
	// logger:        (IN) Client could use this to log errors during merge.
	//
	// Return true on success.
	// All values passed in will be client-specific values. So if this method
	// returns false, it is because client specified bad data or there was
	// internal corruption. The client should assume that this will be treated
	// as an error by the library.
	Merge(key []byte, exval []byte, val []byte, logger *Logger) (suc bool, newval []byte)
}

//export IAssociativeMergeOperatorMerge
func IAssociativeMergeOperatorMerge(cmop unsafe.Pointer, key *C.Slice_t, exval *C.Slice_t, val *C.Slice_t, cnewval *C.String_t, clogger *C.Logger_t) C.bool {
	mop := InterfacesGet(cmop).(IAssociativeMergeOperator)
	logger := clogger.toLogger(false)
	suc, newval := mop.Merge(key.cToBytes(false), exval.cToBytes(false), val.cToBytes(false), logger)
	if suc {
		cnewval.setBytes(newval)
	}
	return toCBool(suc)
}

// Wrap go MergeOperator
type MergeOperator struct {
	mop C.PMergeOperator_t
}

// Release resources
func (mop *MergeOperator) finalize() {
	var cmop *C.PMergeOperator_t= &mop.mop
	C.DeletePMergeOperatorT(cmop, toCBool(false))
}

// C MergeOperator to go MergeOperator
func (cmop *C.PMergeOperator_t) toMergeOperator() (mop *MergeOperator) {
	mop = &MergeOperator{mop: *cmop}	
	runtime.SetFinalizer(mop, finalize)
	return
}

// Return a new MergeOperator that uses IMergeOperator
func NewMergeOperator(itf IMergeOperator) (mop *MergeOperator) {
	var citf unsafe.Pointer = nil

	if nil != itf {
		citf = InterfacesAddReference(itf)
	}
	cmop := C.NewMergeOperator(citf)
	return cmop.toMergeOperator()
}
