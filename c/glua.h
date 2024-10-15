#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stddef.h>
#include <stdarg.h>
#include <stdint.h>
#include <stdbool.h>
#include "result.h"
#include "glua_functions.h"
#include "cross_loader.h"

#define format(fmt, ...)                                        \
    ({                                                          \
        size_t len = snprintf(NULL, 0, fmt, ##__VA_ARGS__) + 1; \
        char *msg = malloc(len);                                \
        sprintf(msg, fmt, ##__VA_ARGS__);                       \
        msg;                                                    \
    })

typedef uintptr_t lua_State;

/*
** basic types
*/
#define LUA_TNONE (-1)

#define LUA_TNIL 0
#define LUA_TBOOLEAN 1
#define LUA_TLIGHTUSERDATA 2
#define LUA_TNUMBER 3
#define LUA_TSTRING 4
#define LUA_TTABLE 5
#define LUA_TFUNCTION 6
#define LUA_TUSERDATA 7
#define LUA_TTHREAD 8

#define LUA_MULTRET (-1)

#define LUA_OK 0
#define LUA_YIELD 1
#define LUA_ERRRUN 2
#define LUA_ERRSYNTAX 3
#define LUA_ERRMEM 4
#define LUA_ERRERR 5

#define LUA_REGISTRYINDEX (-10000)
#define LUA_ENVIRONINDEX (-10001)
#define LUA_GLOBALSINDEX (-10002)
#define lua_upvalueindex(i) (LUA_GLOBALSINDEX - (i))

#define LUA_NUMBER double

typedef LUA_NUMBER lua_Number;
typedef int (*lua_Writer)(lua_State, const void *, size_t, void *);

typedef int (*lua_CFunction)(lua_State L);

typedef struct luaL_Reg
{
    const char *name;
    lua_CFunction func;
} luaL_Reg;

#define LUA_IDSIZE 60

typedef struct lua_Debug
{
    int event;
    const char *name;           /* (n) */
    const char *namewhat;       /* (n) `global', `local', `field', `method' */
    const char *what;           /* (S) `Lua', `C', `main', `tail' */
    const char *source;         /* (S) */
    int currentline;            /* (l) */
    int nups;                   /* (u) number of upvalues */
    int linedefined;            /* (S) */
    int lastlinedefined;        /* (S) */
    char short_src[LUA_IDSIZE]; /* (S) */
    /* private part */
    int i_ci; /* active function */
} lua_Debug;

#define X(return_type, func_name, ...)             \
    extern void *func_name##_ptr;                  \
    typedef return_type (*func_name)(__VA_ARGS__); \
    extern return_type func_name##_wrap(__VA_ARGS__);

GLUA_FUNCTIONS

#undef X

void goLuaCallback(lua_State L, int *, char **);

extern const char *load_lua_shared(void);
extern const char *unload_lua_shared(void);

extern int lua_call_go(lua_State);
extern int luaCFunctionWrapper(void *, lua_State);
extern int lua_debug_getinfo_at(lua_State, int, const char *, lua_Debug *ar);
extern const char *lua_err_argmsg(lua_State, int, const char *);
extern const char *lua_type_error(lua_State, int, const char *);
extern const char *lua_tag_error(lua_State, int, int);

extern Result_Double lua_check_number(lua_State, int);
extern Result_String lua_check_string(lua_State, int, size_t *);
extern const char *lua_check_table(lua_State, int); // returns error instead of using Result
extern const char *lua_check_func(lua_State, int);  // returns error instead of using Result
extern Result_Bool lua_check_bool(lua_State, int);
