// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#include <assert.h>
#include <rocksdb/slice.h>
#include "slice.h"

extern "C" {
#include "_cgo_export.h"
}

DEFINE_C_WRAP_CONSTRUCTOR(CompactionFilterContext)
DEFINE_C_WRAP_DESTRUCTOR(CompactionFilterContext)

DEFINE_C_WRAP_CONSTRUCTOR(CompactionFilter)
DEFINE_C_WRAP_DESTRUCTOR(CompactionFilter)

DEFINE_C_WRAP_CONSTRUCTOR(CompactionFilter_Context)
DEFINE_C_WRAP_DESTRUCTOR(CompactionFilter_Context)

DEFINE_C_WRAP_CONSTRUCTOR(CompactionFilterV2)
DEFINE_C_WRAP_DESTRUCTOR(CompactionFilterV2)

DEFINE_C_WRAP_CONSTRUCTOR_DEC(CompactionFilterV2_SliceVector)
DEFINE_C_WRAP_DESTRUCTOR_DEC(CompactionFilterV2_SliceVector)

DEFINE_C_WRAP_CONSTRUCTOR(PCompactionFilterFactory)
DEFINE_C_WRAP_DESTRUCTOR(PCompactionFilterFactory)

DEFINE_C_WRAP_CONSTRUCTOR(PCompactionFilterFactoryV2)
DEFINE_C_WRAP_DESTRUCTOR(PCompactionFilterFactoryV2)

// C++ wrap class for go ICompactionFilter
// CompactionFilter allows an application to modify/delete a key-value at
// the time of compaction.
class CompactionFilterGo : public CompactionFilter {
public:
    CompactionFilterGo(void* go_cpflt)
        : m_go_cpflt(go_cpflt)
        , m_name(nullptr)
    {
        if (go_cpflt)
        {
            m_name = IFilterPolicyName(go_cpflt);
        }
    }

    // Destructor
    ~CompactionFilterGo()
    {
        if (m_go_cpflt)
        {
            IFilterPolicyRemoveReference(m_go_cpflt);
        }

        if (m_name)
        {
            free(m_name);
        }
    }

    // The compaction process invokes this
    // method for kv that is being compacted. A return value
    // of false indicates that the kv should be preserved in the
    // output of this compaction run and a return value of true
    // indicates that this key-value should be removed from the
    // output of the compaction.  The application can inspect
    // the existing value of the key and make decision based on it.
    //
    // When the value is to be preserved, the application has the option
    // to modify the existing_value and pass it back through new_value.
    // value_changed needs to be set to true in this case.
    //
    // If multithreaded compaction is being used *and* a single CompactionFilter
    // instance was supplied via Options::compaction_filter, this method may be
    // called from different threads concurrently.  The application must ensure
    // that the call is thread-safe.
    //
    // If the CompactionFilter was created by a factory, then it will only ever
    // be used by a single thread that is doing the compaction run, and this
    // call does not need to be thread-safe.  However, multiple filters may be
    // in existence and operating concurrently.
    virtual bool Filter(int level,
                        const Slice& key,
                        const Slice& existing_value,
                        std::string* new_value,
                        bool* value_changed) const
    {
        bool ret = false;
        
        if (m_go_cpflt)
        {
            Slice_t slc_key{&key};
            Slice_t slc_exval{&existing_value};
            String_t str{new_value};
            ret = ICompactionFilterFilter(m_go_cpflt, level, &slc_key, &slc_exval, &str, value_changed);
        }

        return ret;
    }

    // Returns a name that identifies this compaction filter.
    // The name will be printed to LOG file on start up for diagnosis.
    virtual const char* Name() const
    {
        return m_name;
    }

private:
    // Wrapped go IFilterPolicy
    void* m_go_cpflt;

    // The name of the filter policy
    char* m_name;
};

// Return a filter policy from a go filter policy
PFilterPolicy_t NewPFilterPolicy(void* go_cpflt)
{
    PFilterPolicy_t wrap_t;
    wrap_t.rep = new PFilterPolicy(go_cpflt ? new FilterPolicyGo(go_cpflt) : NULL);
    return wrap_t;
}

// Return a new filter policy that uses a bloom filter with approximately
// the specified number of bits per key.
//
// bits_per_key: bits per key in bloom filter. A good value for bits_per_key
// is 10, which yields a filter with ~ 1% false positive rate.
// use_block_based_builder: use block based filter rather than full fiter.
// If you want to builder full filter, it needs to be set to false.
//
// Callers must delete the result after any database that is using the
// result has been closed.
//
// Note: if you are using a custom comparator that ignores some parts
// of the keys being compared, you must not use NewBloomFilterPolicy()
// and must provide your own FilterPolicy that also ignores the
// corresponding parts of the keys.  For example, if the comparator
// ignores trailing spaces, it would be incorrect to use a
// FilterPolicy (like NewBloomFilterPolicy) that does not ignore
// trailing spaces in keys.
PFilterPolicy_t NewPFilterPolicyTRawArgs(int bits_per_key, bool use_block_based_builder)
{
    PFilterPolicy_t wrap_t;
    wrap_t.rep = new PFilterPolicy(const_cast<FilterPolicy *>(NewBloomFilterPolicy(bits_per_key, use_block_based_builder)));
    return wrap_t;
}
