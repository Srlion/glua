#include "c/glua.c"
#include "c/think_queue.c"

// go only includes c files in the same directory as the go file

__attribute__((visibility("default"))) int gmod13_open(lua_State *L);
__attribute__((visibility("default"))) int gmod13_close(lua_State *L);
