// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Class for specifying user-defined functions which perform a
// transformation on a slice.  It is not required that every slice
// belong to the domain and/or range of a function.  Subclasses should
// define InDomain and InRange to determine which slices are in either
// of these sets respectively.

#include <rocksdb/slice_transform.h>

using namespace rocksdb;

#include "sliceTransform.h"

extern "C" {
#include "_cgo_export.h"
}

DEFINE_C_WRAP_CONSTRUCTOR(SliceTransform)
DEFINE_C_WRAP_DESTRUCTOR(SliceTransform)

// C++ wrap class for go ISliceTransform
class SliceTransformGo : public SliceTransform {
public:
    SliceTransformGo(void* go_stf)
        : m_go_stf(go_stf)
        , m_name(nullptr)
    {
        if (go_stf)
        {
            m_name = ISliceTransformName(go_stf);
        }
    }

    // Destructor
    ~SliceTransformGo()
    {
        if (m_go_stf)
        {
            InterfacesRemoveReference(m_go_stf);
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

    // transform a src in domain to a dst in the range
    virtual Slice Transform(const Slice& src) const override
    {
        if (m_go_stf)
        {
            Slice_t slc{const_cast<Slice *>(&src)};
            size_t offset = 0;
            size_t len = 0;
            ISliceTransformTransform(m_go_stf, &slc, &offset, &len);
            return Slice{src.data() + offset, len};
        }

        return Slice{};
    }
    
    // determine whether this is a valid src upon the function applies
    virtual bool InDomain(const Slice& src) const override
    {
        bool ret = false;
        if (m_go_stf)
        {
            Slice_t slc{const_cast<Slice *>(&src)};
            ret = ISliceTransformInDomain(m_go_stf, &slc);
        }
        return ret;
    }

    // determine whether dst=Transform(src) for some src
    virtual bool InRange(const Slice& dst) const override
    {
        bool ret = false;
        if (m_go_stf)
        {
            Slice_t slc{const_cast<Slice *>(&dst)};
            ret = ISliceTransformInRange(m_go_stf, &slc); 
        }
        return ret;
    }

    // Transform(s)=Transform(`prefix`) for any s with `prefix` as a prefix.
    //
    // This function is not used by RocksDB, but for users. If users pass
    // Options by string to RocksDB, they might not know what prefix extractor
    // they are using. This function is to help users can determine:
    //   if they want to iterate all keys prefixing `prefix`, whetherit is
    //   safe to use prefix bloom filter and seek to key `prefix`.
    // If this function returns true, this means a user can Seek() to a prefix
    // using the bloom filter. Otherwise, user needs to skip the bloom filter
    // by setting ReadOptions.total_order_seek = true.
    //
    // Here is an example: Suppose we implement a slice transform that returns
    // the first part of the string after spliting it using deimiter ",":
    // 1. SameResultWhenAppended("abc,") should return true. If aplying prefix
    //    bloom filter using it, all slices matching "abc:.*" will be extracted
    //    to "abc,", so any SST file or memtable containing any of those key
    //    will not be filtered out.
    // 2. SameResultWhenAppended("abc") should return false. A user will not
    //    guaranteed to see all the keys matching "abc.*" if a user seek to "abc"
    //    against a DB with the same setting. If one SST file only contains
    //    "abcd,e", the file can be filtered out and the key will be invisible.
    //
    // i.e., an implementation always returning false is safe.
    virtual bool SameResultWhenAppended(const Slice& prefix) const override
    {
        bool ret = false;
        if (m_go_stf)
        {
            Slice_t slc{const_cast<Slice *>(&prefix)};
            ret = ISliceTransformSameResultWhenAppended(m_go_stf, &slc); 
        }
        return ret;
    }

private:
    // Wrapped go ISliceTransform
    void* m_go_stf;

    // The name of the SliceTransform
    char* m_name;
};

// Return a SliceTransform from a go SliceTransform interface
SliceTransform_t NewSliceTransform(void* go_stf)
{
    SliceTransform_t wrap_t;
    wrap_t.rep = (go_stf ? new SliceTransformGo(go_stf) : NULL);
    return wrap_t;
}

// Create a fixed prefix transform
SliceTransform_t GoNewFixedPrefixTransform(size_t prefix_len)
{
    SliceTransform_t wrap_t;
    wrap_t.rep = const_cast<SliceTransform *>(NewFixedPrefixTransform(prefix_len));
    return wrap_t;
}

// Create a capped prefix transform
SliceTransform_t GoNewCappedPrefixTransform(size_t cap_len)
{
    SliceTransform_t wrap_t;
    wrap_t.rep = const_cast<SliceTransform *>(NewCappedPrefixTransform(cap_len));
    return wrap_t;
}

// Create a noop transform
SliceTransform_t GoNewNoopTransform()
{
    SliceTransform_t wrap_t;
    wrap_t.rep = const_cast<SliceTransform *>(NewNoopTransform());
    return wrap_t;
}
