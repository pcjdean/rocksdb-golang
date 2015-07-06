// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package rocksdb

/*
#cgo LDFLAGS: -lrocksdb -lstdc++ -lz -lrt
#include "db.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type ColumnFamilyHandle struct {
	cfh C.ColumnFamilyHandle_t
}

func (cfh *ColumnFamilyHandle) finalize() {
	var ccfh *C.ColumnFamilyHandle_t = unsafe.Pointer(&cfh.cfh)
	C.DeleteColumnFamilyHandleT(ccfh, false)
}

func (cfh *ColumnFamilyHandle) GetName() string {
	var ptr *C.ColumnFamilyHandle_t = unsafe.Pointer(&cfh.Cfh)
	rstr := cString{C.ColumnFamilyGetName(ptr)}
	return rstr.goString(true);
}
    
func (cfh *ColumnFamilyHandle) GetID() uint32 {
	var ptr *C.ColumnFamilyHandle_t = unsafe.Pointer(&cfh.Cfh)
	return C.ColumnFamilyGetID(ptr)
}

func newColumnFamilyHandleArrayFromCArray(cfh *C.ColumnFamilyHandle_t, sz uint) (cfhs []*ColumnFamilyHandle) {
	defer C.DeleteColumnFamilyHandleTArray(cfh)
	cfhs = make([]*ColumnFamilyHandle, sz)
	for var i = 0; i < sz; i++ {
		cfhs[i] = &ColumnFamilyHandle{cfh: (*[sz]C.String_t)(unsafe.Pointer(cfh))[i]}
		runtime.SetFinalizer(cfhs[i], finalize)
	}
	return
}

func newCArrayFromColumnFamilyHandleArray(cfhs ...*ColumnFamilyHandle) (ccfhs []C.ColumnFamilyHandle_t) {
	var cfhlen int
	if cfhs {
		cfhs.([]*ColumnFamilyHandle)
		cfhlen = len(cfhs)
		ccfhs = make([]C.ColumnFamilyHandle_t, cfhlen)
		for var i = 0; i < cfhlen; i++ {
			ccfhs[i] = cfhs[i].cfh
		}
	}
	return
}

type TablePropertiesCollection struct {
	tpc c.TablePropertiesCollection_t
}

func (tpc *TablePropertiesCollection) finalize() {
	var ctpc *C.TablePropertiesCollection_t = unsafe.Pointer(&tpc.tpc)
	C.DeleteTablePropertiesCollectionT(ctpc, false)
}

// Abstract handle to particular state of a DB.
// A Snapshot is an immutable object and can therefore be safely
// accessed from multiple threads without any external synchronization.
type Snapshot struct {
	snp C.Snapshot_t
}

func (snp *Snapshot) finalize() {
	var csnp *C.Snapshot_t = unsafe.Pointer(&snp.snp)
	C.DeleteSnapshotT(csnp, false)
}

func (snp *Snapshot) GetSequenceNumber() uint64 {
	var csnp *C.Snapshot_t = unsafe.Pointer(&snp.snp)
	return SnapshotGetSequenceNumber(csnp)
}

// A range of keys
type Range struct {
	rng C.Range_t
}

func (rng *Range) finalize() {
	var crng *C.Range_t = unsafe.Pointer(&rng.rng)
	C.DeleteRangeT(crng, false)
}

func newCArrayFromRangeArray(rngs ...*Range) (crngs []C.Range_t) {
	var sz int
	if rngs {
		rngs.([]*Range)
		sz = len(rngs)
		crngs = make([]C.Range_t, sz)
		for var i = 0; i < sz; i++ {
			crngs[i] = rngs[i].rng
		}
	}
	return
}

type ColumnFamilyDescriptor struct {
	cfd C.ColumnFamilyDescriptor_t
}

func (cfd *ColumnFamilyDescriptor) finalize() {
	var ccfd *C.ColumnFamilyDescriptor_t = unsafe.Pointer(&cfd.cfd)
	C.DeleteColumnFamilyDescriptorT(ccfd, false)
}

func newCArrayFromColumnFamilyDescriptorArray(cfds ...*ColumnFamilyDescriptor) (ccfds []C.ColumnFamilyDescriptor_t) {
	var cfdlen int
	if cfds {
		cfds.([]*ColumnFamilyDescriptor)
		cfdlen = len(cfds)
		ccfds = make([]C.ColumnFamilyDescriptor_t, cfdlen)
		for var i = 0; i < cfdlen; i++ {
			ccfds[i] = cfds[i].cfd
		}
	}
	return
}

// A DB is a persistent ordered map from keys to values.
// A DB is safe for concurrent access from multiple threads without
// any external synchronization.
type DB struct {
	db C.DB_t
}

func (db *DB) finalize() {
	var cdb *C.DB_t = unsafe.Pointer(&db.db)
	C.DeleteDBT(cdb, false)
}

// Open the database with the specified "name".
// Stores a pointer to a heap-allocated database in *dbptr and returns
// OK on success.
// Stores nullptr in *dbptr and returns a non-OK status on error.
// Caller should delete *dbptr when it is no longer needed.

// Open DB with column families.
// db_options specify database specific options
// column_families is the vector of all column families in the database,
// containing column family name and options. You need to open ALL column
// families in the database. To get the list of column families, you can use
// ListColumnFamilies(). Also, you can open only a subset of column families
// for read-only access.
// The default column family name is 'default' and it's stored
// in rocksdb::kDefaultColumnFamilyName.
// If everything is OK, handles will on return be the same size
// as column_families --- handles[i] will be a handle that you
// will use to operate on column family column_family[i]
func Open(options *Options, name *string, cfds ...*ColumnFamilyDescriptor) (db *DB, stat *Status, cfhs []*ColumnFamilyHandle) {
	db = &DB{}
	rstr := newCStringFromString(name)
	defer rstr.del()

	ccfds := newCArrayFromColumnFamilyDescriptorArray(cfds...)

	var (
		cdb *C.DB_t = unsafe.Pointer(&db.db)
		opt *C.Options_t = unsafe.Pointer(&options.opt)
		cstr *C.String_t = unsafe.Pointer(&rstr.str)
		cfh *ColumnFamilyHandle
	)

	if ccfds {
		stat = &Status{sta: C.DBOpenWithColumnFamilies(opt, cstr, unsafe.Pointers(&ccfds[0]), len(ccfds), unsafe.Pointer(&cfh), cdb)}
		if stat.Ok() {
			cfhs = newColumnFamilyHandleArrayFromCArray(cfh)
		}
	} else {
		stat = &Status{sta: C.DBOpen(opt, cstr, cdb)}
	}

	if stat.Ok() {
		runtime.SetFinalizer(db, finalize)
	}

	runtime.SetFinalizer(stat, finalize)
	return
}


// Open the database for read only. All DB interfaces
// that modify data, like put/delete, will return error.
// If the db is opened in read only mode, then no compactions
// will happen.
//
// Not supported in ROCKSDB_LITE, in which case the function will
// return Status_t::NotSupported.

// Open the database for read only with column families. When opening DB with
// read only, you can specify only a subset of column families in the
// database that should be opened. However, you always need to specify default
// column family. The default column family name is 'default' and it's stored
// in rocksdb::kDefaultColumnFamilyName
//
// Not supported in ROCKSDB_LITE, in which case the function will
// return Status_t::NotSupported.
func OpenForReadOnly(options *Options, name *string, cfds ...*ColumnFamilyDescriptor, error_if_log_file_exist ...bool) (db *DB, cfhs []*ColumnFamilyHandle, stat *Status) {
	db = &DB{}
	rstr := newCStringFromString(name)
	defer rstr.del()

	ccfds := newCArrayFromColumnFamilyDescriptorArray(cfds...)

	var (
		cdb *C.DB_t = unsafe.Pointer(&db.db)
		opt *C.Options_t = unsafe.Pointer(&options.opt)
		cstr *C.String_t = unsafe.Pointer(&rstr.str)
		cfh *C.ColumnFamilyHandle_t
		cflg C.bool = false
	)

	if error_if_log_file_exist {
		cflg = C.bool(error_if_log_file_exist[0])
	}

	if ccfds {
		stat = &Status{sta: C.DBOpenForReadOnlyWithColumnFamilies(opt, cstr, unsafe.Pointers(&ccfds[0]), len(ccfds), unsafe.Pointer(&cfh), cdb, cflg)}
		if stat.Ok() && cfdlen > 0 {
			cfhs = newColumnFamilyHandleArrayFromCArray(cfh)
		}
	} else {
		stat = &Status{sta: C.DBOpenForReadOnly(opt, cstr, cdb, cflg)}
	}

	if stat.Ok() {
		runtime.SetFinalizer(db, finalize)
	}

	runtime.SetFinalizer(stat, finalize)
	return
}

// ListColumnFamilies will open the DB specified by argument name
// and return the list of all column families in that DB
// through column_families argument. The ordering of
// column families in column_families is unspecified.
func ListColumnFamilies(dbopt *DBOptions, name *string) (cfss []string, stat *Status) {
	rstr := newCStringFromString(name)
	defer rstr.del()

	var (
		opt *C.DBOptions_t = unsafe.Pointers(&dbopt.dbopt)
		cstr *C.String_t = unsafe.Pointer(&rstr.str)
		cfs *C.String_t
		sz C.int
	)

	stat = &Status{sta: C.DBListColumnFamilies(opt, cstr, unsafe.Pointer(&cfs), unsafe.Pointer(&sz))}
	if stat.Ok() && sz > 0 {
		cfss = newStringArrayFromCArray(cfs, sz)
	}
	runtime.SetFinalizer(stat, finalize)
	return
} 


// Create a column_family and return the handle of column family
// through the argument handle.
func (db *DB) CreateColumnFamily(options *ColumnFamilyOptions, colfname *string) (cfd *ColumnFamilyHandle, stat *Status) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		opt *C.DBOptions_t = unsafe.Pointers(&dbopt.dbopt)
		cstr *C.String_t = unsafe.Pointer(&colfname.str)
		ccfd C.ColumnFamilyHandle_t
	)

	stat = &Status{sta: C.DBCreateColumnFamily(cdb, opt, cstr, unsafe.Pointer(&ccfd))}
	if stat.Ok() {
		cfd = &ColumnFamilyHandle{cfh: ccfd}
		runtime.SetFinalizer(cfd, finalize)
	}
	runtime.SetFinalizer(stat, finalize)
	return
}

// Drop a column family specified by column_family handle. This call
// only records a drop record in the manifest and prevents the column
// family from flushing and compacting.
func (db *DB) DropColumnFamily(cfd *ColumnFamilyHandle) (stat *Status) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		ccfd *C.ColumnFamilyHandle_t = unsafe.Pointers(&cfd.cfd) 
	)
	stat = &Status{sta: C.DBDropColumnFamily(cdb, ccfd)}
	runtime.SetFinalizer(stat, finalize)
	return
}

// Set the database entry for "key" to "value".
// If "key" already exists, it will be overwritten.
// Returns OK on success, and a non-OK status on error.
// Note: consider setting options.sync = true.
func (db *DB) Put(options *WriteOptions, key []byte, val []byte, cfd ...*ColumnFamilyHandle) (stat *Status) {
	ckey := newSliceFromBytes(key)
	defer ckey.del()
	cval := newSliceFromBytes(val)
	defer cval.del()

	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cwopt *C.WriteOptions_t = unsafe.Pointers(&options.wopt)
		ccfd *C.ColumnFamilyHandle_t
		cckey *C.Slice_t = unsafe.Pointers(&ckey.slc) 
		ccval *C.Slice_t = unsafe.Pointers(&cval.slc) 
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		stat = &Status{sta: C.DBPutWithColumnFamily(cdb, cwopt, ccfd, cckey, ccval)}
	} else {
		stat = &Status{sta: C.DBPut(cdb, cwopt, cckey, ccval)}
	}
	runtime.SetFinalizer(stat, finalize)
	return
}

// Remove the database entry (if any) for "key".  Returns OK on
// success, and a non-OK status on error.  It is not an error if "key"
// did not exist in the database.
// Note: consider setting options.sync = true.
func (db *DB) Delete(options *WriteOptions, key []byte, cfd ...*ColumnFamilyHandle) (stat *Status) {
	ckey := newSliceFromBytes(key)
	defer ckey.del()

	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cwopt *C.WriteOptions_t = unsafe.Pointers(&options.wopt)
		ccfd *C.ColumnFamilyHandle_t
		cckey *C.Slice_t = unsafe.Pointers(&ckey.slc) 
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		stat = &Status{sta: C.DBDeleteWithColumnFamily(cdb, cwopt, ccfd, cckey)}
	} else {
		stat = &Status{sta: C.DBDelete(cdb, cwopt, cckey)}
	}
	runtime.SetFinalizer(stat, finalize)
	return
}

// Merge the database entry for "key" with "value".  Returns OK on success,
// and a non-OK status on error. The semantics of this operation is
// determined by the user provided merge_operator when opening DB.
// Note: consider setting options.sync = true.
func (db *DB) Merge(options *WriteOptions, key []byte, val []byte, cfd ...*ColumnFamilyHandle) (stat *Status) {
	ckey := newSliceFromBytes(key)
	defer ckey.del()
	cval := newSliceFromBytes(val)
	defer cval.del()

	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cwopt *C.WriteOptions_t = unsafe.Pointers(&options.wopt)
		ccfd *C.ColumnFamilyHandle_t
		cckey *C.Slice_t = unsafe.Pointers(&ckey.slc) 
		ccval *C.Slice_t = unsafe.Pointers(&cval.slc) 
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		stat = &Status{sta: C.DBMergeWithColumnFamily(cdb, cwopt, ccfd, cckey, ccval)}
	} else {
		stat = &Status{sta: C.DBMerge(cdb, cwopt, cckey, ccval)}
	}
	runtime.SetFinalizer(stat, finalize)
	return
}

// Apply the specified updates to the database.
// If `updates` contains no update, WAL will still be synced if
// options.sync=true.
// Returns OK on success, non-OK on failure.
// Note: consider setting options.sync = true.
func (db *DB) Write(options *WriteOptions, wbt *WriteBatch) (stat *Status) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cwopt *C.WriteOptions_t = unsafe.Pointers(&options.wopt)
		cwbt *C.WriteBatch_t = unsafe.Pointers(&wbt.wbt)
	)
	stat = &Status{sta: C.DBWrite(cdb, cwopt, cwbt)}
	runtime.SetFinalizer(stat, finalize)
	return
}

// If the database contains an entry for "key" store the
// corresponding value in *value and return OK.
//
// If there is no entry for "key" leave *value unchanged and return
// a status for which Status_t::IsNotFound() returns true.
//
// May return some other Status_t on an error.
func (db *DB) Get(options *ReadOptions, key []byte, cfd ...*ColumnFamilyHandle) (val string, stat *Status) {
	ckey := newSliceFromBytes(key)
	defer ckey.del()
	cval := C.newCString()

	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cropt *C.ReadOptions_t = unsafe.Pointers(&options.ropt)
		ccfd *C.ColumnFamilyHandle_t
		cckey *C.Slice_t = unsafe.Pointers(&ckey.slc) 
		ccval *C.String_t = unsafe.Pointers(&cval.str) 
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		stat = &Status{sta: C.DBGetWithColumnFamily(cdb, cropt, ccfd, cckey, ccval)}
	} else {
		stat = &Status{sta: C.DBGet(cdb, cropt, cckey, ccval)}
	}
	val = cval.goString(true)
	runtime.SetFinalizer(stat, finalize)
	return
}

// If keys[i] does not exist in the database, then the i'th returned
// status will be one for which Status_t::IsNotFound() is true, and
// (*values)[i] will be set to some arbitrary value (often ""). Otherwise,
// the i'th returned status will have Status_t::ok() true, and (*values)[i]
// will store the value associated with keys[i].
//
// (*values) will always be resized to be the same size as (keys).
// Similarly, the number of returned statuses will be the number of keys.
// Note: keys will not be "de-duplicated". Duplicate keys will return
// duplicate values in order.
func (db *DB) MultiGet(options *ReadOptions, keys [][]byte, cfhs ...*ColumnFamilyHandle) (vals []string, stats []*Status) {
	ckeys := newSlicesFromBytesArray(keys)
	defer ckeys.del()
	cckeys := ckeys.toCArray()
	ccfhs := newCArrayFromColumnFamilyHandleArray(cfhs...)

	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cropt *C.ReadOptions_t = unsafe.Pointers(&options.ropt)
		ccvals *C.String_t
	)

	if ccfhs {
		stats = newStatusArrayFromCArray(C.DBMultiGetWithColumnFamily(cdb, cropt, unsafe.Pointers(&ccfhs[0]), len(cfhs), unsafe.Pointers(&cckeys[0]), len(cckeys), unsafe.Pointers(&ccvals)))
	} else {
		stats = newStatusArrayFromCArray(C.DBMultiGet(cdb, cropt, unsafe.Pointers(&cckeys[0]), len(cckeys), unsafe.Pointers(&ccvals)))
	}
	vals = newStringArrayFromCArray(ccvals, len(keys))
	return
}

// If the key definitely does not exist in the database, then this method
// returns false, else true. If the caller wants to obtain value when the key
// is found in memory, a bool for 'value_found' must be passed. 'value_found'
// will be true on return if value has been set properly.
// This check is potentially lighter-weight than invoking DB::Get(). One way
// to make this lighter weight is to avoid doing any IOs.
// Default implementation here returns true and sets 'value_found' to false
func (db *DB) KeyMayExist(options *ReadOptions, key []byte, cfd ...*ColumnFamilyHandle) (res bool, valfound bool, val string) {
	ckey := newSliceFromBytes(key)
	defer ckey.del()
	cval := C.newCString()

	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cropt *C.ReadOptions_t = unsafe.Pointers(&options.ropt)
		ccfd *C.ColumnFamilyHandle_t
		cckey *C.Slice_t = unsafe.Pointers(&ckey.slc) 
		ccval *C.String_t = unsafe.Pointers(&cval.str) 
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		res = C.DBKeyMayExistWithColumnFamily(cdb, cropt, ccfd, cckey, ccval, unsafe.Pointers(&valfound))
	} else {
		res = C.DBKeyMayExist(cdb, cropt, cckey, ccval, unsafe.Pointers(&valfound))
	}
	val = cval.goString(true)
	return
}

// Return a heap-allocated iterator over the contents of the database.
// The result of NewIterator() is initially invalid (caller must
// call one of the Seek methods on the iterator before using it).
//
// Caller should delete the iterator when it is no longer needed.
// The returned iterator should be deleted before this db is deleted.
func (db *DB) NewIterator(options *ReadOptions, cfd ...*ColumnFamilyHandle) (it *Iterator) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cropt *C.ReadOptions_t = unsafe.Pointers(&options.ropt)
		ccfd *C.ColumnFamilyHandle_t
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		it = &Iterator{it: C.DBNewIteratorWithColumnFamily(cdb, cropt, ccfd), db: db}
	} else {
		it = &Iterator{it: C.DBNewIterator(cdb, cropt), db: db}
	}
	runtime.SetFinalizer(it, finalize)
	return
}

// Returns iterators from a consistent database state across multiple
// column families. Iterators are heap allocated and need to be deleted
// before the db is deleted
func (db *DB) NewIterators(options *ReadOptions, cfhs []*ColumnFamilyHandle) (vals []*Iterator, stat *Status) {
	ccfhs := newCArrayFromColumnFamilyHandleArray(cfhs...)

	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cropt *C.ReadOptions_t = unsafe.Pointers(&options.ropt)
		ccvals *C.Iterator_t
		valsz int
	)

	ccfhs[0].(*C.ColumnFamilyHandle_t)
	stat = &Status{sta: C.DBNewIterators(cdb, cropt, unsafe.Pointers(&ccfhs[0]), len(ccfhs), unsafe.Pointers(&ccvals), unsafe.Pointers(&valsz))}
	vals = newIteratorArrayFromCArray(ccvals, valsz, db)
	runtime.SetFinalizer(stat, finalize)
	return
}

// Return a handle to the current DB state.  Iterators created with
// this handle will all observe a stable snapshot of the current DB
// state.  The caller must call ReleaseSnapshot(result) when the
// snapshot is no longer needed.
//
// nullptr will be returned if the DB fails to take a snapshot or does
// not support snapshot.
func (db *DB) GetSnapshot() (snp *Snapshot) {
	var cdb *C.DB_t = unsafe.Pointers(&db.db)

	snp = &Snapshot{snp: C.DBGetSnapshot(cdb)}
	runtime.SetFinalizer(snp, finalize)
	return
}

// Release a previously acquired snapshot.  The caller must not
// use "snapshot" after this call.
func (db *DB) ReleaseSnapshot(snp *Snapshot) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		csnp *C.Snapshot_t = unsafe.Pointers(&snp.snp)
	)

	C.DBReleaseSnapshot(cdb, csnp)
	return
}

// DB implementations can export properties about their state
// via this method.  If "property" is a valid property understood by this
// DB implementation, fills "*value" with its current value and returns
// true.  Otherwise returns false.
//
//
// Valid property names include:
//
//  "rocksdb.num-files-at-level<N>" - return the number of files at level <N>,
//     where <N> is an ASCII representation of a level number (e.g. "0").
//  "rocksdb.stats" - returns a multi-line string that describes statistics
//     about the internal operation of the DB.
//  "rocksdb.sstables" - returns a multi-line string that describes all
//     of the sstables that make up the db contents.
//  "rocksdb.cfstats"
//  "rocksdb.dbstats"
//  "rocksdb.num-immutable-mem-table"
//  "rocksdb.mem-table-flush-pending"
//  "rocksdb.compaction-pending" - 1 if at least one compaction is pending
//  "rocksdb.background-errors" - accumulated number of background errors
//  "rocksdb.cur-size-active-mem-table"
//  "rocksdb.cur-size-all-mem-tables"
//  "rocksdb.num-entries-active-mem-table"
//  "rocksdb.num-entries-imm-mem-tables"
//  "rocksdb.num-deletes-active-mem-table"
//  "rocksdb.num-deletes-imm-mem-tables"
//  "rocksdb.estimate-num-keys" - estimated keys in the column family
//  "rocksdb.estimate-table-readers-mem" - estimated memory used for reding
//      SST tables, that is not counted as a part of block cache.
//  "rocksdb.is-file-deletions-enabled"
//  "rocksdb.num-snapshots"
//  "rocksdb.oldest-snapshot-time"
//  "rocksdb.num-live-versions" - `version` is an internal data structure.
//      See version_set.h for details. More live versions often mean more SST
//      files are held from being deleted, by iterators or unfinished
//      compactions.
func (db *DB) GetProperty(options *ReadOptions, prop []byte, cfd ...*ColumnFamilyHandle) (val string, res bool) {
	cprop := newSliceFromBytes(prop)
	defer cprop.del()
	cval := C.newCString()

	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cropt *C.ReadOptions_t = unsafe.Pointers(&options.ropt)
		ccfd *C.ColumnFamilyHandle_t
		ccprop *C.Slice_t = unsafe.Pointers(&cprop.slc) 
		ccval *C.String_t = unsafe.Pointers(&cval.str) 
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		res = C.DBGetPropertyWithColumnFamily(cdb, cropt, ccfd, ccprop, ccval)
	} else {
		res = C.DBGetProperty(cdb, cropt, ccprop, ccval)
	}
	val = cval.goString(true)
	return
}

// Similar to GetProperty(), but only works for a subset of properties whose
// return value is an integer. Return the value by integer. Supported
// properties:
//  "rocksdb.num-immutable-mem-table"
//  "rocksdb.mem-table-flush-pending"
//  "rocksdb.compaction-pending"
//  "rocksdb.background-errors"
//  "rocksdb.cur-size-active-mem-table"
//  "rocksdb.cur-size-all-mem-tables"
//  "rocksdb.num-entries-active-mem-table"
//  "rocksdb.num-entries-imm-mem-tables"
//  "rocksdb.num-deletes-active-mem-table"
//  "rocksdb.num-deletes-imm-mem-tables"
//  "rocksdb.estimate-num-keys"
//  "rocksdb.estimate-table-readers-mem"
//  "rocksdb.is-file-deletions-enabled"
//  "rocksdb.num-snapshots"
//  "rocksdb.oldest-snapshot-time"
//  "rocksdb.num-live-versions"
func (db *DB) GetIntProperty(prop []byte, cfd ...*ColumnFamilyHandle) (val uint64, res bool) {
	cprop := newSliceFromBytes(prop)
	defer cprop.del()

	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		ccfd *C.ColumnFamilyHandle_t
		ccprop *C.Slice_t = unsafe.Pointers(&cprop.slc) 
		cval C.uint64_t
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		res = C.DBGetIntPropertyWithColumnFamily(cdb, ccfd, ccprop, unsafe.Pointers(&cval))
	} else {
		res = C.DBGetIntProperty(cdb, ccprop, unsafe.Pointers(&cval))
	}
	val = cval
	return
}

// For each i in [0,n-1], store in "sizes[i]", the approximate
// file system space used by keys in "[range[i].start .. range[i].limit)".
//
// Note that the returned sizes measure file system space usage, so
// if the user data compresses by a factor of ten, the returned
// sizes will be one-tenth the size of the corresponding user data size.
//
// The results may not include the sizes of recently written data.
func (db *DB) GetApproximateSizes(rngs []*Range, cfd ...*ColumnFamilyHandle) (vals []uint64) {
	crngs := newCArrayFromRangeArray(rngs...)

	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		ccfd *C.ColumnFamilyHandle_t
		ccrngs *C.Range_t = unsafe.Pointers(&crngs[0]) 
		sz C.int
		cval *C.uint64_t
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		res = C.DBGetApproximateSizesWithColumnFamily(cdb, ccfd, ccrngs, sz, unsafe.Pointers(&cval))
	} else {
		res = C.DBGetApproximateSizes(cdb, ccrngs, sz, unsafe.Pointers(&cval))
	}
	vals = newUint64ArrayFromCArray(cval)
	return
}

// Compact the underlying storage for the key range [*begin,*end].
// The actual compaction interval might be superset of [*begin, *end].
// In particular, deleted and overwritten versions are discarded,
// and the data is rearranged to reduce the cost of operations
// needed to access the data.  This operation should typically only
// be invoked by users who understand the underlying implementation.
//
// begin==nullptr is treated as a key before all keys in the database.
// end==nullptr is treated as a key after all keys in the database.
// Therefore the following call will compact the entire database:
//    db->CompactRange(nullptr, nullptr);
// Note that after the entire database is compacted, all data are pushed
// down to the last level containing any data. If the total data size
// after compaction is reduced, that level might not be appropriate for
// hosting all the files. In this case, client could set reduce_level
// to true, to move the files back to the minimum level capable of holding
// the data set or a given level (specified by non-negative target_level).
// Compaction outputs should be placed in options.db_paths[target_path_id].
// Behavior is undefined if target_path_id is out of range.
func (db *DB) CompactRange(begin []byte, end []byte, cfd ...*ColumnFamilyHandle, reduce_level ...bool, target_level ...int, target_path_id ...uint32) (stat *Status) {
	cbegin := newSliceFromBytes(begin)
	defer cbegin.del()
	cend := newSliceFromBytes(end)
	defer cend.del()

	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		ccfd *C.ColumnFamilyHandle_t
		credl C.bool 
		ctarl C.int = -1
		ctpi C.uint32 
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}
	if reduce_level {
		credl = C.bool(reduce_level[0])
	}
	if target_level {
		ctarl = C.int(target_level[0])
	}
	if target_path_id {
		ctpi = C.uint32(target_path_id[0])
	}

	if ccfd {
		stat = &Status{sta: C.DBCompactRangeWithColumnFamily(cdb, ccfd, unsafe.Pointers(&cbegin[0]), unsafe.Pointers(&cend[0]), credl, ctarl, ctpi)}
	} else {
		stat = &Status{sta: C.DBCompactRange(cdb, unsafe.Pointers(&cbegin[0]), unsafe.Pointers(&cend[0]), credl, ctarl, ctpi)}
	}
	runtime.SetFinalizer(stat, finalize)
	return
}

func (db *DB) SetOptions(opts []string, cfhs ...*ColumnFamilyHandle) (stat Status) {
	copts := newcStringsFromStringArray(opts)
	defer copts.del()
	ccopts := copts.toCArray()

	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		ccfd *C.ColumnFamilyHandle_t
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		stat = &Status{sta: C.DBSetOptionsWithColumnFamily(cdb, ccfd, unsafe.Pointers(&ccopts[0]), len(ccopts))}
	} else {
		stat = &Status{sta: C.DBSetOptions(cdb, unsafe.Pointers(&ccopts[0]), len(ccopts))}
	}
	return
}

// CompactFiles() inputs a list of files specified by file numbers
// and compacts them to the specified level.  Note that the behavior
// is different from CompactRange in that CompactFiles() will
// perform the compaction job using the CURRENT thread.
//
// @see GetDataBaseMetaData
// @see GetColumnFamilyMetaData
func (db *DB) CompactFiles(options *CompactionOptions, files []string, level int, cfhs ...*ColumnFamilyHandle, path_id ...int) (stat Status) {
	cfiles := newcStringsFromStringArray(opts)
	defer cfiles.del()
	ccfiles := cfiles.toCArray()

	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		ccfd *C.ColumnFamilyHandle_t
		ccopt *C.CompactionOptions_t = unsafe.Pointers(&options.copt)
		cpid C.int = -1
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}
	if path_id {
		cpid = C.int(path_id[0])
	}

	if ccfd {
		stat = &Status{sta: C.DBCompactFilesWithColumnFamily(cdb, ccopt, ccfd, unsafe.Pointers(&ccfiles[0]), len(ccfiles), level, cpid)}
	} else {
		stat = &Status{sta: C.DBCompactFilesWithColumnFamily(cdb, ccopt, unsafe.Pointers(&ccfiles[0]), len(ccfiles), level, cpid)}
	}
	return
}

// Number of levels used for this DB.
func (db *DB) NumberLevels(cfd ...*ColumnFamilyHandle) (level int) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		ccfd *C.ColumnFamilyHandle_t
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		level = C.DBNumberLevelsWithColumnFamily(cdb, ccfd)
	} else {
		level = C.DBNumberLevels(cdb)
	}
	return
}

// Maximum level to which a new compacted memtable is pushed if it
// does not create overlap.
func (db *DB) MaxMemCompactionLevel(cfd ...*ColumnFamilyHandle) (level int) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		ccfd *C.ColumnFamilyHandle_t
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		level = C.DBMaxMemCompactionLevelWithColumnFamily(cdb, ccfd)
	} else {
		level = C.DBMaxMemCompactionLevel(cdb)
	}
	return
}

// Number of files in level-0 that would stop writes.
func (db *DB) Level0StopWriteTrigger(cfd ...*ColumnFamilyHandle) (level int) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		ccfd *C.ColumnFamilyHandle_t
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		level = C.DBLevel0StopWriteTriggerWithColumnFamily(cdb, ccfd)
	} else {
		level = C.DBLevel0StopWriteTrigger(cdb)
	}
	return
}

// Get DB name -- the exact same name that was provided as an argument to
// DB::Open()
func (db *DB) GetName() (name string) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cname C.String_t = C.DBGetName(cdb)
	)
	name = cname.toCString()
	return
}

// Get Env object from the DB
func (db *DB) GetEnv() (env *Env) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
	)
	env = &Env{env: C.DBGetEnv(cdb)}
	runtime.SetFinalizer(env, finalize)
	return
}

// Get DB Options that we use.  During the process of opening the
// column family, the options provided when calling DB::Open() or
// DB::CreateColumnFamily() will have been "sanitized" and transformed
// in an implementation-defined manner.
func (db *DB) GetOptions(cfd ...*ColumnFamilyHandle) (opt *Options) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		ccfd *C.ColumnFamilyHandle_t
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
		copt C.String_t
	}

	if ccfd {
		copt = C.DBGetOptionsWithColumnFamily(cdb, ccfd)
	} else {
		copt = C.DBGetOptions(cdb)
	}

	opt = &Options{opt: copt}
	runtime.SetFinalizer(opt, finalize)
	return
}

func (db *DB) GetDBOptions() (dbopt *DBOptions) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
	)
	dbopt = &Env{dbopt: C.DBGetDBOptions(cdb)}
	runtime.SetFinalizer(dbopt, finalize)
	return
}

// Flush all mem-table data.
func (db *DB) Flush(options *FlushOptions, cfhs ...*ColumnFamilyHandle) (stat Status) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		ccfd *C.ColumnFamilyHandle_t
		cfopt *C.FlushOptions_t = unsafe.Pointers(&options.fopt)
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		stat = &Status{sta: C.DBFlushWithColumnFamily(cdb, cfopt, ccfd)}
	} else {
		stat = &Status{sta: C.DBFlush(cdb, cfopt)}
	}
	return
}

// The sequence number of the most recent transaction.
func (db *DB) GetLatestSequenceNumber() (sqnum SequenceNumber) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
	)
	sqnum = C.DBGetLatestSequenceNumber(cdb)
	return
}

#ifndef ROCKSDB_LITE

// Prevent file deletions. Compactions will continue to occur,
// but no obsolete files will be deleted. Calling this multiple
// times have the same effect as calling it once.
Status_t DBDisableFileDeletions(DB_t* dbptr)
{
    return NewStatusTCopy(dbptr ?
                          &GET_REP(dbptr, DB)->DisableFileDeletions() :
                          &invalid_status);
}

// Allow compactions to delete obsolete files.
// If force == true, the call to EnableFileDeletions() will guarantee that
// file deletions are enabled after the call, even if DisableFileDeletions()
// was called multiple times before.
// If force == false, EnableFileDeletions will only enable file deletion
// after it's been called at least as many times as DisableFileDeletions(),
// enabling the two methods to be called by two threads concurrently without
// synchronization -- i.e., file deletions will be enabled only after both
// threads call EnableFileDeletions()
Status_t DBEnableFileDeletions(DB_t* dbptr, bool force)
{
    return NewStatusTCopy(dbptr ?
                          &GET_REP(dbptr, DB)->EnableFileDeletions(force) :
                          &invalid_status);
}

// GetLiveFiles followed by GetSortedWalFiles can generate a lossless backup

// Retrieve the list of all files in the database. The files are
// relative to the dbname and are not absolute paths. The valid size of the
// manifest file is returned in manifest_file_size. The manifest file is an
// ever growing file, but only the portion specified by manifest_file_size is
// valid for this snapshot.
// Setting flush_memtable to true does Flush before recording the live files.
// Setting flush_memtable to false is useful when we don't want to wait for
// flush which may have to wait for compaction to complete taking an
// indeterminate time.
//
// In case you have multiple column families, even if flush_memtable is true,
// you still need to call GetSortedWalFiles after GetLiveFiles to compensate
// for new data that arrived to already-flushed column families while other
// column families were flushing
Status_t DBGetLiveFiles(DB_t* dbptr,
                        const String_t live_files[],
                        int* n,
                        uint64_t* manifest_file_size,
                        bool flush_memtable)
{
    Status &ret;
    if (dbptr)
    {
        std::vector<std::string> live_files_vec;
        ret = GET_REP(dbptr, DB)->GetLiveFiles(live_files_vec, manifest_file_size, flush_memtable);
        *n = live_files_vec.size();
        live_files = new String_t[*n];
        for (int j = 0; j < *n; j++)
            GET_REP_REF(live_files[j], String) = std::move(live_files_vec[j]);
    }
    else
    {
        ret = invalid_status;
    }
    return NewStatusTCopy(&ret);
}

// Retrieve the sorted list of all wal files with earliest file first
Status_t DBGetSortedWalFiles(DB_t* dbptr, LogFile_t files[], int* n)
{
    Status &ret;
    if (dbptr)
    {
        VectorLogPtr files_vec;
        ret = GET_REP(dbptr, DB)->GetSortedWalFiles(files_vec);
        *n = files_vec.size();
        files = new LogFile_t[*n];
        for (int j = 0; j < *n; j++)
            GET_REP(files[j]) = files_vec[j].release();
    }
    else
    {
        ret = invalid_status;
    }
    return NewStatusTCopy(&ret);
}

// Sets iter to an iterator that is positioned at a write-batch containing
// seq_number. If the sequence number is non existent, it returns an iterator
// at the first available seq_no after the requested seq_no
// Returns Status_t::OK if iterator is valid
// Must set WAL_ttl_seconds or WAL_size_limit_MB to large values to
// use this api, else the WAL files will get
// cleared aggressively and the iterator might keep getting invalid before
// an update is read.
Status_t DBGetUpdatesSince(DB_t* dbptr, SequenceNumber seq_number,
                           TransactionLogIterator_t* iter,
                           const TransactionLogIterator_ReadOptions_t* read_options)
{
    Status &ret;
    if (dbptr)
    {
        unique_ptr<TransactionLogIterator> iter_ptr;
        if (GET_REP(read_options, TransactionLogIterator_ReadOptions) == NULL)
            GET_REP(read_options, TransactionLogIterator_ReadOptions) = &TransactionLogIterator::ReadOptions();
        ret = GET_REP(dbptr, DB)->GetUpdatesSince(seq_number, &iter_ptr, GET_REP_REF(read_options, TransactionLogIterator_ReadOptions));
        GET_REP(iter, TransactionLogIterator) = iter_ptr.release();
    }
    {
        ret = invalid_status;
    }
    return NewStatusTCopy(&ret);
}

// Delete the file name from the db directory and update the internal state to
// reflect that. Supports deletion of sst and log files only. 'name' must be
// path relative to the db directory. eg. 000001.sst, /archive/000003.log
Status_t DBDeleteFile(DB_t* dbptr, String_t* name)
{
    return NewStatusTCopy(dbptr ?
                          &GET_REP(dbptr, DB)->DeleteFile(GET_REP_REF(name, String) :
                          &invalid_status);
}

// Returns a list of all table files with their level, start key
// and end key
void DBGetLiveFilesMetaData(DB_t* dbptr, LiveFileMetaData_t metadata[], int* n)
{
    if (dbptr)
    {
        std::vector<LiveFileMetaData> metadata_vec;
        GET_REP(dbptr, DB)->GetLiveFilesMetaData(&metadata_vec);
        *n = metadata_vec.size();
        metadata = new LiveFileMetaData_t[*n];
        for (int j = 0; j < *n; j++)
            GET_REP_REF(metadata[j], LiveFileMetaData) = metadata_vec[j];
    }
}

// Obtains the meta data of the specified column family of the DB.
// Status_t::NotFound() will be returned if the current DB does not have
// any column family match the specified name.
//
// If cf_name is not specified, then the metadata of the default
// column family will be returned.
void DBGetColumnFamilyMetaDataWithColumnFamily(DB_t* dbptr, 
                                               ColumnFamilyHandle_t* column_family,
                                               ColumnFamilyMetaData_t* metadata)
{
    if (dbptr)
    {
        GET_REP(dbptr, DB)->GetColumnFamilyMetaData(GET_REP(column_family, ColumnFamilyHandle), GET_REP(metadata, ColumnFamilyMetaData));
    }
}

// Get the metadata of the default column family.
void DBGetColumnFamilyMetaData(DB_t* dbptr, 
                               ColumnFamilyMetaData_t* metadata)
{
    DBGetColumnFamilyMetaDataWithColumnFamily(dbptr, &DBDefaultColumnFamily(dbptr), metadata);
}
#endif  // ROCKSDB_LITE

// Sets the globally unique ID created at database creation time by invoking
// Env::GenerateUniqueId(), in identity. Returns Status_t::OK if identity could
// be set properly
Status_t DBGetDbIdentity(DB_t* dbptr, String_t* identity)
{
    return NewStatusTCopy(dbptr ?
                          &GET_REP(dbptr, DB)->GetDbIdentity(GET_REP_REF(identity, String)) :
                          &invalid_status);
}

// Returns default column family handle
ColumnFamilyHandle_t DBDefaultColumnFamily(DB_t* dbptr)
{
    ColumnFamilyHandle_t cf_handle;
    cf_handle.rep = dbptr ? GET_REP(dbptr, DB)->DefaultColumnFamily() : nullptr; 
    return cf_handle;
}

#ifndef ROCKSDB_LITE
Status_t DBGetPropertiesOfAllTablesWithColumnFamily(DB_t* dbptr, 
                                                    ColumnFamilyHandle_t* column_family,
                                                    TablePropertiesCollection_t* props)
{
    return NewStatusTCopy(dbptr ?
                          &GET_REP(dbptr, DB)->GetPropertiesOfAllTables(GET_REP(column_family, ColumnFamilyHandle), GET_REP(props, TablePropertiesCollection)) :
                          &invalid_status);
}

Status_t DBGetPropertiesOfAllTables(DB_t* dbptr, 
                                    TablePropertiesCollection_t* props)
{
    return DBGetPropertiesOfAllTablesWithColumnFamily(dbptr, &DBDefaultColumnFamily(dbptr), props);
}
#endif  // ROCKSDB_LITE

// Destroy the contents of the specified database.
// Be very careful using this method.
Status_t DBDestroyDB(const String_t* name, const Options_t* options)
{
    return NewStatusTCopy(&DestroyDB(GET_REP_REF(name, String), GET_REP_REF(options, Options)));
}

#ifndef ROCKSDB_LITE
// If a DB cannot be opened, you may attempt to call this method to
// resurrect as much of the contents of the database as possible.
// Some data may be lost, so be careful when calling this function
// on a database that contains important information.
Status_t DBRepairDB(const String_t* dbname, const Options_t* options)
{
    return NewStatusTCopy(&RepairDB(GET_REP_REF(dbname, String), GET_REP_REF(options, Options)));
}
#endif
