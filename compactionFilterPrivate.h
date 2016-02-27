// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_COMPACTION_FILTER_PRIVATE_H_
#define GO_ROCKSDB_INCLUDE_COMPACTION_FILTER_PRIVATE_H_

#ifdef __cplusplus
#include <rocksdb/compaction_filter.h>
using namespace rocksdb;

typedef rocksdb::CompactionFilter::Context CompactionFilter_Context;

typedef std::shared_ptr<CompactionFilterFactory> PCompactionFilterFactory;
#endif

#endif  // GO_ROCKSDB_INCLUDE_COMPACTION_FILTER_PRIVATE_H_
