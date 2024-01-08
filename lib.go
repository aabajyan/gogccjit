package main

import "github.com/ebitengine/purego"

type Function uint
type Param uint
type Location uint
type Context uint
type Result uint
type Type uint
type Rvalue uint
type Lvalue uint
type Block uint

type FunctionPtr = *Function
type ParamPtr = *Param
type LocationPtr = *Location
type ContextPtr = *Context
type ResultPtr = *Result
type TypePtr = *Type
type RvaluePtr = *Rvalue
type LvaluePtr = *Lvalue
type BlockPtr = *Block

type BoolOption int
type Types int
type FunctionKind int
type OutputKind int
type Comparison int
type BinaryOp int
type IntOption int
type GlobalKind int

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

var contextAcquire func() *Context
var contextRelease func(ctx *Context)
var contextSetBoolOption func(ctx *Context, opt BoolOption, value bool) uintptr
var contextCompile func(ctx *Context) *Result
var resultRelease func(result *Result)
var contextGetType func(ctx *Context, type_ Types) *Type
var contextNewParam func(ctx *Context, loc *Location, type_ *Type, name string) *Param
var contextNewFunction func(ctx *Context, loc *Location, kind FunctionKind, return_type *Type, name string, num_params int, params []*Param, is_variadic bool) *Function
var contextNewStringLiteral func(ctx *Context, value string) *Rvalue
var functionNewBlock func(fn *Function, name string) *Block
var blockAddEval func(block *Block, loc *Location, rvalue *Rvalue)
var contextNewCall func(ctx *Context, loc *Location, fn *Function, numargs int, args []*Rvalue) *Rvalue
var blockEndWithVoidReturn func(block *Block, loc *Location)
var resultGetCode func(result *Result, name string) uintptr
var paramAsRvalue func(param *Param) *Rvalue
var contextCompileToFile func(ctx *Context, output_kind OutputKind, output_path string)
var contextNewArrayAccess func(ctx *Context, loc *Location, ptr *Rvalue, idx *Rvalue) *Lvalue
var lvalueAsRvalue func(lvalue *Lvalue) *Rvalue
var contextNewComparison func(ctx *Context, loc *Location, op Comparison, lhs *Rvalue, rhs *Rvalue) *Rvalue
var contextNewLocation func(ctx *Context, filename string, line, column int) *Location
var blockAddComment func(block *Block, loc *Location, text string)
var blockAddAssignmentOp func(block *Block, loc *Location, lvalue *Lvalue, op BinaryOp, rvalue *Rvalue)
var contextNewCast func(ctx *Context, loc *Location, rvalue *Rvalue, type_ *Type) *Rvalue
var blockAddAssignment func(block *Block, loc *Location, lvalue *Lvalue, rvalue *Rvalue)
var blockEndWithJump func(block *Block, loc *Location, target *Block)
var blockEndWithConditional func(block *Block, loc *Location, boolval *Rvalue, on_true *Block, on_false *Block)
var contextSetIntOption func(ctx *Context, opt IntOption, value int)
var contextNewArrayType func(ctx *Context, loc *Location, element_type *Type, num_elements int) *Type
var typeGetPointer func(type_ *Type) *Type
var contextZero func(ctx *Context, type_ *Type) *Rvalue
var contextOne func(ctx *Context, type_ *Type) *Rvalue
var contextNewGlobal func(ctx *Context, loc *Location, kind GlobalKind, type_ *Type, name string) *Lvalue
var functionNewLocal func(fn *Function, loc *Location, type_ *Type, name string) *Lvalue
var blockEndWithReturn func(block *Block, loc *Location, rvalue *Rvalue)
var contextGetFirstError func(ctx *Context) string

func init() {
	lib, err := purego.Dlopen("libgccjit.so.0", purego.RTLD_NOW|purego.RTLD_GLOBAL)
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
}

func ContextAcquire() *Context {
	return contextAcquire()
}

func (c *Context) SetBoolOption(opt BoolOption, value bool) {
	contextSetBoolOption(c, opt, value)
}

func (c *Context) GetType(type_ Types) *Type {
	return contextGetType(c, type_)
}

func (c *Context) NewFunction(kind FunctionKind, return_type *Type, name string, params []*Param, isVariadic bool) *Function {
	return contextNewFunction(c, nil, kind, return_type, name, len(params), params, isVariadic)
}

func (c *Context) NewFunctionWithLocation(loc *Location, name string, kind FunctionKind, return_type *Type, params []*Param, isVariadic bool) *Function {
	return contextNewFunction(c, loc, kind, return_type, name, len(params), params, isVariadic)
}

func (c *Context) NewParam(type_ *Type, name string) *Param {
	return contextNewParam(c, nil, type_, name)
}

func (c *Context) NewParamWithLocaction(loc *Location, type_ *Type, name string) *Param {
	return contextNewParam(c, loc, type_, name)
}

func (c *Context) NewBlock(fn *Function, name string) *Block {
	return functionNewBlock(fn, name)
}

func (c *Context) NewCall(fn *Function, args []*Rvalue) *Rvalue {
	return contextNewCall(c, nil, fn, len(args), args)
}

func (c *Context) NewCallWithLocation(loc *Location, fn *Function, args []*Rvalue) *Rvalue {
	return contextNewCall(c, loc, fn, len(args), args)
}

func (c *Context) NewStringLiteral(value string) *Rvalue {
	return contextNewStringLiteral(c, value)
}

func (c *Context) Compile() *Result {
	return contextCompile(c)
}

func (c *Context) Release() {
	contextRelease(c)
}

func (p *Param) AsRvalue() *Rvalue {
	return paramAsRvalue(p)
}

func (b *Block) AddEval(rvalue *Rvalue) {
	blockAddEval(b, nil, rvalue)
}

func (b *Block) AddEvalWithLocation(loc *Location, rvalue *Rvalue) {
	blockAddEval(b, loc, rvalue)
}

func (b *Block) EndWithVoidReturn(loc *Location) {
	blockEndWithVoidReturn(b, loc)
}

func (r *Result) GetCode(name string) uintptr {
	return resultGetCode(r, name)
}

func (r *Result) Release() {
	resultRelease(r)
}
