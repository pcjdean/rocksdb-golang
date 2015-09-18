// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package rocksdb

/*
#include "compactionFilter.h"
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

// CompactionFilter allows an application to modify/delete a key-value at
// the time of compaction.
type ICompactionFilter interface {

	// The compaction process invokes this
	// method for kv that is being compacted. A return value
	// of false indicates that the kv should be preserved in the
	// output of this compaction run and a return value of true
	// indicates that this key-value should be removed from the
	// output of the compaction.  The application can inspect
	// the existing value of the key and make decision based on it.
	//
	// When the value is to be preserved, the application has the option
	// to modify the existing_value and pass it back through new_value.
	// value_changed needs to be set to true in this case.
	//
	// If multithreaded compaction is being used *and* a single CompactionFilter
	// instance was supplied via Options::compaction_filter, this method may be
	// called from different threads concurrently.  The application must ensure
	// that the call is thread-safe.
	//
	// If the CompactionFilter was created by a factory, then it will only ever
	// be used by a single thread that is doing the compaction run, and this
	// call does not need to be thread-safe.  However, multiple filters may be
	// in existence and operating concurrently.
	Filter(level int, key, exval []byte) (newval []byte, valchang bool, removed bool)

	// Returns a name that identifies this compaction filter.
	// The name will be printed to LOG file on start up for diagnosis.
	Name() string
}

// CompactionFilterV2 that buffers kv pairs sharing the same prefix and let
// application layer to make individual decisions for all the kv pairs in the
// buffer.
type ICompactionFilterV2 interface {

	// The compaction process invokes this method for all the kv pairs
	// sharing the same prefix. It is a "roll-up" version of CompactionFilter.
	//
	// Each entry in the return vector indicates if the corresponding kv should
	// be preserved in the output of this compaction run. The application can
	// inspect the existing values of the keys and make decision based on it.
	//
	// When a value is to be preserved, the application has the option
	// to modify the entry in existing_values and pass it back through an entry
	// in new_values. A corresponding values_changed entry needs to be set to
	// true in this case. Note that the new_values vector contains only changed
	// values, i.e. new_values.size() <= values_changed.size().
	//
	Filter(level int, keys, exvals [][]byte) (newvals [][]byte, valchangs []bool, removed []bool)

	// Returns a name that identifies this compaction filter.
	// The name will be printed to LOG file on start up for diagnosis.
	Name() string
}

// Each compaction will create a new CompactionFilter allowing the
// application to know about different compactions
type ICompactionFilterFactory interface {

	// Create a ICompactionFilter
	CreateCompactionFilter(context *CompactionFilter_Context) *ICompactionFilter

	// Returns a name that identifies this compaction filter factory.
	Name() string
}

// Each compaction will create a new CompactionFilterV2
//
// CompactionFilterFactoryV2 enables application to specify a prefix and use
// CompactionFilterV2 to filter kv-pairs in batches. Each batch contains all
// the kv-pairs sharing the same prefix.
//
// This is useful for applications that require grouping kv-pairs in
// compaction filter to make a purge/no-purge decision. For example, if the
// key prefix is user id and the rest of key represents the type of value.
// This batching filter will come in handy if the application's compaction
// filter requires knowledge of all types of values for any user id.
//
type ICompactionFilterFactoryV2 interface {

	// Create a ICreateCompactionFilterV2
	CreateCompactionFilterV2(context *CompactionFilterContext) *ICreateCompactionFilterV2

	// Returns a name that identifies this compaction filter factory.
	Name() string

	// Return the prefix extractor
	GetPrefixExtractor() *SliceTransform

	// Set the prefix extractor
	SetPrefixExtractor(pextr *SliceTransform)
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