// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_FILTERPOLICY_H_
#define GO_ROCKSDB_INCLUDE_FILTERPOLICY_H_

#ifdef __cplusplus
#include <rocksdb/filter_policy.h>
using namespace rocksdb;
#endif

#include "types.h"

#ifdef __cplusplus
typedef std::shared_ptr<FilterPolicy> PFilterPolicy;
#endif

#ifdef __cplusplus
extern "C" {
#endif

DEFINE_C_WRAP_STRUCT(PFilterPolicy)
DEFINE_C_WRAP_CONSTRUCTOR_DEC(PFilterPolicy)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(PFilterPolicy)
DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(PFilterPolicy, int, bool)
DEFINE_C_WRAP_DESTRUCTOR_DEC(PFilterPolicy)

// Return a filter policy from a go filter policy
PFilterPolicy_t NewPFilterPolicy(void* go_flp);


#ifdef __cplusplus
}  /* end extern "C" */
#endif

#endif  // GO_ROCKSDB_INCLUDE_FILTERPOLICY_H_
