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
}

// Release resources
func (btop *BlockBasedTableOptions) finalize() {
	var cbtop *C.BlockBasedTableOptions_t= &btop.btop
	C.DeleteBlockBasedTableOptionsT(cbtop, toCBool(false))
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
}

// Release resources
func (ptop *PlainTableOptions) finalize() {
	var cptop *C.PlainTableOptions_t= &ptop.ptop
	C.DeletePlainTableOptionsT(cptop, toCBool(false))
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

// Create default plain table factory.
func (ptop *PlainTableOptions) NewPlainTableFactory() *TableFactory {
	var cptop *C.PlainTableOptions_t= &ptop.ptop
	ctbf := C.NewPlainTableFactory(cptop)
	return ctbf.toTableFactory()
}

// Wrap go CuckooTableOptions
type CuckooTableOptions struct {
	ctop C.CuckooTableOptions_t
}

// Release resources
func (ctop *CuckooTableOptions) finalize() {
	var cctop *C.CuckooTableOptions_t= &ctop.ctop
	C.DeleteCuckooTableOptionsT(cctop, toCBool(false))
}

// C CuckooTableOptions to go CuckooTableOptions
func (cctop *C.CuckooTableOptions_t) toCuckooTableOptions() (ctop *CuckooTableOptions) {
	ctop = &CuckooTableOptions{ctop: *cctop}	
	runtime.SetFinalizer(ctop, finalize)
	return
}

// Create default CuckooTableOptions
func NewCuckooTableOptions() *CuckooTableOptions {
	cctop := C.NewCuckooTableOptionsTDefault()
	return cctop.toCuckooTableOptions()
}
