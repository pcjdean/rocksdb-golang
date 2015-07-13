// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_STATUS_H_
#define GO_ROCKSDB_INCLUDE_STATUS_H_

#include "types.h"
#include "string.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(Status)

extern bool StatusOk(Status_t *stat);
extern bool StatusIsNotFound(Status_t *stat);
extern bool StatusIsCorruption(Status_t *stat);
extern bool StatusIsNotSupported(Status_t *stat);
extern bool StatusIsInvalidArgument(Status_t *stat);
extern bool StatusIsInvalidArgument(Status_t *stat);
extern bool StatusIsMergeInProgress(Status_t *stat);
extern bool StatusIsIncomplete(Status_t *stat);
extern bool StatusIsShutdownInProgress(Status_t *stat);
extern bool StatusIsTimedOut(Status_t *stat);
extern bool StatusIsAborted(Status_t *stat);
extern bool StatusIsBusy(Status_t *stat);
extern String_t StatusToString(Status_t *stat);

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_STATUS_H_
