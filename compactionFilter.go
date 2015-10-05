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
)

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
	Filter(level int, key, exval []byte) (newval []byte, valchg bool, removed bool)

	// Returns a name that identifies this compaction filter.
	// The name will be printed to LOG file on start up for diagnosis.
	Name() string
}

// Wrap functions for ICompactionFilter

//export ICompactionFilterName
func ICompactionFilterName(ccpf unsafe.Pointer) *C.char {
	cpf := InterfacesGet(ccpf).(ICompactionFilter)
	return C.CString(cpf.Name())
}

//export ICompactionFilterFilter
func ICompactionFilterFilter(ccpf unsafe.Pointer, level C.int, key, exval *C.Slice_t, newval *C.String_t, valchg *C.bool) C.bool {
	cpf := InterfacesGet(ccpf).(ICompactionFilter)
	gnewval, gvalchg, removed := cpf.Filter(int(level), key.cToBytes(false), exval.cToBytes(false))
	*valchg = toCBool(gvalchg)
	if gvalchg {
		newval.setBytes(gnewval)
	}
	return toCBool(removed)
}

// Wrap go CompactionFilter
type CompactionFilter struct {
	cpf C.CompactionFilter_t
	// True if the CompactionFilter is closed
	closed bool
}

// Release resources
func (cpf *CompactionFilter) finalize() {
	if !cpf.closed {
		cpf.closed = true
		var ccpf *C.CompactionFilter_t= &cpf.cpf
		C.DeleteCompactionFilterT(ccpf, toCBool(false))
	}
}

// Close the @cpf
func (cpf *CompactionFilter) Close() {
	runtime.SetFinalizer(cpf, nil)
	cpf.finalize()
}

// C CompactionFilter to go CompactionFilter
func (ccpf *C.CompactionFilter_t) toCompactionFilter() (cpf *CompactionFilter) {
	cpf = &CompactionFilter{cpf: *ccpf}	
	runtime.SetFinalizer(cpf, finalize)
	return
}

// Return a new default CompactionFilter
func NewDefaultCompactionFilter() (cpf *CompactionFilter) {
	cpf = &CompactionFilter{cpf: C.CompactionFilter_t{nil}}	
	return
}

// Return a new CompactionFilter that uses ICompactionFilter
func NewCompactionFilter(itf ICompactionFilter) (cpf *CompactionFilter) {
	var iftp unsafe.Pointer = nil

	if nil != itf {
		iftp =InterfacesAddReference(itf)
	}
	ccpf := C.NewCompactionFilter(iftp)
	return ccpf.toCompactionFilter()
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
	Filter(level int, keys, exvals [][]byte) (newvals [][]byte, valchgs []bool, removed []bool)

	// Returns a name that identifies this compaction filter.
	// The name will be printed to LOG file on start up for diagnosis.
	Name() string
}

// Wrap functions for ICompactionFilterV2

//export ICompactionFilterV2Name
func ICompactionFilterV2Name(ccpf unsafe.Pointer) *C.char {
	cpf := InterfacesGet(ccpf).(ICompactionFilterV2)
	return C.CString(cpf.Name())
}

//export ICompactionFilterV2Filter
func ICompactionFilterV2Filter(ccpf unsafe.Pointer, level C.int, keys, exvals *C.SliceVector_t, newvals *C.StringVector_t, valchgs *C.BoolVector_t, removeds *C.BoolVector_t) {
	cpf := InterfacesGet(ccpf).(ICompactionFilterV2)
	gnewvals, gvalchgs, gremoveds := cpf.Filter(int(level), keys.toBytesArray(), exvals.toBytesArray())
	newvals.setBytesArray(gnewvals)
	valchgs.setBoolArray(gvalchgs)
	removeds.setBoolArray(gremoveds)
}

// Wrap go CompactionFilter_Context
type CompactionFilter_Context struct {
	cfc C.CompactionFilter_Context_t
}

// C CompactionFilter_Context to go CompactionFilter_Context
func (ccfc *C.CompactionFilter_Context_t) toCompactionFilter_Context() (cfc *CompactionFilter_Context) {
	cfc = &CompactionFilter_Context{cfc: *ccfc}	
	return
}

// Each compaction will create a new CompactionFilter allowing the
// application to know about different compactions
type ICompactionFilterFactory interface {

	// Create a ICompactionFilter
	CreateCompactionFilter(context *CompactionFilter_Context) ICompactionFilter

	// Returns a name that identifies this compaction filter factory.
	Name() string
}

// Wrap functions for ICompactionFilterFactory

//export ICompactionFilterFactoryName
func ICompactionFilterFactoryName(ccpf unsafe.Pointer) *C.char {
	cpf := InterfacesGet(ccpf).(ICompactionFilterFactory)
	return C.CString(cpf.Name())
}

//export ICompactionFilterFactoryCreateCompactionFilter
func ICompactionFilterFactoryCreateCompactionFilter(ccpf unsafe.Pointer, context *C.CompactionFilter_Context_t) (filter unsafe.Pointer) {
	cpf := InterfacesGet(ccpf).(ICompactionFilterFactory)
	filter = InterfacesAddReference(cpf.CreateCompactionFilter(context.toCompactionFilter_Context()))
	return
}

// Wrap go CompactionFilterFactory
type CompactionFilterFactory struct {
	cff C.PCompactionFilterFactory_t
}

// Release resources
func (cff *CompactionFilterFactory) finalize() {
	var ccff *C.PCompactionFilterFactory_t= &cff.cff
	C.DeletePCompactionFilterFactoryT(ccff, toCBool(false))
}

// C CompactionFilterFactory to go CompactionFilterFactory
func (ccff *C.PCompactionFilterFactory_t) toCompactionFilterFactory() (cff *CompactionFilterFactory) {
	cff = &CompactionFilterFactory{cff: *ccff}	
	runtime.SetFinalizer(cff, finalize)
	return
}

// Return a new default CompactionFilterFactory
func NewDefaultCompactionFilterFactory() (cff *CompactionFilterFactory) {
	cff = &CompactionFilterFactory{cff: C.NewPCompactionFilterFactoryTDefault()}	
	runtime.SetFinalizer(cff, finalize)
	return
}

// Return a new CompactionFilterFactory that uses ICompactionFilterFactory
func NewCompactionFilterFactory(itf ICompactionFilterFactory) (cff *CompactionFilterFactory) {
	var iftp unsafe.Pointer = nil

	if nil != itf {
		iftp =InterfacesAddReference(itf)
	}
	ccff := C.NewPCompactionFilterFactory(iftp)
	return ccff.toCompactionFilterFactory()
}

// Wrap go CompactionFilterContext
type CompactionFilterContext struct {
	cfc C.CompactionFilterContext_t
}

// C CompactionFilterContext_t to go CompactionFilterContext
func (ccfc *C.CompactionFilterContext_t) toCompactionFilterContext() (cfc *CompactionFilterContext) {
	cfc = &CompactionFilterContext{cfc: *ccfc}	
	return
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
	CreateCompactionFilterV2(context *CompactionFilterContext) ICompactionFilterV2

	// Returns a name that identifies this compaction filter factory.
	Name() string

	// Return the prefix extractor
	// This is not virtual function in c++ head.
	// It's not called with the newly set prefix extractor by 
	// SetPrefixExtractor
	GetPrefixExtractor() ISliceTransform

	// Set the prefix extractor
	// There is no effect once underlying CompactionFilterFactoryV2 is created.
	// See notes on GetPrefixExtractor(). Create a new CompactionFilterFactoryV2
	// if we want to use a new ISliceTransform. Set ISliceTransform somewhere 
	// internaly to keep @prextrc from garbage collected.
	SetPrefixExtractor(prextrc ISliceTransform)
}

// Wrap functions for ICompactionFilterFactoryV2

//export ICompactionFilterFactoryV2Name
func ICompactionFilterFactoryV2Name(ccpf unsafe.Pointer) *C.char {
	cpf := InterfacesGet(ccpf).(ICompactionFilterFactoryV2)
	return C.CString(cpf.Name())
}

//export ICompactionFilterFactoryV2CreateCompactionFilterV2
func ICompactionFilterFactoryV2CreateCompactionFilterV2(ccpf unsafe.Pointer, context *C.CompactionFilterContext_t) (filter unsafe.Pointer) {
	cpf := InterfacesGet(ccpf).(ICompactionFilterFactoryV2)
	filter = InterfacesAddReference(cpf.CreateCompactionFilterV2(context.toCompactionFilterContext()))
	return
}

//export ICompactionFilterFactoryV2GetPrefixExtractor
func ICompactionFilterFactoryV2GetPrefixExtractor(ccpf unsafe.Pointer) (cstf C.SliceTransform_t) {
	cpf := InterfacesGet(ccpf).(ICompactionFilterFactoryV2)
	stf := NewSliceTransform(cpf.GetPrefixExtractor())
	cstf = stf.stf
	return
}

// Wrap go CompactionFilterFactoryV2
type CompactionFilterFactoryV2 struct {
	cff C.PCompactionFilterFactoryV2_t
}

// Release resources
func (cff *CompactionFilterFactoryV2) finalize() {
	var ccff *C.PCompactionFilterFactoryV2_t= &cff.cff
	C.DeletePCompactionFilterFactoryV2T(ccff, toCBool(false))
}

// C CompactionFilterFactoryV2 to go CompactionFilterFactoryV2
func (ccff *C.PCompactionFilterFactoryV2_t) toCompactionFilterFactoryV2() (cff *CompactionFilterFactoryV2) {
	cff = &CompactionFilterFactoryV2{cff: *ccff}	
	runtime.SetFinalizer(cff, finalize)
	return
}

// Return a new default CompactionFilterFactoryV2
func NewDefaultCompactionFilterFactoryV2() (cff *CompactionFilterFactoryV2) {
	cff = &CompactionFilterFactoryV2{cff: C.NewPCompactionFilterFactoryV2TDefault()}	
	runtime.SetFinalizer(cff, finalize)
	return
}

// Return a new CompactionFilterFactoryV2 that uses ICompactionFilterFactoryV2
func NewCompactionFilterFactoryV2(itf ICompactionFilterFactoryV2, sitf ISliceTransform) (cff *CompactionFilterFactoryV2) {
	var iftp unsafe.Pointer = nil

	if nil != itf {
		iftp =InterfacesAddReference(itf)
		itf.SetPrefixExtractor(sitf)
	}
	// SliceTransform of CompactionFilterFactoryV2 will be initialised
	// from GetPrefixExtractor of ICompactionFilterFactoryV2
	ccff := C.NewPCompactionFilterFactoryV2(iftp, nil)
	return ccff.toCompactionFilterFactoryV2()
}
