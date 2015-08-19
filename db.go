// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package rocksdb

/*
#cgo CXXFLAGS: -std=c++11
#cgo LDFLAGS: -lrocksdb -lstdc++ -lz -lrt
#include "db.h"
*/
import "C"

import (
	"runtime"
)

type TablePropertiesCollection struct {
	tpc C.TablePropertiesCollection_t
}

func (tpc *TablePropertiesCollection) finalize() {
	var ctpc *C.TablePropertiesCollection_t = &tpc.tpc
	C.DeleteTablePropertiesCollectionT(ctpc, toCBool(false))
}

func (ctpc *C.TablePropertiesCollection_t) toTablePropertiesCollection() (tpc *TablePropertiesCollection) {
	tpc = &TablePropertiesCollection{tpc: *ctpc}	
	runtime.SetFinalizer(tpc, finalize)
	return
}

// A range of keys
type Range struct {
	startSlc *cSlice
	limitSlc *cSlice
	rng C.Range_t
}

func (rng *Range) finalize() {
	rng.startSlc.del()
	rng.limitSlc.del()
	var crng *C.Range_t = &rng.rng
	C.DeleteRangeT(crng, toCBool(false))
}

func NewRange(start, limit []byte) (rng *Range) {
	startSlc := newSliceFromBytes(start)
	limitSlc := newSliceFromBytes(limit)
	rng = &Range{rng: C.NewRangeTArgs(&startSlc.slc, &limitSlc.slc), startSlc: startSlc, limitSlc: limitSlc}
	runtime.SetFinalizer(rng, finalize)
	return
}

func newCArrayFromRangeArray(rngs ...*Range) (crngs []C.Range_t) {
	var sz int
	if rngs != nil {
		sz = len(rngs)
		crngs = make([]C.Range_t, sz)
		for i := 0; i < sz; i++ {
			crngs[i] = rngs[i].rng
		}
	}
	return
}

type ColumnFamilyDescriptor struct {
	cfd C.ColumnFamilyDescriptor_t
}

func (cfd *ColumnFamilyDescriptor) finalize() {
	var ccfd *C.ColumnFamilyDescriptor_t = &cfd.cfd
	C.DeleteColumnFamilyDescriptorT(ccfd, toCBool(false))
}

func newCArrayFromColumnFamilyDescriptorArray(cfds ...interface{}) (ccfds []C.ColumnFamilyDescriptor_t) {
	var cfdlen int
	if cfds != nil {
		cfdlen = len(cfds)
		n := 0

		for i := 0; i < cfdlen; i++ {
			v, ok := cfds[i].(*ColumnFamilyDescriptor)
			if ok {
				if cfds == nil {
					ccfds = make([]C.ColumnFamilyDescriptor_t, cfdlen)
				}
				ccfds[i] = v.cfd
				n++
			}
		}

		if 0 < n && cfdlen != n {
			ccfds = ccfds[:n]
		}
	}
	return
}

func newCArrayFromColumnFamilyHandleInterface(cfhs ...interface{}) (ccfhs []C.ColumnFamilyHandle_t) {
	var cfhlen int
	if cfhs != nil {
		cfhlen = len(cfhs)
		n := 0

		for i := 0; i < cfhlen; i++ {
			v, ok := cfhs[i].(*ColumnFamilyHandle)
			if ok {
				if cfhs == nil {
					ccfhs = make([]C.ColumnFamilyHandle_t, cfhlen)
				}
				ccfhs[i] = v.cfh
				n++
			}
		}

		if 0 < n && cfhlen != n {
			ccfhs = ccfhs[:n]
		}
	}
	return
}

// A DB is a persistent ordered map from keys to values.
// A DB is safe for concurrent access from multiple threads without
// any external synchronization.
type DB struct {
	db C.DB_t
	closed bool
}

func (db *DB) finalize() {
	if !db.closed {
		db.closed = true
		var cdb *C.DB_t = &db.db
		C.DeleteDBT(cdb, toCBool(false))
	}
}

func (db *DB) Close() {
	runtime.SetFinalizer(db, nil)
	db.finalize()
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

	var ccfds []C.ColumnFamilyDescriptor_t

	if cfds != nil {
		s := make([]interface{}, len(cfds))
		for i, v := range cfds {
			s[i] = v
		}
		ccfds = newCArrayFromColumnFamilyDescriptorArray(s...)
	}

	var (
		cdb *C.DB_t = &db.db
		opt *C.Options_t = &options.opt
		cstr *C.String_t = &rstr.str
		cfh *C.ColumnFamilyHandle_t
	)

	if ccfds != nil {
		cstat := C.DBOpenWithColumnFamilies(opt, cstr, &ccfds[0], C.int(len(ccfds)), &cfh, cdb)
		stat = cstat.toStatus()
		if stat.Ok() {
			cfhs = newColumnFamilyHandleArrayFromCArray(cfh, uint(len(ccfds)))
		}
	} else {
		cstat := C.DBOpen(opt, cstr, cdb)
		stat = cstat.toStatus()
	}

	if stat.Ok() {
		runtime.SetFinalizer(db, finalize)
	}

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
func OpenForReadOnly(options *Options, name *string, cfds ...interface{}) (db *DB, cfhs []*ColumnFamilyHandle, stat *Status) {
	db = &DB{}
	rstr := newCStringFromString(name)
	defer rstr.del()

	ccfds := newCArrayFromColumnFamilyDescriptorArray(cfds...)

	var (
		cdb *C.DB_t = &db.db
		opt *C.Options_t = &options.opt
		cstr *C.String_t = &rstr.str
		cfh *C.ColumnFamilyHandle_t
		cflg C.bool = toCBool(false)
	)

	if cfds != nil {
		n:= len(cfds)
		v, ok := cfds[n].(bool)
		if ok {
			cflg = toCBool(v)
		}
	}

	if ccfds != nil {
		cfdlen := len(ccfds)
		cstat := C.DBOpenForReadOnlyWithColumnFamilies(opt, cstr, &ccfds[0], C.int(cfdlen), &cfh, cdb, cflg)
		stat = cstat.toStatus()
		if stat.Ok() && cfdlen > 0 {
			cfhs = newColumnFamilyHandleArrayFromCArray(cfh, uint(cfdlen))
		}
	} else {
		cstat := C.DBOpenForReadOnly(opt, cstr, cdb, cflg)
		stat = cstat.toStatus()
	}

	if stat.Ok() {
		runtime.SetFinalizer(db, finalize)
	}
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
		opt *C.DBOptions_t = &dbopt.dbopt
		cstr *C.String_t = &rstr.str
		cfs *C.String_t
		sz C.int
	)

	cstat := C.DBListColumnFamilies(opt, cstr, &cfs, &sz)
	stat = cstat.toStatus()
	if stat.Ok() && sz > 0 {
		cfss = newStringArrayFromCArray(cfs, uint(sz))
	}
	return
} 


// Create a column_family and return the handle of column family
// through the argument handle.
func (db *DB) CreateColumnFamily(options *ColumnFamilyOptions, colfname *string) (cfd *ColumnFamilyHandle, stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	cstr := newCStringFromString(colfname)
	var (
		cdb *C.DB_t = &db.db
		opt *C.ColumnFamilyOptions_t = &options.cfopt
		ccstr *C.String_t = &cstr.str
		ccfd C.ColumnFamilyHandle_t
	)

	cstat := C.DBCreateColumnFamily(cdb, opt, ccstr, &ccfd)
	stat = cstat.toStatus()
	if stat.Ok() {
		cfd = ccfd.toColumnFamilyHandle()
	}
	return
}

// Drop a column family specified by column_family handle. This call
// only records a drop record in the manifest and prevents the column
// family from flushing and compacting.
func (db *DB) DropColumnFamily(cfh *ColumnFamilyHandle) (stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	var (
		cdb *C.DB_t = &db.db
		ccfh *C.ColumnFamilyHandle_t = &cfh.cfh 
	)
	cstat := C.DBDropColumnFamily(cdb, ccfh)
	stat = cstat.toStatus()
	return
}

// Set the database entry for "key" to "value".
// If "key" already exists, it will be overwritten.
// Returns OK on success, and a non-OK status on error.
// Note: consider setting options.sync = true.
func (db *DB) Put(options *WriteOptions, key, val []byte, cfh ...*ColumnFamilyHandle) (stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	ckey := newSliceFromBytes(key)
	defer ckey.del()
	cval := newSliceFromBytes(val)
	defer cval.del()

	var (
		cdb *C.DB_t = &db.db
		cwopt *C.WriteOptions_t = &options.wopt
		ccfh *C.ColumnFamilyHandle_t
		cckey *C.Slice_t = &ckey.slc 
		ccval *C.Slice_t = &cval.slc 
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	var cstat C.Status_t
	if ccfh != nil {
		cstat = C.DBPutWithColumnFamily(cdb, cwopt, ccfh, cckey, ccval)
	} else {
		cstat = C.DBPut(cdb, cwopt, cckey, ccval)
	}
	stat = cstat.toStatus()
	return
}

// Remove the database entry (if any) for "key".  Returns OK on
// success, and a non-OK status on error.  It is not an error if "key"
// did not exist in the database.
// Note: consider setting options.sync = true.
func (db *DB) Delete(options *WriteOptions, key []byte, cfh ...*ColumnFamilyHandle) (stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	ckey := newSliceFromBytes(key)
	defer ckey.del()

	var (
		cdb *C.DB_t = &db.db
		cwopt *C.WriteOptions_t = &options.wopt
		ccfh *C.ColumnFamilyHandle_t
		cckey *C.Slice_t = &ckey.slc 
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	var cstat C.Status_t
	if ccfh != nil {
		cstat = C.DBDeleteWithColumnFamily(cdb, cwopt, ccfh, cckey)
	} else {
		cstat = C.DBDelete(cdb, cwopt, cckey)
	}
	stat = cstat.toStatus()
	return
}

// Merge the database entry for "key" with "value".  Returns OK on success,
// and a non-OK status on error. The semantics of this operation is
// determined by the user provided merge_operator when opening DB.
// Note: consider setting options.sync = true.
func (db *DB) Merge(options *WriteOptions, key, val []byte, cfh ...*ColumnFamilyHandle) (stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	ckey := newSliceFromBytes(key)
	defer ckey.del()
	cval := newSliceFromBytes(val)
	defer cval.del()

	var (
		cdb *C.DB_t = &db.db
		cwopt *C.WriteOptions_t = &options.wopt
		ccfh *C.ColumnFamilyHandle_t
		cckey *C.Slice_t = &ckey.slc 
		ccval *C.Slice_t = &cval.slc 
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	var cstat C.Status_t
	if ccfh != nil {
		cstat = C.DBMergeWithColumnFamily(cdb, cwopt, ccfh, cckey, ccval)
	} else {
		cstat = C.DBMerge(cdb, cwopt, cckey, ccval)
	}
	stat = cstat.toStatus()
	return
}

// Apply the specified updates to the database.
// If `updates` contains no update, WAL will still be synced if
// options.sync=true.
// Returns OK on success, non-OK on failure.
// Note: consider setting options.sync = true.
func (db *DB) Write(options *WriteOptions, wbt *WriteBatch) (stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	var (
		cdb *C.DB_t = &db.db
		cwopt *C.WriteOptions_t = &options.wopt
		cwbt *C.WriteBatch_t = &wbt.wbt
	)

	defer wbt.mutex.Unlock()
	wbt.mutex.Lock()
	cstat := C.DBWrite(cdb, cwopt, cwbt)
	stat = cstat.toStatus()
	return
}

// If the database contains an entry for "key" store the
// corresponding value in *value and return OK.
//
// If there is no entry for "key" leave *value unchanged and return
// a status for which Status_t::IsNotFound() returns true.
//
// May return some other Status_t on an error.
func (db *DB) Get(options *ReadOptions, key []byte, cfh ...*ColumnFamilyHandle) (val []byte, stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	ckey := newSliceFromBytes(key)
	defer ckey.del()
	cval := newCString()

	var (
		cdb *C.DB_t = &db.db
		cropt *C.ReadOptions_t = &options.ropt
		ccfh *C.ColumnFamilyHandle_t
		cckey *C.Slice_t = &ckey.slc 
		ccval *C.String_t = &cval.str 
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	var cstat C.Status_t
	if ccfh != nil {
		cstat = C.DBGetWithColumnFamily(cdb, cropt, ccfh, cckey, ccval)
	} else {
		cstat = C.DBGet(cdb, cropt, cckey, ccval)
	}
	stat = cstat.toStatus()
	val = cval.goBytes(true)
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
func (db *DB) MultiGet(options *ReadOptions, keys [][]byte, cfhs ...*ColumnFamilyHandle) (vals [][]byte, stats []*Status) {
	if db.closed {
		stats = make([]*Status, len(keys))
		stat := NewDBClosedStatus()
		for i, _ := range stats {
			stats[i] = stat	
		}
		return
	}

	ckeys := cSlicePtrAry(newSlicesFromBytesArray(keys))
	defer ckeys.del()
	cckeys := ckeys.toCArray()
	ccfhs := newCArrayFromColumnFamilyHandleArray(cfhs...)

	var (
		cdb *C.DB_t = &db.db
		cropt *C.ReadOptions_t = &options.ropt
		ccvals *C.String_t
	)
	
	n := len(cckeys)
	if ccfhs != nil {
		stats = newStatusArrayFromCArray(C.DBMultiGetWithColumnFamily(cdb, cropt, &ccfhs[0], C.int(len(cfhs)), &cckeys[0], C.int(n), &ccvals), uint(n))
	} else {
		stats = newStatusArrayFromCArray(C.DBMultiGet(cdb, cropt, &cckeys[0], C.int(n), &ccvals), uint(n))
	}
	vals = newBytesFromCArray(ccvals, uint(n))
	return
}

// If the key definitely does not exist in the database, then this method
// returns false, else true. If the caller wants to obtain value when the key
// is found in memory, a bool for 'value_found' must be passed. 'value_found'
// will be true on return if value has been set properly.
// This check is potentially lighter-weight than invoking DB::Get(). One way
// to make this lighter weight is to avoid doing any IOs.
// Default implementation here returns true and sets 'value_found' to false
func (db *DB) KeyMayExist(options *ReadOptions, key []byte, cfh ...*ColumnFamilyHandle) (res bool, valfound bool, val string) {
	if db.closed {
		return
	}

	ckey := newSliceFromBytes(key)
	defer ckey.del()
	cval := newCString()

	var (
		cdb *C.DB_t = &db.db
		cropt *C.ReadOptions_t = &options.ropt
		ccfh *C.ColumnFamilyHandle_t
		cckey *C.Slice_t = &ckey.slc 
		ccval *C.String_t = &cval.str 
		cvalfound C.bool
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	if ccfh != nil {
		res = C.DBKeyMayExistWithColumnFamily(cdb, cropt, ccfh, cckey, ccval, &cvalfound).toBool()
	} else {
		res = C.DBKeyMayExist(cdb, cropt, cckey, ccval, &cvalfound).toBool()
	}
	valfound = cvalfound.toBool()
	val = cval.goString(true)
	return
}

// Return a heap-allocated iterator over the contents of the database.
// The result of NewIterator() is initially invalid (caller must
// call one of the Seek methods on the iterator before using it).
//
// Caller should delete the iterator when it is no longer needed.
// The returned iterator should be deleted before this db is deleted.
func (db *DB) NewIterator(options *ReadOptions, cfh ...*ColumnFamilyHandle) (it *Iterator) {
	if db.closed {
		return
	}

	var (
		cdb *C.DB_t = &db.db
		cropt *C.ReadOptions_t = &options.ropt
		ccfh *C.ColumnFamilyHandle_t
		cit C.Iterator_t
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	if ccfh != nil {
		cit = C.DBNewIteratorWithColumnFamily(cdb, cropt, ccfh)
	} else {
		cit = C.DBNewIterator(cdb, cropt)
	}
	it = cit.toIterator(db)
	return
}

// Returns iterators from a consistent database state across multiple
// column families. Iterators are heap allocated and need to be deleted
// before the db is deleted
func (db *DB) NewIterators(options *ReadOptions, cfhs []*ColumnFamilyHandle) (vals []*Iterator, stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	ccfhs := newCArrayFromColumnFamilyHandleArray(cfhs...)

	var (
		cdb *C.DB_t = &db.db
		cropt *C.ReadOptions_t = &options.ropt
		ccvals *C.Iterator_t
		valsz C.int
	)

	cstat := C.DBNewIterators(cdb, cropt, &ccfhs[0], C.int(len(ccfhs)), &ccvals, &valsz)
	stat = cstat.toStatus()
	vals = newIteratorArrayFromCArray(ccvals, uint(valsz), db)
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
	if db.closed {
		return
	}

	var cdb *C.DB_t = &db.db
	var csnp C.Snapshot_t = C.DBGetSnapshot(cdb)

	snp = csnp.toSnapshot(db)
	return
}

// Release a previously acquired snapshot.  The caller must not
// use "snapshot" after this call.
func (db *DB) ReleaseSnapshot(snp *Snapshot) {
	if db.closed {
		return
	}

	if snp.db != db {
		panic("ReleaseSnapshot error!")
	}
	snp.db = nil

	var (
		cdb *C.DB_t = &db.db
		csnp *C.Snapshot_t = &snp.snp
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
func (db *DB) GetProperty(prop []byte, cfh ...*ColumnFamilyHandle) (val string, res bool) {
	if db.closed {
		return
	}

	cprop := newSliceFromBytes(prop)
	defer cprop.del()
	cval := newCString()

	var (
		cdb *C.DB_t = &db.db
		ccfh *C.ColumnFamilyHandle_t
		ccprop *C.Slice_t = &cprop.slc 
		ccval *C.String_t = &cval.str 
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	if ccfh != nil {
		res = C.DBGetPropertyWithColumnFamily(cdb, ccfh, ccprop, ccval).toBool()
	} else {
		res = C.DBGetProperty(cdb, ccprop, ccval).toBool()
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
func (db *DB) GetIntProperty(prop []byte, cfh ...*ColumnFamilyHandle) (val uint64, res bool) {
	if db.closed {
		return
	}

	cprop := newSliceFromBytes(prop)
	defer cprop.del()

	var (
		cdb *C.DB_t = &db.db
		ccfh *C.ColumnFamilyHandle_t
		ccprop *C.Slice_t = &cprop.slc 
		cval C.uint64_t
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	if ccfh != nil {
		res = C.DBGetIntPropertyWithColumnFamily(cdb, ccfh, ccprop, &cval).toBool()
	} else {
		res = C.DBGetIntProperty(cdb, ccprop, &cval).toBool()
	}
	val = uint64(cval)
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
func (db *DB) GetApproximateSizes(rngs []*Range, cfh ...*ColumnFamilyHandle) (vals []uint64) {
	if db.closed {
		return
	}

	crngs := newCArrayFromRangeArray(rngs...)

	var (
		cdb *C.DB_t = &db.db
		ccfh *C.ColumnFamilyHandle_t
		ccrngs *C.Range_t = &crngs[0] 
		sz C.int = C.int(len(crngs))
		cval []C.uint64_t = make([]C.uint64_t, sz)
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	if ccfh != nil {
		C.DBGetApproximateSizesWithColumnFamily(cdb, ccfh, ccrngs, sz, &cval[0])
	} else {
		C.DBGetApproximateSizes(cdb, ccrngs, sz, &cval[0])
	}
	vals = newUint64ArrayFromCArray(&cval)
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
func (db *DB) CompactRange(begin []byte, end []byte, cfhs ...interface{}) (stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	cbegin := newSliceFromBytes(begin)
	defer cbegin.del()
	cend := newSliceFromBytes(end)
	defer cend.del()
	var ccfhs []C.ColumnFamilyHandle_t

	if cfhs != nil {
		ccfhs = newCArrayFromColumnFamilyHandleInterface(cfhs...)
	}

	var (
		cdb *C.DB_t = &db.db
		credl C.bool 
		ctarl C.int = -1
		ctpi C.uint32_t 
	)

	if cfhs != nil {
		n:= len(cfhs)
		v, ok := cfhs[n].(uint32)
		if ok {
			ctpi = C.uint32_t(v)
			n--
		}

		if 0 < n {
			v, ok := cfhs[n].(int)
			if ok {
				ctarl = C.int(v)
				n--
			}

			if 0 < n {
				v, ok := cfhs[n].(bool)
				if ok {
					credl = toCBool(v)
				}
			}
		}
	}

	var cstat C.Status_t
	if ccfhs != nil {
		cstat = C.DBCompactRangeWithColumnFamily(cdb, &ccfhs[0], &cbegin.slc, &cend.slc, credl, ctarl, ctpi)
	} else {
		cstat = C.DBCompactRange(cdb, &cbegin.slc, &cend.slc, credl, ctarl, ctpi)
	}
	stat = cstat.toStatus()
	return
}

func (db *DB) SetOptions(opts []string, cfhs ...*ColumnFamilyHandle) (stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	copts := cStringPtrAry(newcStringsFromStringArray(opts))
	defer copts.del()
	ccopts := copts.toCArray()

	var (
		cdb *C.DB_t = &db.db
		ccfh *C.ColumnFamilyHandle_t
	)

	if cfhs != nil {
		ccfh = &cfhs[0].cfh
	}

	var cstat C.Status_t
	if ccfh != nil {
		cstat = C.DBSetOptionsWithColumnFamily(cdb, ccfh, &ccopts[0], C.int(len(ccopts)))
	} else {
		cstat = C.DBSetOptions(cdb, &ccopts[0], C.int(len(ccopts)))
	}
	stat = cstat.toStatus()
	return
}

// CompactFiles() inputs a list of files specified by file numbers
// and compacts them to the specified level.  Note that the behavior
// is different from CompactRange in that CompactFiles() will
// perform the compaction job using the CURRENT thread.
//
// @see GetDataBaseMetaData
// @see GetColumnFamilyMetaData
func (db *DB) CompactFiles(options *CompactionOptions, files []string, level int, cfhs ...interface{}) (stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	cfiles := cStringPtrAry(newcStringsFromStringArray(files))
	defer cfiles.del()
	ccfiles := cfiles.toCArray()

	var (
		cdb *C.DB_t = &db.db
		ccfh *C.ColumnFamilyHandle_t
		ccopt *C.CompactionOptions_t = &options.copt
		cpid C.int = -1
	)

	if cfhs != nil {
		n:= len(cfhs)
		v, ok := cfhs[n].(int)
		if ok {
			cpid = C.int(v)
			n--
		}

		if 0 < n {
			v, ok := cfhs[n].(*ColumnFamilyHandle)
			if ok {
				ccfh = &v.cfh
				n--
			}
		}
	}

	var cstat C.Status_t
	if ccfh != nil {
		cstat = C.DBCompactFilesWithColumnFamily(cdb, ccopt, ccfh, &ccfiles[0], C.int(len(ccfiles)), C.int(level), cpid)
	} else {
		cstat = C.DBCompactFiles(cdb, ccopt, &ccfiles[0], C.int(len(ccfiles)), C.int(level), cpid)
	}
	stat = cstat.toStatus()
	return
}

// Number of levels used for this DB.
func (db *DB) NumberLevels(cfh ...*ColumnFamilyHandle) (level int) {
	if db.closed {
		return
	}

	var (
		cdb *C.DB_t = &db.db
		ccfh *C.ColumnFamilyHandle_t
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	if ccfh != nil {
		level = int(C.DBNumberLevelsWithColumnFamily(cdb, ccfh))
	} else {
		level = int(C.DBNumberLevels(cdb))
	}
	return
}

// Maximum level to which a new compacted memtable is pushed if it
// does not create overlap.
func (db *DB) MaxMemCompactionLevel(cfh ...*ColumnFamilyHandle) (level int) {
	if db.closed {
		return
	}

	var (
		cdb *C.DB_t = &db.db
		ccfh *C.ColumnFamilyHandle_t
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	if ccfh != nil {
		level = int(C.DBMaxMemCompactionLevelWithColumnFamily(cdb, ccfh))
	} else {
		level = int(C.DBMaxMemCompactionLevel(cdb))
	}
	return
}

// Number of files in level-0 that would stop writes.
func (db *DB) Level0StopWriteTrigger(cfh ...*ColumnFamilyHandle) (level int) {
	if db.closed {
		return
	}

	var (
		cdb *C.DB_t = &db.db
		ccfh *C.ColumnFamilyHandle_t
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	if ccfh != nil {
		level = int(C.DBLevel0StopWriteTriggerWithColumnFamily(cdb, ccfh))
	} else {
		level = int(C.DBLevel0StopWriteTrigger(cdb))
	}
	return
}

// Get DB name -- the exact same name that was provided as an argument to
// DB::Open()
func (db *DB) GetName() (name string) {
	if db.closed {
		return
	}

	var (
		cdb *C.DB_t = &db.db
		cname C.String_t = C.DBGetName(cdb)
	)
	name = cname.cToString()
	return
}

// Get Env object from the DB
func (db *DB) GetEnv() (env *Env) {
	if db.closed {
		return
	}

	var (
		cdb *C.DB_t = &db.db
		cenv C.Env_t = C.DBGetEnv(cdb)
	)
	env = cenv.toEnv()
	return
}

// Get DB Options that we use.  During the process of opening the
// column family, the options provided when calling DB::Open() or
// DB::CreateColumnFamily() will have been "sanitized" and transformed
// in an implementation-defined manner.
func (db *DB) GetOptions(cfh ...*ColumnFamilyHandle) (opt *Options) {
	if db.closed {
		return
	}

	var (
		cdb *C.DB_t = &db.db
		ccfh *C.ColumnFamilyHandle_t
		copt C.Options_t
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	if ccfh != nil {
		copt = C.DBGetOptionsWithColumnFamily(cdb, ccfh)
	} else {
		copt = C.DBGetOptions(cdb)
	}

	opt = copt.toOptions()
	return
}

func (db *DB) GetDBOptions() (dbopt *DBOptions) {
	if db.closed {
		return
	}

	var (
		cdb *C.DB_t = &db.db
		cdbopt C.DBOptions_t = C.DBGetDBOptions(cdb)
	)
	dbopt = cdbopt.toDBOptions()
	return
}

// Flush all mem-table data.
func (db *DB) Flush(options *FlushOptions, cfhs ...*ColumnFamilyHandle) (stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	var (
		cdb *C.DB_t = &db.db
		ccfh *C.ColumnFamilyHandle_t
		cfopt *C.FlushOptions_t = &options.fopt
	)

	if cfhs != nil {
		ccfh = &cfhs[0].cfh
	}

	var cstat C.Status_t
	if ccfh != nil {
		cstat = C.DBFlushWithColumnFamily(cdb, cfopt, ccfh)
	} else {
		cstat = C.DBFlush(cdb, cfopt)
	}
	stat = cstat.toStatus()
	return
}

// The sequence number of the most recent transaction.
func (db *DB) GetLatestSequenceNumber() (sqnum SequenceNumber) {
	if db.closed {
		return
	}

	var (
		cdb *C.DB_t = &db.db
	)
	sqnum = SequenceNumber(C.DBGetLatestSequenceNumber(cdb))
	return
}

// Sets the globally unique ID created at database creation time by invoking
// Env::GenerateUniqueId(), in identity. Returns Status_t::OK if identity could
// be set properly
func (db *DB) GetDbIdentity() (id string, stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	var (
		cdb *C.DB_t = &db.db
		cid C.String_t
	)
	cstat := C.DBGetDbIdentity(cdb, &cid)
	stat = cstat.toStatus()
	id = cid.cToString()
	return
}

// Returns default column family handle
func (db *DB) DefaultColumnFamily() (cfh *ColumnFamilyHandle) {
	if db.closed {
		return
	}

	var (
		cdb *C.DB_t = &db.db
		ccfh C.ColumnFamilyHandle_t = C.DBDefaultColumnFamily(cdb)
	)
	cfh = ccfh.toColumnFamilyHandle()
	return
}

// Destroy the contents of the specified database.
// Be very careful using this method.
func DestroyDB(opt *Options, name *string) (stat *Status) {
	cname := newCStringFromString(name)

	var (
		ccname *C.String_t = &cname.str
		copt *C.Options_t = &opt.opt
	)

	cstat := C.DBDestroyDB(ccname, copt)
	stat = cstat.toStatus()
	return
}
