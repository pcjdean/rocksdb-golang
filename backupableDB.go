// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package rocksdb

/*
#include "backupableDB.h"
*/
import "C"

import (
	"runtime"
)

// Wrap go BackupableDBOptions
type BackupableDBOptions struct {
	bdbop C.BackupableDBOptions_t
	// true if bdbop is deleted
	closed bool
}

// Release resources
func (bdbop *BackupableDBOptions) finalize() {
	if !bdbop.closed {
		bdbop.closed = true
		var cbdbop *C.BackupableDBOptions_t = &bdbop.bdbop
		C.DeleteBackupableDBOptionsT(cbdbop, toCBool(false))
	}
}

// Close the @BackupableDBOptions
func (bdbop *BackupableDBOptions) Close() {
	runtime.SetFinalizer(bdbop, nil)
	bdbop.finalize()
}

// C BackupableDBOptions to go BackupableDBOptions
func (cbdbop *C.BackupableDBOptions_t) toBackupableDBOptions() (bdbop *BackupableDBOptions) {
	bdbop = &BackupableDBOptions{bdbop: *cbdbop}	
	runtime.SetFinalizer(bdbop, finalize)
	return
}

// Create a new BackupableDBOptions with a directory path @dpath
func NewBackupableDBOptions(dpath *string) *BackupableDBOptions {
	cstr := newCStringFromString(dpath)
	defer cstr.del()
	cbdbop := C.NewBackupableDBOptionsTArgs(&cstr.str)
	return cbdbop.toBackupableDBOptions()
}

// Wrap go RestoreOptions
type RestoreOptions struct {
	rsop C.RestoreOptions_t
	// true if rsop is deleted
	closed bool
}

// Release resources
func (rsop *RestoreOptions) finalize() {
	if !rsop.closed {
		rsop.closed = true
		var crsop *C.RestoreOptions_t = &rsop.rsop
		C.DeleteRestoreOptionsT(crsop, toCBool(false))
	}
}

// Close the @RestoreOptions
func (rsop *RestoreOptions) Close() {
	runtime.SetFinalizer(rsop, nil)
	rsop.finalize()
}

// C RestoreOptions to go RestoreOptions
func (crsop *C.RestoreOptions_t) toRestoreOptions() (rsop *RestoreOptions) {
	rsop = &RestoreOptions{rsop: *crsop}	
	runtime.SetFinalizer(rsop, finalize)
	return
}

// If true, restore won't overwrite the existing log files in wal_dir. It will
// also move all log files from archive directory to wal_dir. Use this option
// in combination with BackupableDBOptions::backup_log_files = false for
// persisting in-memory databases.
// Default: false
func (rsop *RestoreOptions) SetKeepLogFile(val bool) {
	var crsop *C.RestoreOptions_t = &rsop.rsop
	C.RestoreOptions_set_keep_log_files(crsop, toCBool(val))
}

// Create a new RestoreOptions with a directory path @dpath
func NewRestoreOptions() *RestoreOptions {
	crsop := C.NewRestoreOptionsTDefault()
	return crsop.toRestoreOptions()
}

// Wrap go BackupEngine
type BackupEngine struct {
	beg C.BackupEngine_t
	// true if beg is deleted
	closed bool
}

// Release resources
func (beg *BackupEngine) finalize() {
	if !beg.closed {
		beg.closed = true
		var cbeg *C.BackupEngine_t = &beg.beg
		C.DeleteBackupEngineT(cbeg, toCBool(false))
	}
}

// Close the @BackupEngine
func (beg *BackupEngine) Close() {
	runtime.SetFinalizer(beg, nil)
	beg.finalize()
}

// C BackupEngine to go BackupEngine
func (cbeg *C.BackupEngine_t) toBackupEngine() (beg *BackupEngine) {
	beg = &BackupEngine{beg: *cbeg}	
	runtime.SetFinalizer(beg, finalize)
	return
}

// Create a new backup engine from db_env and options.
func BackupEngineOpen (env *Env, options *BackupableDBOptions) (beg *BackupEngine, stat *Status) {
	beg = &BackupEngine{}
	cstat := C.BackupEngineOpen(&env.env, &options.bdbop, &beg.beg)
	stat = cstat.toStatus()
	return
}

// Create a new backup for db.
func (beg *BackupEngine) CreateNewBackup(db *DB, flush_before_backup ...interface{}) (stat *Status) {
	var fbb bool = false
	if nil != flush_before_backup && 0 < len(flush_before_backup) {
		fbb = flush_before_backup[0].(bool)
	}
	cstat := C.BackupEngineCreateNewBackup(&beg.beg, &db.db, toCBool(fbb))
	stat = cstat.toStatus()
	return
}

// Restoring DB from backup is NOT safe when there is another BackupEngine
// running that might call DeleteBackup() or PurgeOldBackups(). It is caller's
// responsibility to synchronize the operation, i.e. don't delete the backup
// when you're restoring from it
func (beg *BackupEngine) RestoreDBFromBackup(id uint32, dbdir *string, waldir *string, rsop *RestoreOptions) (stat *Status) {
	cdbd := newCStringFromString(dbdir)
	defer cdbd.del()
	cwald := newCStringFromString(waldir)
	defer cwald.del()
	cstat := C.BackupEngineRestoreDBFromBackup(&beg.beg, C.BackupID(id), &cdbd.str, &cwald.str, &rsop.rsop)
	stat = cstat.toStatus()
	return
}

func (beg *BackupEngine) RestoreDBFromLatestBackup(dbdir *string, waldir *string, rsop *RestoreOptions) (stat *Status) {
	cdbd := newCStringFromString(dbdir)
	defer cdbd.del()
	cwald := newCStringFromString(waldir)
	defer cwald.del()
	cstat := C.BackupEngineRestoreDBFromLatestBackup(&beg.beg, &cdbd.str, &cwald.str, &rsop.rsop)
	stat = cstat.toStatus()
	return
}
