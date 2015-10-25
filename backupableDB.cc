// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#include <rocksdb/utilities/backupable_db.h>

using namespace rocksdb;

#include "backupableDB.h"

DEFINE_C_WRAP_CONSTRUCTOR(BackupableDBOptions)
DEFINE_C_WRAP_DESTRUCTOR(BackupableDBOptions)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS(BackupableDBOptions, String)

DEFINE_C_WRAP_CONSTRUCTOR(RestoreOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(RestoreOptions)
DEFINE_C_WRAP_DESTRUCTOR(RestoreOptions)

// If true, restore won't overwrite the existing log files in wal_dir. It will
// also move all log files from archive directory to wal_dir. Use this option
// in combination with BackupableDBOptions::backup_log_files = false for
// persisting in-memory databases.
// Default: false
DEFINE_C_WRAP_SETTER(RestoreOptions, keep_log_files, bool)


DEFINE_C_WRAP_CONSTRUCTOR(BackupEngine)
DEFINE_C_WRAP_DESTRUCTOR(BackupEngine)

// Create a new backup engine from db_env and options.
Status_t BackupEngineOpen(Env_t* db_env, const BackupableDBOptions_t* options, BackupEngine_t* backup_engine_ptr)
{
    assert(db_env != NULL);
    assert(options != NULL);
    assert(backup_engine_ptr != NULL);
    BackupEngine** backup_engine = GET_REP_ADDR(backup_engine_ptr, BackupEngine);
    Status stat = BackupEngine::Open(GET_REP(db_env, Env), GET_REP_REF(options, BackupableDBOptions), backup_engine);
    return NewStatusTCopy(&stat);
}

// Create a new backup for db.
Status_t BackupEngineCreateNewBackup(BackupEngine_t* backup_engine_ptr, DB_t* db, bool flush_before_backup)
{
    assert(db != NULL);
    assert(backup_engine_ptr != NULL);
    assert(GET_REP(backup_engine_ptr, BackupEngine) != NULL);
    Status stat = GET_REP(backup_engine_ptr, BackupEngine)->CreateNewBackup(GET_REP(db, DB), flush_before_backup);
    return NewStatusTCopy(&stat);
}

// Restoring DB from backup is NOT safe when there is another BackupEngine
// running that might call DeleteBackup() or PurgeOldBackups(). It is caller's
// responsibility to synchronize the operation, i.e. don't delete the backup
// when you're restoring from it
Status_t BackupEngineRestoreDBFromBackup(BackupEngine_t* backup_engine_ptr, BackupID backup_id, String_t* db_dir, String_t* wal_dir, RestoreOptions_t* restore_options)
{
    assert(db_dir != NULL);
    assert(GET_REP(db_dir, String) != NULL);
    assert(wal_dir != NULL);
    assert(GET_REP(wal_dir, String) != NULL);
    assert(backup_engine_ptr != NULL);
    assert(GET_REP(backup_engine_ptr, BackupEngine) != NULL);
    Status stat = GET_REP(backup_engine_ptr, BackupEngine)->RestoreDBFromBackup(backup_id,
                                                                                GET_REP_REF(db_dir, String),
                                                                                GET_REP_REF(wal_dir, String),
                                                                                NULL == restore_options ? RestoreOptions() : GET_REP_REF(restore_options, RestoreOptions));
    return NewStatusTCopy(&stat);
}

Status_t BackupEngineRestoreDBFromLatestBackup(BackupEngine_t* backup_engine_ptr, String_t* db_dir, String_t* wal_dir, RestoreOptions_t* restore_options)
{
    assert(db_dir != NULL);
    assert(GET_REP(db_dir, String) != NULL);
    assert(wal_dir != NULL);
    assert(GET_REP(wal_dir, String) != NULL);
    assert(backup_engine_ptr != NULL);
    assert(GET_REP(backup_engine_ptr, BackupEngine) != NULL);
    Status stat = GET_REP(backup_engine_ptr, BackupEngine)->RestoreDBFromLatestBackup(GET_REP_REF(db_dir, String),
                                                                                      GET_REP_REF(wal_dir, String),
                                                                                      NULL == restore_options ? RestoreOptions() : GET_REP_REF(restore_options, RestoreOptions));
    return NewStatusTCopy(&stat);
}

DEFINE_C_WRAP_CONSTRUCTOR(BackupableDB)
DEFINE_C_WRAP_DESTRUCTOR(BackupableDB)

DEFINE_C_WRAP_CONSTRUCTOR(RestoreBackupableDB)
DEFINE_C_WRAP_DESTRUCTOR(RestoreBackupableDB)
