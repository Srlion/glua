#ifdef _WIN32
#include <windows.h>
#define LIB_HANDLE HMODULE
#define LOAD_LIBRARY(lib) LoadLibrary(lib)
#define GET_FUNCTION(lib, func) GetProcAddress(lib, func)
#define CLOSE_LIBRARY(lib) FreeLibrary(lib)
#define GET_LOAD_ERROR() GetLastError()
#else
#include <dlfcn.h>
#define LIB_HANDLE void *
#define LOAD_LIBRARY(lib) dlopen(lib, RTLD_LAZY)
#define GET_FUNCTION(lib, func) dlsym(lib, func)
#define CLOSE_LIBRARY(lib) dlclose(lib)
#define GET_LOAD_ERROR() dlerror()
#endif
