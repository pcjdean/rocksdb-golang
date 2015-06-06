// Copyright (c) 2013, Facebook, Inc.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.
// Copyright (c) 2011 The LevelDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file. See the AUTHORS file for names of contributors.

#ifndef STORAGE_ROCKSDB_INCLUDE_DB_H_
#define STORAGE_ROCKSDB_INCLUDE_DB_H_

#include <stdint.h>
#include <stdio.h>
#include <memory>
#include <vector>
#include <string>
#include <unordered_map>
#include "metadata.h"
#include "version.h"
#include "iterator.h"
#include "options.h"
#include "types.h"
#include "transaction_log.h"
#include "listener.h"
#include "thread_status.h"

namespace rocksdb {

struct Options;
struct DBOptions;
struct ColumnFamilyOptions;
struct ReadOptions;
struct WriteOptions;
struct FlushOptions;
struct CompactionOptions;
struct TableProperties;
class WriteBatch;
class Env;
class EventListener;

using std::unique_ptr;

DEFINE_C_WRAP_STRUCT(ColumnFamilyHandle)

DEFINE_C_WRAP_STRUCT(String, std::string)

extern const std::string kDefaultColumnFamilyName;

DEFINE_C_WRAP_STRUCT(ColumnFamilyDescriptor)

static const int kMajorVersion = __ROCKSDB_MAJOR__;
static const int kMinorVersion = __ROCKSDB_MINOR__;

// Abstract handle to particular state of a DB.
// A Snapshot is an immutable object and can therefore be safely
// accessed from multiple threads without any external synchronization.
DEFINE_C_WRAP_STRUCT(Snapshot)

// A range of keys
struct Range {
  Slice start;          // Included in the range
  Slice limit;          // Not included in the range

  Range() { }
  Range(const Slice& s, const Slice& l) : start(s), limit(l) { }
};

// A collections of table properties objects, where
//  key: is the table's file name.
//  value: the table properties object of the given table.
typedef std::unordered_map<std::string, std::shared_ptr<const TableProperties>>
    TablePropertiesCollection;

// A DB is a persistent ordered map from keys to values.
// A DB is safe for concurrent access from multiple threads without
// any external synchronization.
DEFINE_C_WRAP_STRUCT(DB)

#ifndef ROCKSDB_LITE
  struct Properties {
    static const std::string kNumFilesAtLevelPrefix;
    static const std::string kStats;
    static const std::string kSSTables;
    static const std::string kCFStats;
    static const std::string kDBStats;
    static const std::string kNumImmutableMemTable;
    static const std::string kMemTableFlushPending;
    static const std::string kCompactionPending;
    static const std::string kBackgroundErrors;
    static const std::string kCurSizeActiveMemTable;
    static const std::string kCurSizeAllMemTables;
    static const std::string kNumEntriesActiveMemTable;
    static const std::string kNumEntriesImmMemTables;
    static const std::string kNumDeletesActiveMemTable;
    static const std::string kNumDeletesImmMemTables;
    static const std::string kEstimateNumKeys;
    static const std::string kEstimateTableReadersMem;
    static const std::string kIsFileDeletionsEnabled;
    static const std::string kNumSnapshots;
    static const std::string kOldestSnapshotTime;
    static const std::string kNumLiveVersions;
  };
#endif /* ROCKSDB_LITE */

// Destroy the contents of the specified database.
// Be very careful using this method.
Status DestroyDB(const std::string& name, const Options& options);

#ifndef ROCKSDB_LITE
// If a DB cannot be opened, you may attempt to call this method to
// resurrect as much of the contents of the database as possible.
// Some data may be lost, so be careful when calling this function
// on a database that contains important information.
Status RepairDB(const std::string& dbname, const Options& options);
#endif

}  // namespace rocksdb

#endif  // STORAGE_ROCKSDB_INCLUDE_DB_H_
