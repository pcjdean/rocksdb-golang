// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// +build !lite

package rocksdb

/*
#cgo CFLAGS: -DPNG_DEBUG=1
#cgo LDFLAGS: -lrocksdb -lstdc++ -lz -lrt
#include "db.h"
*/
import "C"

// Prevent file deletions. Compactions will continue to occur,
// but no obsolete files will be deleted. Calling this multiple
// times have the same effect as calling it once.
func (db *DB) DisableFileDeletions() (stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	var (
		cdb *C.DB_t = &db.db
	)

	cstat := C.DBDisableFileDeletions(cdb)
	stat = cstat.toStatus()
	return
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
func (db *DB) EnableFileDeletions(force ...bool) (stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	var (
		cdb *C.DB_t = &db.db
		cforce C.bool = toCBool(true)
	)

	if force != nil {
		cforce = toCBool(force[0])
	}

	cstat := C.DBEnableFileDeletions(cdb, cforce)
	stat = cstat.toStatus()
	return
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
func (db *DB) GetLiveFiles(flush_memtable ...bool) (files []string, fileSz uint64, stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	var (
		cdb *C.DB_t = &db.db
		cflhmem C.bool = toCBool(true)
		n C.int
		cfiles *C.String_t
		cfsz C.uint64_t
	)

	if flush_memtable != nil {
		cflhmem = toCBool(flush_memtable[0])
	}

	cstat := C.DBGetLiveFiles(cdb, &cfiles, &n, &cfsz, cflhmem)
	stat = cstat.toStatus()
	fileSz = uint64(cfsz)
	files = newStringArrayFromCArray(cfiles, uint(n))
	return
}

// Retrieve the sorted list of all wal files with earliest file first
func (db *DB) GetSortedWalFiles() (files []*LogFile, stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	var (
		cdb *C.DB_t = &db.db
		n C.int
		cfiles *C.LogFile_t
	)

	cstat := C.DBGetSortedWalFiles(cdb, &cfiles, &n)
	stat = cstat.toStatus()
	files = newLogFileArrayFromCArray(cfiles, uint(n))
	return
}

// Sets iter to an iterator that is positioned at a write-batch containing
// seq_number. If the sequence number is non existent, it returns an iterator
// at the first available seq_no after the requested seq_no
// Returns Status_t::OK if iterator is valid
// Must set WAL_ttl_seconds or WAL_size_limit_MB to large values to
// use this api, else the WAL files will get
// cleared aggressively and the iterator might keep getting invalid before
// an update is read.
func (db *DB) GetUpdatesSince(sqn SequenceNumber, tranropt ...TransactionLogIteratorReadOptions) (it *TransactionLogIterator, stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	var (
		cdb *C.DB_t = &db.db
		cit C.TransactionLogIterator_t
		ctranropt *C.TransactionLogIterator_ReadOptions_t
	)

	if tranropt != nil {
		ctranropt = &tranropt[0].tranropt	
	}

	cstat := C.DBGetUpdatesSince(cdb, C.SequenceNumber(sqn), &cit, ctranropt)
	stat = cstat.toStatus()
	it = cit.toTransactionLogIterator()
	return
}

// Delete the file name from the db directory and update the internal state to
// reflect that. Supports deletion of sst and log files only. 'name' must be
// path relative to the db directory. eg. 000001.sst, /archive/000003.log
func (db *DB) DeleteFile(name string) (stat *Status) {
	var (
		cdb *C.DB_t = &db.db
		cname C.String_t
	)
	cstat := C.DBDeleteFile(cdb, &cname)
	stat = cstat.toStatus()
	name = cname.cToString()
	return
}

// Returns a list of all table files with their level, start key
// and end key
func (db *DB) GetLiveFilesMetaData() (lfmds []*LiveFileMetaData) {
	if db.closed {
		return
	}

	var (
		cdb *C.DB_t = &db.db
		clfmds *C.LiveFileMetaData_t
		valsz C.int
	)

	C.DBGetLiveFilesMetaData(cdb, &clfmds, &valsz)
	lfmds = newLiveFileMetaDataArrayFromCArray(clfmds, uint(valsz))
	return
}

// Obtains the meta data of the specified column family of the DB.
// Status_t::NotFound() will be returned if the current DB does not have
// any column family match the specified name.
//
// If cf_name is not specified, then the metadata of the default
// column family will be returned.
func (db *DB) GetColumnFamilyMetaData(cfh ...*ColumnFamilyHandle) (md *ColumnFamilyMetaData) {
	if db.closed {
		return
	}

	var (
		cdb *C.DB_t = &db.db
		ccfh *C.ColumnFamilyHandle_t
		ccfmd C.ColumnFamilyMetaData_t
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	if ccfh != nil {
		C.DBGetColumnFamilyMetaDataWithColumnFamily(cdb, ccfh, &ccfmd)
	} else {
		C.DBGetColumnFamilyMetaData(cdb, &ccfmd)
	}
	md = ccfmd.toColumnFamilyMetaData()
	return
}

func (db *DB) GetPropertiesOfAllTables(cfh ...*ColumnFamilyHandle) (tpc *TablePropertiesCollection, stat *Status) {
	if db.closed {
		stat = NewDBClosedStatus()
		return
	}

	var (
		cdb *C.DB_t = &db.db
		ccfh *C.ColumnFamilyHandle_t
		ctpc C.TablePropertiesCollection_t
	)

	if cfh != nil {
		ccfh = &cfh[0].cfh
	}

	var cstat C.Status_t
	if ccfh != nil {
		cstat = C.DBGetPropertiesOfAllTablesWithColumnFamily(cdb, ccfh, &ctpc)
	} else {
		cstat = C.DBGetPropertiesOfAllTables(cdb, &ctpc)
	}
	stat = cstat.toStatus()
	tpc = ctpc.toTablePropertiesCollection()
	return
}

// If a DB cannot be opened, you may attempt to call this method to
// resurrect as much of the contents of the database as possible.
// Some data may be lost, so be careful when calling this function
// on a database that contains important information.
func RepairDB(opt *Options, name *string) (stat *Status) {
	cname := newCStringFromString(name)

	var (
		ccname *C.String_t = &cname.str
		copt *C.Options_t = &opt.opt
	)

	cstat := C.DBRepairDB(ccname, copt)
	stat = cstat.toStatus()
	return
}
