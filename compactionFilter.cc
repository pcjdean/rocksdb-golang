// Copyright (c) 2015, Dean ChaoJun Pan.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

#include <assert.h>
#include <rocksdb/slice_transform.h>
#include "compactionfilterPrivate.h"
#include "compactionfilter.h"

extern "C" {
#include "_cgo_export.h"
}

DEFINE_C_WRAP_CONSTRUCTOR(CompactionFilterContext)
DEFINE_C_WRAP_DESTRUCTOR(CompactionFilterContext)

DEFINE_C_WRAP_CONSTRUCTOR(CompactionFilter)
DEFINE_C_WRAP_DESTRUCTOR(CompactionFilter)

DEFINE_C_WRAP_CONSTRUCTOR(CompactionFilter_Context)
DEFINE_C_WRAP_DESTRUCTOR(CompactionFilter_Context)

DEFINE_C_WRAP_CONSTRUCTOR(PCompactionFilterFactory)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(PCompactionFilterFactory)
DEFINE_C_WRAP_DESTRUCTOR(PCompactionFilterFactory)

// C++ wrap class for go ICompactionFilter
// CompactionFilter allows an application to modify/delete a key-value at
// the time of compaction.
class CompactionFilterGo : public CompactionFilter {
public:
    CompactionFilterGo(void* go_cpf)
        : m_go_cpf(go_cpf)
        , m_name(nullptr)
    {
        if (go_cpf)
        {
            m_name = ICompactionFilterName(go_cpf);
        }
    }

    // Destructor
    ~CompactionFilterGo()
    {
        if (m_go_cpf)
        {
            InterfacesRemoveReference(m_go_cpf);
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
                        bool* value_changed) const override
    {
        bool ret = false;
        
        if (m_go_cpf)
        {
            Slice_t slc_key{const_cast<Slice *>(&key)};
            Slice_t slc_exval{const_cast<Slice *>(&existing_value)};
            String_t str{new_value};
            ret = ICompactionFilterFilter(m_go_cpf, level, &slc_key, &slc_exval, &str, value_changed);
        }

        return ret;
    }

    // Returns a name that identifies this compaction filter.
    // The name will be printed to LOG file on start up for diagnosis.
    virtual const char* Name() const override
    {
        return m_name;
    }

private:
    // Wrapped go ICompactionFilter
    void* m_go_cpf;

    // The name of the compaction filter
    char* m_name;
};

// Return a CompactionFilter from a go ICompactionFilter
CompactionFilter_t NewCompactionFilter(void* go_cpf)
{
    CompactionFilter_t wrap_t;
    wrap_t.rep = (go_cpf ? new CompactionFilterGo(go_cpf) : NULL);
    return wrap_t;
}

// Each compaction will create a new CompactionFilter allowing the
// application to know about different compactions
class CompactionFilterFactoryGo : public CompactionFilterFactory {
public:
    CompactionFilterFactoryGo(void* go_cpfac)
        : m_go_cpfac(go_cpfac)
        , m_name(nullptr)
    {
        if (go_cpfac)
        {
            m_name = ICompactionFilterFactoryName(go_cpfac);
        }
    }

    // Destructor
    ~CompactionFilterFactoryGo()
    {
        if (m_go_cpfac)
        {
            InterfacesRemoveReference(m_go_cpfac);
        }

        if (m_name)
        {
            free(m_name);
        }
    }

    virtual std::unique_ptr<CompactionFilter> CreateCompactionFilter(
        const CompactionFilter::Context& context)
    {
        std::unique_ptr<CompactionFilter> ret;
        
        if (m_go_cpfac)
        {
            CompactionFilter_Context_t cxt{const_cast<CompactionFilter::Context *>(&context)};
            ret.reset(new CompactionFilterGo(ICompactionFilterFactoryCreateCompactionFilter(m_go_cpfac, &cxt)));
        }

        return ret;
    }

    // Returns a name that identifies this compaction filter factory.
    virtual const char* Name() const
    {
        return m_name;
    }

private:
    // Wrapped go ICompactionFilterFactory
    void* m_go_cpfac;

    // The name of the CompactionFilterFactory
    char* m_name;
};

// Return a CompactionFilterFactory from a go ICompactionFilterFactory
PCompactionFilterFactory_t NewPCompactionFilterFactory(void* go_cpflt)
{
    PCompactionFilterFactory_t wrap_t;
    wrap_t.rep = new PCompactionFilterFactory(go_cpflt ? new CompactionFilterFactoryGo(go_cpflt) : NULL);
    return wrap_t;
}
