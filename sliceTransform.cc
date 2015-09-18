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
#include "sliceTransform.h"

using namespace rocksdb;

DEFINE_C_WRAP_STRUCT(SliceTransform)
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
            m_name = ISliceTransformName(go_flp);
        }
    }

    // Destructor
    ~SliceTransformGo()
    {
        if (m_go_stf)
        {
            ISliceTransformRemoveReference(m_go_stf);
        }

        if (m_name)
        {
            free(m_name);
        }
    }

private:
    // Wrapped go IFilterPolicy
    void* m_go_stf;

    // The name of the filter policy
    char* m_name;
};

// Return a SliceTransform from a go SliceTransform interface
SliceTransform_t SliceTransformNewSliceTransform(void* go_stf)
{
    SliceTransform_t wrap_t;
    wrap_t.rep = new SliceTransform(go_stf ? new SliceTransformGo(go_stf) : NULL);
    return wrap_t;
}

// Create a fixed prefix transform
SliceTransform_t SliceTransformNewFixedPrefixTransform(size_t prefix_len)
{
    SliceTransform_t wrap_t;
    wrap_t.rep = new NewFixedPrefixTransform(prefix_len);
    return wrap_t;
}

// Create a capped prefix transform
SliceTransform_t SliceTransformNewCappedPrefixTransform(size_t cap_len)
{
    SliceTransform_t wrap_t;
    wrap_t.rep = new NewCappedPrefixTransform(cap_len);
    return wrap_t;
}

// Create a noop transform
SliceTransform_t SliceTransformNewNoopTransform()
{
    SliceTransform_t wrap_t;
    wrap_t.rep = new NewNoopTransform();
    return wrap_t;
}
