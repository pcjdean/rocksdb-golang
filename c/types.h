// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_TYPES_H_
#define GO_ROCKSDB_INCLUDE_TYPES_H_

#define GET_MACRO2(_1,_2,NAME,...) NAME
#define GET_MACRO3(_1,_2,_3,NAME,...) NAME

#define GET_REP(x,y) ((y##*)x->rep)
#define GET_REP_REF(x,y) (*GET_REP(x, y))

#define DEFINE_C_WRAP_STRUCT(x) typedef struct x##_t   \
                              {                      \
                                  void* rep;          \
                              } x##_t;

#define DEFINE_C_WRAP_CONSTRUCTOR1(x) inline x##_t New##x##T(x##* ptr) \
                                     { \
                                          x##_t wrap_t; \
                                          wrap_t.rep = (void*)ptr;  \
                                          return wrap_t; \
                                     } \

#define DEFINE_C_WRAP_CONSTRUCTOR2(x,y) inline x##_t New##x##T(y##* ptr)  \
                                     { \
                                          x##_t wrap_t; \
                                          wrap_t.rep = (void*)ptr;  \
                                          return wrap_t; \
                                     } \

#define DEFINE_C_WRAP_CONSTRUCTOR(...) GET_MACRO2(__VA_ARGS__, DEFINE_C_WRAP_CONSTRUCTOR2, DEFINE_C_WRAP_CONSTRUCTOR1)(__VA_ARGS__)

#define DEFINE_C_WRAP_CONSTRUCTOR_COPY1(x) inline x##_t New##x##TCopy(x##* ptr) \
                                     { \
                                          x##_t wrap_t; \
                                          wrap_t.rep = (void*)new x(*ptr); \
                                          return wrap_t; \
                                     } \

#define DEFINE_C_WRAP_CONSTRUCTOR_COPY2(x,y) inline x##_t New##x##TCopy(y##* ptr)  \
                                     { \
                                          x##_t wrap_t; \
                                          wrap_t.rep = (void*)new y(*ptr);  \
                                          return wrap_t; \
                                     } \

#define DEFINE_C_WRAP_CONSTRUCTOR_COPY(...) GET_MACRO2(__VA_ARGS__, DEFINE_C_WRAP_CONSTRUCTOR_COPY2, DEFINE_C_WRAP_CONSTRUCTOR_COPY1)(__VA_ARGS__)

#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS0(x) inline x##_t New##x##TArgs(x##_t* ptr) \
                                     { \
                                          x##_t wrap_t; \
                                          wrap_t.rep = (void*)new x(GET_REP_REF(ptr)); \
                                          return wrap_t; \
                                     } \

#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS1(x,a) inline x##_t New##x##TArgs(a##_t* ptr_a) \
                                     { \
                                          x##_t wrap_t; \
                                          wrap_t.rep = (void*)new x(GET_REP_REF(ptr_a)); \
                                          return wrap_t; \
                                     } \

#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS2(x,a,b) inline x##_t New##x##TArgs(a##_t* ptr_a, b##_t* ptr_b) \
                                     { \
                                          x##_t wrap_t; \
                                          wrap_t.rep = (void*)new x(GET_REP_REF(ptr_a), GET_REP_REF(ptr_b)); \
                                          return wrap_t; \
                                     } \

#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS(...) GET_MACRO3(__VA_ARGS__, DEFINE_C_WRAP_CONSTRUCTOR_ARGS2, DEFINE_C_WRAP_CONSTRUCTOR_ARGS1, DEFINE_C_WRAP_CONSTRUCTOR_ARGS0)(__VA_ARGS__)

#define DEFINE_C_WRAP_DESTRUCTOR1(x) inline void Delete##x##T(x##_t* ptr) \
                                     { \
                                          if (ptr) \
                                          {        \
                                              x* rep = GET_REP(ptr,x); \
                                              if (rep) \
                                                  delete rep;   \
                                              delete ptr; \
                                          }        \
                                     } \

#define DEFINE_C_WRAP_DESTRUCTOR1(x,y) inline void Delete##x##T(x##_t* ptr) \
                                     { \
                                          if (ptr) \
                                          {        \
                                              y* rep = GET_REP(ptr,y); \
                                              if (rep) \
                                                  delete rep;   \
                                              delete ptr; \
                                          }        \
                                     } \

#define DEFINE_C_WRAP_DESTRUCTOR(...) GET_MACRO2(__VA_ARGS__, DEFINE_C_WRAP_DESTRUCTOR2, DEFINE_C_WRAP_DESTRUCTOR1)(__VA_ARGS__)


#endif //  GO_ROCKSDB_INCLUDE_TYPES_H_
