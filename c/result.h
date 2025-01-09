#pragma once

#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#define DECLARE_RESULT_TYPE(TypeName, ValueType, Default)           \
    typedef struct                                                  \
    {                                                               \
        ValueType value;                                            \
        const char *err;                                            \
    } Result_##TypeName;                                            \
                                                                    \
    /* Function to construct a success result */                    \
    static inline Result_##TypeName Ok_##TypeName(ValueType value)  \
    {                                                               \
        Result_##TypeName res;                                      \
        res.value = value;                                          \
        res.err = NULL;                                             \
        return res;                                                 \
    }                                                               \
                                                                    \
    /* Function to construct an error result */                     \
    static inline Result_##TypeName Err_##TypeName(const char *err) \
    {                                                               \
        Result_##TypeName res;                                      \
        res.value = Default;                                        \
        res.err = err;                                              \
        return res;                                                 \
    }

// not assigning a default value makes go crash?
DECLARE_RESULT_TYPE(Int, int, 0)
DECLARE_RESULT_TYPE(Float, float, 0)
DECLARE_RESULT_TYPE(Double, double, 0)
DECLARE_RESULT_TYPE(Pointer, void *, NULL)
DECLARE_RESULT_TYPE(String, const char *, NULL)
DECLARE_RESULT_TYPE(Bool, bool, false)
