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

DEFINE_C_WRAP_CONSTRUCTOR(CompactionFilterV2)
DEFINE_C_WRAP_DESTRUCTOR(CompactionFilterV2)

DEFINE_C_WRAP_CONSTRUCTOR(PCompactionFilterFactory)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(PCompactionFilterFactory)
DEFINE_C_WRAP_DESTRUCTOR(PCompactionFilterFactory)

DEFINE_C_WRAP_CONSTRUCTOR(PCompactionFilterFactoryV2)
DEFINE_C_WRAP_CONSTRUCTOR_DEFAULT(PCompactionFilterFactoryV2)
DEFINE_C_WRAP_DESTRUCTOR(PCompactionFilterFactoryV2)

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

// CompactionFilterV2 that buffers kv pairs sharing the same prefix and let
// application layer to make individual decisions for all the kv pairs in the
// buffer.
class CompactionFilterV2Go : public CompactionFilterV2 {
public:
    CompactionFilterV2Go(void* go_cpfv2)
        : m_go_cpfv2(go_cpfv2)
        , m_name(nullptr)
    {
        if (go_cpfv2)
        {
            m_name = ICompactionFilterV2Name(go_cpfv2);
        }
    }

    // Destructor
    ~CompactionFilterV2Go()
    {
        if (m_go_cpfv2)
        {
            InterfacesRemoveReference(m_go_cpfv2);
        }

        if (m_name)
        {
            free(m_name);
        }
    }

    // The compaction process invokes this method for all the kv pairs
    // sharing the same prefix. It is a "roll-up" version of CompactionFilter.
    //
    // Each entry in the return vector indicates if the corresponding kv should
    // be preserved in the output of this compaction run. The application can
    // inspect the existing values of the keys and make decision based on it.
    //
    // When a value is to be preserved, the application has the option
    // to modify the entry in existing_values and pass it back through an entry
    // in new_values. A corresponding values_changed entry needs to be set to
    // true in this case. Note that the new_values vector contains only changed
    // values, i.e. new_values.size() <= values_changed.size().
    //
    virtual std::vector<bool> Filter(int level,
                                     const SliceVector& keys,
                                     const SliceVector& existing_values,
                                     std::vector<std::string>* new_values,
                                     std::vector<bool>* values_changed) const
    {
        std::vector<bool> ret;
        
        if (m_go_cpfv2)
        {
            SliceVector_t slcv_keys{const_cast<SliceVector *>(&keys)};
            SliceVector_t slcv_exvals{const_cast<SliceVector *>(&existing_values)};
            StringVector_t strs{new_values};
            BoolVector_t valchgs{values_changed};
            BoolVector_t rets{&ret};
            ICompactionFilterV2Filter(m_go_cpfv2, level, &slcv_keys, &slcv_exvals, &strs, &valchgs, &rets);
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
    // Wrapped go ICompactionFilterV2
    void* m_go_cpfv2;

    // The name of the CompactionFilterV2
    char* m_name;
};

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

// Each compaction will create a new CompactionFilterV2
//
// CompactionFilterFactoryV2 enables application to specify a prefix and use
// CompactionFilterV2 to filter kv-pairs in batches. Each batch contains all
// the kv-pairs sharing the same prefix.
//
// This is useful for applications that require grouping kv-pairs in
// compaction filter to make a purge/no-purge decision. For example, if the
// key prefix is user id and the rest of key represents the type of value.
// This batching filter will come in handy if the application's compaction
// filter requires knowledge of all types of values for any user id.
//
class CompactionFilterFactoryV2Go : public CompactionFilterFactoryV2 {
public:
    CompactionFilterFactoryV2Go(void* go_cpfacv2, const SliceTransform* prefix_extractor)
        : CompactionFilterFactoryV2(prefix_extractor)
        , m_go_cpfacv2(go_cpfacv2)
        , m_name(nullptr)
    {
        if (go_cpfacv2)
        {
            m_name = ICompactionFilterFactoryV2Name(go_cpfacv2);
            SliceTransform_t stf = ICompactionFilterFactoryV2GetPrefixExtractor(go_cpfacv2);
            SetPrefixExtractor(GET_REP(&stf, SliceTransform));
        }
    }

    // Delete prefix_extractor if not NULL
    void DelPrefixExtractor() {
        const SliceTransform* prefix_extractor = GetPrefixExtractor();
        if (prefix_extractor)
        {
            delete prefix_extractor;
            CompactionFilterFactoryV2::SetPrefixExtractor(nullptr);
        }
    }

    // Destructor
    ~CompactionFilterFactoryV2Go()
    {
        if (m_go_cpfacv2)
        {
            InterfacesRemoveReference(m_go_cpfacv2);
        }
        
        DelPrefixExtractor();
        
        if (m_name)
        {
            free(m_name);
        }
    }

    virtual std::unique_ptr<CompactionFilterV2> CreateCompactionFilterV2(
        const CompactionFilterContext& context)
    {
        std::unique_ptr<CompactionFilterV2> ret;
        
        if (m_go_cpfacv2)
        {
            CompactionFilterContext_t cxt{const_cast<CompactionFilterContext *>(&context)};
            ret.reset(new CompactionFilterV2Go(ICompactionFilterFactoryV2CreateCompactionFilterV2(m_go_cpfacv2, &cxt)));
        }

        return ret;
    }

    // Delete old prefix_extractor before setting the new one.
    void SetPrefixExtractor(const SliceTransform* prefix_extractor) {
        DelPrefixExtractor();
        CompactionFilterFactoryV2::SetPrefixExtractor(prefix_extractor);
    }

    // Returns a name that identifies this compaction filter factory.
    virtual const char* Name() const
    {
        return m_name;
    }

private:
    // Wrapped go CompactionFilterFactoryV2
    void* m_go_cpfacv2;

    // The name of the CompactionFilterFactoryV2
    char* m_name;
};

// Return a CompactionFilterFactoryV2 from a go ICompactionFilterFactoryV2
PCompactionFilterFactoryV2_t NewPCompactionFilterFactoryV2(void* go_cpflt, void* go_stf)
{
    PCompactionFilterFactoryV2_t wrap_t;
    SliceTransform_t stf{nullptr};
    if (go_stf)
    {
        stf = NewSliceTransform(go_stf);
    }
    
    wrap_t.rep = new PCompactionFilterFactoryV2(go_cpflt ? new CompactionFilterFactoryV2Go(go_cpflt, GET_REP(&stf, SliceTransform)) : NULL);
    return wrap_t;
}
