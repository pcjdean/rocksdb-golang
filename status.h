// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_STATUS_H_
#define GO_ROCKSDB_INCLUDE_STATUS_H_

#include "types.h"
#include "cstring.h"

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(Status)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(Status)
DEFINE_C_WRAP_DESTRUCTOR_DEC(Status)
DEFINE_C_WRAP_CONSTRUCTOR_COPY_DEC(Status)
DEFINE_C_WRAP_DESTRUCTOR_ARRAY_DEC(Status)

bool StatusOk(Status_t *stat);
bool StatusIsNotFound(Status_t *stat);
bool StatusIsCorruption(Status_t *stat);
bool StatusIsNotSupported(Status_t *stat);
bool StatusIsInvalidArgument(Status_t *stat);
bool StatusIsInvalidArgument(Status_t *stat);
bool StatusIsMergeInProgress(Status_t *stat);
bool StatusIsIncomplete(Status_t *stat);
bool StatusIsShutdownInProgress(Status_t *stat);
bool StatusIsTimedOut(Status_t *stat);
bool StatusIsAborted(Status_t *stat);
bool StatusIsBusy(Status_t *stat);
String_t StatusToString(Status_t *stat);

#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_STATUS_H_
