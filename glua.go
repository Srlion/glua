package glua

// #include "c/glua.h"
import "C"
import (
	"errors"
	"fmt"
	"runtime/cgo"
	"strconv"
	"unsafe"
)

func LoadLuaShared() *string {
	err := C.load_lua_shared()
	if err != nil {
		errStr := C.GoString(err)
		return &errStr
	}

	return nil
}

func UnloadLuaShared() {
	C.unload_lua_shared()
}

func GetLuaSharedPath() string {
	return C.GoString(C.get_lua_shared_path())
}

// Creates a new Lua state.
//
// # Example
//
//	L := glua.NewState()
func NewState() State {
	return State(C.luaL_newstate_wrap())
}

// Creates a new coroutine
//
// # Example
//
//	L := glua.NewState()
//	L.NewCoroutine()
func (L State) NewCoroutine() State {
	return State(C.lua_newthread_wrap(L.c()))
}

// Returns the index of the top element in the stack. Because indices start at 1, this result is equal to the
// number of elements in the stack (and so 0 means an empty stack).
func (L State) GetTop() int {
	return int(C.lua_gettop_wrap(L.c()))
}

/*
Sets the stack top to the given index.

If the new top is larger than the old one, the new elements are filled with nil.

If index is 0, then all stack elements are removed.

# Example

	L.PushString("Hello, world!")

	fmt.Println("Stack size before SetTop:", L.GetTop())

	L.SetTop(0)

	fmt.Println("Stack size after SetTop:", L.GetTop())
*/
func (L State) SetTop(idx int) {
	C.lua_settop_wrap(L.c(), C.int(idx))
}

/*
Pushes a copy of the element at the given index onto the stack.

# Example

	// Push a string onto the stack
	L.PushString("Hello, world!")

	// Duplicate the value on top of the stack (the string)
	L.PushValue(-1)

	// The stack now has two copies of the string
	L.DumpStack()
*/
func (L State) PushValue(idx int) {
	C.lua_pushvalue_wrap(L.c(), C.int(idx))
}

/*
Removes the element at the given valid index, shifting down the elements above this index to fill the gap.

Cannot be called with a pseudo-index, because a pseudo-index is not an actual stack position.

# Example

	// Push a string onto the stack
	L.PushString("Hello, world!")

	// Push a number onto the stack
	glua.PushNumber(L, 123)

	// Remove the string
	L.Remove(-2)

	// The stack now has only the number
	L.DumpStack()
*/
func (L State) Remove(idx int) {
	C.lua_remove_wrap(L.c(), C.int(idx))
}

/*
Moves the top element into the given valid index, shifting up the elements above this index to open space.

Cannot be called with a pseudo-index, because a pseudo-index is not an actual stack position.

# Example

	L.PushString("first")
	L.PushString("second")
	L.PushString("third")

	// Stack is:
	// 3. third
	// 2. second
	// 1. first

	// Insert the top element ("third") at position 1
	L.Insert(1)

	// Stack is:
	// 3. second
	// 2. first
	// 1. third
*/
func (L State) Insert(idx int) {
	C.lua_insert_wrap(L.c(), C.int(idx))
}

/*
Moves the top element into the given position (and pops it), without shifting any
element (therefore replacing the value at the given position).

# Example

	L.PushString("original") // stack: 1. "original"
	L.PushString("new") 	// stack: 2. "new", 1. "original"

	// Replace the value at position 1 with the top element ("new")
	L.Replace(1)

	// Stack is:
	// 1. new
*/
func (L State) Replace(idx int) {
	C.lua_replace_wrap(L.c(), C.int(idx))
}

/*
Ensures that there are at least extra free stack slots in the stack.

It returns false if it cannot grow the stack to that size.

This function never shrinks the stack; if the stack is already larger than the new size, it is left unchanged.

# Example

	// Ensure there are at least 10 free stack slots
	if !L.CheckStack(10) {
		fmt.Println("Couldn't grow the stack")
	}
*/
func (L State) CheckStack(extra int) bool {
	return C.lua_checkstack_wrap(L.c(), C.int(extra)) != 0
}

/*
Returns the type of the value in the given index.

It returns LUA_TNONE for a non-valid index (That is, an index to an "empty" stack position).

The types returned by this function are:

  - LUA_TNIL
  - LUA_TBOOLEAN
  - LUA_TLIGHTUSERDATA
  - LUA_TNUMBER
  - LUA_TSTRING
  - LUA_TTABLE
  - LUA_TFUNCTION
  - LUA_TUSERDATA
  - LUA_TTHREAD

# Example

	L.PushString("Hello, world!")
	fmt.Println(L.TypeName(L.Type(-1)))
*/
func (L State) Type(idx int) int {
	return int(C.lua_type_wrap(L.c(), C.int(idx)))
}

/*
Returns the name of the type encoded by the value typeid.

# Example

	L.PushString("Hello, world!")
	fmt.Println(L.TypeName(L.Type(-1)))
*/
func (L State) TypeName(typeid int) string {
	return C.GoString(C.lua_typename_wrap(L.c(), C.int(typeid)))
}

/*
Returns true if the two values in acceptable indices index1 and index2 are equal,
following the semantics of the Lua == operator (that is, may call metamethods).
Otherwise returns false. Also returns false if any of the indices is non valid.

# Example

	L.PushString("Hello, world!")
	L.PushString("Hello, world!")

	if L.AreEqual(-1, -2) {
		fmt.Println("The strings are equal")
	}
*/
func (L State) AreEqual(idx1, idx2 int) bool {
	return C.lua_equal_wrap(L.c(), C.int(idx1), C.int(idx2)) != 0
}

/*
Returns true if the two values in acceptable indices index1 and index2 are primitively equal
(that is, without calling metamethods).
Otherwise returns false. Also returns false if any of the indices are non valid.

# Example

	L.PushString("Hello, world!")
	L.PushString("Hello, world!")

	if L.AreRawEqual(-1, -2) {
		fmt.Println("The strings are equal")
	}
*/
func (L State) AreRawEqual(idx1, idx2 int) bool {
	return C.lua_rawequal_wrap(L.c(), C.int(idx1), C.int(idx2)) != 0
}

/*
Returns true if the value at acceptable index index1 is smaller than the value at
acceptable index index2, following the semantics of the Lua < operator (that is, may call metamethods).
Otherwise returns false. Also returns false if any of the indices is non valid.

# Example

	glua.PushNumber(L, 123)
	glua.PushNumber(L, 456)

	if L.IsLessThan(-1, -2) {
		fmt.Println("123 is less than 456")
	}
*/
func (L State) IsLessThan(idx1, idx2 int) bool {
	return C.lua_lessthan_wrap(L.c(), C.int(idx1), C.int(idx2)) != 0
}

/*
Returns the number at the given index.

The value must be a number or a string convertible to a number, otherwise it returns 0.

# Example

	glua.PushNumber(L, 123.456)
	fmt.Println(L.GetNumber(-1))
*/
func (L State) GetNumber(idx int) LUA_NUMBER {
	return LUA_NUMBER(C.lua_tonumber_wrap(L.c(), C.int(idx)))
}

/*
Returns the boolean value at the given index.

If the value is a number or a string that is convertible to a number, it returns true for non-zero numbers.
It also returns 0 when called with a non-valid index.

If you want to accept only true or false, use IsBool.

# Example

	L.PushBool(true)
	fmt.Println(L.GetBool(-1))
*/
func (L State) GetBool(idx int) bool {
	return C.lua_toboolean_wrap(L.c(), C.int(idx)) != 0
}

/*
Returns the string at the given index in the Lua stack.

# Example

	L.PushString("Hello, world!")
	str := L.GetString(-1)
	if str != nil {
		fmt.Println(*str)
	}
*/
func (L State) GetString(idx int) *string {
	if !L.IsString(idx) {
		return nil
	}

	size := C.size_t(0)
	str := C.lua_tolstring_wrap(L.c(), C.int(idx), &size)
	if str == nil {
		return nil
	}

	result := goStringN(str, size)
	return &result
}

/*
Returns the binary string at the given index in the Lua stack.

# Example

	L.PushString("Hello, world!")
	str := L.GetBinaryString(-1)
	if str != nil {
		fmt.Println(string(str))
	}
*/
func (L State) GetBinaryString(idx int) []byte {
	if !L.IsString(idx) {
		return nil
	}

	size := C.size_t(0)
	str := C.lua_tolstring_wrap(L.c(), C.int(idx), &size)
	if str == nil {
		return nil
	}

	return goBytes(unsafe.Pointer(str), size)
}

func (L State) GetLength(idx int) int {
	return int(C.lua_objlen_wrap(L.c(), C.int(idx)))
}

func (L State) GetFunction(idx int) (func(State) int, error) {
	ptr := C.lua_tocfunction_wrap(L.c(), C.int(idx))
	if ptr == nil {
		return nil, errors.New("not a function")
	}

	return func(L State) int {
		return int(C.luaCFunctionWrapper(ptr, L.c()))
	}, nil
}

/*
Returns the userdata at the given index.

If the value at the given index is not a userdata, it returns nil.

If the value is a light userdata, it returns the pointer.
*/
func (L State) GetUserData(idx int, metatable *string) cgo.Handle {
	if !L.IsUserData(idx) {
		var metaMessage string
		if metatable != nil {
			metaMessage = " of type: " + *metatable
		}
		panic("expected a userdata" + metaMessage)
	}

	if metatable != nil {
		L.GetMetatable(idx)
		L.GetMetatableByName(*metatable)

		res := L.AreRawEqual(-1, -2)
		L.PopN(2)

		if !res {
			panic("expected a userdata of type: " + *metatable)
		}
	}

	ud := C.lua_touserdata_wrap(L.c(), C.int(idx))
	if ud == nil {
		panic("invalid userdata pointer")
	}

	handle := *(*cgo.Handle)(ud)

	return handle
}

func (L State) GetLightUserData(idx int) unsafe.Pointer {
	if !L.IsUserData(idx) {
		panic("expected a light userdata")
	}

	return C.lua_touserdata_wrap(L.c(), C.int(idx))
}

/*
Returns the thread at the given index.

If the value at the given index is not a thread, it returns nil.
*/
func (L State) GetThread(idx int) State {
	return State(C.lua_tothread_wrap(L.c(), C.int(idx)))
}

/*
Gets a pointer to the value at the given index.
*/
func (L State) GetPointer(idx int) unsafe.Pointer {
	return unsafe.Pointer(C.lua_topointer_wrap(L.c(), C.int(idx)))
}

/*
Pushes a nil value onto the stack.
*/
func (L State) PushNil() {
	C.lua_pushnil_wrap(L.c())
}

/*
Pushes a boolean value onto the stack.
*/
func (L State) PushBool(b bool) {
	if b {
		C.lua_pushboolean_wrap(L.c(), 1)
	} else {
		C.lua_pushboolean_wrap(L.c(), 0)
	}
}

/*
Pushes a string onto the stack.
*/
func (L State) PushString(str string) {
	if len(str) == 0 {
		C.lua_pushlstring_wrap(L.c(), nil, 0)
		return
	}

	strPtr := unsafe.Pointer(&[]byte(str)[0])
	C.lua_pushlstring_wrap(L.c(), (*C.char)(strPtr), C.size_t(len(str)))
}

/*
Pushes a string onto the stack with a given length.
*/
func (L State) PushBinaryString(data []byte) {
	if len(data) == 0 {
		C.lua_pushlstring_wrap(L.c(), nil, 0)
		return
	}

	strPtr := unsafe.Pointer(&data[0])
	C.lua_pushlstring_wrap(L.c(), (*C.char)(strPtr), C.size_t(len(data)))
}

// this is not same as lua_pushfstring, it just mimics the behavior
func (L State) PushFString(fmtstr string, args ...any) {
	if len(args) == 0 {
		L.PushString(fmtstr)
		return
	}

	L.PushString(fmt.Sprintf(fmtstr, args...))
}

/*
Pushes a light userdata onto the stack.

Light userdata is a pointer that is not managed by Lua.
*/
func (L State) PushLightUserData(p unsafe.Pointer) {
	C.lua_pushlightuserdata_wrap(L.c(), p)
}

/*
Pushes the current thread (i.e, the coroutine) onto the stack,
and returns whether the thread is the main thread or not.

Returns 1 if the thread is the main thread, otherwise 0.
*/
func (L State) PushThread() int {
	return int(C.lua_pushthread_wrap(L.c()))
}

/*
Pushes onto the stack the value t[k], where t is the table at the given index and k is the value at the top of the stack.

This function pops the key from the stack.

# Example

	L.NewTable()
		L.PushString("message")
		L.PushString("Hello, world!")
	L.SetTable(-3)

	L.PushString("message")
	L.GetTable(-2)
	fmt.Println(L.GetString(-1))
*/
func (L State) GetTable(idx int) {
	C.lua_gettable_wrap(L.c(), C.int(idx))
}

/*
Pushes onto the stack the value t[k], where t is the table at the given index and k is the value at the top of the stack.

# Example

	L.NewTable()
		L.PushString("message")
		L.PushString("Hello, world!")
	L.SetTable(-3)

	L.GetField(-1, "message")
	fmt.Println(L.GetString(-1))
*/
func (L State) GetField(idx int, key string) {
	cKey := CStr(key)
	defer cKey.free()

	C.lua_getfield_wrap(L.c(), C.int(idx), cKey.c)
}

/*
Pushes onto the stack the value of the global name.

# Example

	L.GetGlobal("print")
*/
func (L State) GetGlobal(name string) {
	L.GetField(LUA_GLOBALSINDEX, name)
}

/*
Similar to GetTable, but does not perform any metamethods.
*/
func (L State) RawGet(idx int) {
	C.lua_rawget_wrap(L.c(), C.int(idx))
}

/*
Pushes onto the stack the value t[n], where t is the table at the given index.

The access is raw, that is, it does not invoke metamethods.

# Example

	L.NewTable()
		L.PushString("Hello, world!")
	L.RawSetI(-2, 1)

	L.RawGetI(-1, 1)
	fmt.Println(L.GetString(-1))
*/
func (L State) RawGetI(idx int, n int) {
	C.lua_rawgeti_wrap(L.c(), C.int(idx), C.int(n))
}

func (L State) CreateTable(narr, nrec int) {
	C.lua_createtable_wrap(L.c(), C.int(narr), C.int(nrec))
}

func (L State) NewTable() {
	L.CreateTable(0, 0)
}

/*
Creates a new userdata and pushes it onto the stack.

It returns a cgo.Handle that can be used to retrieve the value.

You need to call handle.Delete() when __gc is called.

# Example

	type MyStruct struct {
		Message string
	}

	myStruct := &MyStruct{"Hello, world!"}

	h := L.NewUserData(myStruct, nil)
*/
func (L State) NewUserData(value any, metatable *string) cgo.Handle {
	const goUserDataSize = C.size_t(unsafe.Sizeof(uintptr(0)))

	h := cgo.NewHandle(value)

	ptr := C.lua_newuserdata_wrap(L.c(), goUserDataSize)

	if metatable != nil {
		L.GetMetatableByName(*metatable)
		if L.Type(-1) != LUA_TTABLE {
			panic("metatable not found")
		}
		L.SetMetatable(-2)
	}

	*(*cgo.Handle)(ptr) = h

	return h
}

/*
Pushes onto the stack the metatable of the value at the given index.

If the value does not have a metatable, the function returns 0 and pushes nothing.
*/
func (L State) GetMetatable(idx int) int {
	return int(C.lua_getmetatable_wrap(L.c(), C.int(idx)))
}

func (L State) GetMetatableByName(name string) {
	L.GetField(LUA_REGISTRYINDEX, name)
}

/*
Pushes onto the stack the environment table of the value at the given index.
*/
func (L State) GetFenv(idx int) {
	C.lua_getfenv_wrap(L.c(), C.int(idx))
}

/*
Does the equivalent of t[k] = v, where t is the table at the given index and v is the value at the top of the stack, and k is the value just below the top.

This function pops the key and the value from the stack.

# Example

	L.NewTable()
		L.PushString("message")
		L.PushString("Hello, world!")
	L.SetTable(-3)

	L.SetGlobal("myTable")
	L.RunString("print(myTable.message)")
*/
func (L State) SetTable(idx int) {
	C.lua_settable_wrap(L.c(), C.int(idx))
}

/*
Does the equivalent of t[k] = v, where t is the table at the given index and v is the value at the top of the stack.

This function pops the value from the stack.

# Example

	L.NewTable()

	L.PushString("Hello, world!")
	L.SetField(-2, "message")

	L.SetGlobal("myTable")

	L.RunString("print(myTable.message)")
*/
func (L State) SetField(idx int, key string) {
	cKey := CStr(key)
	defer cKey.free()

	C.lua_setfield_wrap(L.c(), C.int(idx), cKey.c)
}

/*
Similar to SetTable, but does not perform any metamethods.

# Example

	L.NewTable()
		L.PushString("message")
		L.PushString("Hello, world!")
	L.RawSet(-3)

	L.SetGlobal("myTable")
	L.RunString("print(myTable.message)")
*/
func (L State) RawSet(idx int) {
	C.lua_rawset_wrap(L.c(), C.int(idx))
}

/*
Does the equivalent of t[n] = v, where t is the table at the given index and v is the value at the top of the stack.

This function pops the value from the stack.

The assignment is raw, that is, it does not invoke metamethods.

# Example

	L.NewTable()
		L.PushString("Hello, world!")
	L.RawSetI(-2, 1)

	L.SetGlobal("myTable")
	L.RunString("print(myTable[1])")
*/
func (L State) RawSetI(idx int, n int) {
	C.lua_rawseti_wrap(L.c(), C.int(idx), C.int(n))
}

/*
Sets the metatable for the object at the given index.

# Example

	L.NewTable()

	L.NewTable()
	L.SetMetatable(-2)

	L.SetGlobal("myTable")

	L.RunString("print(getmetatable(myTable))")
*/
func (L State) SetMetatable(idx int) {
	C.lua_setmetatable_wrap(L.c(), C.int(idx))
}

/*
Pops a table from the stack and sets it as the new environment for the value at the given index.

If the value at the given index is neither a function nor a thread nor a userdata, it returns 0.
Otherwise, it returns 1.

# You cannot set the environment of a C function. It will return 1 but won't work.

# Example

	L.RunString("function myfunc() print(a) end")
	L.GetGlobal("myfunc")

	L.NewTable()

	L.PushString("a")
	glua.PushNumber(L, 123)
	L.SetTable(-3)

	L.PushString("print")
	L.GetGlobal("print")
	L.SetTable(-3)

	L.SetFEnv(-2)

	L.Call(0, 0)
*/
func (L State) SetFEnv(idx int) int {
	return int(C.lua_setfenv_wrap(L.c(), C.int(idx)))
}

/*
Calls a function.

nargs is the number of arguments in the stack.

nresults is the number of results to be returned.

# Example

	L.CompileString("print('Hello, world!')")
	L.Call(0, 0)
*/
func (L State) Call(nargs, nresults int) {
	C.lua_call_wrap(L.c(), C.int(nargs), C.int(nresults))
}

/*
Calls a function (which is on top of the Lua stack) in protected mode.

If there are no errors, PCall returns LUA_OK.

If errFunc is 0, the original error message is returned on the stack.

If errFunc is a valid stack index, it acts as an error handler function, and the error message is returned on top of the stack.

# Examples

1- no error handler

	err := L.CompileString("doesntexist()")
	if err != nil {
		fmt.Println(err)
		return 0
	}

	err = L.PCall(0, 0, 0)
	if err != nil {
		fmt.Println(err)
		return 0
	}

----

2- with error handler

	// we use PushOneTimeGoFunc because we only want to call this function once, it will be unregistered after it's called
	L.PushOneTimeGoFunc(func(L glua.State) int {
		fmt.Println(L.GetErrorString())
		return 0
	})
	errFuncIdx := L.GetTop()

	L.CompileString("doesntexist()")
	err := L.PCall(0, 0, errFuncIdx)
	if err != 0 {
		// error handler already printed the error
		return 0
	}
*/
func (L State) PCall(nargs, nresults, errfunc int) error {
	status := C.lua_pcall_wrap(L.c(), C.int(nargs), C.int(nresults), C.int(errfunc))
	if status != LUA_OK {
		return errors.New(L.GetErrorMessage(int(status)))
	}

	return nil
}

/*
Calls a function in protected mode.

If there are errors it returns false and prints the error message.

If there are no errors it returns true.

# Example

	err := L.RunString("doesntexist()")
	if err != nil {
		fmt.Println(err)
	}

	L.TryCall(0, 0)
*/
func (L State) TryCall(nargs, nresults int) bool {
	if err := L.PCall(nargs, nresults, 0); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (L State) CPCall(funcPtr, ud unsafe.Pointer) int {
	return int(C.lua_cpcall_wrap(L.c(), funcPtr, ud))
}

// TODO lua_yield
// TODO lua_resume
// TODO lua_status

/*
Opens all standard Lua libraries into the given Lua state.

# Example

	L.OpenLibs()
*/
func (L State) OpenLibs() {
	C.luaL_openlibs_wrap(L.c())
}

/*
Calls a metamethod.

If the object at index obj has a metatable with a field e,
this function calls it, passing the object as its argument.

It returns 1 and pushes the call's return value onto the stack.

If there is no metatable or field e, it returns 0 without
pushing any value.

# Example

	L.RunString(`
		myObject = {}
		mt = { __tostring = function() return 'Hello from __tostring!' end }
		setmetatable(myObject, mt)
	`)

	L.GetGlobal("myObject")

	if L.CallMeta(-1, "__tostring") == 1 {
		fmt.Println(L.GetString(-1))
	} else {
		fmt.Println("No metatable or field __tostring")
	}
*/
func (L State) CallMeta(objIdx int, e string) int {
	cEvent := CStr(e)
	defer cEvent.free()

	status := C.luaL_callmeta_wrap(L.c(), C.int(objIdx), cEvent.c)
	return int(status)
}

/*
If the registry already has the key tname, returns false. Otherwise,
creates a new table to be used as a metatable for userdata, adds it to the registry with key tname, and returns true.

In both cases pushes onto the stack the final value associated with tname in the registry.
*/
func (L State) NewMetaTable(name string) bool {
	cName := CStr(name)
	defer cName.free()

	return C.luaL_newmetatable_wrap(L.c(), cName.c) != 0
}

/*
Creates and returns a reference in the registry for the object at the top of the stack.

It pops the object from the stack.

# Example

	L.PushString("Hello, world!")
	ref := L.CreateRef()

	if L.FromRef(ref) {
		fmt.Println(L.GetString(-1))
	}
*/
func (L State) CreateRef() int {
	return int(C.luaL_ref_wrap(L.c(), LUA_REGISTRYINDEX))
}

/*
Gets the value associated with ref in the registry and pushes it onto the stack.

If the reference is invalid/nil, it returns false and does not push anything.

# Example

	L.PushString("Hello, world!")
	ref := L.CreateRef()

	if L.FromRef(ref) {
		fmt.Println(L.GetString(-1))
	}
*/
func (L State) FromRef(ref int) bool {
	if ref == LUA_REFNIL || ref == LUA_NOREF {
		return false
	}

	L.RawGetI(LUA_REGISTRYINDEX, ref)

	return true
}

/*
Deletes the reference ref from the registry.

If ref is LUA_REFNIL or LUA_NOREF, this function does nothing.

# Example

	L.PushString("Hello, world!")
	ref := L.CreateRef()

	if L.FromRef(ref) {
		fmt.Println(L.GetString(-1))
	}

	L.DeleteRef(ref)
*/
func (L State) DeleteRef(ref int) {
	if ref == LUA_REFNIL || ref == LUA_NOREF {
		return
	}

	C.luaL_unref_wrap(L.c(), LUA_REGISTRYINDEX, C.int(ref))
}

/*
An alias for DeleteRef.
*/
func (L State) Unref(ref int) {
	L.DeleteRef(ref)
}

// TODO luaL_findtable
// TODO lua_getstack
// TODO lua_getinfo

/*
Compiles a buffer into Lua code and pushes a function onto the stack that, when called, executes it.

It does not automatically run the function.

Returns LUA_OK if successful, otherwise an error code.

# Example

	buf := []byte("print('Hello, world!')")

	err := L.CompileBuffer(buf, 16, "example.lua")
	if err != nil {
		fmt.Println(err)
	}

	L.Call(0, 0)
*/
func (L State) CompileBuffer(code []byte, size uint, name string) int {
	cName, cBuf := CStr(name), CByt(code)
	defer cName.free()
	defer cBuf.Free()

	lua_error_code := C.luaL_loadbuffer_wrap(L.c(), cBuf.c, cBuf.size, cName.c)
	if lua_error_code != LUA_OK {
		return int(lua_error_code)
	}

	return LUA_OK
}

/*
Same as CompileBuffer, but with an additional mode parameter.

mode can be

  - "b" Treat the script as a binary chunk.

  - "t" Treat the script as a text chunk.

  - "bt" (default) Accepts both text and binary chunks.

# Example

	buf := []byte("print('Hello, world!')")
	err := L.LoadBufferX(buf, "example.lua", "t")
	if err != nil {
		fmt.Println(err)
	}
*/
func (L State) CompileBufferX(buf []byte, name, mode string) error {
	cName, cMode, cBuf := CStr(name), CStr(mode), CByt(buf)
	defer cName.free()
	defer cMode.free()
	defer cBuf.Free()

	status := C.luaL_loadbufferx_wrap(L.c(), cBuf.c, cBuf.size, cName.c, cMode.c)
	if status != LUA_OK {
		return errors.New(L.GetErrorMessage(int(status)))
	}

	return nil
}

/*
Compiles a string into Lua code and pushes a function onto the stack that, when called, executes it.

# It does not automatically run the function.

# Example

	err := L.CompileString("print('Hello, world!')")
	if err != nil {
		fmt.Println(err)
		return 0
	}

	L.Call(0, 0)
*/
func (L State) CompileString(s string) error {
	cS := CStr(s)
	defer cS.free()

	status := C.luaL_loadstring_wrap(L.c(), cS.c)
	if status != LUA_OK {
		return errors.New(L.GetErrorMessage(int(status)))
	}

	return nil
}

/*
Compiles a file into Lua code and pushes a function onto the stack that, when called, executes it.

It does not run the chunk.

# Example

	err := L.CompileFile("example.lua")
	if err != nil {
		fmt.Println(err)
	}
*/
func (L State) CompileFile(name string) error {
	cName := CStr(name)
	defer cName.free()

	status := C.luaL_loadfile_wrap(L.c(), cName.c)
	if status != 0 {
		return errors.New(L.GetErrorMessage(int(status)))
	}

	return nil
}

/*
Checks whether the value at the given index is a number.

If it is, it returns the number, otherwise it throws a Lua error.

# Example

	L.PushString("123")
	num := L.CheckNumber(-1) // will throw an error
*/
func (L State) CheckNumber(arg int) LUA_NUMBER {
	resNumber := C.lua_check_number(L.c(), C.int(arg))
	if resNumber.err != nil {
		panic(C.GoString(resNumber.err))
	}
	return LUA_NUMBER(resNumber.value)
}

/*
Checks whether the value at the given index is a string.

If it is, it returns the string, otherwise it throws a Lua error.

# Example

	L.PushString("Hello, world!")
	str := L.CheckString(-1)
*/
func (L State) CheckString(arg int) string {
	size := C.size_t(0)
	resString := C.lua_check_string(L.c(), C.int(arg), &size)
	if resString.err != nil {
		panic(C.GoString(resString.err))
	}
	return goStringN(resString.value, size)
}

/*
Checks whether the value at the given index is a binary string.

If it is, it returns the binary string, otherwise it throws a Lua error.

# Example

	L.PushString("Hello, world!")
	str := L.CheckBinaryString(-1)
*/
func (L State) CheckBinaryString(arg int) []byte {
	size := C.size_t(0)
	resString := C.lua_check_string(L.c(), C.int(arg), &size)
	if resString.err != nil {
		panic(C.GoString(resString.err))
	}
	return goBytes(unsafe.Pointer(resString.value), size)
}

/*
Checks whether the value at the given index is a boolean.

If it is, it returns the boolean, otherwise it throws a Lua error.

# Example

	L.PushBool(true)
	b := L.CheckBool(-1)
*/
func (L State) CheckBool(arg int) bool {
	resBool := C.lua_check_bool(L.c(), C.int(arg))
	if resBool.err != nil {
		panic(C.GoString(resBool.err))
	}
	return bool(resBool.value)
}

/*
Checks whether the value at the given index is a table.

It throws a Lua error if it is not.

# Example

	L.NewTable()
	L.CheckTable(-1)
*/
func (L State) CheckTable(arg int) {
	err := C.lua_check_table(L.c(), C.int(arg))
	if err != nil {
		panic(C.GoString(err))
	}
}

func (L State) CheckFunc(arg int) {
	err := C.lua_check_func(L.c(), C.int(arg))
	if err != nil {
		panic(C.GoString(err))
	}
}

func (L State) PopN(n int) {
	L.SetTop(-n - 1)
}

func (L State) Pop() {
	L.PopN(1)
}

func (L State) IsFunc(idx int) bool {
	return L.Type(idx) == LUA_TFUNCTION
}

func (L State) IsNil(idx int) bool {
	return L.Type(idx) == LUA_TNIL
}

func (L State) IsBool(idx int) bool {
	return L.Type(idx) == LUA_TBOOLEAN
}

func (L State) IsNumber(idx int) bool {
	return L.Type(idx) == LUA_TNUMBER
}

func (L State) IsString(idx int) bool {
	return L.Type(idx) == LUA_TSTRING
}

func (L State) IsUserData(idx int) bool {
	switch L.Type(idx) {
	case LUA_TUSERDATA, LUA_TLIGHTUSERDATA:
		return true
	}

	return false
}

func (L State) IsThread(idx int) bool {
	return L.Type(idx) == LUA_TTHREAD
}

func (L State) IsNone(idx int) bool {
	return L.Type(idx) == LUA_TNONE
}

func (L State) IsNoneOrNil(idx int) bool {
	return L.Type(idx) <= 0
}

func (L State) IsTable(idx int) bool {
	return L.Type(idx) == LUA_TTABLE
}

func (L State) SetGlobal(name string) {
	cName := CStr(name)
	defer cName.free()

	C.lua_setfield_wrap(L.c(), C.int(LUA_GLOBALSINDEX), cName.c)
}

/*
Runs a string as Lua code.

# Example

	err := L.RunString("print('Hello, world!')")
	if err != nil {
		fmt.Println(err)
		return 0
	}
*/
func (L State) RunString(str string) error {
	err := L.CompileString(str)
	if err != nil {
		return err
	}

	err = L.PCall(0, LUA_MULTRET, 0)
	if err != nil {
		return err
	}

	return nil
}

func (L State) GetErrorString() *string {
	return L.GetString(-1)
}

func (L State) GetCallingFileName() string {
	fileNameCStr := C.lua_get_calling_file_name(L.c())
	if fileNameCStr == nil {
		return ""
	}

	fileName := C.GoString(fileNameCStr)

	// free the C string
	C.free(unsafe.Pointer(fileNameCStr))

	return fileName
}

func (L State) GetErrorMessage(errorCode int) string {
	errorMessage := func(defaultMsg string) string {
		if err := L.GetErrorString(); err != nil {
			return defaultMsg + ": " + *err
		}
		return defaultMsg
	}

	switch errorCode {
	case LUA_ERRMEM:
		return "out of memory"
	case LUA_ERRERR:
		return "failed to run error handler"
	case LUA_ERRSYNTAX:
		return errorMessage("syntax error")
	case LUA_ERRRUN:
		return errorMessage("runtime error")
	case LUA_ERRFILE:
		return errorMessage("file error")
	default:
		return "unknown error code: " + strconv.Itoa(errorCode)
	}
}

func (L State) DumpStack() {
	top := L.GetTop()
	fmt.Printf("=== Stack size: %v ===\n", top)
	for i := 1; i <= top; i++ {
		t := L.Type(i)
		switch t {
		case LUA_TSTRING:
			fmt.Printf("(string) %v => %q\n", i, *L.GetString(i))
		case LUA_TNUMBER:
			fmt.Printf("(number) %v => %v\n", i, L.GetNumber(i))
		case LUA_TBOOLEAN:
			fmt.Printf("(bool) %v => %v\n", i, L.GetBool(i))
		case LUA_TTABLE:
			fmt.Printf("(table) %v\n", i)
		case LUA_TFUNCTION:
			fmt.Printf("(function) %v\n", i)
		default:
			fmt.Printf("%v: %v:\n", i, L.TypeName(i))
		}
	}
	println("===")
}

// func (L State) TypeError(arg int, tname string) string {
// 	C.luaL_typerror_wrap(L.c(), C.int(arg), C.CString(tname))
// }

// func (L State) ErrArgMsg(arg int, extraMsg string) {
// 	fname := "?"
// 	var namewhat string
//
// 	if
// }
