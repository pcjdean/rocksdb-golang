// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// A Cache is an interface that maps keys to values.  It has internal
// synchronization and may be safely accessed concurrently from
// multiple threads.  It may automatically evict entries to make room
// for new entries.  Values have a specified charge against the cache
// capacity.  For example, a cache where the values are variable
// length strings, may use the length of the string as the charge for
// the string.
//
// A builtin cache implementation with a least-recently-used eviction
// policy is provided.  Clients may use their own implementations if
// they want something more sophisticated (like scan-resistance, a
// custom eviction policy, variable cache sizing, etc.)

package rocksdb

/*
#include "table.h"
*/
import "C"

import (
	"runtime"
)

// Wrap go TableFactory
type TableFactory struct {
	tbf C.PTableFactory_t
}

// Release resources
func (tbf *TableFactory) finalize() {
	var ctbf *C.PTableFactory_t= &tbf.tbf
	C.DeletePTableFactoryT(ctbf, toCBool(false))
}

// C TableFactory to go TableFactory
func (ctbf *C.PTableFactory_t) toTableFactory() (tbf *TableFactory) {
	tbf = &TableFactory{tbf: *ctbf}	
	runtime.SetFinalizer(tbf, finalize)
	return
}

// Wrap go BlockBasedTableOptions
type BlockBasedTableOptions struct {
	btop C.BlockBasedTableOptions_t
	// true if btop is deleted
	closed bool
}

// Release resources
func (btop *BlockBasedTableOptions) finalize() {
	if !btop.closed {
		btop.closed = true
		var cbtop *C.BlockBasedTableOptions_t= &btop.btop
		C.DeleteBlockBasedTableOptionsT(cbtop, toCBool(false))
	}
}

// Close the @btop
func (btop *BlockBasedTableOptions) Close() {
	runtime.SetFinalizer(btop, nil)
	btop.finalize()
}

// C BlockBasedTableOptions to go BlockBasedTableOptions
func (cbtop *C.BlockBasedTableOptions_t) toBlockBasedTableOptions() (btop *BlockBasedTableOptions) {
	btop = &BlockBasedTableOptions{btop: *cbtop}	
	runtime.SetFinalizer(btop, finalize)
	return
}

// Create default BlockBasedTableOptions
func NewBlockBasedTableOptions() *BlockBasedTableOptions {
	cbtop := C.NewBlockBasedTableOptionsTDefault()
	return cbtop.toBlockBasedTableOptions()
}

// If non-NULL use the specified cache for blocks.
// If NULL, rocksdb will automatically create and use an 8MB internal cache.
func (btop *BlockBasedTableOptions) SetBlockCache(cache *Cache) {
	var cbtop *C.BlockBasedTableOptions_t= &btop.btop
	C.BlockBasedTableOptions_set_block_cache(cbtop, &cache.cache)
}

// If non-nullptr, use the specified filter policy to reduce disk reads.
// Many applications will benefit from passing the result of
// NewBloomFilterPolicy() here.
func (btop *BlockBasedTableOptions) SetFilterPolicy(flp *FilterPolicy) {
	var cbtop *C.BlockBasedTableOptions_t= &btop.btop
	if nil == flp {
		def := C.NewPFilterPolicyTDefault()
		flp = def.toFilterPolicy()
	}
	C.BlockBasedTableOptions_set_filter_policy(cbtop, &flp.flp)
}

// Create default block based table factory.
func (btop *BlockBasedTableOptions) NewBlockBasedTableFactory() *TableFactory {
	var cbtop *C.BlockBasedTableOptions_t= &btop.btop
	ctbf := C.NewBlockBasedTableFactory(cbtop)
	return ctbf.toTableFactory()
}

// Wrap go PlainTableOptions
type PlainTableOptions struct {
	ptop C.PlainTableOptions_t
	// true if ptop is deleted
	closed bool
}

// Release resources
func (ptop *PlainTableOptions) finalize() {
	if !ptop.closed {
		ptop.closed = true
		var cptop *C.PlainTableOptions_t= &ptop.ptop
		C.DeletePlainTableOptionsT(cptop, toCBool(false))
	}
}

// Close the @ptop
func (ptop *PlainTableOptions) Close() {
	runtime.SetFinalizer(ptop, nil)
	ptop.finalize()
}

// C PlainTableOptions to go PlainTableOptions
func (cptop *C.PlainTableOptions_t) toPlainTableOptions() (ptop *PlainTableOptions) {
	ptop = &PlainTableOptions{ptop: *cptop}	
	runtime.SetFinalizer(ptop, finalize)
	return
}

// Create default PlainTableOptions
func NewPlainTableOptions() *PlainTableOptions {
	cptop := C.NewPlainTableOptionsTDefault()
	return cptop.toPlainTableOptions()
}

// @user_key_len: plain table has optimization for fix-sized keys, which can
//                be specified via user_key_len.  Alternatively, you can pass
//                `kPlainTableVariableLength` if your keys have variable
//                lengths.
func (ptop *PlainTableOptions) SetUserKeyLen(keylen uint32) {
	var cptop *C.PlainTableOptions_t= &ptop.ptop
	C.PlainTableOptions_set_user_key_len(cptop, C.uint32_t(keylen))
}

// @bloom_bits_per_key: the number of bits used for bloom filer per prefix.
//                      You may disable it by passing a zero.
func (ptop *PlainTableOptions) SetBloomBitsPerKey(bits int) {
	var cptop *C.PlainTableOptions_t= &ptop.ptop
	C.PlainTableOptions_set_bloom_bits_per_key(cptop, C.int(bits))
}

// @hash_table_ratio: the desired utilization of the hash table used for
//                    prefix hashing.
//                    hash_table_ratio = number of prefixes / #buckets in the
//                    hash table
func (ptop *PlainTableOptions) SetHashTableRatio(ratio float64) {
	var cptop *C.PlainTableOptions_t= &ptop.ptop
	C.PlainTableOptions_set_hash_table_ratio(cptop, C.double(ratio))
}

// @index_sparseness: inside each prefix, need to build one index record for
//                    how many keys for binary search inside each hash bucket.
//                    For encoding type kPrefix, the value will be used when
//                    writing to determine an interval to rewrite the full
//                    key. It will also be used as a suggestion and satisfied
//                    when possible.
func (ptop *PlainTableOptions) SetIndexSparseness(sparseness uint64) {
	var cptop *C.PlainTableOptions_t= &ptop.ptop
	C.PlainTableOptions_set_index_sparseness(cptop, C.size_t(sparseness))
}

// @huge_page_tlb_size: if <=0, allocate hash indexes and blooms from malloc.
//                      Otherwise from huge page TLB. The user needs to
//                      reserve huge pages for it to be allocated, like:
//                          sysctl -w vm.nr_hugepages=20
//                      See linux doc Documentation/vm/hugetlbpage.txt
func (ptop *PlainTableOptions) SetHugePageTlbSize(sz uint64) {
	var cptop *C.PlainTableOptions_t= &ptop.ptop
	C.PlainTableOptions_set_huge_page_tlb_size(cptop, C.size_t(sz))
}

// @encoding_type: how to encode the keys. See enum EncodingType above for
//                 the choices. The value will determine how to encode keys
//                 when writing to a new SST file. This value will be stored
//                 inside the SST file which will be used when reading from
//                 the file, which makes it possible for users to choose
//                 different encoding type when reopening a DB. Files with
//                 different encoding types can co-exist in the same DB and
//                 can be read.
func (ptop *PlainTableOptions) SetEncodingType(etype byte) {
	var cptop *C.PlainTableOptions_t= &ptop.ptop
	C.PlainTableOptions_set_encoding_type(cptop, C.char(etype))
}

// @full_scan_mode: mode for reading the whole file one record by one without
//                  using the index.
func (ptop *PlainTableOptions) SetFullScanMode(mode bool) {
	var cptop *C.PlainTableOptions_t= &ptop.ptop
	C.PlainTableOptions_set_full_scan_mode(cptop, toCBool(mode))
}

// @store_index_in_file: compute plain table index and bloom filter during
//                       file building and store it in file. When reading
//                       file, index will be mmaped instead of recomputation.
func (ptop *PlainTableOptions) SetStoreIndexInFile(store bool) {
	var cptop *C.PlainTableOptions_t= &ptop.ptop
	C.PlainTableOptions_set_store_index_in_file(cptop, toCBool(store))
}

// Create plain table factory.
func (ptop *PlainTableOptions) NewPlainTableFactory() *TableFactory {
	var cptop *C.PlainTableOptions_t= &ptop.ptop
	ctbf := C.NewPlainTableFactory(cptop)
	return ctbf.toTableFactory()
}

// Wrap go CuckooTableOptions
type CuckooTableOptions struct {
	ctop C.CuckooTableOptions_t
	// true if ctop is deleted
	closed bool
}

// Release resources
func (ctop *CuckooTableOptions) finalize() {
	if !ctop.closed {
		ctop.closed = true
		var cctop *C.CuckooTableOptions_t= &ctop.ctop
		C.DeleteCuckooTableOptionsT(cctop, toCBool(false))
	}
}

// Close the @ctop
func (ctop *CuckooTableOptions) Close() {
	runtime.SetFinalizer(ctop, nil)
	ctop.finalize()
}

// C CuckooTableOptions to go CuckooTableOptions
func (cctop *C.CuckooTableOptions_t) toCuckooTableOptions() (ctop *CuckooTableOptions) {
	ctop = &CuckooTableOptions{ctop: *cctop}	
	runtime.SetFinalizer(ctop, finalize)
	return
}

// Setter methods for CuckooTableOptions
// Determines the utilization of hash tables. Smaller values
// result in larger hash tables with fewer collisions.
func (ctop *CuckooTableOptions) SetHashTableRatio(ratio float64) {
	var cctop *C.CuckooTableOptions_t= &ctop.ctop
	C.CuckooTableOptions_set_hash_table_ratio(cctop, C.double(ratio))
}

// A property used by builder to determine the depth to go to
// to search for a path to displace elements in case of
// collision. See Builder.MakeSpaceForKey method. Higher
// values result in more efficient hash tables with fewer
// lookups but take more time to build.
func (ctop *CuckooTableOptions) SetMaxSearchDepth(depth uint32) {
	var cctop *C.CuckooTableOptions_t= &ctop.ctop
	C.CuckooTableOptions_set_max_search_depth(cctop, C.uint32_t(depth))
}

// In case of collision while inserting, the builder
// attempts to insert in the next cuckoo_block_size
// locations before skipping over to the next Cuckoo hash
// function. This makes lookups more cache friendly in case
// of collisions.
func (ctop *CuckooTableOptions) SetCuckooBlockSize(sz uint32) {
	var cctop *C.CuckooTableOptions_t= &ctop.ctop
	C.CuckooTableOptions_set_cuckoo_block_size(cctop, C.uint32_t(sz))
}

// If this option is enabled, user key is treated as uint64_t and its value
// is used as hash value directly. This option changes builder's behavior.
// Reader ignore this option and behave according to what specified in table
// property.
func (ctop *CuckooTableOptions) SetIdentityAsFirstHash(identity bool) {
	var cctop *C.CuckooTableOptions_t= &ctop.ctop
	C.CuckooTableOptions_set_identity_as_first_hash(cctop, toCBool(identity))
}

// If this option is set to true, module is used during hash calculation.
// This often yields better space efficiency at the cost of performance.
// If this optino is set to false, # of entries in table is constrained to be
// power of two, and bit and is used to calculate hash, which is faster in
// general.
func (ctop *CuckooTableOptions) SetUseModuleHash(use bool) {
	var cctop *C.CuckooTableOptions_t= &ctop.ctop
	C.CuckooTableOptions_set_use_module_hash(cctop, toCBool(use))
}

// Create a cuckoo table factory.
func (ctop *CuckooTableOptions) NewCuckooTableFactory() *TableFactory {
	var cctop *C.CuckooTableOptions_t= &ctop.ctop
	ctbf := C.NewCuckooTableFactory(cctop)
	return ctbf.toTableFactory()
}

// Create default CuckooTableOptions
func NewCuckooTableOptions() *CuckooTableOptions {
	cctop := C.NewCuckooTableOptionsTDefault()
	return cctop.toCuckooTableOptions()
}
