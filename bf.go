package main

import (
	"fmt"
	"os"
)

const MAX_OPEN_PARENS = 20

type bfCompiler struct {
	filename     string
	line, column int
	ctx          *Context
	void_type    *Type
	int_type     *Type
	byte_type    *Type
	array_type   *Type

	func_getchar *Function
	func_putchar *Function
	func_main    *Function

	curblock *Block

	int_zero   *Rvalue
	int_one    *Rvalue
	byte_zero  *Rvalue
	byte_one   *Rvalue
	data_cells *Lvalue
	idx        *Lvalue

	num_open_parens int

	paren_test  []*Block
	paren_body  []*Block
	paren_after []*Block
}

func (c *bfCompiler) fatalError(msg string) {
	fmt.Printf("%s:%d:%d: %s\n", c.filename, c.line, c.column, msg)
	os.Exit(1)
}

func (c *bfCompiler) getCurrentData(loc *Location) *Lvalue {
	return contextNewArrayAccess(
		c.ctx,
		loc,
		lvalueAsRvalue(c.data_cells),
		lvalueAsRvalue(c.idx),
	)
}

func (c *bfCompiler) currentDataIsZero(loc *Location) *Rvalue {
	return contextNewComparison(
		c.ctx,
		loc,
		COMPARISON_EQ,
		lvalueAsRvalue(c.getCurrentData(loc)),
		c.byte_zero,
	)
}

func (c *bfCompiler) compileChar(ch byte) {
	loc := contextNewLocation(c.ctx, c.filename, c.line, c.column)

	switch ch {
	case '>':
		blockAddComment(c.curblock, loc, "'>': idx += 1;")
		blockAddAssignmentOp(
			c.curblock,
			loc,
			c.idx,
			BINARY_OP_PLUS,
			c.int_one,
		)
	case '<':
		blockAddComment(c.curblock, loc, "'<': idx -= 1;")
		blockAddAssignmentOp(
			c.curblock,
			loc,
			c.idx,
			BINARY_OP_MINUS,
			c.int_one,
		)
	case '+':
		blockAddComment(c.curblock, loc, "'+': data[idx] += 1;")
		blockAddAssignmentOp(
			c.curblock,
			loc,
			c.getCurrentData(loc),
			BINARY_OP_PLUS,
			c.byte_one,
		)
	case '-':
		blockAddComment(c.curblock, loc, "'-': data[idx] -= 1;")
		blockAddAssignmentOp(
			c.curblock,
			loc,
			c.getCurrentData(loc),
			BINARY_OP_MINUS,
			c.byte_one,
		)
	case '.':
		arg := contextNewCast(
			c.ctx,
			loc,
			lvalueAsRvalue(c.getCurrentData(loc)),
			c.int_type,
		)

		call := contextNewCall(
			c.ctx,
			loc,
			c.func_putchar,
			1,
			[]*Rvalue{arg},
		)

		blockAddComment(c.curblock, loc, "'.': putchar(data[idx]);")
		blockAddEval(c.curblock, loc, call)
	case ',':
		call := contextNewCall(
			c.ctx,
			loc,
			c.func_getchar,
			0,
			[]*Rvalue{},
		)

		blockAddComment(c.curblock, loc, "',': data[idx] = getchar();")
		blockAddAssignment(
			c.curblock,
			loc,
			c.getCurrentData(loc),
			contextNewCast(
				c.ctx,
				loc,
				call,
				c.byte_type,
			),
		)
	case '[':
		loop_test := functionNewBlock(c.func_main, "loop_test")
		on_zero := functionNewBlock(c.func_main, "on_zero")
		on_non_zero := functionNewBlock(c.func_main, "on_non_zero")

		if c.num_open_parens >= MAX_OPEN_PARENS {
			c.fatalError("too many open parens")
		}

		blockEndWithJump(
			c.curblock,
			loc,
			loop_test,
		)

		blockAddComment(
			loop_test,
			loc,
			"'['",
		)

		blockEndWithConditional(
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
		blockAddComment(c.curblock, loc, "']':")

		if c.num_open_parens == 0 {
			c.fatalError("mismatching parens")
		}
		c.num_open_parens -= 1
		blockEndWithJump(
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

func make_main(ctx *Context) *Function {
	int_type := contextGetType(ctx, TYPE_INT)
	char_ptr_ptr_type := typeGetPointer(contextGetType(ctx, TYPE_CONST_CHAR_PTR))

	param_argc := contextNewParam(ctx, nil, int_type, "argc")
	param_argv := contextNewParam(ctx, nil, char_ptr_ptr_type, "argv")

	main_func := contextNewFunction(
		ctx,
		nil,
		FUNCTION_EXPORTED,
		int_type,
		"main",
		2,
		[]*Param{param_argc, param_argv},
		false,
	)

	return main_func
}

func compile_bf(filename string) {
	c := bfCompiler{
		filename:    filename,
		paren_test:  make([]*Block, MAX_OPEN_PARENS),
		paren_body:  make([]*Block, MAX_OPEN_PARENS),
		paren_after: make([]*Block, MAX_OPEN_PARENS),
	}

	code, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	c.line = 1

	if c.ctx = contextAcquire(); c.ctx == nil {
		panic("failed to acquire context")
	}

	defer contextRelease(c.ctx)

	contextSetIntOption(c.ctx, INT_OPTION_OPTIMIZATION_LEVEL, 3)
	contextSetBoolOption(c.ctx, BOOL_OPTION_DUMP_INITIAL_GIMPLE, false)
	contextSetBoolOption(c.ctx, BOOL_OPTION_DEBUGINFO, true)
	contextSetBoolOption(c.ctx, BOOL_OPTION_DUMP_EVERYTHING, false)
	contextSetBoolOption(c.ctx, BOOL_OPTION_KEEP_INTERMEDIATES, false)

	c.void_type = contextGetType(c.ctx, TYPE_VOID)
	c.int_type = contextGetType(c.ctx, TYPE_INT)
	c.byte_type = contextGetType(c.ctx, TYPE_UNSIGNED_CHAR)
	c.array_type = contextNewArrayType(c.ctx, nil, c.byte_type, 30000)

	c.func_getchar = contextNewFunction(
		c.ctx,
		nil,
		FUNCTION_IMPORTED,
		c.int_type,
		"getchar",
		0,
		[]*Param{},
		false,
	)

	param_c := contextNewParam(c.ctx, nil, c.int_type, "c")
	c.func_putchar = contextNewFunction(
		c.ctx,
		nil,
		FUNCTION_IMPORTED,
		c.void_type,
		"putchar",
		1,
		[]*Param{param_c},
		false,
	)

	c.func_main = make_main(c.ctx)
	c.curblock = functionNewBlock(c.func_main, "main")
	c.int_zero = contextZero(c.ctx, c.int_type)
	c.int_one = contextOne(c.ctx, c.int_type)
	c.byte_zero = contextZero(c.ctx, c.byte_type)
	c.byte_one = contextOne(c.ctx, c.byte_type)
	c.data_cells = contextNewGlobal(c.ctx, nil, GLOBAL_INTERNAL, c.array_type, "data_cells")
	c.idx = functionNewLocal(c.func_main, nil, c.int_type, "idx")

	blockAddComment(c.curblock, nil, "idx = 0;")
	blockAddAssignment(c.curblock, nil, c.idx, c.int_zero)

	c.num_open_parens = 0

	for _, ch := range code {
		c.compileChar(ch)
	}

	blockEndWithReturn(c.curblock, nil, c.int_zero)

	contextCompileToFile(c.ctx, OUTPUT_KIND_EXECUTABLE, "a.out")

	strerr := contextGetFirstError(c.ctx)
	if strerr != "" {
		println(strerr)
	}
}
