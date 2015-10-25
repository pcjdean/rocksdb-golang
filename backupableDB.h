// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_BACKUPABLE_DB_H_
#define GO_ROCKSDB_INCLUDE_BACKUPABLE_DB_H_

#include "types.h"
#include "db.h"

typedef uint32_t BackupID;

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(BackupableDBOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(BackupableDBOptions)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS_DEC(BackupableDBOptions, String)
DEFINE_C_WRAP_DESTRUCTOR_DEC(BackupableDBOptions)

DEFINE_C_WRAP_STRUCT(RestoreOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(RestoreOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(RestoreOptions)
DEFINE_C_WRAP_DESTRUCTOR_DEC(RestoreOptions)
// Get/Set methods
DEFINE_C_WRAP_SETTER_DEC(RestoreOptions, keep_log_files, bool)

// Please see the documentation in BackupableDB and RestoreBackupableDB
DEFINE_C_WRAP_STRUCT(BackupEngine)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(BackupEngine)
DEFINE_C_WRAP_DESTRUCTOR_DEC(BackupEngine)

// Create a new backup engine from db_env and options.
Status_t BackupEngineOpen(Env_t* db_env, const BackupableDBOptions_t* options, BackupEngine_t* backup_engine_ptr);

// Create a new backup for db.
Status_t BackupEngineCreateNewBackup(BackupEngine_t* backup_engine_ptr, DB_t* db, bool flush_before_backup);

// Restoring DB from backup is NOT safe when there is another BackupEngine
// running that might call DeleteBackup() or PurgeOldBackups(). It is caller's
// responsibility to synchronize the operation, i.e. don't delete the backup
// when you're restoring from it
Status_t BackupEngineRestoreDBFromBackup(BackupEngine_t* backup_engine_ptr, BackupID backup_id, String_t* db_dir, String_t* wal_dir, RestoreOptions_t* restore_options);
Status_t BackupEngineRestoreDBFromLatestBackup(BackupEngine_t* backup_engine_ptr, String_t* db_dir, String_t* wal_dir, RestoreOptions_t* restore_options);


// BackupableDBOptions have to be the same as the ones used in a previous
// incarnation of the DB
DEFINE_C_WRAP_STRUCT(BackupableDB)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(BackupableDB)
DEFINE_C_WRAP_DESTRUCTOR_DEC(BackupableDB)

// Use this class to access information about backups and restore from them
DEFINE_C_WRAP_STRUCT(RestoreBackupableDB)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(RestoreBackupableDB)
DEFINE_C_WRAP_DESTRUCTOR_DEC(RestoreBackupableDB)

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_BACKUPABLE_DB_H_
