// Copyright (c) 2013, Facebook, Inc.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#ifndef STORAGE_ROCKSDB_INCLUDE_TYPES_H_
#define STORAGE_ROCKSDB_INCLUDE_TYPES_H_

#define GET_REP(x) (x->rep)
#define GET_REP_REF(x) (*GET_REP(x))

#define DEFINE_C_WRAP_STRUCT(x) typedef struct x##_t   \
                              {                      \
                                  x##* rep;          \
                              } x##_t;

#define DEFINE_C_WRAP_STRUCT(x, y) typedef struct x##_t  \
                              {                      \
                                  y##* rep;          \
                              } x##_t;


#define DEFINE_C_WRAP_CONSTRUCTOR(x) inline x##_t New##x##T(x##* ptr) \
                                     { \
                                          x##_t wrap_t; \
                                          wrap_t.rep = ptr; \
                                          return wrap_t; \
                                     } \

#endif //  STORAGE_ROCKSDB_INCLUDE_TYPES_H_
