package main

import "github.com/ebitengine/purego"

func greet() {
	ctx := gcc_jit_context_acquire()
	if ctx == 0 {
		panic("ctx is null")
	}

	defer gcc_jit_context_release(ctx)

	gcc_jit_context_set_bool_option(ctx, GCC_JIT_BOOL_OPTION_DEBUGINFO, 0)

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

	res := gcc_jit_context_compile(ctx)
	if res == 0 {
		panic("res is null")
	}
	defer gcc_jit_result_release(res)

	var greet func(name string)
	ptr := gcc_jit_result_get_code(res, "greet")
	purego.RegisterFunc(&greet, ptr)

	greet("world")
}
