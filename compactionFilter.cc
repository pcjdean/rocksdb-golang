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

DEFINE_C_WRAP_CONSTRUCTOR(PCompactionFilterFactory)
DEFINE_C_WRAP_DESTRUCTOR(PCompactionFilterFactory)

DEFINE_C_WRAP_CONSTRUCTOR(PCompactionFilterFactoryV2)
DEFINE_C_WRAP_DESTRUCTOR(PCompactionFilterFactoryV2)

// C++ wrap class for go ICompactionFilter
// CompactionFilter allows an application to modify/delete a key-value at
// the time of compaction.
class CompactionFilterGo : public CompactionFilter {
public:
    FilterPolicyGo(void* go_flp)
        : m_go_flp(go_flp)
        , m_name(nullptr)
    {
        if (go_flp)
        {
            m_name = IFilterPolicyName(go_flp);
        }
    }

    // Destructor
    ~FilterPolicyGo()
    {
        if (m_go_flp)
        {
            IFilterPolicyRemoveReference(m_go_flp);
        }

        if (m_name)
        {
            free(m_name);
        }
    }

    // Return the name of this policy.  Note that if the filter encoding
    // changes in an incompatible way, the name returned by this method
    // must be changed.  Otherwise, old incompatible filters may be
    // passed to methods of this type.
    virtual const char* Name() const
    {
        return m_name;
    }

    // keys[0,n-1] contains a list of keys (potentially with duplicates)
    // that are ordered according to the user supplied comparator.
    // Append a filter that summarizes keys[0,n-1] to *dst.
    //
    // Warning: do not change the initial contents of *dst.  Instead,
    // append the newly constructed filter to *dst.
    virtual void CreateFilter(const Slice* keys, int n, std::string* dst) const
    {
        Slice_t* slcs = new Slice_t[n];
        assert(slcs != NULL);
        for (int j = 0; j < n; j++)
        {
            slcs[j].rep = const_cast<Slice *>(&keys[j]);
        }

        if (m_go_flp)
        {
            String_t str = IFilterPolicyCreateFilter(m_go_flp, slcs, n);
            dst->append(GET_REP_REF(&str, String));
            DeleteStringT(&str, false);
        }

        if (slcs)
        {
            delete[] slcs;
        }
    }

    // "filter" contains the data appended by a preceding call to
    // CreateFilter() on this class.  This method must return true if
    // the key was in the list of keys passed to CreateFilter().
    // This method may return true or false if the key was not on the
    // list, but it should aim to return false with a high probability.
    virtual bool KeyMayMatch(const Slice& key, const Slice& filter) const
    {
        bool ret = false;

        if (m_go_flp)
        {
            Slice_t keyslc{ const_cast<Slice *>(&key) };
            Slice_t filterslc{ const_cast<Slice *>(&filter) };
            ret = IFilterPolicyKeyMayMatch(m_go_flp, &keyslc, &filterslc);
        }

        return ret;
    }

    // Get the FilterBitsBuilder, which is ONLY used for full filter block
    // It contains interface to take individual key, then generate filter
    virtual FilterBitsBuilder* GetFilterBitsBuilder()
    {
        // TODO
        return nullptr;
    }

    // Get the FilterBitsReader, which is ONLY used for full filter block
    // It contains interface to tell if key can be in filter
    // The input slice should NOT be deleted by FilterPolicy
    virtual FilterBitsReader* GetFilterBitsReader(const Slice& contents)
    {
        // TODO
        return nullptr;
    }

private:
    // Wrapped go IFilterPolicy
    void* m_go_flp;

    // The name of the filter policy
    char* m_name;
};

// Return a filter policy from a go filter policy
PFilterPolicy_t NewPFilterPolicy(void* go_flp)
{
    PFilterPolicy_t wrap_t;
    wrap_t.rep = new PFilterPolicy(go_flp ? new FilterPolicyGo(go_flp) : NULL);
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
