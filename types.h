// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_TYPES_H_
#define GO_ROCKSDB_INCLUDE_TYPES_H_

#define GET_MACRO3(_1,_2,_3,NAME,...) NAME

#define GET_REP(x,y) ((y##*)x->rep)
#define GET_REP_REF(x,y) (*GET_REP(x, y))

#define DEFINE_C_WRAP_STRUCT(x) typedef struct x##_t   \
    {                                                  \
        void* rep;                                     \
    } x##_t;                                           \
    extern inline void Delete##x##T(x##_t* ptr, bool self);    
    

// Used internally by the C/C++ code
#define DEFINE_C_WRAP_CONSTRUCTOR(x) inline x##_t New##x##T(x##* ptr) \
    { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)ptr;            \
            return wrap_t;                      \
    }                                           \

// Used internally by the C/C++ code
#define DEFINE_C_WRAP_CONSTRUCTOR_COPY(x) inline x##_t New##x##TCopy(x##* ptr) \
    { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x(*ptr);    \
            return wrap_t;                      \
    }                                           

// Used internally by the C/C++ code
#define DEFINE_C_WRAP_CONSTRUCTOR_MOVE(x) inline x##_t New##x##TMove(x##* ptr) \
    { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x(std::move(*ptr)); \
            return wrap_t;                              \
    }                                                   

// Used externally by the calling code
#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS0(x) inline x##_t New##x##TArgs() \
    { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x();        \
            return wrap_t;                      \
    }                                           

#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS1(x,a) inline x##_t New##x##TArgs(a##_t* ptr_a) \
    { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x(GET_REP_REF(ptr_a, a));   \
            return wrap_t;                                      \
    }                                                           

#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS2(x,a,b) inline x##_t New##x##TArgs(a##_t* ptr_a, b##_t* ptr_b) \
    { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x(GET_REP_REF(ptr_a, a), GET_REP_REF(ptr_b, b)); \
            return wrap_t;                                              \
    }                                                                   

// Used externally by the calling code
#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS(...) GET_MACRO3(__VA_ARGS__, DEFINE_C_WRAP_CONSTRUCTOR_ARGS2, DEFINE_C_WRAP_CONSTRUCTOR_ARGS1, DEFINE_C_WRAP_CONSTRUCTOR_ARGS0)(__VA_ARGS__)

// Used externally by the calling code
#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS0(x) DEFINE_C_WRAP_CONSTRUCTOR_ARGS0(x)

#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS1(x,a) inline x##_t New##x##TRawArgs(a _a) \
    { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x(_a);      \
            return wrap_t;                      \
    }                                           

#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS2(x,a,b) inline x##_t New##x##TRawArgs(a _a, b _b) \
    { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x(_a, _b);  \
            return wrap_t;                      \
    }                                           

// Used internally by the C/C++ code
#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(...) GET_MACRO3(__VA_ARGS__, DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS2, DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS1, DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS0)(__VA_ARGS__)

// Used externally by the calling code. But it's implemented inernally.
#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT0(x) DEFINE_C_WRAP_CONSTRUCTOR_ARGS0(x)

#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT1(x,a) inline x##_t New##x##TDefault() \
    { \
        return New##x##TRawArgs(a);             \
    }                                           

#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT2(x,a,b) inline x##_t New##x##TDefault() \
    { \
        return New##x##TRawArgs(a, b);          \
    }                                           

// Used externally by the calling code
#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(...) GET_MACRO3(__VA_ARGS__, DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT2, DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT1, DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT0)(__VA_ARGS__)

// Used externally by the calling code
#define DEFINE_C_WRAP_DESTRUCTOR(x) inline void Delete##x##T(x##_t* ptr, bool self) \
    { \
        if (ptr) \
        {                                       \
            x* rep = GET_REP(ptr,x);            \
            if (rep)                            \
                delete rep;                     \
            if (self)                           \
                delete ptr;                     \
        }                                       \
    } 


#endif //  GO_ROCKSDB_INCLUDE_TYPES_H_
