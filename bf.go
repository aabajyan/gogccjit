package main

import (
	"fmt"
	"os"
)

const MAX_OPEN_PARENS = 20

type bfCompiler struct {
	filename     string
	line, column int
	ctx          gcc_jit_context
	void_type    gcc_jit_type
	int_type     gcc_jit_type
	byte_type    gcc_jit_type
	array_type   gcc_jit_type

	func_getchar gcc_jit_function
	func_putchar gcc_jit_function
	func_main    gcc_jit_function

	curblock gcc_jit_block

	int_zero   gcc_jit_rvalue
	int_one    gcc_jit_rvalue
	byte_zero  gcc_jit_rvalue
	byte_one   gcc_jit_rvalue
	data_cells gcc_jit_lvalue
	idx        gcc_jit_lvalue

	num_open_parens int

	paren_test  []gcc_jit_block
	paren_body  []gcc_jit_block
	paren_after []gcc_jit_block
}

func (c *bfCompiler) fatalError(msg string) {
	fmt.Printf("%s:%d:%d: %s\n", c.filename, c.line, c.column, msg)
	os.Exit(1)
}

func (c *bfCompiler) getCurrentData(loc gcc_jit_location) gcc_jit_lvalue {
	return gcc_jit_context_new_array_access(
		c.ctx,
		loc,
		gcc_jit_lvalue_as_rvalue(c.data_cells),
		gcc_jit_lvalue_as_rvalue(c.idx),
	)
}

func (c *bfCompiler) currentDataIsZero(loc gcc_jit_location) gcc_jit_rvalue {
	return gcc_jit_context_new_comparison(
		c.ctx,
		loc,
		GCC_JIT_COMPARISON_EQ,
		gcc_jit_lvalue_as_rvalue(c.getCurrentData(loc)),
		c.byte_zero,
	)
}

func (c *bfCompiler) compileChar(ch byte) {
	loc := gcc_jit_context_new_location(c.ctx, c.filename, c.line, c.column)

	switch ch {
	case '>':
		gcc_jit_block_add_comment(c.curblock, loc, "'>': idx += 1;")
		gcc_jit_block_add_assignment_op(
			c.curblock,
			loc,
			c.idx,
			GCC_JIT_BINARY_OP_PLUS,
			c.int_one,
		)
	case '<':
		gcc_jit_block_add_comment(c.curblock, loc, "'<': idx -= 1;")
		gcc_jit_block_add_assignment_op(
			c.curblock,
			loc,
			c.idx,
			GCC_JIT_BINARY_OP_MINUS,
			c.int_one,
		)
	case '+':
		gcc_jit_block_add_comment(c.curblock, loc, "'+': data[idx] += 1;")
		gcc_jit_block_add_assignment_op(
			c.curblock,
			loc,
			c.getCurrentData(loc),
			GCC_JIT_BINARY_OP_PLUS,
			c.byte_one,
		)
	case '-':
		gcc_jit_block_add_comment(c.curblock, loc, "'-': data[idx] -= 1;")
		gcc_jit_block_add_assignment_op(
			c.curblock,
			loc,
			c.getCurrentData(loc),
			GCC_JIT_BINARY_OP_MINUS,
			c.byte_one,
		)
	case '.':
		arg := gcc_jit_context_new_cast(
			c.ctx,
			loc,
			gcc_jit_lvalue_as_rvalue(c.getCurrentData(loc)),
			c.int_type,
		)

		call := gcc_jit_context_new_call(
			c.ctx,
			loc,
			c.func_putchar,
			1,
			[]gcc_jit_rvalue{arg},
		)

		gcc_jit_block_add_comment(c.curblock, loc, "'.': putchar(data[idx]);")
		gcc_jit_block_add_eval(c.curblock, loc, call)
	case ',':
		call := gcc_jit_context_new_call(
			c.ctx,
			loc,
			c.func_getchar,
			0,
			[]gcc_jit_rvalue{},
		)

		gcc_jit_block_add_comment(c.curblock, loc, "',': data[idx] = getchar();")
		gcc_jit_block_add_assignment(
			c.curblock,
			loc,
			c.getCurrentData(loc),
			gcc_jit_context_new_cast(
				c.ctx,
				loc,
				call,
				c.byte_type,
			),
		)
	case '[':
		loop_test := gcc_jit_function_new_block(c.func_main, "loop_test")
		on_zero := gcc_jit_function_new_block(c.func_main, "on_zero")
		on_non_zero := gcc_jit_function_new_block(c.func_main, "on_non_zero")

		if c.num_open_parens >= MAX_OPEN_PARENS {
			c.fatalError("too many open parens")
		}

		gcc_jit_block_end_with_jump(
			c.curblock,
			loc,
			loop_test,
		)

		gcc_jit_block_add_comment(
			loop_test,
			loc,
			"'['",
		)

		gcc_jit_block_end_with_conditional(
			loop_test,
			loc,
			c.currentDataIsZero(loc),
			on_zero,
			on_non_zero,
		)

		c.paren_test[c.num_open_parens] = loop_test
		c.paren_body[c.num_open_parens] = on_non_zero
		c.paren_after[c.num_open_parens] = on_zero
		c.num_open_parens += 1
		c.curblock = on_non_zero
	case ']':
		gcc_jit_block_add_comment(c.curblock, loc, "']':")

		if c.num_open_parens == 0 {
			c.fatalError("mismatching parens")
		}
		c.num_open_parens -= 1
		gcc_jit_block_end_with_jump(
			c.curblock,
			loc,
			c.paren_test[c.num_open_parens],
		)
		c.curblock = c.paren_after[c.num_open_parens]
	case '\n':
		c.line += 1
		c.column = 0
	}

	if ch != '\n' {
		c.column += 1
	}
}

func make_main(ctx gcc_jit_context) gcc_jit_function {
	int_type := gcc_jit_context_get_type(ctx, GCC_JIT_TYPE_INT)
	char_ptr_ptr_type := gcc_jit_type_get_pointer(gcc_jit_context_get_type(ctx, GCC_JIT_TYPE_CONST_CHAR_PTR))

	param_argc := gcc_jit_context_new_param(ctx, 0, int_type, "argc")
	param_argv := gcc_jit_context_new_param(ctx, 0, char_ptr_ptr_type, "argv")

	main_func := gcc_jit_context_new_function(
		ctx,
		0,
		GCC_JIT_FUNCTION_EXPORTED,
		int_type,
		"main",
		2,
		[]gcc_jit_param{param_argc, param_argv},
		0,
	)

	return main_func
}

func compile_bf(filename string) {
	c := bfCompiler{
		filename:    filename,
		paren_test:  make([]gcc_jit_block, MAX_OPEN_PARENS),
		paren_body:  make([]gcc_jit_block, MAX_OPEN_PARENS),
		paren_after: make([]gcc_jit_block, MAX_OPEN_PARENS),
	}

	code, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	c.line = 1

	if c.ctx = gcc_jit_context_acquire(); c.ctx == 0 {
		panic("failed to acquire context")
	}

	defer gcc_jit_context_release(c.ctx)

	gcc_jit_context_set_int_option(c.ctx, GCC_JIT_INT_OPTION_OPTIMIZATION_LEVEL, 3)
	gcc_jit_context_set_bool_option(c.ctx, GCC_JIT_BOOL_OPTION_DUMP_INITIAL_GIMPLE, false)
	gcc_jit_context_set_bool_option(c.ctx, GCC_JIT_BOOL_OPTION_DEBUGINFO, true)
	gcc_jit_context_set_bool_option(c.ctx, GCC_JIT_BOOL_OPTION_DUMP_EVERYTHING, false)
	gcc_jit_context_set_bool_option(c.ctx, GCC_JIT_BOOL_OPTION_KEEP_INTERMEDIATES, false)

	c.void_type = gcc_jit_context_get_type(c.ctx, GCC_JIT_TYPE_VOID)
	c.int_type = gcc_jit_context_get_type(c.ctx, GCC_JIT_TYPE_INT)
	c.byte_type = gcc_jit_context_get_type(c.ctx, GCC_JIT_TYPE_UNSIGNED_CHAR)
	c.array_type = gcc_jit_context_new_array_type(c.ctx, 0, c.byte_type, 30000)

	c.func_getchar = gcc_jit_context_new_function(
		c.ctx,
		0,
		GCC_JIT_FUNCTION_IMPORTED,
		c.int_type,
		"getchar",
		0,
		[]gcc_jit_param{},
		0,
	)

	param_c := gcc_jit_context_new_param(c.ctx, 0, c.int_type, "c")
	c.func_putchar = gcc_jit_context_new_function(
		c.ctx,
		0,
		GCC_JIT_FUNCTION_IMPORTED,
		c.void_type,
		"putchar",
		1,
		[]gcc_jit_param{param_c},
		0,
	)

	c.func_main = make_main(c.ctx)
	c.curblock = gcc_jit_function_new_block(c.func_main, "main")
	c.int_zero = gcc_jit_context_zero(c.ctx, c.int_type)
	c.int_one = gcc_jit_context_one(c.ctx, c.int_type)
	c.byte_zero = gcc_jit_context_zero(c.ctx, c.byte_type)
	c.byte_one = gcc_jit_context_one(c.ctx, c.byte_type)
	c.data_cells = gcc_jit_context_new_global(c.ctx, 0, GCC_JIT_GLOBAL_INTERNAL, c.array_type, "data_cells")
	c.idx = gcc_jit_function_new_local(c.func_main, 0, c.int_type, "idx")

	gcc_jit_block_add_comment(c.curblock, 0, "idx = 0;")
	gcc_jit_block_add_assignment(c.curblock, 0, c.idx, c.int_zero)

	c.num_open_parens = 0

	for _, ch := range code {
		c.compileChar(ch)
	}

	gcc_jit_block_end_with_return(c.curblock, 0, c.int_zero)

	gcc_jit_context_compile_to_file(c.ctx, GCC_JIT_OUTPUT_KIND_EXECUTABLE, "a.out")

	strerr := gcc_jit_context_get_first_error(c.ctx)
	if strerr != "" {
		println(strerr)
	}
}
