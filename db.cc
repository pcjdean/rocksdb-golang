// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#include <assert.h>
#include <string.h>
#include <stdio.h>
#include <rocksdb/db.h>
#include "db.h"

using namespace rocksdb;

DEFINE_C_WRAP_CONSTRUCTOR(TablePropertiesCollection)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(TablePropertiesCollection)
DEFINE_C_WRAP_DESTRUCTOR(TablePropertiesCollection)

// A range of keys
DEFINE_C_WRAP_CONSTRUCTOR(Range)
DEFINE_C_WRAP_CONSTRUCTOR_ARGS(Range, Slice, Slice)
DEFINE_C_WRAP_DESTRUCTOR(Range)

DEFINE_C_WRAP_CONSTRUCTOR_ARGS(ColumnFamilyDescriptor, String, ColumnFamilyOptions)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(ColumnFamilyDescriptor, String, ColumnFamilyOptions)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(ColumnFamilyDescriptor, kDefaultColumnFamilyName, ColumnFamilyOptions())
DEFINE_C_WRAP_DESTRUCTOR(ColumnFamilyDescriptor)
// A DB is a persistent ordered map from keys to values.
// A DB is safe for concurrent access from multiple threads without
// any external synchronization.
DEFINE_C_WRAP_DESTRUCTOR(DB)

// Open the database with the specified "name".
// Stores a pointer to a heap-allocated database in *dbptr and returns
// OK on success.
// Stores nullptr in *dbptr and returns a non-OK status on error.
// Caller should delete *dbptr when it is no longer needed.
Status_t DBOpen(const Options_t* options,
                const String_t* name,
                DB_t* dbptr)
{
    assert(dbptr != NULL);
    DB** rdbptr = GET_REP_ADDR(dbptr, DB);
    assert(rdbptr != NULL);
    assert(GET_REP(options, Options) != NULL);
    assert(GET_REP(name, String) != NULL);
    Status stat = DB::Open(GET_REP_REF(options, Options), GET_REP_REF(name, String), rdbptr);
    return NewStatusTCopy(&stat);
}


// Open the database for read only. All DB interfaces
// that modify data, like put/delete, will return error.
// If the db is opened in read only mode, then no compactions
// will happen.
//
// Not supported in ROCKSDB_LITE, in which case the function will
// return Status_t::NotSupported.
Status_t DBOpenForReadOnly(const Options_t* options,
                           const String_t* name, DB_t* dbptr,
                           bool error_if_log_file_exist)
{
    assert(dbptr != NULL);
    DB** rdbptr = GET_REP_ADDR(dbptr, DB);
    assert(rdbptr != NULL);
    assert(GET_REP(options, Options) != NULL);
    assert(GET_REP(name, String) != NULL);
    Status stat = DB::OpenForReadOnly(GET_REP_REF(options, Options), GET_REP_REF(name, String),  rdbptr, error_if_log_file_exist);
    return NewStatusTCopy(&stat);
}

// Open the database for read only with column families. When opening DB with
// read only, you can specify only a subset of column families in the
// database that should be opened. However, you always need to specify default
// column family. The default column family name is 'default' and it's stored
// in rocksdb::kDefaultColumnFamilyName
//
// Not supported in ROCKSDB_LITE, in which case the function will
// return Status_t::NotSupported.
Status_t DBOpenForReadOnlyWithColumnFamilies(const Options_t* options,
                                             const String_t* name,
                                             const ColumnFamilyDescriptor_t column_families[],
                                             const int size_col,
                                             ColumnFamilyHandle_t **handles,
                                             DB_t* dbptr, bool error_if_log_file_exist)
{
    std::vector<ColumnFamilyDescriptor> column_families_vec = std::vector<ColumnFamilyDescriptor>(size_col);
    for (int i = 0; i < size_col; i++)
        column_families_vec.push_back(*(ColumnFamilyDescriptor*)column_families[i].rep);
    std::vector<ColumnFamilyHandle*> handles_vec;
    assert(dbptr != NULL);
    DB** rdbptr = GET_REP_ADDR(dbptr, DB);
    assert(rdbptr != NULL);
    assert(GET_REP(options, Options) != NULL);
    assert(GET_REP(name, String) != NULL);
    Status stat = DB::OpenForReadOnly(GET_REP_REF(options, Options), GET_REP_REF(name, String), column_families_vec, &handles_vec, rdbptr, error_if_log_file_exist);
    Status_t ret = NewStatusTCopy(&stat);
    assert(handles_vec.size() == size_col);
    *handles = new ColumnFamilyHandle_t[size_col];
    for (int j = 0; j < size_col; j++)
    {
        ColumnFamilyHandle** rcfh = GET_REP_ADDR(&(*handles)[j], ColumnFamilyHandle);
        *rcfh = handles_vec[j];
    }
    return ret;
}

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
Status_t DBOpenWithColumnFamilies(const Options_t* options, const String_t* name,
                                  const ColumnFamilyDescriptor_t column_families[], const int size_col,
                                  ColumnFamilyHandle_t **handles, DB_t* dbptr)
{
    std::vector<ColumnFamilyDescriptor> column_families_vec = std::vector<ColumnFamilyDescriptor>(size_col);
    for (int i = 0; i < size_col; i++)
        column_families_vec.push_back(*(ColumnFamilyDescriptor*)column_families[i].rep);
    std::vector<ColumnFamilyHandle*> handles_vec;
    assert(dbptr != NULL);
    DB** rdbptr = GET_REP_ADDR(dbptr, DB);
    assert(rdbptr != NULL);
    assert(GET_REP(options, Options) != NULL);
    assert(GET_REP(name, String) != NULL);
    Status stat = DB::Open(GET_REP_REF(options, Options), GET_REP_REF(name, String), column_families_vec, &handles_vec, rdbptr);
    Status_t ret = NewStatusTCopy(&stat);
    assert(handles_vec.size() == size_col);
    *handles = new ColumnFamilyHandle_t[size_col];
    for (int j = 0; j < size_col; j++)
    {
        ColumnFamilyHandle** rcfh = GET_REP_ADDR(&(*handles)[j], ColumnFamilyHandle);
        *rcfh = handles_vec[j];
    }
    return ret;
}

// ListColumnFamilies will open the DB specified by argument name
// and return the list of all column families in that DB
// through column_families argument. The ordering of
// column families in column_families is unspecified.
Status_t DBListColumnFamilies(DBOptions_t* db_options,
                              const String_t* name,
                              String_t **column_families, int* size_col)
{
    std::vector<std::string> column_families_vec;
    assert(GET_REP(db_options, Options) != NULL);
    assert(GET_REP(name, String) != NULL);
    Status stat = DB::ListColumnFamilies(GET_REP_REF(db_options, Options), GET_REP_REF(name, String), &column_families_vec);
    Status_t ret = NewStatusTCopy(&stat);
    *size_col = column_families_vec.size();
    *column_families = new String_t[*size_col];
    for (int j = 0; j < *size_col; j++)
    {
        (*column_families)[j].rep = new String();
        String* rstr = GET_REP(&(*column_families)[j], String);
        assert(rstr != NULL);
        *rstr = std::move(column_families_vec[j]);
    }
    return ret;
}


// Create a column_family and return the handle of column family
// through the argument handle.
Status_t DBCreateColumnFamily(const DB_t* dbptr, const ColumnFamilyOptions_t* options,
                              const String_t* column_family_name,
                              ColumnFamilyHandle_t* handle)
{
    assert(handle != NULL);
    ColumnFamilyHandle** rcfh = GET_REP_ADDR(handle, ColumnFamilyHandle);
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    assert(GET_REP(options, ColumnFamilyOptions) != NULL);
    assert(GET_REP(column_family_name, String) != NULL);
    Status stat = GET_REP(dbptr, DB)->CreateColumnFamily(GET_REP_REF(options, ColumnFamilyOptions), GET_REP_REF(column_family_name, String), rcfh);
    return NewStatusTCopy(dbptr ?
                          &stat :
                          // TODO create a const constructor
                          const_cast<Status*>(&invalid_status));
}

// Drop a column family specified by column_family handle. This call
// only records a drop record in the manifest and prevents the column
// family from flushing and compacting.
Status_t DBDropColumnFamily(const DB_t* dbptr, const ColumnFamilyHandle_t* column_family)
{
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
    Status stat = GET_REP(dbptr, DB)->DropColumnFamily(GET_REP(column_family, ColumnFamilyHandle));
    return NewStatusTCopy(dbptr ?
                          &stat :
                          const_cast<Status*>(&invalid_status));
}

// Set the database entry for "key" to "value".
// If "key" already exists, it will be overwritten.
// Returns OK on success, and a non-OK status on error.
// Note: consider setting options.sync = true.
Status_t DBPutWithColumnFamily(const DB_t* dbptr, const WriteOptions_t* options,
                               const ColumnFamilyHandle_t* column_family,
                               const Slice_t* key,
                               const Slice_t* value)
{
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    assert(GET_REP(options, WriteOptions) != NULL);
    assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
    assert(GET_REP(key, Slice) != NULL);
    assert(GET_REP(value, Slice) != NULL);
    Status stat = GET_REP(dbptr, DB)->Put(GET_REP_REF(options, WriteOptions), GET_REP(column_family, ColumnFamilyHandle), GET_REP_REF(key, Slice), GET_REP_REF(value, Slice));
    return NewStatusTCopy(dbptr ?
                          &stat :
                          const_cast<Status*>(&invalid_status));
}

Status_t DBPut(const DB_t* dbptr, const WriteOptions_t* options,
               const Slice_t* key,
               const Slice_t* value)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBPutWithColumnFamily(dbptr, options, &column_family, key, value);
}

// Remove the database entry (if any) for "key".  Returns OK on
// success, and a non-OK status on error.  It is not an error if "key"
// did not exist in the database.
// Note: consider setting options.sync = true.
Status_t DBDeleteWithColumnFamily(const DB_t* dbptr, const WriteOptions_t* options,
                                  const ColumnFamilyHandle_t* column_family,
                                  const Slice_t* key)
{
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    assert(GET_REP(options, WriteOptions) != NULL);
    assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
    assert(GET_REP(key, Slice) != NULL);
    Status stat = GET_REP(dbptr, DB)->Delete(GET_REP_REF(options, WriteOptions), GET_REP(column_family, ColumnFamilyHandle), GET_REP_REF(key, Slice));
    return NewStatusTCopy(dbptr ?
                          &stat :
                          const_cast<Status*>(&invalid_status));
}

Status_t DBDelete(const DB_t* dbptr, const WriteOptions_t* options,
                  const Slice_t* key)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBDeleteWithColumnFamily(dbptr, options, &column_family, key);
}

// Merge the database entry for "key" with "value".  Returns OK on success,
// and a non-OK status on error. The semantics of this operation is
// determined by the user provided merge_operator when opening DB.
// Note: consider setting options.sync = true.
Status_t DBMergeWithColumnFamily(const DB_t* dbptr, const WriteOptions_t* options,
                                 const ColumnFamilyHandle_t* column_family,
                                 const Slice_t* key,
                                 const Slice_t* value)
{
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    assert(GET_REP(options, WriteOptions) != NULL);
    assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
    assert(GET_REP(key, Slice) != NULL);
    assert(GET_REP(value, Slice) != NULL);
    Status stat = GET_REP(dbptr, DB)->Merge(GET_REP_REF(options, WriteOptions), GET_REP(column_family, ColumnFamilyHandle), GET_REP_REF(key, Slice), GET_REP_REF(value, Slice));
    return NewStatusTCopy(dbptr ?
                          &stat :
                          const_cast<Status*>(&invalid_status));
}

Status_t DBMerge(const DB_t* dbptr, const WriteOptions_t* options,
                 const Slice_t* key,
                 const Slice_t* value)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBMergeWithColumnFamily(dbptr, options, &column_family, key, value);
}

// Apply the specified updates to the database.
// If `updates` contains no update, WAL will still be synced if
// options.sync=true.
// Returns OK on success, non-OK on failure.
// Note: consider setting options.sync = true.
Status_t DBWrite(const DB_t* dbptr, const WriteOptions_t* options, WriteBatch_t* updates)
{
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    assert(GET_REP(options, WriteOptions) != NULL);
    assert(GET_REP(updates, WriteBatch) != NULL);
    Status stat = GET_REP(dbptr, DB)->Write(GET_REP_REF(options, WriteOptions), GET_REP(updates, WriteBatch));
    return NewStatusTCopy(dbptr ?
                          &stat :
                          const_cast<Status*>(&invalid_status));
}

// If the database contains an entry for "key" store the
// corresponding value in *value and return OK.
//
// If there is no entry for "key" leave *value unchanged and return
// a status for which Status_t::IsNotFound() returns true.
//
// May return some other Status_t on an error.
Status_t DBGetWithColumnFamily(const DB_t* dbptr, const ReadOptions_t* options,
                               const ColumnFamilyHandle_t* column_family,
                               const Slice_t* key,
                               const String_t* value)
{
    Status ret;
    if (dbptr)
    {
        std::string str_val;
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(options, ReadOptions) != NULL);
        assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
        assert(GET_REP(key, Slice) != NULL);
        ret = GET_REP(dbptr, DB)->Get(GET_REP_REF(options, ReadOptions), GET_REP(column_family, ColumnFamilyHandle), GET_REP_REF(key, Slice), &str_val);
        if (!str_val.empty())
            GET_REP_REF(value, String) = std::move(str_val);
        
    }
    else
        ret = invalid_status;
    return NewStatusTCopy(&ret);
}

Status_t DBGet(const DB_t* dbptr, const ReadOptions_t* options,
             const Slice_t* key,
             const String_t* value)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBGetWithColumnFamily(dbptr, options, &column_family, key, value);
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
Status_t* DBMultiGetWithColumnFamily(const DB_t* dbptr, const ReadOptions_t* options,
                                     const ColumnFamilyHandle_t column_families[],
                                     const int size_col,
                                     const Slice_t keys[],
                                     const int size_keys,
                                     String_t** values)
{
    Status_t* ret;
    if (dbptr)
    {
        std::vector<ColumnFamilyHandle*> column_families_vec = std::vector<ColumnFamilyHandle*>(size_col);
        for (int i = 0; i < size_col; i++)
            column_families_vec.push_back((ColumnFamilyHandle*)column_families[i].rep);
        std::vector<Slice> keys_vec = std::vector<Slice>(size_col);
        for (int i = 0; i < size_keys; i++)
            keys_vec.push_back(*(Slice*)(keys[i].rep));
        std::vector<std::string> values_vec;
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(options, ReadOptions) != NULL);
        std::vector<Status> ret_vec = GET_REP(dbptr, DB)->MultiGet(GET_REP_REF(options, ReadOptions), column_families_vec, keys_vec, &values_vec);
        assert(values_vec.size() == size_keys);
        assert(ret_vec.size() == size_keys);
        *values = new String_t[size_keys];
        ret = new Status_t[size_keys];
        for (int j = 0; j < size_keys; j++)
        {
            (*values)[j].rep = new String();
            GET_REP_REF(&(*values)[j], String) = std::move(values_vec[j]);
            ret[j] = NewStatusTCopy(&ret_vec[j]);
        }
    }
    else
    {
        ret = new Status_t();
        ret->rep = new Status(invalid_status);
    }
    return ret;
}

Status_t* DBMultiGet(const DB_t* dbptr, const ReadOptions_t* options,
                     const Slice_t keys[],
                     const int size_keys,
                     String_t** values)
{
    ColumnFamilyHandle_t *column_families = new ColumnFamilyHandle_t[size_keys];
    std::fill_n(column_families, size_keys, DBDefaultColumnFamily(dbptr));
    Status_t* ret = DBMultiGetWithColumnFamily(dbptr, options, column_families, 1, keys, size_keys, values);
    delete []column_families;
    return ret;
}

// If the key definitely does not exist in the database, then this method
// returns false, else true. If the caller wants to obtain value when the key
// is found in memory, a bool for 'value_found' must be passed. 'value_found'
// will be true on return if value has been set properly.
// This check is potentially lighter-weight than invoking DB::Get(). One way
// to make this lighter weight is to avoid doing any IOs.
// Default implementation here returns true and sets 'value_found' to false
bool DBKeyMayExistWithColumnFamily(const DB_t* dbptr, const ReadOptions_t* options,
                                   const ColumnFamilyHandle_t* column_family,
                                   const Slice_t* key,
                                   String_t* value,
                                   bool* value_found)
{
    std::string val_str;
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    assert(GET_REP(options, ReadOptions) != NULL);
    assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
    assert(GET_REP(key, Slice) != NULL);
    bool ret = GET_REP(dbptr, DB)->KeyMayExist(GET_REP_REF(options, ReadOptions), GET_REP(column_family, ColumnFamilyHandle), GET_REP_REF(key, Slice), &val_str, value_found);
    GET_REP_REF(value, String) = std::move(val_str);
    return ret;
}

bool DBKeyMayExist(const DB_t* dbptr, const ReadOptions_t* options,
                   const Slice_t* key,
                   String_t* value,
                   bool* value_found)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBKeyMayExistWithColumnFamily(dbptr, options, &column_family, key, value, value_found);
}

// Return a heap-allocated iterator over the contents of the database.
// The result of NewIterator() is initially invalid (caller must
// call one of the Seek methods on the iterator before using it).
//
// Caller should delete the iterator when it is no longer needed.
// The returned iterator should be deleted before this db is deleted.
Iterator_t DBNewIteratorWithColumnFamily(const DB_t* dbptr, const ReadOptions_t* options,
                                         const ColumnFamilyHandle_t* column_family)
{
    assert(GET_REP(options, ReadOptions) != NULL);
    assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
    return NewIteratorT((dbptr && GET_REP(dbptr, DB)) ? GET_REP(dbptr, DB)->NewIterator(GET_REP_REF(options, ReadOptions), GET_REP(column_family, ColumnFamilyHandle)) : nullptr);
}
    
Iterator_t DBNewIterator(const DB_t* dbptr, const ReadOptions_t* options)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBNewIteratorWithColumnFamily(dbptr, options, &column_family);
}

// Returns iterators from a consistent database state across multiple
// column families. Iterators are heap allocated and need to be deleted
// before the db is deleted
Status_t DBNewIterators(const DB_t* dbptr, const ReadOptions_t* options,
                        const ColumnFamilyHandle_t column_families[],
                        const int size_col,
                        Iterator_t** values,
                        int *val_sz)
{
    Status ret;
    if (dbptr)
    {
        std::vector<ColumnFamilyHandle*> column_families_vec = std::vector<ColumnFamilyHandle*>(size_col);
        for (int i = 0; i < size_col; i++)
            column_families_vec.push_back(GET_REP(&column_families[i], ColumnFamilyHandle));
        std::vector<Iterator*> values_vec;
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(options, ReadOptions) != NULL);
        ret = GET_REP(dbptr, DB)->NewIterators(GET_REP_REF(options, ReadOptions), column_families_vec, &values_vec);
        *val_sz = values_vec.size();
        *values = new Iterator_t[*val_sz];
        for (int j = 0; j < *val_sz; j++)
        {
            (*values)[j].rep = values_vec[j];
        }
    }
    else
    {
        *val_sz = 0;
        ret = invalid_status;
    }
    return NewStatusTCopy(&ret);
}

// Return a handle to the current DB state.  Iterators created with
// this handle will all observe a stable snapshot of the current DB
// state.  The caller must call ReleaseSnapshot(result) when the
// snapshot is no longer needed.
//
// nullptr will be returned if the DB fails to take a snapshot or does
// not support snapshot.
Snapshot_t DBGetSnapshot(const DB_t* dbptr)
{
    return NewSnapshotT((dbptr && GET_REP(dbptr, DB)) ? const_cast<Snapshot *>(GET_REP(dbptr, DB)->GetSnapshot()) : nullptr);
}

// Release a previously acquired snapshot.  The caller must not
// use "snapshot" after this call.
void DBReleaseSnapshot(const DB_t* dbptr, const Snapshot_t* snapshot)
{
    if (dbptr)
    {
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(snapshot, Snapshot) != NULL);
        GET_REP(dbptr, DB)->ReleaseSnapshot(GET_REP(snapshot, Snapshot));
    }
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

bool DBGetPropertyWithColumnFamily(const DB_t* dbptr, const ColumnFamilyHandle_t* column_family,
                                   const Slice_t* property, String_t* value)
{
    if (dbptr)
    {
        std::string str_val;
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
        assert(GET_REP(property, Slice) != NULL);
        bool ret = GET_REP(dbptr, DB)->GetProperty( GET_REP(column_family, ColumnFamilyHandle), GET_REP_REF(property, Slice), &str_val);
        assert(GET_REP(value, String) != NULL);
        GET_REP_REF(value, String) = std::move(str_val);
        return ret;
    }
    else
        return false;
}
    
bool DBGetProperty(const DB_t* dbptr, const Slice_t* property, String_t* value)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBGetPropertyWithColumnFamily(dbptr, &column_family, property, value);
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
bool DBGetIntPropertyWithColumnFamily(const DB_t* dbptr, 
                                      const ColumnFamilyHandle_t* column_family,
                                      const Slice_t* property, uint64_t* value)
{
    if (dbptr)
    {
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
        assert(GET_REP(property, Slice) != NULL);
        return GET_REP(dbptr, DB)->GetIntProperty(GET_REP(column_family, ColumnFamilyHandle), GET_REP_REF(property, Slice), value);
    }
    else
        return false;
}

bool DBGetIntProperty(const DB_t* dbptr, 
                      const Slice_t* property, uint64_t* value)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBGetIntPropertyWithColumnFamily(dbptr, &column_family, property, value);
}

// For each i in [0,n-1], store in "sizes[i]", the approximate
// file system space used by keys in "[range[i].start .. range[i].limit)".
//
// Note that the returned sizes measure file system space usage, so
// if the user data compresses by a factor of ten, the returned
// sizes will be one-tenth the size of the corresponding user data size.
//
// The results may not include the sizes of recently written data.
void DBGetApproximateSizesWithColumnFamily(const DB_t* dbptr, 
                                           const ColumnFamilyHandle_t* column_family,
                                           const Range_t* range, int n,
                                           uint64_t* sizes)
{
    if (dbptr)
    {
        Range* range_ary = new Range[n];
        for (int i = 0; i < n; i++)
            range_ary[i] = GET_REP_REF(&range[i], Range);
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
        GET_REP(dbptr, DB)->GetApproximateSizes(GET_REP(column_family, ColumnFamilyHandle), range_ary, n, sizes);
        delete []range_ary;
    }
}

void DBGetApproximateSizes(const DB_t* dbptr, 
                           const Range_t* range, int n,
                           uint64_t* sizes)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    DBGetApproximateSizesWithColumnFamily(dbptr, &column_family, range, n, sizes);
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
Status_t DBCompactRangeWithColumnFamily(const DB_t* dbptr, 
                                        const ColumnFamilyHandle_t* column_family,
                                        const Slice_t* begin, const Slice_t* end,
                                        bool reduce_level, int target_level,
                                        uint32_t target_path_id)
{
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
    assert(begin != NULL);
    assert(end != NULL);
    Status stat = GET_REP(dbptr, DB)->CompactRange(GET_REP(column_family, ColumnFamilyHandle), GET_REP(begin, Slice), GET_REP(end, Slice), reduce_level, target_level, target_path_id);
    return NewStatusTCopy(dbptr ?
                          &stat :
                          const_cast<Status*>(&invalid_status));
}

Status_t DBCompactRange(const DB_t* dbptr, 
                        const Slice_t* begin, const Slice_t* end,
                        bool reduce_level, int target_level,
                        uint32_t target_path_id)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBCompactRangeWithColumnFamily(dbptr, &column_family, begin, end, reduce_level, target_level, target_path_id);
}

Status_t DBSetOptionsWithColumnFamily(const DB_t* dbptr, 
                                      const ColumnFamilyHandle_t* column_family,
                                      const String_t new_options[],
                                      int n)
{
    std::unordered_map<std::string, std::string> new_options_map;
    for (int i = 0; i < n; i++)
    {
        new_options_map[std::move(GET_REP_REF(&new_options[i], String))] = std::move(GET_REP_REF(&new_options[++i], String));
    }
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
    Status stat = GET_REP(dbptr, DB)->SetOptions(GET_REP(column_family, ColumnFamilyHandle), new_options_map);
    return NewStatusTCopy(&stat);
}

Status_t DBSetOptions(const DB_t* dbptr, 
                      const String_t new_options[],
                      const int n)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBSetOptionsWithColumnFamily(dbptr, &column_family, new_options, n);
}

// CompactFiles() inputs a list of files specified by file numbers
// and compacts them to the specified level.  Note that the behavior
// is different from CompactRange in that CompactFiles() will
// perform the compaction job using the CURRENT thread.
//
// @see GetDataBaseMetaData
// @see GetColumnFamilyMetaData
Status_t DBCompactFilesWithColumnFamily(const DB_t* dbptr, 
                                        const CompactionOptions_t* compact_options,
                                        const ColumnFamilyHandle_t* column_family,
                                        const String_t input_file_names[],
                                        const int n,
                                        const int output_level, const int output_path_id)
{
    Status ret;
    if (dbptr)
    {
        std::vector<std::string> input_file_names_vec;
        for (int i = 0; i < n; i++)
            input_file_names_vec.push_back(std::move(GET_REP_REF(&input_file_names[i], String)));
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(compact_options, CompactionOptions) != NULL);
        assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
        ret = GET_REP(dbptr, DB)->CompactFiles(GET_REP_REF(compact_options, CompactionOptions), GET_REP(column_family, ColumnFamilyHandle), input_file_names_vec, output_level, output_path_id);
    }
    else
    {
        ret = invalid_status;
    }
    return NewStatusTCopy(&ret);
}

Status_t DBCompactFiles(const DB_t* dbptr, 
                        const CompactionOptions_t* compact_options,
                        const String_t input_file_names[],
                        const int n,
                        const int output_level, const int output_path_id)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBCompactFilesWithColumnFamily(dbptr, compact_options, &column_family,
                                          input_file_names, n, output_level, output_path_id);
}

// Number of levels used for this DB.
int DBNumberLevelsWithColumnFamily(const DB_t* dbptr, 
                                   const ColumnFamilyHandle_t* column_family)
{
    int ret = 0;
    if (dbptr)
    {
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
        ret = GET_REP(dbptr, DB)->NumberLevels(GET_REP(column_family, ColumnFamilyHandle));
    }
    return ret;
}

int DBNumberLevels(const DB_t* dbptr)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBNumberLevelsWithColumnFamily(dbptr, &column_family);
}

// Maximum level to which a new compacted memtable is pushed if it
// does not create overlap.
int DBMaxMemCompactionLevelWithColumnFamily(const DB_t* dbptr, 
                                            const ColumnFamilyHandle_t* column_family)
{
    int ret = 0;
    if (dbptr)
    {
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
        ret = GET_REP(dbptr, DB)->MaxMemCompactionLevel(GET_REP(column_family, ColumnFamilyHandle));
    }
    return ret;
}

int DBMaxMemCompactionLevel(const DB_t* dbptr)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBMaxMemCompactionLevelWithColumnFamily(dbptr, &column_family);
}

// Number of files in level-0 that would stop writes.
int DBLevel0StopWriteTriggerWithColumnFamily(const DB_t* dbptr, 
                                             const ColumnFamilyHandle_t* column_family)
{
    int ret = 0;
    if (dbptr)
    {
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
        ret = GET_REP(dbptr, DB)->Level0StopWriteTrigger(GET_REP(column_family, ColumnFamilyHandle));
    }
    return ret;
}

int DBLevel0StopWriteTrigger(const DB_t* dbptr)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBLevel0StopWriteTriggerWithColumnFamily(dbptr, &column_family);
}

// Get DB name -- the exact same name that was provided as an argument to
// DB::Open()
String_t DBGetName(const DB_t* dbptr)
{
    if (dbptr)
    {
        assert(GET_REP(dbptr, DB) != NULL);
        return NewStringTMove(const_cast<String *>(&GET_REP(dbptr, DB)->GetName()));
    }
    else
        return NewStringT(nullptr);
}

// Get Env object from the DB
Env_t DBGetEnv(const DB_t* dbptr)
{
    return NewEnvT((dbptr && GET_REP(dbptr, DB)) ? GET_REP(dbptr, DB)->GetEnv() : nullptr);
}

// Get DB Options that we use.  During the process of opening the
// column family, the options provided when calling DB::Open() or
// DB::CreateColumnFamily() will have been "sanitized" and transformed
// in an implementation-defined manner.
Options_t DBGetOptionsWithColumnFamily(const DB_t* dbptr, 
                                       const ColumnFamilyHandle_t* column_family)
{
    if (dbptr)
    {
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
        const Options &options = GET_REP(dbptr, DB)->GetOptions(GET_REP(column_family, ColumnFamilyHandle));
        return NewOptionsTRawArgs(options, options);
    }
    else
        return NewOptionsT(nullptr);
}

Options_t DBGetOptions(const DB_t* dbptr)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBGetOptionsWithColumnFamily(dbptr, &column_family);
}

DBOptions_t DBGetDBOptions(const DB_t* dbptr)
{
    if (dbptr)
    {
        assert(GET_REP(dbptr, DB) != NULL);
        const DBOptions &dboptions = GET_REP(dbptr, DB)->GetDBOptions();
        const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
        assert(GET_REP(&column_family, ColumnFamilyHandle) != NULL);
        const Options &options = GET_REP(dbptr, DB)->GetOptions(GET_REP(&column_family, ColumnFamilyHandle));
        Options new_options(dboptions, options);
        return NewDBOptionsTRawArgs(new_options);
    }
    else
        return NewDBOptionsT(nullptr);
}

// Flush all mem-table data.
Status_t DBFlushWithColumnFamily(const DB_t* dbptr, 
                                 const FlushOptions_t* options,
                                 const ColumnFamilyHandle_t* column_family)
{
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    assert(GET_REP(options, FlushOptions) != NULL);
    assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
    Status stat = GET_REP(dbptr, DB)->Flush(GET_REP_REF(options, FlushOptions), GET_REP(column_family, ColumnFamilyHandle));
    return NewStatusTCopy(dbptr ?
                          &stat :
                          const_cast<Status*>(&invalid_status));
}

Status_t DBFlush(const DB_t* dbptr, 
                 const FlushOptions_t* options)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBFlushWithColumnFamily(dbptr, options, &column_family);
}

// The sequence number of the most recent transaction.
SequenceNumber DBGetLatestSequenceNumber(const DB_t* dbptr)
{
    SequenceNumber ret = -1;
    if (dbptr)
    {
        assert(GET_REP(dbptr, DB) != NULL);
        ret = GET_REP(dbptr, DB)->GetLatestSequenceNumber();
    }
    return ret;
}

#ifndef ROCKSDB_LITE

// Prevent file deletions. Compactions will continue to occur,
// but no obsolete files will be deleted. Calling this multiple
// times have the same effect as calling it once.
Status_t DBDisableFileDeletions(const DB_t* dbptr)
{
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    Status stat = GET_REP(dbptr, DB)->DisableFileDeletions();
    return NewStatusTCopy(dbptr ?
                          &stat :
                          const_cast<Status*>(&invalid_status));
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
Status_t DBEnableFileDeletions(const DB_t* dbptr, bool force)
{
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    Status stat = GET_REP(dbptr, DB)->EnableFileDeletions(force);
    return NewStatusTCopy(dbptr ?
                          &stat :
                          const_cast<Status*>(&invalid_status));
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
Status_t DBGetLiveFiles(const DB_t* dbptr,
                        String_t **live_files,
                        int* n,
                        uint64_t* manifest_file_size,
                        bool flush_memtable)
{
    Status ret;
    if (dbptr)
    {
        std::vector<std::string> live_files_vec;
        assert(GET_REP(dbptr, DB) != NULL);
        ret = GET_REP(dbptr, DB)->GetLiveFiles(live_files_vec, manifest_file_size, flush_memtable);
        *n = live_files_vec.size();
        *live_files = new String_t[*n];
        for (int j = 0; j < *n; j++)
        {
            (*live_files)[j].rep = new String();
            GET_REP_REF(&(*live_files)[j], String) = std::move(live_files_vec[j]);
        }
    }
    else
    {
        ret = invalid_status;
    }
    return NewStatusTCopy(&ret);
}

// Retrieve the sorted list of all wal files with earliest file first
Status_t DBGetSortedWalFiles(const DB_t* dbptr, LogFile_t **files, int* n)
{
    Status ret;
    if (dbptr)
    {
        VectorLogPtr files_vec;
        assert(GET_REP(dbptr, DB) != NULL);
        ret = GET_REP(dbptr, DB)->GetSortedWalFiles(files_vec);
        *n = files_vec.size();
        *files = new LogFile_t[*n];
        for (int j = 0; j < *n; j++)
            (*files)[j].rep = files_vec[j].release();
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
Status_t DBGetUpdatesSince(const DB_t* dbptr, SequenceNumber seq_number,
                           TransactionLogIterator_t* iter,
                           const TransactionLogIterator_ReadOptions_t* read_options)
{
    Status ret;
    if (dbptr)
    {
        unique_ptr<TransactionLogIterator> iter_ptr;
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(read_options, TransactionLogIterator_ReadOptions) != NULL);
        ret = GET_REP(dbptr, DB)->GetUpdatesSince(seq_number, &iter_ptr, GET_REP_REF(read_options, TransactionLogIterator_ReadOptions));
        iter->rep = iter_ptr.release();
    }
    {
        ret = invalid_status;
    }
    return NewStatusTCopy(&ret);
}

// Delete the file name from the db directory and update the internal state to
// reflect that. Supports deletion of sst and log files only. 'name' must be
// path relative to the db directory. eg. 000001.sst, /archive/000003.log
Status_t DBDeleteFile(const DB_t* dbptr, String_t* name)
{
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    assert(GET_REP(name, String) != NULL);
    Status stat = GET_REP(dbptr, DB)->DeleteFile(GET_REP_REF(name, String));
    return NewStatusTCopy(dbptr ?
                          &stat :
                          const_cast<Status*>(&invalid_status));
}

// Returns a list of all table files with their level, start key
// and end key
void DBGetLiveFilesMetaData(const DB_t* dbptr, LiveFileMetaData_t** metadata, int* n)
{
    if (dbptr)
    {
        std::vector<LiveFileMetaData> metadata_vec;
        assert(GET_REP(dbptr, DB) != NULL);
        GET_REP(dbptr, DB)->GetLiveFilesMetaData(&metadata_vec);
        *n = metadata_vec.size();
        *metadata = new LiveFileMetaData_t[*n];
        for (int j = 0; j < *n; j++)
        {
            (*metadata)[j].rep = new LiveFileMetaData();
            GET_REP_REF(&(*metadata)[j], LiveFileMetaData) = metadata_vec[j];
        }
    }
}

// Obtains the meta data of the specified column family of the DB.
// Status_t::NotFound() will be returned if the current DB does not have
// any column family match the specified name.
//
// If cf_name is not specified, then the metadata of the default
// column family will be returned.
void DBGetColumnFamilyMetaDataWithColumnFamily(const DB_t* dbptr, 
                                               const ColumnFamilyHandle_t* column_family,
                                               ColumnFamilyMetaData_t* metadata)
{
    if (dbptr)
    {
        assert(GET_REP(dbptr, DB) != NULL);
        assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
        assert(GET_REP(metadata, ColumnFamilyMetaData) != NULL);
        GET_REP(dbptr, DB)->GetColumnFamilyMetaData(GET_REP(column_family, ColumnFamilyHandle), GET_REP(metadata, ColumnFamilyMetaData));
    }
}

// Get the metadata of the default column family.
void DBGetColumnFamilyMetaData(const DB_t* dbptr, 
                               ColumnFamilyMetaData_t* metadata)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    DBGetColumnFamilyMetaDataWithColumnFamily(dbptr, &column_family, metadata);
}
#endif  // ROCKSDB_LITE

// Sets the globally unique ID created at database creation time by invoking
// Env::GenerateUniqueId(), in identity. Returns Status_t::OK if identity could
// be set properly
Status_t DBGetDbIdentity(const DB_t* dbptr, String_t* identity)
{
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    assert(GET_REP(identity, String) != NULL);
    Status stat = GET_REP(dbptr, DB)->GetDbIdentity(GET_REP_REF(identity, String));
    return NewStatusTCopy(dbptr ?
                          &stat :
                          const_cast<Status*>(&invalid_status));
}

// Returns default column family handle
ColumnFamilyHandle_t DBDefaultColumnFamily(const DB_t* dbptr)
{
    ColumnFamilyHandle_t cf_handle;
    cf_handle.rep = (dbptr && GET_REP(dbptr, DB)) ?
        GET_REP(dbptr, DB)->DefaultColumnFamily() :
        nullptr; 
    return cf_handle;
}

#ifndef ROCKSDB_LITE
Status_t DBGetPropertiesOfAllTablesWithColumnFamily(const DB_t* dbptr, 
                                                    const ColumnFamilyHandle_t* column_family,
                                                    TablePropertiesCollection_t* props)
{
    assert(dbptr != NULL);
    assert(GET_REP(dbptr, DB) != NULL);
    assert(GET_REP(column_family, ColumnFamilyHandle) != NULL);
    assert(GET_REP(props, TablePropertiesCollection) != NULL);
    Status stat = GET_REP(dbptr, DB)->GetPropertiesOfAllTables(GET_REP(column_family, ColumnFamilyHandle), GET_REP(props, TablePropertiesCollection));
    return NewStatusTCopy(dbptr ?
                          &stat :
                          const_cast<Status*>(&invalid_status));
}

Status_t DBGetPropertiesOfAllTables(const DB_t* dbptr, 
                                    TablePropertiesCollection_t* props)
{
    const ColumnFamilyHandle_t column_family = DBDefaultColumnFamily(dbptr);
    return DBGetPropertiesOfAllTablesWithColumnFamily(dbptr, &column_family, props);
}
#endif  // ROCKSDB_LITE

// Destroy the contents of the specified database.
// Be very careful using this method.
Status_t DBDestroyDB(const String_t* name, const Options_t* options)
{
    assert(GET_REP(name, String) != NULL);
    assert(GET_REP(options, Options) != NULL);
    Status stat = DestroyDB(GET_REP_REF(name, String), GET_REP_REF(options, Options));
    return NewStatusTCopy(&stat);
}

#ifndef ROCKSDB_LITE
// If a DB cannot be opened, you may attempt to call this method to
// resurrect as much of the contents of the database as possible.
// Some data may be lost, so be careful when calling this function
// on a database that contains important information.
Status_t DBRepairDB(const String_t* dbname, const Options_t* options)
{
    assert(GET_REP(dbname, String) != NULL);
    assert(GET_REP(options, Options) != NULL);
    Status stat = RepairDB(GET_REP_REF(dbname, String), GET_REP_REF(options, Options));
    return NewStatusTCopy(&stat);
}
#endif
