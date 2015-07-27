// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#ifndef GO_ROCKSDB_INCLUDE_TYPES_H_
#define GO_ROCKSDB_INCLUDE_TYPES_H_

#include <stdint.h>

typedef uint64_t SequenceNumber;
typedef uint64_t size_t;

#ifndef __cplusplus
typedef char bool;

enum bool_t {
    false, true
};
    
#endif

#define GET_MACRO3(_1,_2,_3,NAME,...) NAME

#define GET_REP(x,y) ((y*)((x)->rep))
#define GET_REP_REF(x,y) (*GET_REP(x, y))

#define DEFINE_C_WRAP_STRUCT(x) typedef struct x##_t   \
    {                                                  \
        void* rep;                                     \
    } x##_t;     

    

// Used internally by the C/C++ code
#define DEFINE_C_WRAP_CONSTRUCTOR_DEC_R(x) x##_t New##x##T(void* ptr)
#define DEFINE_C_WRAP_CONSTRUCTOR_BODY(x) { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)ptr;            \
            return wrap_t;                      \
    }                                           

#define DEFINE_C_WRAP_CONSTRUCTOR_DEC(x) DEFINE_C_WRAP_CONSTRUCTOR_DEC_R(x);
#define DEFINE_C_WRAP_CONSTRUCTOR(x) DEFINE_C_WRAP_CONSTRUCTOR_DEC_R(x) \
    DEFINE_C_WRAP_CONSTRUCTOR_BODY(x)



// Used internally by the C/C++ code
#define DEFINE_C_WRAP_CONSTRUCTOR_COPY_DEC_R(x)  x##_t New##x##TCopy(void* ptr)
#define DEFINE_C_WRAP_CONSTRUCTOR_COPY_BODY(x) { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x(*(x*)ptr);      \
            return wrap_t;                      \
    }                                           

#define DEFINE_C_WRAP_CONSTRUCTOR_COPY_DEC(x) DEFINE_C_WRAP_CONSTRUCTOR_COPY_DEC_R(x);
#define DEFINE_C_WRAP_CONSTRUCTOR_COPY(x) DEFINE_C_WRAP_CONSTRUCTOR_COPY_DEC_R(x) \
    DEFINE_C_WRAP_CONSTRUCTOR_COPY_BODY(x)



// Used internally by the C/C++ code
#define DEFINE_C_WRAP_CONSTRUCTOR_MOVE_DEC_R(x) x##_t New##x##TMove(void* ptr)
#define DEFINE_C_WRAP_CONSTRUCTOR_MOVE_BODY(x) { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x(std::move(*(x*)ptr)); \
            return wrap_t;                              \
    }                                                   

#define DEFINE_C_WRAP_CONSTRUCTOR_MOVE_DEC(x) DEFINE_C_WRAP_CONSTRUCTOR_MOVE_DEC_R(x);
#define DEFINE_C_WRAP_CONSTRUCTOR_MOVE(x) DEFINE_C_WRAP_CONSTRUCTOR_MOVE_DEC_R(x) \
    DEFINE_C_WRAP_CONSTRUCTOR_MOVE_BODY(x)



// Used externally by the calling code
#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS0_DEC_R(x) x##_t New##x##TArgs()
#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS0_BODY(x) { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x();        \
            return wrap_t;                      \
    }                                           

#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS0_DEC(x) DEFINE_C_WRAP_CONSTRUCTOR_ARGS0_DEC_R(x);
#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS0(x) DEFINE_C_WRAP_CONSTRUCTOR_ARGS0_DEC_R(x) \
    DEFINE_C_WRAP_CONSTRUCTOR_ARGS0_BODY(x)



#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS1_DEC_R(x,a) x##_t New##x##TArgs(a##_t* ptr_a)
#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS1_BODY(x,a) { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x(GET_REP_REF(ptr_a, a));   \
            return wrap_t;                                      \
    }      

#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS1_DEC(x,a) DEFINE_C_WRAP_CONSTRUCTOR_ARGS1_DEC_R(x,a);
#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS1(x,a) DEFINE_C_WRAP_CONSTRUCTOR_ARGS1_DEC_R(x,a) \
    DEFINE_C_WRAP_CONSTRUCTOR_ARGS1_BODY(x,a)

                                                     

#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS2_DEC_R(x,a,b) x##_t New##x##TArgs(a##_t* ptr_a, b##_t* ptr_b)
#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS2_BODY(x,a,b) { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x(GET_REP_REF(ptr_a, a), GET_REP_REF(ptr_b, b)); \
            return wrap_t;                                              \
    }                                                                   

#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS2_DEC(x,a,b) DEFINE_C_WRAP_CONSTRUCTOR_ARGS2_DEC_R(x,a,b);
#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS2(x,a,b) DEFINE_C_WRAP_CONSTRUCTOR_ARGS2_DEC_R(x,a,b) \
    DEFINE_C_WRAP_CONSTRUCTOR_ARGS2_BODY(x,a,b)



// Used externally by the calling code
#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS_DEC(...) GET_MACRO3(__VA_ARGS__, DEFINE_C_WRAP_CONSTRUCTOR_ARGS2_DEC, DEFINE_C_WRAP_CONSTRUCTOR_ARGS1_DEC, DEFINE_C_WRAP_CONSTRUCTOR_ARGS0_DEC)(__VA_ARGS__)
#define DEFINE_C_WRAP_CONSTRUCTOR_ARGS(...) GET_MACRO3(__VA_ARGS__, DEFINE_C_WRAP_CONSTRUCTOR_ARGS2, DEFINE_C_WRAP_CONSTRUCTOR_ARGS1, DEFINE_C_WRAP_CONSTRUCTOR_ARGS0)(__VA_ARGS__)



// Used internally by the calling code
#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS0_DEC_R(x)  x##_t New##x##TRawArgs()
#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS0_DEC(x) DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS0_DEC_R(x);
#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS0(x) DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS0_DEC_R(x) \
    DEFINE_C_WRAP_CONSTRUCTOR_ARGS0_BODY(x)



#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS1_DEC_R(x,a) x##_t New##x##TRawArgs(a _a)
#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS1_BODY(x,a) { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x(_a);      \
            return wrap_t;                      \
    } 
                                          
#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS1_DEC(x,a) DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS1_DEC_R(x,a);
#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS1(x,a) DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS1_DEC_R(x,a) \
    DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS1_BODY(x,a)



#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS2_DEC_R(x,a,b) x##_t New##x##TRawArgs(a _a, b _b)
#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS2_BODY(x,a,b)  { \
        x##_t wrap_t; \
            wrap_t.rep = (void*)new x(_a, _b);  \
            return wrap_t;                      \
    } 
                                          
#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS2_DEC(x,a,b) DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS2_DEC_R(x,a,b);
#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS2(x,a,b) DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS2_DEC_R(x,a,b) \
    DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS2_BODY(x,a,b)



// Used internally by the C/C++ code
#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS_DEC(...) GET_MACRO3(__VA_ARGS__, DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS2_DEC, DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS1_DEC, DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS0_DEC)(__VA_ARGS__)
#define DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS(...) GET_MACRO3(__VA_ARGS__, DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS2, DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS1, DEFINE_C_WRAP_CONSTRUCTOR_RAW_ARGS0)(__VA_ARGS__)

// Used externally by the calling code. But it's implemented inernally.
#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT0_DEC_R(x) x##_t New##x##TDefault()
#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT0_DEC(x) DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT0_DEC_R(x);
#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT0(x) DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT0_DEC_R(x) \
    DEFINE_C_WRAP_CONSTRUCTOR_ARGS0_BODY(x)



#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT1_DEC_R(x,a)  x##_t New##x##TDefault()
#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT1_BODY(x,a) { \
        return New##x##TRawArgs(a);             \
    }    
                                       
#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT1_DEC(x,a) DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT1_DEC_R(x,a);
#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT1(x,a) DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT1_DEC_R(x,a) \
    DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT1_BODY(x,a)



#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT2_DEC_R(x,a,b) x##_t New##x##TDefault()
#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT2_BODY(x,a,b) { \
        return New##x##TRawArgs(a, b);          \
    } 
                                          
#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT2_DEC(x,a,b) DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT2_DEC_R(x,a,b);
#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT2(x,a,b) DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT2_DEC_R(x,a,b) \
    DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT2_BODY(x,a,b)

// Used externally by the calling code
#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT_DEC(...) GET_MACRO3(__VA_ARGS__, DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT2_DEC, DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT1_DEC, DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT0_DEC)(__VA_ARGS__)
#define DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(...) GET_MACRO3(__VA_ARGS__, DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT2, DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT1, DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT0)(__VA_ARGS__)



// Used externally by the calling code
#define DEFINE_C_WRAP_DESTRUCTOR_DEC_R(x)  void Delete##x##T(x##_t* ptr, bool self)
#define DEFINE_C_WRAP_DESTRUCTOR_BODY(x) { \
        if (ptr) \
        {                                       \
            x* rep = GET_REP(ptr,x);            \
            if (rep)                            \
                delete rep;                     \
            if (self)                           \
                delete ptr;                     \
        }                                       \
    } 

#define DEFINE_C_WRAP_DESTRUCTOR_DEC(x) DEFINE_C_WRAP_DESTRUCTOR_DEC_R(x);
#define DEFINE_C_WRAP_DESTRUCTOR(x) DEFINE_C_WRAP_DESTRUCTOR_DEC_R(x) \
    DEFINE_C_WRAP_DESTRUCTOR_BODY(x)



// Used externally by the calling code
#define DEFINE_C_WRAP_DESTRUCTOR_ARRAY_DEC_R(x) void Delete##x##TArray(x##_t* ptr)
#define DEFINE_C_WRAP_DESTRUCTOR_ARRAY_BODY(x) { \
        if (ptr) \
        {                                       \
            delete[] (x##_t*)ptr;                   \
        }                                       \
    } 

#define DEFINE_C_WRAP_DESTRUCTOR_ARRAY_DEC(x) DEFINE_C_WRAP_DESTRUCTOR_ARRAY_DEC_R(x);
#define DEFINE_C_WRAP_DESTRUCTOR_ARRAY(x) DEFINE_C_WRAP_DESTRUCTOR_ARRAY_DEC_R(x) \
    DEFINE_C_WRAP_DESTRUCTOR_ARRAY_BODY(x)

#ifndef __cplusplus

static inline size_t intToSizeT(int i) 
{
    return (size_t)i;
}

#endif

#endif //  GO_ROCKSDB_INCLUDE_TYPES_H_
