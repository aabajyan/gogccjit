package main

import (
	"github.com/ebitengine/purego"
)

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
type gcc_jit_block uintptr

const (
	GCC_JIT_BOOL_OPTION_DEBUGINFO gcc_jit_bool_option = iota
	GCC_JIT_BOOL_OPTION_DUMP_INITIAL_TREE
	GCC_JIT_BOOL_OPTION_DUMP_INITIAL_GIMPLE
	GCC_JIT_BOOL_OPTION_DUMP_GENERATED_CODE
	GCC_JIT_BOOL_OPTION_DUMP_SUMMARY
	GCC_JIT_BOOL_OPTION_DUMP_EVERYTHING
	GCC_JIT_BOOL_OPTION_SELFCHECK_GC
	GCC_JIT_BOOL_OPTION_KEEP_INTERMEDIATES
	GCC_JIT_NUM_BOOL_OPTIONS
)

const (
	GCC_JIT_FUNCTION_EXPORTED gcc_jit_function_kind = iota
	GCC_JIT_FUNCTION_INTERNAL
	GCC_JIT_FUNCTION_IMPORTED
	GCC_JIT_FUNCTION_ALWAYS_INLINE
)

const (
	GCC_JIT_TYPE_VOID gcc_jit_types = iota
	GCC_JIT_TYPE_VOID_PTR
	GCC_JIT_TYPE_BOOL
	GCC_JIT_TYPE_CHAR
	GCC_JIT_TYPE_SIGNED_CHAR
	GCC_JIT_TYPE_UNSIGNED_CHAR
	GCC_JIT_TYPE_SHORT
	GCC_JIT_TYPE_UNSIGNED_SHORT
	GCC_JIT_TYPE_INT
	GCC_JIT_TYPE_UNSIGNED_INT
	GCC_JIT_TYPE_LONG
	GCC_JIT_TYPE_UNSIGNED_LONG
	GCC_JIT_TYPE_LONG_LONG
	GCC_JIT_TYPE_UNSIGNED_LONG_LONG
	GCC_JIT_TYPE_FLOAT
	GCC_JIT_TYPE_DOUBLE
	GCC_JIT_TYPE_LONG_DOUBLE
	GCC_JIT_TYPE_CONST_CHAR_PTR
	GCC_JIT_TYPE_SIZE_T
	GCC_JIT_TYPE_FILE_PTR
	GCC_JIT_TYPE_COMPLEX_FLOAT
	GCC_JIT_TYPE_COMPLEX_DOUBLE
	GCC_JIT_TYPE_COMPLEX_LONG_DOUBLE
	GCC_JIT_TYPE_UINT8_T
	GCC_JIT_TYPE_UINT16_T
	GCC_JIT_TYPE_UINT32_T
	GCC_JIT_TYPE_UINT64_T
	GCC_JIT_TYPE_UINT128_T
	GCC_JIT_TYPE_INT8_T
	GCC_JIT_TYPE_INT16_T
	GCC_JIT_TYPE_INT32_T
	GCC_JIT_TYPE_INT64_T
	GCC_JIT_TYPE_INT128_T
)

var gcc_jit_context_acquire func() gcc_jit_context
var gcc_jit_context_release func(ctx gcc_jit_context)
var gcc_jit_context_set_bool_option func(ctx gcc_jit_context, opt gcc_jit_bool_option, value int) uintptr
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

func create_code(ctx gcc_jit_context) {
	void_type := gcc_jit_context_get_type(ctx, GCC_JIT_TYPE_VOID)
	const_char_type := gcc_jit_context_get_type(ctx, GCC_JIT_TYPE_CONST_CHAR_PTR)
	param_name := gcc_jit_context_new_param(ctx, 0, const_char_type, "param")
	fn := gcc_jit_context_new_function(
		ctx,
		0,
		GCC_JIT_FUNCTION_EXPORTED,
		void_type,
		"greet",
		1,
		[]gcc_jit_param{param_name},
		0,
	)

	param_format := gcc_jit_context_new_param(ctx, 0, const_char_type, "format")
	printf_func := gcc_jit_context_new_function(
		ctx,
		0,
		GCC_JIT_FUNCTION_IMPORTED,
		gcc_jit_context_get_type(ctx, GCC_JIT_TYPE_INT),
		"printf",
		1,
		[]gcc_jit_param{param_format},
		1,
	)

	block := gcc_jit_function_new_block(fn, "entry")
	gcc_jit_block_add_eval(block, 0, gcc_jit_context_new_call(
		ctx,
		0,
		printf_func,
		2,
		[]gcc_jit_rvalue{gcc_jit_context_new_string_literal(ctx, "Hello %s from GO!\n"), gcc_jit_param_as_rvalue(param_name)},
	))

	gcc_jit_block_end_with_void_return(block, 0)
}

func main() {
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

	ctx := gcc_jit_context_acquire()
	if ctx == 0 {
		panic("ctx is null")
	}

	gcc_jit_context_set_bool_option(ctx, GCC_JIT_BOOL_OPTION_DEBUGINFO, 0)

	create_code(ctx)

	res := gcc_jit_context_compile(ctx)
	if res == 0 {
		panic("res is null")
	}
	defer gcc_jit_result_release(res)

	var greet func(name string)
	ptr := gcc_jit_result_get_code(res, "greet")
	purego.RegisterFunc(&greet, ptr)

	greet("world")

	defer gcc_jit_context_release(ctx)

}
