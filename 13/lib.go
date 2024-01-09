package gccjit

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/ebitengine/purego"
)

type (
	Timer    uint
	Object   uint
	Function struct{ Object }
	Location struct{ Object }
	Context  struct{ Object }
	Result   struct{ Object }
	Type     struct{ Object }
	Rvalue   struct{ Object }
	Lvalue   struct{ Rvalue }
	Param    struct{ Lvalue }
	Block    struct{ Object }
	Field    struct{ Object }
	Struct   struct{ Type }
)

type (
	TimerPtr    = *Timer
	ObjectPtr   = *Object
	FunctionPtr = *Function
	ParamPtr    = *Param
	LocationPtr = *Location
	ContextPtr  = *Context
	ResultPtr   = *Result
	TypePtr     = *Type
	RvaluePtr   = *Rvalue
	LvaluePtr   = *Lvalue
	BlockPtr    = *Block
	FieldPtr    = *Field
	StructPtr   = *Struct
)

type (
	StrOption    int
	BoolOption   int
	Types        int
	FunctionKind int
	OutputKind   int
	Comparison   int
	BinaryOp     int
	IntOption    int
	GlobalKind   int
)

const (
	GCC_JIT_STR_OPTION_PROGNAME StrOption = iota
	GCC_JIT_NUM_STR_OPTIONS
)

const (
	BOOL_OPTION_DEBUGINFO BoolOption = iota
	BOOL_OPTION_DUMP_INITIAL_TREE
	BOOL_OPTION_DUMP_INITIAL_GIMPLE
	BOOL_OPTION_DUMP_GENERATED_CODE
	BOOL_OPTION_DUMP_SUMMARY
	BOOL_OPTION_DUMP_EVERYTHING
	BOOL_OPTION_SELFCHECK_GC
	BOOL_OPTION_KEEP_INTERMEDIATES
	NUM_BOOL_OPTIONS
)

const (
	FUNCTION_EXPORTED FunctionKind = iota
	FUNCTION_INTERNAL
	FUNCTION_IMPORTED
	FUNCTION_ALWAYS_INLINE
)

const (
	TYPE_VOID Types = iota
	TYPE_VOID_PTR
	TYPE_BOOL
	TYPE_CHAR
	TYPE_SIGNED_CHAR
	TYPE_UNSIGNED_CHAR
	TYPE_SHORT
	TYPE_UNSIGNED_SHORT
	TYPE_INT
	TYPE_UNSIGNED_INT
	TYPE_LONG
	TYPE_UNSIGNED_LONG
	TYPE_LONG_LONG
	TYPE_UNSIGNED_LONG_LONG
	TYPE_FLOAT
	TYPE_DOUBLE
	TYPE_LONG_DOUBLE
	TYPE_CONST_CHAR_PTR
	TYPE_SIZE_T
	TYPE_FILE_PTR
	TYPE_COMPLEX_FLOAT
	TYPE_COMPLEX_DOUBLE
	TYPE_COMPLEX_LONG_DOUBLE
	TYPE_UINT8_T
	TYPE_UINT16_T
	TYPE_UINT32_T
	TYPE_UINT64_T
	TYPE_UINT128_T
	TYPE_INT8_T
	TYPE_INT16_T
	TYPE_INT32_T
	TYPE_INT64_T
	TYPE_INT128_T
)

const (
	OUTPUT_KIND_ASSEMBLER OutputKind = iota
	OUTPUT_KIND_OBJECT_FILE
	OUTPUT_KIND_DYNAMIC_LIBRARY
	OUTPUT_KIND_EXECUTABLE
)

const (
	COMPARISON_EQ Comparison = iota
	COMPARISON_NE
	COMPARISON_LT
	COMPARISON_LE
	COMPARISON_GT
	COMPARISON_GE
)

const (
	BINARY_OP_PLUS BinaryOp = iota
	BINARY_OP_MINUS
	BINARY_OP_MULT
	BINARY_OP_DIVIDE
	BINARY_OP_MODULO
	BINARY_OP_BITWISE_AND
	BINARY_OP_BITWISE_XOR
	BINARY_OP_BITWISE_OR
	BINARY_OP_LOGICAL_AND
	BINARY_OP_LOGICAL_OR
	BINARY_OP_LSHIFT
	BINARY_OP_RSHIFT
)

const (
	INT_OPTION_OPTIMIZATION_LEVEL IntOption = iota
	NUM_INT_OPTIONS
)

const (
	GLOBAL_EXPORTED GlobalKind = iota
	GLOBAL_INTERNAL
	GLOBAL_IMPORTED
)

var (
	contextAcquire                       func() *Context
	contextRelease                       func(ctx *Context)
	contextSetStrOption                  func(ctx *Context, opt StrOption, value string)
	contextSetBoolOption                 func(ctx *Context, opt BoolOption, value bool)
	contextCompile                       func(ctx *Context) *Result
	resultRelease                        func(result *Result)
	contextGetType                       func(ctx *Context, typ Types) *Type
	contextNewParam                      func(ctx *Context, loc *Location, typ *Type, name string) *Param
	contextNewFunction                   func(ctx *Context, loc *Location, kind FunctionKind, return_type *Type, name string, numParams int, params []*Param, isVariadic bool) *Function
	contextNewStringLiteral              func(ctx *Context, value string) *Rvalue
	functionNewBlock                     func(fn *Function, name string) *Block
	blockAddEval                         func(block *Block, loc *Location, rvalue *Rvalue)
	contextNewCall                       func(ctx *Context, loc *Location, fn *Function, numargs int, args []*Rvalue) *Rvalue
	blockEndWithVoidReturn               func(block *Block, loc *Location)
	resultGetCode                        func(result *Result, name string) uintptr
	paramAsRvalue                        func(param *Param) *Rvalue
	contextCompileToFile                 func(ctx *Context, outputKind OutputKind, outputPath string)
	contextNewArrayAccess                func(ctx *Context, loc *Location, ptr *Rvalue, idx *Rvalue) *Lvalue
	lvalueAsRvalue                       func(lvalue *Lvalue) *Rvalue
	contextNewComparison                 func(ctx *Context, loc *Location, op Comparison, lhs *Rvalue, rhs *Rvalue) *Rvalue
	contextNewLocation                   func(ctx *Context, filename string, line, column int) *Location
	blockAddComment                      func(block *Block, loc *Location, text string)
	blockAddAssignmentOp                 func(block *Block, loc *Location, lvalue *Lvalue, op BinaryOp, rvalue *Rvalue)
	contextNewCast                       func(ctx *Context, loc *Location, rvalue *Rvalue, typ *Type) *Rvalue
	blockAddAssignment                   func(block *Block, loc *Location, lvalue *Lvalue, rvalue *Rvalue)
	blockEndWithJump                     func(block *Block, loc *Location, target *Block)
	blockEndWithConditional              func(block *Block, loc *Location, boolval *Rvalue, onTrue *Block, onFalse *Block)
	contextSetIntOption                  func(ctx *Context, opt IntOption, value int)
	contextNewArrayType                  func(ctx *Context, loc *Location, elementType *Type, numElements int) *Type
	typeGetPointer                       func(typ *Type) *Type
	contextZero                          func(ctx *Context, typ *Type) *Rvalue
	contextOne                           func(ctx *Context, typ *Type) *Rvalue
	contextNewGlobal                     func(ctx *Context, loc *Location, kind GlobalKind, typ *Type, name string) *Lvalue
	functionNewLocal                     func(fn *Function, loc *Location, typ *Type, name string) *Lvalue
	blockEndWithReturn                   func(block *Block, loc *Location, rvalue *Rvalue)
	contextGetFirstError                 func(ctx *Context) string
	contextGetLastError                  func(ctx *Context) string
	contextDumpToFile                    func(ctx *Context, path string, updateLocations bool)
	contextDumpReproducerToFile          func(ctx *Context, path string)
	contextSetBoolAllowUnreachableBlocks func(ctx *Context, value bool)
	contextSetBoolPrintErrorsToStderr    func(ctx *Context, value bool)
	contextSetBoolUseExternalDriver      func(ctx *Context, value bool)
	contextAddCommandLineOption          func(ctx *Context, optname string)
	contextNewField                      func(ctx *Context, loc *Location, typ *Type, name string) *Field
	contextNewStructType                 func(ctx *Context, loc *Location, name string, numFields int, fields []*Field) *Struct
	rvalueDereferenceField               func(ptr *Rvalue, loc *Location, field *Field) *Lvalue
	structAsType                         func(structType *Struct) *Type
	contextNewRvalueFromInt              func(ctx *Context, typ *Type, value int) *Rvalue
	contextNewRvalueFromLong             func(ctx *Context, typ *Type, value int64) *Rvalue
	contextNewRvalueFromPtr              func(ctx *Context, typ *Type, value uintptr) *Rvalue
	contextNewFunctionPtrType            func(ctx *Context, loc *Location, returnType *Type, numParams int, paramTypes []*Type, isVariadic bool) *Type
	contextNewCallThroughPtr             func(ctx *Context, loc *Location, fnPtr *Rvalue, numArgs int, args []*Rvalue) *Rvalue
	lvalueAccessField                    func(structOrUnion *Lvalue, loc *Location, field *Field) *Lvalue
	contextNewBitCast                    func(ctx *Context, loc *Location, rvalue *Rvalue, typ *Type) *Rvalue
	lvalueGetAddress                     func(lvalue *Lvalue, loc *Location) *Rvalue
	rvalueDereference                    func(rvalue *Rvalue, loc *Location) *Lvalue
	typeIsBool                           func(typ *Type) bool
	typeIsPointer                        func(typ *Type) bool
	typeIsIntegral                       func(typ *Type) bool
	typeIsStruct                         func(typ *Type) bool
	typeUnqualified                      func(typ *Type) *Type
	typeGetConst                         func(typ *Type) *Type
	typeGetVolatile                      func(typ *Type) *Type
	typeGetSize                          func(typ *Type) uint64
	objectGetContext                     func(obj *Object) *Context
	objectGetDebugString                 func(obj *Object) string
	timerNew                             func() *Timer
	timerRelease                         func(t *Timer)
	timerPush                            func(t *Timer, itemName string)
	timerPop                             func(t *Timer, itemName string)
	contextSetTimer                      func(ctx *Context, t *Timer)
	contextGetTimer                      func(ctx *Context) *Timer
	versionMajor                         func() int
	versionMinor                         func() int
	versionPatchLevel                    func() int
)

func getLibrary() string {
	switch runtime.GOOS {
	case "linux":
		return "libgccjit.so.0"
	case "darwin":
		// FIXME: Replace hardcoded path with something else
		return "/opt/homebrew/lib/gcc/current/libgccjit.0.dylib"
	case "windows":
		return "libgccjit-0.dll"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

func init() {
	lib, err := loadLibrary(getLibrary())
	if err != nil {
		panic(err)
	}

	purego.RegisterLibFunc(&contextAcquire, lib, "gcc_jit_context_acquire")
	purego.RegisterLibFunc(&contextRelease, lib, "gcc_jit_context_release")
	purego.RegisterLibFunc(&contextSetBoolOption, lib, "gcc_jit_context_set_bool_option")
	purego.RegisterLibFunc(&contextCompile, lib, "gcc_jit_context_compile")
	purego.RegisterLibFunc(&resultRelease, lib, "gcc_jit_result_release")
	purego.RegisterLibFunc(&contextGetType, lib, "gcc_jit_context_get_type")
	purego.RegisterLibFunc(&contextNewParam, lib, "gcc_jit_context_new_param")
	purego.RegisterLibFunc(&contextNewFunction, lib, "gcc_jit_context_new_function")
	purego.RegisterLibFunc(&contextNewStringLiteral, lib, "gcc_jit_context_new_string_literal")
	purego.RegisterLibFunc(&functionNewBlock, lib, "gcc_jit_function_new_block")
	purego.RegisterLibFunc(&blockAddEval, lib, "gcc_jit_block_add_eval")
	purego.RegisterLibFunc(&contextNewCall, lib, "gcc_jit_context_new_call")
	purego.RegisterLibFunc(&blockEndWithVoidReturn, lib, "gcc_jit_block_end_with_void_return")
	purego.RegisterLibFunc(&resultGetCode, lib, "gcc_jit_result_get_code")
	purego.RegisterLibFunc(&paramAsRvalue, lib, "gcc_jit_param_as_rvalue")
	purego.RegisterLibFunc(&contextCompileToFile, lib, "gcc_jit_context_compile_to_file")
	purego.RegisterLibFunc(&contextNewArrayAccess, lib, "gcc_jit_context_new_array_access")
	purego.RegisterLibFunc(&lvalueAsRvalue, lib, "gcc_jit_lvalue_as_rvalue")
	purego.RegisterLibFunc(&contextNewComparison, lib, "gcc_jit_context_new_comparison")
	purego.RegisterLibFunc(&contextNewLocation, lib, "gcc_jit_context_new_location")
	purego.RegisterLibFunc(&blockAddComment, lib, "gcc_jit_block_add_comment")
	purego.RegisterLibFunc(&blockAddAssignmentOp, lib, "gcc_jit_block_add_assignment_op")
	purego.RegisterLibFunc(&contextNewCast, lib, "gcc_jit_context_new_cast")
	purego.RegisterLibFunc(&blockAddAssignment, lib, "gcc_jit_block_add_assignment")
	purego.RegisterLibFunc(&blockEndWithJump, lib, "gcc_jit_block_end_with_jump")
	purego.RegisterLibFunc(&blockEndWithConditional, lib, "gcc_jit_block_end_with_conditional")
	purego.RegisterLibFunc(&contextSetIntOption, lib, "gcc_jit_context_set_int_option")
	purego.RegisterLibFunc(&contextNewArrayType, lib, "gcc_jit_context_new_array_type")
	purego.RegisterLibFunc(&typeGetPointer, lib, "gcc_jit_type_get_pointer")
	purego.RegisterLibFunc(&contextZero, lib, "gcc_jit_context_zero")
	purego.RegisterLibFunc(&contextOne, lib, "gcc_jit_context_one")
	purego.RegisterLibFunc(&contextNewGlobal, lib, "gcc_jit_context_new_global")
	purego.RegisterLibFunc(&functionNewLocal, lib, "gcc_jit_function_new_local")
	purego.RegisterLibFunc(&blockEndWithReturn, lib, "gcc_jit_block_end_with_return")
	purego.RegisterLibFunc(&contextGetFirstError, lib, "gcc_jit_context_get_first_error")
	purego.RegisterLibFunc(&contextGetLastError, lib, "gcc_jit_context_get_last_error")
	purego.RegisterLibFunc(&contextDumpToFile, lib, "gcc_jit_context_dump_to_file")
	purego.RegisterLibFunc(&contextDumpReproducerToFile, lib, "gcc_jit_context_dump_reproducer_to_file")
	purego.RegisterLibFunc(&contextSetStrOption, lib, "gcc_jit_context_set_str_option")
	purego.RegisterLibFunc(&contextSetBoolAllowUnreachableBlocks, lib, "gcc_jit_context_set_bool_allow_unreachable_blocks")
	purego.RegisterLibFunc(&contextSetBoolPrintErrorsToStderr, lib, "gcc_jit_context_set_bool_print_errors_to_stderr")
	purego.RegisterLibFunc(&contextSetBoolUseExternalDriver, lib, "gcc_jit_context_set_bool_use_external_driver")
	purego.RegisterLibFunc(&contextAddCommandLineOption, lib, "gcc_jit_context_add_command_line_option")
	purego.RegisterLibFunc(&contextNewField, lib, "gcc_jit_context_new_field")
	purego.RegisterLibFunc(&contextNewStructType, lib, "gcc_jit_context_new_struct_type")
	purego.RegisterLibFunc(&rvalueDereferenceField, lib, "gcc_jit_rvalue_dereference_field")
	purego.RegisterLibFunc(&structAsType, lib, "gcc_jit_struct_as_type")
	purego.RegisterLibFunc(&contextNewRvalueFromInt, lib, "gcc_jit_context_new_rvalue_from_int")
	purego.RegisterLibFunc(&contextNewRvalueFromLong, lib, "gcc_jit_context_new_rvalue_from_long")
	purego.RegisterLibFunc(&contextNewRvalueFromPtr, lib, "gcc_jit_context_new_rvalue_from_ptr")
	purego.RegisterLibFunc(&contextNewFunctionPtrType, lib, "gcc_jit_context_new_function_ptr_type")
	purego.RegisterLibFunc(&contextNewCallThroughPtr, lib, "gcc_jit_context_new_call_through_ptr")
	purego.RegisterLibFunc(&lvalueAccessField, lib, "gcc_jit_lvalue_access_field")
	purego.RegisterLibFunc(&contextNewBitCast, lib, "gcc_jit_context_new_bitcast")
	purego.RegisterLibFunc(&lvalueGetAddress, lib, "gcc_jit_lvalue_get_address")
	purego.RegisterLibFunc(&rvalueDereference, lib, "gcc_jit_rvalue_dereference")
	purego.RegisterLibFunc(&typeIsBool, lib, "gcc_jit_type_is_bool")
	purego.RegisterLibFunc(&typeIsPointer, lib, "gcc_jit_type_is_pointer")
	purego.RegisterLibFunc(&typeIsIntegral, lib, "gcc_jit_type_is_integral")
	purego.RegisterLibFunc(&typeIsStruct, lib, "gcc_jit_type_is_struct")
	purego.RegisterLibFunc(&typeUnqualified, lib, "gcc_jit_type_unqualified")
	purego.RegisterLibFunc(&typeGetConst, lib, "gcc_jit_type_get_const")
	purego.RegisterLibFunc(&typeGetVolatile, lib, "gcc_jit_type_get_volatile")
	purego.RegisterLibFunc(&typeGetSize, lib, "gcc_jit_type_get_size")
	purego.RegisterLibFunc(&objectGetContext, lib, "gcc_jit_object_get_context")
	purego.RegisterLibFunc(&objectGetDebugString, lib, "gcc_jit_object_get_debug_string")
	purego.RegisterLibFunc(&timerNew, lib, "gcc_jit_timer_new")
	purego.RegisterLibFunc(&timerRelease, lib, "gcc_jit_timer_release")
	purego.RegisterLibFunc(&timerPush, lib, "gcc_jit_timer_push")
	purego.RegisterLibFunc(&timerPop, lib, "gcc_jit_timer_pop")
	purego.RegisterLibFunc(&contextSetTimer, lib, "gcc_jit_context_set_timer")
	purego.RegisterLibFunc(&contextGetTimer, lib, "gcc_jit_context_get_timer")
	purego.RegisterLibFunc(&versionMajor, lib, "gcc_jit_version_major")
	purego.RegisterLibFunc(&versionMinor, lib, "gcc_jit_version_minor")
	purego.RegisterLibFunc(&versionPatchLevel, lib, "gcc_jit_version_patchlevel")
}

func VersionMajor() int {
	return versionMajor()
}

func VersionMinor() int {
	return versionMinor()
}

func VersionPatchLevel() int {
	return versionPatchLevel()
}

func TimerNew() *Timer {
	return timerNew()
}

func (t *Timer) Release() {
	timerRelease(t)
}

func (t *Timer) Push(name string) {
	timerPush(t, name)
}

func (t *Timer) Pop(name string) {
	timerPop(t, name)
}

func ContextAcquire() *Context {
	return contextAcquire()
}

func (o *Object) GetContext() *Context {
	return objectGetContext(o)
}

func (o *Object) GetDebugString() string {
	return objectGetDebugString(o)
}

func (c *Context) SetTimer(t *Timer) {
	contextSetTimer(c, t)
}

func (c *Context) GetTimer() *Timer {
	return contextGetTimer(c)
}

func (c *Context) SetBoolOption(opt BoolOption, value bool) {
	contextSetBoolOption(c, opt, value)
}

func (c *Context) SetIntOption(opt IntOption, value int) {
	contextSetIntOption(c, opt, value)
}

func (c *Context) SetStrOption(opt StrOption, value string) {
	contextSetStrOption(c, opt, value)
}

func (c *Context) SetBoolAllowUnreachableBlocks(value bool) {
	contextSetBoolAllowUnreachableBlocks(c, value)
}

func (c *Context) SetBoolPrintErrorsToStderr(value bool) {
	contextSetBoolPrintErrorsToStderr(c, value)
}

func (c *Context) SetBoolUseExternalDriver(value bool) {
	contextSetBoolUseExternalDriver(c, value)
}

func (c *Context) AddCommandLineOption(optname string) {
	contextAddCommandLineOption(c, optname)
}

func (c *Context) GetType(typ Types) *Type {
	return contextGetType(c, typ)
}

func (c *Context) GetArrayType(loc *Location, elementType *Type, numElements int) *Type {
	return contextNewArrayType(c, loc, elementType, numElements)
}

func (c *Context) NewFunctionPtrType(loc *Location, returnType *Type, paramTypes []*Type, isVariadic bool) *Type {
	return contextNewFunctionPtrType(c, loc, returnType, len(paramTypes), paramTypes, isVariadic)
}

func (c *Context) NewStructType(loc *Location, name string, fields []*Field) *Struct {
	return contextNewStructType(c, loc, name, len(fields), fields)
}

func (c *Context) NewFunction(loc *Location, kind FunctionKind, return_type *Type, name string, params []*Param, isVariadic bool) *Function {
	return contextNewFunction(c, loc, kind, return_type, name, len(params), params, isVariadic)
}

func (c *Context) NewParam(loc *Location, typ *Type, name string) *Param {
	return contextNewParam(c, loc, typ, name)
}

func (c *Context) NewBlock(fn *Function, name string) *Block {
	return functionNewBlock(fn, name)
}

func (c *Context) NewCall(loc *Location, fn *Function, args []*Rvalue) *Rvalue {
	return contextNewCall(c, loc, fn, len(args), args)
}

func (c *Context) NewCallThroughPtr(loc *Location, ptr *Rvalue, args []*Rvalue) *Rvalue {
	return contextNewCallThroughPtr(c, loc, ptr, len(args), args)
}

func (c *Context) NewStringLiteral(value string) *Rvalue {
	return contextNewStringLiteral(c, value)
}

func (c *Context) NewArrayAccess(loc *Location, ptr *Rvalue, idx *Rvalue) *Lvalue {
	return contextNewArrayAccess(c, loc, ptr, idx)
}

func (c *Context) NewNewComparison(loc *Location, op Comparison, lhs *Rvalue, rhs *Rvalue) *Rvalue {
	return contextNewComparison(c, loc, op, lhs, rhs)
}

func (c *Context) NewLocation(filename string, line, column int) *Location {
	return contextNewLocation(c, filename, line, column)
}

func (c *Context) NewCast(loc *Location, rvalue *Rvalue, typ *Type) *Rvalue {
	return contextNewCast(c, loc, rvalue, typ)
}

func (c *Context) NewBitCast(loc *Location, rvalue *Rvalue, typ *Type) *Rvalue {
	return contextNewBitCast(c, loc, rvalue, typ)
}

func (c *Context) NewGlobal(loc *Location, kind GlobalKind, typ *Type, name string) *Lvalue {
	return contextNewGlobal(c, loc, kind, typ, name)
}

func (c *Context) NewRValueFromInt(typ *Type, value int) *Rvalue {
	return contextNewRvalueFromInt(c, typ, value)
}

func (c *Context) NewRValueFromLong(typ *Type, value int64) *Rvalue {
	return contextNewRvalueFromLong(c, typ, value)
}

func (c *Context) NewRvalueFromPtr(typ *Type, value uintptr) *Rvalue {
	return contextNewRvalueFromPtr(c, typ, value)
}

func (c *Context) NewField(loc *Location, typ *Type, name string) *Field {
	return contextNewField(c, loc, typ, name)
}

func (c *Context) Zero(typ *Type) *Rvalue {
	return contextZero(c, typ)
}

func (c *Context) One(typ *Type) *Rvalue {
	return contextOne(c, typ)
}

func (c *Context) DumpToFile(path string, updateLocations bool) {
	contextDumpToFile(c, path, updateLocations)
}

func (c *Context) DumpReproducerToFile(path string) {
	contextDumpReproducerToFile(c, path)
}

func (c *Context) Compile() (*Result, error) {
	res := contextCompile(c)
	if res == nil {
		return nil, errors.New(c.GetLastError())
	}

	return res, nil
}

func (c *Context) GetFirstError() string {
	return contextGetFirstError(c)
}

func (c *Context) GetLastError() string {
	return contextGetLastError(c)
}

func (c *Context) CompileToFile(outputKind OutputKind, outputPath string) {
	contextCompileToFile(c, outputKind, outputPath)
}

func (c *Context) Release() {
	contextRelease(c)
}

func (p *Param) AsRvalue() *Rvalue {
	return paramAsRvalue(p)
}

func (b *Block) AddEval(loc *Location, rvalue *Rvalue) {
	blockAddEval(b, loc, rvalue)
}

func (b *Block) EndWithVoidReturn(loc *Location) {
	blockEndWithVoidReturn(b, loc)
}

func (b *Block) AddComment(loc *Location, text string) {
	blockAddComment(b, loc, text)
}

func (b *Block) AddAssignmentOp(loc *Location, lvalue *Lvalue, op BinaryOp, rvalue *Rvalue) {
	blockAddAssignmentOp(b, loc, lvalue, op, rvalue)
}

func (b *Block) AddAssignment(loc *Location, lvalue *Lvalue, rvalue *Rvalue) {
	blockAddAssignment(b, loc, lvalue, rvalue)
}

func (b *Block) EndWithJump(loc *Location, target *Block) {
	blockEndWithJump(b, loc, target)
}

func (b *Block) EndWithConditional(loc *Location, boolval *Rvalue, onTrue *Block, on_false *Block) {
	blockEndWithConditional(b, loc, boolval, onTrue, on_false)
}

func (b *Block) EndWithReturn(loc *Location, rvalue *Rvalue) {
	blockEndWithReturn(b, loc, rvalue)
}

func (r *Result) GetCode(name string) uintptr {
	return resultGetCode(r, name)
}

func (r *Result) RegisterFunc(name string, fn any) {
	ptr := r.GetCode(name)
	purego.RegisterFunc(fn, ptr)
}

func (r *Result) Release() {
	resultRelease(r)
}

func (l *Lvalue) GetAddress(loc *Location) *Rvalue {
	return lvalueGetAddress(l, loc)
}

func (l *Lvalue) AsRvalue() *Rvalue {
	return lvalueAsRvalue(l)
}

func (l *Lvalue) AccessField(loc *Location, field *Field) *Lvalue {
	return lvalueAccessField(l, loc, field)
}

func (r *Rvalue) DereferenceField(loc *Location, field *Field) *Lvalue {
	return rvalueDereferenceField(r, loc, field)
}

func (r *Rvalue) Dereference(loc *Location) *Lvalue {
	return rvalueDereference(r, loc)
}

func (f *Function) NewBlock(name string) *Block {
	return functionNewBlock(f, name)
}

func (f *Function) NewLocal(loc *Location, typ *Type, name string) *Lvalue {
	return functionNewLocal(f, loc, typ, name)
}

func (t *Type) GetPointer() *Type {
	return typeGetPointer(t)
}

func (t *Type) GetConst() *Type {
	return typeGetConst(t)
}

func (t *Type) GetVolatile() *Type {
	return typeGetVolatile(t)
}

func (t *Type) GetSize() uint64 {
	return typeGetSize(t)
}

func (t *Type) IsBool() bool {
	return typeIsBool(t)
}

func (t *Type) IsPointer() bool {
	return typeIsPointer(t)
}

func (t *Type) IsIntegral() bool {
	return typeIsIntegral(t)
}

func (t *Type) IsStruct() bool {
	return typeIsStruct(t)
}

func (t *Type) Unqualified() *Type {
	return typeUnqualified(t)
}

func (t *Struct) AsType() *Type {
	return structAsType(t)
}
