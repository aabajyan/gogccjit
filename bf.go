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
	intType      *Type
	byteType     *Type
	array_type   *Type

	funcGetchar *Function
	funcPutchar *Function
	funcMain    *Function

	curblock *Block

	intZero   *Rvalue
	intOne    *Rvalue
	byteZero  *Rvalue
	byteOne   *Rvalue
	dataCells *Lvalue
	idx       *Lvalue

	numOpenParens int

	parenTest  []*Block
	parenBody  []*Block
	parenAfter []*Block
}

func (c *bfCompiler) fatalError(msg string) {
	fmt.Printf("%s:%d:%d: %s\n", c.filename, c.line, c.column, msg)
	os.Exit(1)
}

func (c *bfCompiler) getCurrentData(loc *Location) *Lvalue {
	return c.ctx.NewArrayAccess(
		loc,
		c.dataCells.AsRvalue(),
		c.idx.AsRvalue(),
	)
}

func (c *bfCompiler) currentDataIsZero(loc *Location) *Rvalue {
	return c.ctx.NewNewComparison(
		loc,
		COMPARISON_EQ,
		c.getCurrentData(loc).AsRvalue(),
		c.byteZero,
	)
}

func (c *bfCompiler) compileChar(ch byte) {
	loc := c.ctx.NewLocation(c.filename, c.line, c.column)

	switch ch {
	case '>':
		c.curblock.AddComment(loc, "'>': idx += 1;")
		c.curblock.AddAssignmentOp(
			loc,
			c.idx,
			BINARY_OP_PLUS,
			c.intOne,
		)
	case '<':
		c.curblock.AddComment(loc, "'<': idx -= 1;")
		c.curblock.AddAssignmentOp(
			loc,
			c.idx,
			BINARY_OP_MINUS,
			c.intOne,
		)
	case '+':
		c.curblock.AddComment(loc, "'+': data[idx] += 1;")
		c.curblock.AddAssignmentOp(
			loc,
			c.getCurrentData(loc),
			BINARY_OP_PLUS,
			c.byteOne,
		)
	case '-':
		c.curblock.AddComment(loc, "'-': data[idx] -= 1;")
		c.curblock.AddAssignmentOp(
			loc,
			c.getCurrentData(loc),
			BINARY_OP_MINUS,
			c.byteOne,
		)
	case '.':
		arg := c.ctx.NewCast(
			loc,
			c.getCurrentData(loc).AsRvalue(),
			c.intType,
		)

		call := c.ctx.NewCall(
			loc,
			c.funcPutchar,
			[]*Rvalue{arg},
		)

		c.curblock.AddComment(loc, "'.': putchar(data[idx]);")
		c.curblock.AddEval(loc, call)
	case ',':
		call := c.ctx.NewCall(
			loc,
			c.funcGetchar,
			[]*Rvalue{},
		)

		c.curblock.AddComment(loc, "',': data[idx] = getchar();")
		c.curblock.AddAssignment(
			loc,
			c.getCurrentData(loc),
			c.ctx.NewCast(
				loc,
				call,
				c.byteType,
			),
		)
	case '[':
		loopTest := c.funcMain.NewBlock("loopTest")
		onZero := c.funcMain.NewBlock("onZero")
		onNonZero := c.funcMain.NewBlock("onNonZero")

		if c.numOpenParens >= MAX_OPEN_PARENS {
			c.fatalError("too many open parens")
		}

		c.curblock.EndWithJump(loc, loopTest)
		loopTest.AddComment(loc, "'[':")
		loopTest.EndWithConditional(
			loc,
			c.currentDataIsZero(loc),
			onZero,
			onNonZero,
		)

		c.parenTest[c.numOpenParens] = loopTest
		c.parenBody[c.numOpenParens] = onNonZero
		c.parenAfter[c.numOpenParens] = onZero
		c.numOpenParens += 1
		c.curblock = onNonZero
	case ']':
		c.curblock.AddComment(loc, "']':")
		if c.numOpenParens == 0 {
			c.fatalError("mismatching parens")
		}

		c.numOpenParens -= 1
		c.curblock.EndWithJump(loc, c.parenTest[c.numOpenParens])
		c.curblock = c.parenAfter[c.numOpenParens]
	case '\n':
		c.line += 1
		c.column = 0
	}

	if ch != '\n' {
		c.column += 1
	}
}

func makeMain(ctx *Context) *Function {
	intType := ctx.GetType(TYPE_INT)
	charPtrPtrType := ctx.GetType(TYPE_CONST_CHAR_PTR).GetPointer()

	paramArgc := ctx.NewParam(nil, intType, "argc")
	paramArgv := ctx.NewParam(nil, charPtrPtrType, "argv")
	mainFunc := ctx.NewFunction(
		nil,
		FUNCTION_EXPORTED,
		intType,
		"main",
		[]*Param{paramArgc, paramArgv},
		false,
	)

	return mainFunc
}

func compile_bf(filename string) {
	c := bfCompiler{
		filename:   filename,
		parenTest:  make([]*Block, MAX_OPEN_PARENS),
		parenBody:  make([]*Block, MAX_OPEN_PARENS),
		parenAfter: make([]*Block, MAX_OPEN_PARENS),
	}

	code, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	c.line = 1

	if c.ctx = ContextAcquire(); c.ctx == nil {
		panic("failed to acquire context")
	}

	defer c.ctx.Release()

	c.ctx.SetIntOption(INT_OPTION_OPTIMIZATION_LEVEL, 3)
	c.ctx.SetBoolOption(BOOL_OPTION_DUMP_INITIAL_GIMPLE, false)
	c.ctx.SetBoolOption(BOOL_OPTION_DEBUGINFO, true)
	c.ctx.SetBoolOption(BOOL_OPTION_DUMP_EVERYTHING, false)
	c.ctx.SetBoolOption(BOOL_OPTION_KEEP_INTERMEDIATES, false)

	c.void_type = c.ctx.GetType(TYPE_VOID)
	c.intType = c.ctx.GetType(TYPE_INT)
	c.byteType = c.ctx.GetType(TYPE_UNSIGNED_CHAR)
	c.array_type = c.ctx.GetArrayType(nil, c.byteType, 30000)

	c.funcGetchar = c.ctx.NewFunction(
		nil,
		FUNCTION_IMPORTED,
		c.intType,
		"getchar",
		[]*Param{},
		false,
	)

	paramC := c.ctx.NewParam(nil, c.intType, "c")
	c.funcPutchar = c.ctx.NewFunction(
		nil,
		FUNCTION_IMPORTED,
		c.void_type,
		"putchar",
		[]*Param{paramC},
		false,
	)

	c.funcMain = makeMain(c.ctx)
	c.curblock = c.funcMain.NewBlock("main")
	c.intZero = c.ctx.Zero(c.intType)
	c.intOne = c.ctx.One(c.intType)
	c.byteZero = c.ctx.Zero(c.byteType)
	c.byteOne = c.ctx.One(c.byteType)
	c.dataCells = c.ctx.NewGlobal(nil, GLOBAL_INTERNAL, c.array_type, "dataCells")
	c.idx = c.funcMain.NewLocal(nil, c.intType, "idx")

	c.curblock.AddComment(nil, "idx = 0;")
	c.curblock.AddAssignment(nil, c.idx, c.intZero)

	c.numOpenParens = 0

	for _, ch := range code {
		c.compileChar(ch)
	}

	c.curblock.EndWithReturn(nil, c.intZero)
	c.ctx.CompileToFile(OUTPUT_KIND_EXECUTABLE, "a.out")

	if strerr := c.ctx.GetFirstError(); strerr != "" {
		println(strerr)
	}
}
