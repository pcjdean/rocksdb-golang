// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// A Comparator object provides a total order across slices that are
// used as keys in an sstable or a database.  A Comparator implementation
// must be thread-safe since rocksdb may invoke its methods concurrently
// from multiple threads.

#include <rocksdb/comparator.h>

using namespace rocksdb;

#include "comparator.h"

extern "C" {
#include "_cgo_export.h"
}

DEFINE_C_WRAP_CONSTRUCTOR(Comparator)
DEFINE_C_WRAP_DESTRUCTOR(Comparator)

// C++ wrap class for go IComparator
class ComparatorGo : public Comparator {
public:
    ComparatorGo(void* go_cmp)
        : m_go_cmp(go_cmp)
        , m_name(nullptr)
    {
        if (go_cmp)
        {
            m_name = IComparatorName(go_cmp);
        }
    }

    // Destructor
    ~ComparatorGo()
    {
        if (m_go_cmp)
        {
            InterfacesRemoveReference(m_go_cmp);
        }

        if (m_name)
        {
            free(m_name);
        }
    }

    // Return the name of this transformation.
    virtual const char* Name() const override
    {
        return m_name;
    }

    // Three-way comparison.  Returns value:
    //   < 0 iff "a" < "b",
    //   == 0 iff "a" == "b",
    //   > 0 iff "a" > "b"
    virtual int Compare(const Slice& a, const Slice& b) const override
    {
        int ret;
        if (m_go_cmp)
        {
            Slice_t a_slc{const_cast<Slice *>(&a)};
            Slice_t b_slc{const_cast<Slice *>(&b)};
            ret = IComparatorCompare(m_go_cmp, &a_slc, &b_slc); 
        }
        else
        {
            ret = a.compare(b);
        }

        return ret;
    }

    // Advanced functions: these are used to reduce the space requirements
    // for internal data structures like index blocks.

    // If *start < limit, changes *start to a short string in [start,limit).
    // Simple comparator implementations may return with *start unchanged,
    // i.e., an implementation of this method that does nothing is correct.
    virtual void FindShortestSeparator(std::string* start, const Slice& limit) const override
    {
        if (m_go_cmp)
        {
            String_t start_str{start};
            Slice_t limit_slc{const_cast<Slice *>(&limit)};
            size_t sz = 0;
            char* ret = IComparatorFindShortestSeparator(m_go_cmp, &start_str, &limit_slc, &sz); 
            if (ret)
            {
                start->assign(ret, sz);
                free(ret);
            }
        }
    }

    // Changes *key to a short string >= *key.
    // Simple comparator implementations may return with *key unchanged,
    // i.e., an implementation of this method that does nothing is correct.
    virtual void FindShortSuccessor(std::string* key) const override
    {
        if (m_go_cmp)
        {
            String_t key_str{key};
            size_t sz = 0;
            char* ret = IComparatorFindShortSuccessor(m_go_cmp, &key_str, &sz); 
            if (ret)
            {
                key->assign(ret, sz);
                free(ret);
            }
        }
    }

private:
    // Wrapped go IComparator
    void* m_go_cmp;

    // The name of the Comparator
    char* m_name;
};

// Return a Comparator from a go Comparator interface
Comparator_t NewComparator(void* go_cmp)
{
    Comparator_t wrap_t;
    wrap_t.rep = (go_cmp ? new ComparatorGo(go_cmp) : NULL);
    return wrap_t;
}

// Return a builtin comparator that uses lexicographic byte-wise
// ordering.  The result remains the property of this module and
// must not be deleted.
Comparator_t GoBytewiseComparator()
{
    Comparator_t wrap_t;
    wrap_t.rep = const_cast<Comparator *>(BytewiseComparator());
    return wrap_t;
}

// Return a builtin comparator that uses reverse lexicographic byte-wise
// ordering.
Comparator_t GoReverseBytewiseComparator()
{
    Comparator_t wrap_t;
    wrap_t.rep = const_cast<Comparator *>(ReverseBytewiseComparator());
    return wrap_t;
}
