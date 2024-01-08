package main

import (
	"fmt"
	"os"

	gccjit "github.com/aabajyan/gogccjit/13"
)

const MAX_OPEN_PARENS = 20

type bfCompiler struct {
	filename     string
	line, column int
	ctx          *gccjit.Context
	void_type    *gccjit.Type
	intType      *gccjit.Type
	byteType     *gccjit.Type
	array_type   *gccjit.Type

	funcGetchar *gccjit.Function
	funcPutchar *gccjit.Function
	funcMain    *gccjit.Function

	curblock *gccjit.Block

	intZero   *gccjit.Rvalue
	intOne    *gccjit.Rvalue
	byteZero  *gccjit.Rvalue
	byteOne   *gccjit.Rvalue
	dataCells *gccjit.Lvalue
	idx       *gccjit.Lvalue

	numOpenParens int

	parenTest  []*gccjit.Block
	parenBody  []*gccjit.Block
	parenAfter []*gccjit.Block
}

func (c *bfCompiler) fatalError(msg string) {
	fmt.Printf("%s:%d:%d: %s\n", c.filename, c.line, c.column, msg)
	os.Exit(1)
}

func (c *bfCompiler) getCurrentData(loc *gccjit.Location) *gccjit.Lvalue {
	return c.ctx.NewArrayAccess(
		loc,
		c.dataCells.AsRvalue(),
		c.idx.AsRvalue(),
	)
}

func (c *bfCompiler) currentDataIsZero(loc *gccjit.Location) *gccjit.Rvalue {
	return c.ctx.NewNewComparison(
		loc,
		gccjit.COMPARISON_EQ,
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
			gccjit.BINARY_OP_PLUS,
			c.intOne,
		)
	case '<':
		c.curblock.AddComment(loc, "'<': idx -= 1;")
		c.curblock.AddAssignmentOp(
			loc,
			c.idx,
			gccjit.BINARY_OP_MINUS,
			c.intOne,
		)
	case '+':
		c.curblock.AddComment(loc, "'+': data[idx] += 1;")
		c.curblock.AddAssignmentOp(
			loc,
			c.getCurrentData(loc),
			gccjit.BINARY_OP_PLUS,
			c.byteOne,
		)
	case '-':
		c.curblock.AddComment(loc, "'-': data[idx] -= 1;")
		c.curblock.AddAssignmentOp(
			loc,
			c.getCurrentData(loc),
			gccjit.BINARY_OP_MINUS,
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
			[]*gccjit.Rvalue{arg},
		)

		c.curblock.AddComment(loc, "'.': putchar(data[idx]);")
		c.curblock.AddEval(loc, call)
	case ',':
		call := c.ctx.NewCall(
			loc,
			c.funcGetchar,
			[]*gccjit.Rvalue{},
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

func makeMain(ctx *gccjit.Context) *gccjit.Function {
	intType := ctx.GetType(gccjit.TYPE_INT)
	charPtrPtrType := ctx.GetType(gccjit.TYPE_CONST_CHAR_PTR).GetPointer()

	paramArgc := ctx.NewParam(nil, intType, "argc")
	paramArgv := ctx.NewParam(nil, charPtrPtrType, "argv")
	mainFunc := ctx.NewFunction(
		nil,
		gccjit.FUNCTION_EXPORTED,
		intType,
		"main",
		[]*gccjit.Param{paramArgc, paramArgv},
		false,
	)

	return mainFunc
}

func main() {
	if len(os.Args) < 2 {
		panic("usage go run ./examples/bf/main.go ./examples/bf/example.bf")
	}

	filename := os.Args[1]

	c := bfCompiler{
		filename:   filename,
		parenTest:  make([]*gccjit.Block, MAX_OPEN_PARENS),
		parenBody:  make([]*gccjit.Block, MAX_OPEN_PARENS),
		parenAfter: make([]*gccjit.Block, MAX_OPEN_PARENS),
	}

	code, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	c.line = 1

	if c.ctx = gccjit.ContextAcquire(); c.ctx == nil {
		panic("failed to acquire context")
	}

	defer c.ctx.Release()

	c.ctx.SetIntOption(gccjit.INT_OPTION_OPTIMIZATION_LEVEL, 3)
	c.ctx.SetBoolOption(gccjit.BOOL_OPTION_DUMP_INITIAL_GIMPLE, false)
	c.ctx.SetBoolOption(gccjit.BOOL_OPTION_DEBUGINFO, true)
	c.ctx.SetBoolOption(gccjit.BOOL_OPTION_DUMP_EVERYTHING, false)
	c.ctx.SetBoolOption(gccjit.BOOL_OPTION_KEEP_INTERMEDIATES, false)

	c.void_type = c.ctx.GetType(gccjit.TYPE_VOID)
	c.intType = c.ctx.GetType(gccjit.TYPE_INT)
	c.byteType = c.ctx.GetType(gccjit.TYPE_UNSIGNED_CHAR)
	c.array_type = c.ctx.GetArrayType(nil, c.byteType, 30000)

	c.funcGetchar = c.ctx.NewFunction(
		nil,
		gccjit.FUNCTION_IMPORTED,
		c.intType,
		"getchar",
		[]*gccjit.Param{},
		false,
	)

	paramC := c.ctx.NewParam(nil, c.intType, "c")
	c.funcPutchar = c.ctx.NewFunction(
		nil,
		gccjit.FUNCTION_IMPORTED,
		c.void_type,
		"putchar",
		[]*gccjit.Param{paramC},
		false,
	)

	c.funcMain = makeMain(c.ctx)
	c.curblock = c.funcMain.NewBlock("main")
	c.intZero = c.ctx.Zero(c.intType)
	c.intOne = c.ctx.One(c.intType)
	c.byteZero = c.ctx.Zero(c.byteType)
	c.byteOne = c.ctx.One(c.byteType)
	c.dataCells = c.ctx.NewGlobal(nil, gccjit.GLOBAL_INTERNAL, c.array_type, "dataCells")
	c.idx = c.funcMain.NewLocal(nil, c.intType, "idx")

	c.curblock.AddComment(nil, "idx = 0;")
	c.curblock.AddAssignment(nil, c.idx, c.intZero)

	c.numOpenParens = 0

	for _, ch := range code {
		c.compileChar(ch)
	}

	c.curblock.EndWithReturn(nil, c.intZero)
	c.ctx.CompileToFile(gccjit.OUTPUT_KIND_EXECUTABLE, "a.out")

	if strerr := c.ctx.GetFirstError(); strerr != "" {
		println(strerr)
	}
}
