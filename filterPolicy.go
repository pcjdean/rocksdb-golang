// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// A database can be configured with a custom FilterPolicy object.
// This object is responsible for creating a small filter from a set
// of keys.  These filters are stored in rocksdb and are consulted
// automatically by rocksdb to decide whether or not to read some
// information from disk. In many cases, a filter can cut down the
// number of disk seeks form a handful to a single disk seek per
// DB::Get() call.
//
// Most people will want to use the builtin bloom filter support (see
// NewBloomFilterPolicy() below).

package rocksdb

/*
#include "filterPolicy.h"
#include "slice.h"
#include "cstring.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
	"log"
	"sync"
)

const (
	// The initial size of callbackIFilterPolicy
	initialIFilterPolicySize = 100
)

var (
	// Map to keep all the IFilterPolicy callbacks from garbage collected
	callbackIFilterPolicy map[unsafe.Pointer]IFilterPolicy

	// Mutext to protect callbackIFilterPolicy
	callbackIFilterPolicyMutex sync.Mutex = sync.Mutex{}
)

//export IFilterPolicyRemoveReference
// Remove interface citf from the callbackIFilterPolicy 
// to leave citf garbage collected
func IFilterPolicyRemoveReference(citf unsafe.Pointer) {
	defer callbackIFilterPolicyMutex.Unlock()
	callbackIFilterPolicyMutex.Lock()
	if nil != callbackIFilterPolicy {
		delete(callbackIFilterPolicy, citf)
	} else {
		log.Println("IFilterPolicyRemoveReference: callbackIFilterPolicy is not created!")
	}
}

// Get interface itf from the callbackIFilterPolicy
// with citf as the key
func IFilterPolicyGet(citf unsafe.Pointer) (itf IFilterPolicy) {
	defer callbackIFilterPolicyMutex.Unlock()
	callbackIFilterPolicyMutex.Lock()
	if nil != callbackIFilterPolicy {
		itf = callbackIFilterPolicy[citf]
	} else {
		log.Println("IFilterPolicyGet: callbackIFilterPolicy is not created!")
	}
	return
}

// Add interface itf to the callbackIFilterPolicy to keep itf alive
// Return the key of the IFilterPolicy in map callbackIFilterPolicy
func IFilterPolicyAddReference(itf IFilterPolicy) (citf unsafe.Pointer) {
	defer callbackIFilterPolicyMutex.Unlock()
	callbackIFilterPolicyMutex.Lock()
	if nil == callbackIFilterPolicy {
		callbackIFilterPolicy = make(map[unsafe.Pointer]IFilterPolicy, initialIFilterPolicySize)
	}
	citf = unsafe.Pointer(&itf)
	callbackIFilterPolicy[citf] = itf
	return
}

// A class that takes a bunch of keys, then generates filter
type IFilterBitsBuilder interface {

	// Add Key to filter, you could use any way to store the key.
	// Such as: storing hashes or original keys
	// Keys are in sorted order and duplicated keys are possible.
	AddKey(key []byte)

	// Generate the filter using the keys that are added
	// The return value of this function would be the filter bits,
	// The ownership of actual data is set to buf
	Finish() []byte
}

// A class that checks if a key can be in filter
// It should be initialized by Slice generated by BitsBuilder
type IFilterBitsReader interface {

	// Check if the entry match the bits in filter
	MayMatch(entry []byte) bool
}

// go interface IFilterPolicy
// We add a new format of filter block called full filter block
// This new interface gives you more space of customization
//
// For the full filter block, you can plug in your version by implement
// the IFilterBitsBuilder and IFilterBitsReader
//
// There are two sets of interface in FilterPolicy
// Set 1: CreateFilter, KeyMayMatch: used for blockbased filter
// Set 2: GetFilterBitsBuilder, GetFilterBitsReader, they are used for
// full filter.
// Set 1 MUST be implemented correctly, Set 2 is optional
// RocksDB would first try using functions in Set 2. if they return nullptr,
// it would use Set 1 instead.
// You can choose filter type in NewBloomFilterPolicy
type IFilterPolicy interface {
	// Return the name of this policy.  Note that if the filter encoding
	// changes in an incompatible way, the name returned by this method
	// must be changed.  Otherwise, old incompatible filters may be
	// passed to methods of this type.
	Name() string

	// keys[0,n-1] contains a list of keys (potentially with duplicates)
	// that are ordered according to the user supplied comparator.
	// Append a filter that summarizes keys[0,n-1] to a internal *dst.
	CreateFilter(keys [][]byte) []byte

	// "filter" contains the data appended by a preceding call to
	// CreateFilter() on this class.  This method must return true if
	// the key was in the list of keys passed to CreateFilter().
	// This method may return true or false if the key was not on the
	// list, but it should aim to return false with a high probability.
	KeyMayMatch(key, filter []byte) bool

	// Get the FilterBitsBuilder, which is ONLY used for full filter block
	// It contains interface to take individual key, then generate filter
	GetFilterBitsBuilder() *IFilterBitsBuilder

	// Get the FilterBitsReader, which is ONLY used for full filter block
	// It contains interface to tell if key can be in filter
	// The input slice should NOT be deleted by FilterPolicy
	GetFilterBitsReader() *IFilterBitsReader
}

// Wrap functions for IFilterPolicy

//export IFilterPolicyName
func IFilterPolicyName(cflp unsafe.Pointer) *C.char {
	flp := IFilterPolicyGet(cflp)
	return C.CString(flp.Name())
}

//export IFilterPolicyCreateFilter
func IFilterPolicyCreateFilter(cflp unsafe.Pointer, ckeys *C.Slice_t, sz C.int) C.String_t {
	flp := IFilterPolicyGet(cflp)
	keys := newBytesFromCSliceArray(ckeys, uint(sz), false, false)
	filter := string(flp.CreateFilter(keys))
	str := newCStringFromString(&filter)
	return str.str
}

//export IFilterPolicyKeyMayMatch
func IFilterPolicyKeyMayMatch(cflp unsafe.Pointer, key, filter *C.Slice_t) C.bool {
	flp := IFilterPolicyGet(cflp)
	return toCBool(flp.KeyMayMatch(key.cToBytes(false), filter.cToBytes(false)))
}

//export IFilterPolicyGetFilterBitsBuilder
func IFilterPolicyGetFilterBitsBuilder(cflp unsafe.Pointer) unsafe.Pointer {
	flp := IFilterPolicyGet(cflp)
	return unsafe.Pointer(flp.GetFilterBitsBuilder())
}

//export IFilterPolicyGetFilterBitsReader
func IFilterPolicyGetFilterBitsReader(cflp unsafe.Pointer) unsafe.Pointer {
	flp := IFilterPolicyGet(cflp)
	return unsafe.Pointer(flp.GetFilterBitsReader())
}

// Wrap go FilterPolicy
type FilterPolicy struct {
	flp C.PFilterPolicy_t
}

// Release resources
func (flp *FilterPolicy) finalize() {
	var cflp *C.PFilterPolicy_t= &flp.flp
	C.DeletePFilterPolicyT(cflp, toCBool(false))
}

// C filterPolicy to go filterPolicy
func (cflp *C.PFilterPolicy_t) toFilterPolicy() (flp *FilterPolicy) {
	flp = &FilterPolicy{flp: *cflp}	
	runtime.SetFinalizer(flp, finalize)
	return
}

// Return a new filter policy that uses IFilterPolicy
func NewFilterPolicy(itf IFilterPolicy) (flp *FilterPolicy) {
	var iftp unsafe.Pointer = nil

	if nil != itf {
		iftp =IFilterPolicyAddReference(itf)
	}
	cflp := C.NewPFilterPolicy(iftp)
	return cflp.toFilterPolicy()
}

// Return a new filter policy that uses a bloom filter with approximately
// the specified number of bits per key.
//
// bits_per_key: bits per key in bloom filter. A good value for bits_per_key
// is 10, which yields a filter with ~ 1% false positive rate.
// use_block_based_builder: use block based filter rather than full fiter.
// If you want to builder full filter, it needs to be set to false.
//
// Callers must delete the result after any database that is using the
// result has been closed.
//
// Note: if you are using a custom comparator that ignores some parts
// of the keys being compared, you must not use NewBloomFilterPolicy()
// and must provide your own FilterPolicy that also ignores the
// corresponding parts of the keys.  For example, if the comparator
// ignores trailing spaces, it would be incorrect to use a
// FilterPolicy (like NewBloomFilterPolicy) that does not ignore
// trailing spaces in keys.
func NewBloomFilterPolicy(bitsPerKey int, useBlockBasedBuilder ...bool) (flp *FilterPolicy) {
	var blockBased bool = true

	if useBlockBasedBuilder != nil {
		blockBased = useBlockBasedBuilder[0]
	}

	cflp := C.NewPFilterPolicyTRawArgs(C.int(bitsPerKey), toCBool(blockBased))
	return cflp.toFilterPolicy()
}
