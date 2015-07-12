// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// +build !lite

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

// Prevent file deletions. Compactions will continue to occur,
// but no obsolete files will be deleted. Calling this multiple
// times have the same effect as calling it once.
func (db *DB) DisableFileDeletions() (stat *Status) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
	)

	stat = C.DBDisableFileDeletions(cdb).toStatus()
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
func (db *DB) EnableFileDeletions() (force ...bool) (stat *Status)  {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cforce C.bool = true
	)

	if force {
		cforce = force
	}

	stat = C.DBEnableFileDeletions(cdb, cforce).toStatus()
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
func (db *DB) GetLiveFiles() (flush_memtable ...bool) (files []string, fileSz uint64, stat *Status) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cflhmem C.bool = true
		n C.int
		cfiles *C.String_t
	)

	if flush_memtable {
		cflhmem = flush_memtable
	}

	stat = C.DBGetLiveFiles(cdb, unsafe.Pointers(&cfiles), unsafe.Pointers(&n), unsafe.Pointers(&fileSz), cflhmem).toStatus()
	files = newStringArrayFromCArray(cfiles, n)
	return
}

// Retrieve the sorted list of all wal files with earliest file first
func (db *DB) GetSortedWalFiles() (files []*LogFile, stat *Status) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		n C.int
		cfiles *C.LogFile_t
	)

	stat = C.DBGetSortedWalFiles(cdb, unsafe.Pointers(&cfiles), unsafe.Pointers(&n)).toStatus()
	files = newLogFileArrayFromCArray(cfiles, n)
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
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cit C.TransactionLogIterator_t
		ctranropt *C.TransactionLogIterator_ReadOptions_t
	)

	if tranropt {
		ctranropt = unsafe.Pointers(&tranropt.tranropt)	
	}

	stat = C.DBGetUpdatesSince(cdb, sqn, unsafe.Pointers(&cit), ctranropt).toStatus()
	it = cit.toTransactionLogIterator()
	return
}

// Delete the file name from the db directory and update the internal state to
// reflect that. Supports deletion of sst and log files only. 'name' must be
// path relative to the db directory. eg. 000001.sst, /archive/000003.log
func (db *DB) DeleteFile() (name string, stat *Status) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		cname C.String_t
	)
	stat = C.DBDeleteFile(cdb, unsafe.Pointers(&cname)).toStatus()
	name = cname.cToString()
	return
}

// Returns a list of all table files with their level, start key
// and end key
func (db *DB) GetLiveFilesMetaData() (lfmds []*LiveFileMetaData) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		clfmds *C.LiveFileMetaData_t
		valsz int
	)

	C.DBGetLiveFilesMetaData(cdb, unsafe.Pointers(&cclfmds), len(valsz))
	vals = newLiveFileMetaDataArrayFromCArray(cclfmds, valsz)
	return
}

// Obtains the meta data of the specified column family of the DB.
// Status_t::NotFound() will be returned if the current DB does not have
// any column family match the specified name.
//
// If cf_name is not specified, then the metadata of the default
// column family will be returned.
func (db *DB) GetColumnFamilyMetaData(cfd ...*ColumnFamilyHandle) (md *ColumnFamilyMetaData) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		ccfd *C.ColumnFamilyHandle_t
		ccfmd C.ColumnFamilyMetaData_t
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		level = C.DBGetColumnFamilyMetaDataWithColumnFamily(cdb, ccfd, unsafe.Pointers(&ccfmd))
	} else {
		level = C.DBGetColumnFamilyMetaData(cdb, unsafe.Pointers(&ccfmd))
	}
	md = ccfmd.toColumnFamilyMetaData()
	return
}

func (db *DB) GetPropertiesOfAllTables(cfd ...*ColumnFamilyHandle) (tpc *TablePropertiesCollection, stat *Status) {
	var (
		cdb *C.DB_t = unsafe.Pointers(&db.db)
		ccfd *C.ColumnFamilyHandle_t
		ctpc C.TablePropertiesCollection_t
	)

	if cfd {
		cfd[0].(*ColumnFamilyHandle)
		ccfd = unsafe.Pointers(&cfd[0].cfd)
	}

	if ccfd {
		stat = C.DBGetPropertiesOfAllTablesWithColumnFamily(cdb, ccfd, unsafe.Pointers(&ctpc)).toStatus()
	} else {
		stat = C.DBGetPropertiesOfAllTables(cdb, unsafe.Pointers(&ctpc)).toStatus()
	}
	tpc = ctpc.toTablePropertiesCollection()
	return
}

// If a DB cannot be opened, you may attempt to call this method to
// resurrect as much of the contents of the database as possible.
// Some data may be lost, so be careful when calling this function
// on a database that contains important information.
func RepairDB(name string, opt Options) (stat *Status) {
	cname := newCStringFromString(name)

	var (
		ccname *C.String_t = unsafe.Pointers(&cname.str)
		copt C.Options_t = unsafe.Pointers(&opt.opt)
	)

	stat = C.DBRepairDB(ccname, copt).toStatus()
	return
}
