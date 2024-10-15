#include "glua.h"

void *LuaShared(void)
{
    const char **paths = NULL;
    size_t path_count = 0;
    void *handle = NULL;

#ifdef _WIN32
#ifdef _WIN64
    static const char *win_paths[] = {
        "lua_shared.dll",
        "bin/win64/lua_shared.dll"};
#else
    static const char *win_paths[] = {
        "lua_shared.dll",
        "garrysmod/bin/lua_shared.dll",
        "bin/lua_shared.dll"};
#endif
    paths = win_paths;
    path_count = sizeof(win_paths) / sizeof(win_paths[0]);

#else
#if defined(__x86_64__) || defined(_M_X64)
    static const char *linux_paths[] = {
        "lua_shared.so",
        "bin/linux64/lua_shared.so"};
#else
    static const char *linux_paths[] = {
        "lua_shared_srv.so",
        "garrysmod/bin/lua_shared_srv.so",
        "bin/linux32/lua_shared.so"};
#endif
    paths = linux_paths;
    path_count = sizeof(linux_paths) / sizeof(linux_paths[0]);
#endif

    // Attempt to load the library from each path
    for (size_t i = 0; i < path_count; ++i)
    {
        handle = LOAD_LIBRARY(paths[i]);
        if (handle != NULL)
        {
            return handle; // Successfully loaded
        }
    }

    return NULL; // Failed to load
}

LIB_HANDLE hModule = NULL;

const char *load_lua_shared()
{
    hModule = LuaShared();
    if (hModule == NULL)
    {
        return "Failed to load Lua shared library";
    }

#define X(return_type, func_name, ...)                             \
    func_name##_ptr = GET_FUNCTION(hModule, #func_name);           \
    if (func_name##_ptr == NULL)                                   \
    {                                                              \
        return format("Failed to load function '%s'", #func_name); \
    }

    GLUA_FUNCTIONS

#undef X

    return NULL;
}

const char *unload_lua_shared()
{
    if (hModule == NULL)
    {
        return "Lua shared library is not loaded";
    }

    CLOSE_LIBRARY(hModule);
    hModule = NULL;

    return NULL;
}

// https://stackoverflow.com/a/45043324
#define CAT(A, B) A##B
#define SELECT(NAME, NUM) CAT(NAME##_, NUM)
#define COMPOSE(NAME, ARGS) NAME ARGS

#define GET_COUNT(_0, _1, _2, _3, _4, _5, _6 /* ad nauseam */, COUNT, ...) COUNT
#define EXPAND() , , , , , , // 6 commas (or 7 empty tokens)
#define VA_SIZE(...) COMPOSE(GET_COUNT, (EXPAND __VA_ARGS__(), 0, 6, 5, 4, 3, 2, 1))

#define VA_SELECT(NAME, return_type, func_name, ...) SELECT(NAME, VA_SIZE(__VA_ARGS__))(return_type, func_name, __VA_ARGS__)

#define X(return_type, func_name, ...) VA_SELECT(X_ARG, return_type, func_name, __VA_ARGS__)

#define X_ARG_0(return_type, func_name, _void) \
    void *func_name##_ptr;                     \
    return_type func_name##_wrap()             \
    {                                          \
        return ((func_name)func_name##_ptr)(); \
    }

#define X_ARG_1(return_type, func_name, t1)        \
    void *func_name##_ptr;                         \
    return_type func_name##_wrap(t1 arg1)          \
    {                                              \
        return ((func_name)func_name##_ptr)(arg1); \
    }

#define X_ARG_2(return_type, func_name, t1, t2)          \
    void *func_name##_ptr;                               \
    return_type func_name##_wrap(t1 arg1, t2 arg2)       \
    {                                                    \
        return ((func_name)func_name##_ptr)(arg1, arg2); \
    }

#define X_ARG_3(return_type, func_name, t1, t2, t3)            \
    void *func_name##_ptr;                                     \
    return_type func_name##_wrap(t1 arg1, t2 arg2, t3 arg3)    \
    {                                                          \
        return ((func_name)func_name##_ptr)(arg1, arg2, arg3); \
    }

#define X_ARG_4(return_type, func_name, t1, t2, t3, t4)              \
    void *func_name##_ptr;                                           \
    return_type func_name##_wrap(t1 arg1, t2 arg2, t3 arg3, t4 arg4) \
    {                                                                \
        return ((func_name)func_name##_ptr)(arg1, arg2, arg3, arg4); \
    }

#define X_ARG_5(return_type, func_name, t1, t2, t3, t4, t5)                   \
    void *func_name##_ptr;                                                    \
    return_type func_name##_wrap(t1 arg1, t2 arg2, t3 arg3, t4 arg4, t5 arg5) \
    {                                                                         \
        return ((func_name)func_name##_ptr)(arg1, arg2, arg3, arg4, arg5);    \
    }

#define X_ARG_6(return_type, func_name, t1, t2, t3, t4, t5, t6)                        \
    void *func_name##_ptr;                                                             \
    return_type func_name##_wrap(t1 arg1, t2 arg2, t3 arg3, t4 arg4, t5 arg5, t6 arg6) \
    {                                                                                  \
        return ((func_name)func_name##_ptr)(arg1, arg2, arg3, arg4, arg5, arg6);       \
    }

GLUA_FUNCTIONS

#undef CAT
#undef SELECT
#undef COMPOSE
#undef GET_COUNT
#undef EXPAND
#undef VA_SIZE
#undef VA_SELECT

#undef X
#undef X_ARG_0
#undef X_ARG_1
#undef X_ARG_2
#undef X_ARG_3
#undef X_ARG_4
#undef X_ARG_5

int lua_call_go(lua_State L)
{
    int result = 0;
    char *err = NULL;

    goLuaCallback(L, &result, &err);

    if (err != NULL)
    {
        lua_pushstring_wrap(L, err);
        free(err);
        lua_error_wrap(L);
        return 0; // unreachable
    }

    return result;
}

int luaCFunctionWrapper(void *f, lua_State L)
{
    return ((int (*)(lua_State))(f))(L);
}

int lua_debug_getinfo_at(lua_State L, int level, const char *what, lua_Debug *ar)
{
    if (lua_getstack_wrap(L, level, ar) != 0 && lua_getinfo_wrap(L, what, ar) != 0)
    {
        return 1;
    }
    return 0;
}

const char *lua_err_argmsg(lua_State L, int narg, const char *extramsg)
{
    const char *fname = "?";
    const char *namewhat = NULL;

    lua_Debug ar;
    if (lua_debug_getinfo_at(L, 0, "n", &ar) == 1)
    {
        if (ar.name != NULL)
        {
            fname = ar.name;
        }

        if (ar.namewhat != NULL)
        {
            namewhat = ar.namewhat;
        }
    }

    if (narg < 0 && narg > LUA_REGISTRYINDEX)
    {
        narg = lua_gettop_wrap(L) + narg + 1;
    }

    if (namewhat != NULL && strcmp(namewhat, "method") == 0)
    {
        narg--;
        if (narg == 0)
        {
            return format("calling '%s' on bad self (%s)", fname, extramsg);
        }
    }

    return format("bad argument #%d to '%s' (%s)", narg, fname, extramsg);
}

const char *lua_type_error(lua_State L, int narg, const char *tname)
{
    const char *type_name = lua_typename_wrap(L, lua_type_wrap(L, narg));
    const char *err = format("%s expected, got %s", tname, type_name);
    return lua_err_argmsg(L, narg, err);
}

const char *lua_tag_error(lua_State L, int narg, int tag)
{
    return lua_type_error(L, narg, lua_typename_wrap(L, tag));
}

Result_Double lua_check_number(lua_State L, int narg)
{
    if (lua_type_wrap(L, narg) != LUA_TNUMBER)
    {
        return Err_Double(lua_tag_error(L, narg, LUA_TNUMBER));
    }

    return Ok_Double(lua_tonumber_wrap(L, narg));
}

Result_String lua_check_string(lua_State L, int narg, size_t *len)
{
    if (lua_type_wrap(L, narg) != LUA_TSTRING)
    {
        return Err_String(lua_tag_error(L, narg, LUA_TSTRING));
    }

    return Ok_String(lua_tolstring_wrap(L, narg, len));
}

Result_Bool lua_check_bool(lua_State L, int narg)
{
    if (lua_type_wrap(L, narg) != LUA_TBOOLEAN)
    {
        return Err_Bool(lua_tag_error(L, narg, LUA_TBOOLEAN));
    }

    return Ok_Bool(lua_toboolean_wrap(L, narg));
}

const char *lua_check_table(lua_State L, int narg)
{
    if (lua_type_wrap(L, narg) != LUA_TTABLE)
    {
        return lua_tag_error(L, narg, LUA_TTABLE);
    }

    return NULL;
}

const char *lua_check_func(lua_State L, int narg)
{
    if (lua_type_wrap(L, narg) != LUA_TFUNCTION)
    {
        return lua_tag_error(L, narg, LUA_TFUNCTION);
    }

    return NULL;
}
