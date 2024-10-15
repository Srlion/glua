#define GLUA_FUNCTIONS                                                                    \
    /* state manipulation */                                                              \
    X(lua_State, luaL_newstate)                                                           \
    X(lua_State, lua_newthread, lua_State)                                                \
    /* basic stack manipulation */                                                        \
    X(int, lua_gettop, lua_State)                                                         \
    X(void, lua_settop, lua_State, int)                                                   \
    X(void, lua_pushvalue, lua_State, int)                                                \
    X(void, lua_remove, lua_State, int)                                                   \
    X(void, lua_insert, lua_State, int)                                                   \
    X(void, lua_replace, lua_State, int)                                                  \
    X(int, lua_checkstack, lua_State, int)                                                \
    /* access functions (stack -> C) */                                                   \
    X(int, lua_type, lua_State, int)                                                      \
    X(const char *, lua_typename, lua_State, int)                                         \
    X(int, lua_equal, lua_State, int, int)                                                \
    X(int, lua_rawequal, lua_State, int, int)                                             \
    X(int, lua_lessthan, lua_State, int, int)                                             \
    X(double, lua_tonumber, lua_State, int)                                               \
    X(int, lua_toboolean, lua_State, int)                                                 \
    X(const char *, lua_tolstring, lua_State, int, size_t *)                              \
    X(size_t, lua_objlen, lua_State, int)                                                 \
    X(void *, lua_tocfunction, lua_State, int)                                            \
    X(void *, lua_touserdata, lua_State, int)                                             \
    X(lua_State, lua_tothread, lua_State, int)                                            \
    X(void *, lua_topointer, lua_State, int)                                              \
    /* push functions (C -> stack) */                                                     \
    X(void, lua_pushnil, lua_State)                                                       \
    X(void, lua_pushnumber, lua_State, double)                                            \
    X(void, lua_pushlstring, lua_State, const char *, size_t)                             \
    X(void, lua_pushstring, lua_State, const char *)                                      \
    X(void, lua_pushcclosure, lua_State, void *, int)                                     \
    X(void, lua_pushboolean, lua_State, int)                                              \
    X(void, lua_pushlightuserdata, lua_State, void *)                                     \
    X(int, lua_pushthread, lua_State)                                                     \
    /* get functions (Lua -> stack) */                                                    \
    X(void, lua_gettable, lua_State, int)                                                 \
    X(void, lua_getfield, lua_State, int, const char *)                                   \
    X(void, lua_rawget, lua_State, int)                                                   \
    X(void, lua_rawgeti, lua_State, int, int)                                             \
    X(void, lua_createtable, lua_State, int, int)                                         \
    X(void *, lua_newuserdata, lua_State, size_t)                                         \
    X(int, lua_getmetatable, lua_State, int)                                              \
    X(void, lua_getfenv, lua_State, int)                                                  \
    /* set functions (stack -> Lua) */                                                    \
    X(void, lua_settable, lua_State, int)                                                 \
    X(void, lua_setfield, lua_State, int, const char *)                                   \
    X(void, lua_rawset, lua_State, int)                                                   \
    X(void, lua_rawseti, lua_State, int, int)                                             \
    X(void, lua_setmetatable, lua_State, int)                                             \
    X(int, lua_setfenv, lua_State, int)                                                   \
    /* `load' and `call' functions (load and run Lua code) */                             \
    X(void, lua_call, lua_State, int, int)                                                \
    X(int, lua_pcall, lua_State, int, int, int)                                           \
    X(int, lua_cpcall, lua_State, void *, void *)                                         \
    /* coroutine functions */                                                             \
    X(int, lua_yield, lua_State, int)                                                     \
    X(int, lua_resume_real, lua_State, int)                                               \
    X(int, lua_status, lua_State)                                                         \
    /* miscellaneous functions */                                                         \
    X(int, lua_error, lua_State)                                                          \
    X(int, lua_next, lua_State, int)                                                      \
    X(void, lua_concat, lua_State, int)                                                   \
    X(void, luaL_openlibs, lua_State)                                                     \
    X(int, luaL_callmeta, lua_State, int, const char *)                                   \
    X(int, luaL_newmetatable, lua_State, const char *)                                    \
    X(int, luaL_ref, lua_State, int)                                                      \
    X(void, luaL_unref, lua_State, int, int)                                              \
    X(int, luaL_loadbuffer, lua_State, const char *, size_t, const char *)                \
    X(int, luaL_loadbufferx, lua_State, const char *, size_t, const char *, const char *) \
    X(int, luaL_loadstring, lua_State, const char *)                                      \
    X(int, luaL_loadfile, lua_State, const char *)                                        \
    X(const char *, luaL_findtable, lua_State, int, const char *, int)                    \
    /* Functions to be called by the debugger in specific events */                       \
    X(int, lua_getstack, lua_State, int, lua_Debug *)                                     \
    X(int, lua_getinfo, lua_State, const char *, lua_Debug *)
