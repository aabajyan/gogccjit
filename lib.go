package main

import "github.com/ebitengine/purego"

type gcc_jit_function uintptr
type gcc_jit_param uintptr
type gcc_jit_location uintptr
type gcc_jit_context uintptr
type gcc_jit_result uintptr
type gcc_jit_bool_option int
type gcc_jit_types int
type gcc_jit_function_kind int
type gcc_jit_type uintptr
type gcc_jit_rvalue uintptr
type gcc_jit_lvalue uintptr
type gcc_jit_block uintptr
type gcc_jit_output_kind int
type gcc_jit_comparison int
type gcc_jit_binary_op int
type gcc_jit_int_option int
type gcc_jit_global_kind int

const (
	BOOL_OPTION_DEBUGINFO gcc_jit_bool_option = iota
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
	FUNCTION_EXPORTED gcc_jit_function_kind = iota
	FUNCTION_INTERNAL
	FUNCTION_IMPORTED
	FUNCTION_ALWAYS_INLINE
)

const (
	TYPE_VOID gcc_jit_types = iota
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
	OUTPUT_KIND_ASSEMBLER gcc_jit_output_kind = iota
	OUTPUT_KIND_OBJECT_FILE
	OUTPUT_KIND_DYNAMIC_LIBRARY
	OUTPUT_KIND_EXECUTABLE
)

const (
	COMPARISON_EQ gcc_jit_comparison = iota
	COMPARISON_NE
	COMPARISON_LT
	COMPARISON_LE
	COMPARISON_GT
	COMPARISON_GE
)

const (
	BINARY_OP_PLUS gcc_jit_binary_op = iota
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
	INT_OPTION_OPTIMIZATION_LEVEL gcc_jit_int_option = iota
	NUM_INT_OPTIONS
)

const (
	GLOBAL_EXPORTED gcc_jit_global_kind = iota
	GLOBAL_INTERNAL
	GLOBAL_IMPORTED
)

var gcc_jit_context_acquire func() gcc_jit_context
var gcc_jit_context_release func(ctx gcc_jit_context)
var gcc_jit_context_set_bool_option func(ctx gcc_jit_context, opt gcc_jit_bool_option, value bool) uintptr
var gcc_jit_context_compile func(ctx gcc_jit_context) gcc_jit_result
var gcc_jit_result_release func(ctx gcc_jit_result)
var gcc_jit_context_get_type func(ctx gcc_jit_context, type_ gcc_jit_types) gcc_jit_type
var gcc_jit_context_new_param func(ctx gcc_jit_context, loc gcc_jit_location, type_ gcc_jit_type, name string) gcc_jit_param
var gcc_jit_context_new_function func(ctx gcc_jit_context, loc gcc_jit_location, kind gcc_jit_function_kind, return_type gcc_jit_type, name string, num_params int, params []gcc_jit_param, is_variadic int) gcc_jit_function
var gcc_jit_context_new_string_literal func(ctx gcc_jit_context, value string) gcc_jit_rvalue
var gcc_jit_function_new_block func(fn gcc_jit_function, name string) gcc_jit_block
var gcc_jit_block_add_eval func(block gcc_jit_block, loc gcc_jit_location, rvalue gcc_jit_rvalue)
var gcc_jit_context_new_call func(ctx gcc_jit_context, loc gcc_jit_location, fn gcc_jit_function, numargs int, args []gcc_jit_rvalue) gcc_jit_rvalue
var gcc_jit_block_end_with_void_return func(block gcc_jit_block, loc gcc_jit_location)
var gcc_jit_result_get_code func(result gcc_jit_result, name string) uintptr
var gcc_jit_param_as_rvalue func(param gcc_jit_param) gcc_jit_rvalue
var gcc_jit_context_compile_to_file func(ctx gcc_jit_context, output_kind gcc_jit_output_kind, output_path string)
var gcc_jit_context_new_array_access func(ctx gcc_jit_context, loc gcc_jit_location, ptr gcc_jit_rvalue, idx gcc_jit_rvalue) gcc_jit_lvalue
var gcc_jit_lvalue_as_rvalue func(lvalue gcc_jit_lvalue) gcc_jit_rvalue
var gcc_jit_context_new_comparison func(ctx gcc_jit_context, loc gcc_jit_location, op gcc_jit_comparison, lhs gcc_jit_rvalue, rhs gcc_jit_rvalue) gcc_jit_rvalue
var gcc_jit_context_new_location func(ctx gcc_jit_context, filename string, line, column int) gcc_jit_location
var gcc_jit_block_add_comment func(block gcc_jit_block, loc gcc_jit_location, text string)
var gcc_jit_block_add_assignment_op func(block gcc_jit_block, loc gcc_jit_location, lvalue gcc_jit_lvalue, op gcc_jit_binary_op, rvalue gcc_jit_rvalue)
var gcc_jit_context_new_cast func(ctx gcc_jit_context, loc gcc_jit_location, rvalue gcc_jit_rvalue, type_ gcc_jit_type) gcc_jit_rvalue
var gcc_jit_block_add_assignment func(block gcc_jit_block, loc gcc_jit_location, lvalue gcc_jit_lvalue, rvalue gcc_jit_rvalue)
var gcc_jit_block_end_with_jump func(block gcc_jit_block, loc gcc_jit_location, target gcc_jit_block)
var gcc_jit_block_end_with_conditional func(block gcc_jit_block, loc gcc_jit_location, boolval gcc_jit_rvalue, on_true gcc_jit_block, on_false gcc_jit_block)
var gcc_jit_context_set_int_option func(ctx gcc_jit_context, opt gcc_jit_int_option, value int)
var gcc_jit_context_new_array_type func(ctx gcc_jit_context, loc gcc_jit_location, element_type gcc_jit_type, num_elements int) gcc_jit_type
var gcc_jit_type_get_pointer func(type_ gcc_jit_type) gcc_jit_type
var gcc_jit_context_zero func(ctx gcc_jit_context, type_ gcc_jit_type) gcc_jit_rvalue
var gcc_jit_context_one func(ctx gcc_jit_context, type_ gcc_jit_type) gcc_jit_rvalue
var gcc_jit_context_new_global func(ctx gcc_jit_context, loc gcc_jit_location, kind gcc_jit_global_kind, type_ gcc_jit_type, name string) gcc_jit_lvalue
var gcc_jit_function_new_local func(fn gcc_jit_function, loc gcc_jit_location, type_ gcc_jit_type, name string) gcc_jit_lvalue
var gcc_jit_block_end_with_return func(block gcc_jit_block, loc gcc_jit_location, rvalue gcc_jit_rvalue)
var gcc_jit_context_get_first_error func(ctx gcc_jit_context) string

func init() {
	lib, err := purego.Dlopen("libgccjit.so.0", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}

	purego.RegisterLibFunc(&gcc_jit_context_acquire, lib, "gcc_jit_context_acquire")
	purego.RegisterLibFunc(&gcc_jit_context_release, lib, "gcc_jit_context_release")
	purego.RegisterLibFunc(&gcc_jit_context_set_bool_option, lib, "gcc_jit_context_set_bool_option")
	purego.RegisterLibFunc(&gcc_jit_context_compile, lib, "gcc_jit_context_compile")
	purego.RegisterLibFunc(&gcc_jit_result_release, lib, "gcc_jit_result_release")
	purego.RegisterLibFunc(&gcc_jit_context_get_type, lib, "gcc_jit_context_get_type")
	purego.RegisterLibFunc(&gcc_jit_context_new_param, lib, "gcc_jit_context_new_param")
	purego.RegisterLibFunc(&gcc_jit_context_new_function, lib, "gcc_jit_context_new_function")
	purego.RegisterLibFunc(&gcc_jit_context_new_string_literal, lib, "gcc_jit_context_new_string_literal")
	purego.RegisterLibFunc(&gcc_jit_function_new_block, lib, "gcc_jit_function_new_block")
	purego.RegisterLibFunc(&gcc_jit_block_add_eval, lib, "gcc_jit_block_add_eval")
	purego.RegisterLibFunc(&gcc_jit_context_new_call, lib, "gcc_jit_context_new_call")
	purego.RegisterLibFunc(&gcc_jit_block_end_with_void_return, lib, "gcc_jit_block_end_with_void_return")
	purego.RegisterLibFunc(&gcc_jit_result_get_code, lib, "gcc_jit_result_get_code")
	purego.RegisterLibFunc(&gcc_jit_param_as_rvalue, lib, "gcc_jit_param_as_rvalue")
	purego.RegisterLibFunc(&gcc_jit_context_compile_to_file, lib, "gcc_jit_context_compile_to_file")
	purego.RegisterLibFunc(&gcc_jit_context_new_array_access, lib, "gcc_jit_context_new_array_access")
	purego.RegisterLibFunc(&gcc_jit_lvalue_as_rvalue, lib, "gcc_jit_lvalue_as_rvalue")
	purego.RegisterLibFunc(&gcc_jit_context_new_comparison, lib, "gcc_jit_context_new_comparison")
	purego.RegisterLibFunc(&gcc_jit_context_new_location, lib, "gcc_jit_context_new_location")
	purego.RegisterLibFunc(&gcc_jit_block_add_comment, lib, "gcc_jit_block_add_comment")
	purego.RegisterLibFunc(&gcc_jit_block_add_assignment_op, lib, "gcc_jit_block_add_assignment_op")
	purego.RegisterLibFunc(&gcc_jit_context_new_cast, lib, "gcc_jit_context_new_cast")
	purego.RegisterLibFunc(&gcc_jit_block_add_assignment, lib, "gcc_jit_block_add_assignment")
	purego.RegisterLibFunc(&gcc_jit_block_end_with_jump, lib, "gcc_jit_block_end_with_jump")
	purego.RegisterLibFunc(&gcc_jit_block_end_with_conditional, lib, "gcc_jit_block_end_with_conditional")
	purego.RegisterLibFunc(&gcc_jit_context_set_int_option, lib, "gcc_jit_context_set_int_option")
	purego.RegisterLibFunc(&gcc_jit_context_new_array_type, lib, "gcc_jit_context_new_array_type")
	purego.RegisterLibFunc(&gcc_jit_type_get_pointer, lib, "gcc_jit_type_get_pointer")
	purego.RegisterLibFunc(&gcc_jit_context_zero, lib, "gcc_jit_context_zero")
	purego.RegisterLibFunc(&gcc_jit_context_one, lib, "gcc_jit_context_one")
	purego.RegisterLibFunc(&gcc_jit_context_new_global, lib, "gcc_jit_context_new_global")
	purego.RegisterLibFunc(&gcc_jit_function_new_local, lib, "gcc_jit_function_new_local")
	purego.RegisterLibFunc(&gcc_jit_block_end_with_return, lib, "gcc_jit_block_end_with_return")
	purego.RegisterLibFunc(&gcc_jit_context_get_first_error, lib, "gcc_jit_context_get_first_error")

}
