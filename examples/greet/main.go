package main

import (
	gccjit "github.com/aabajyan/gogccjit/13"
)

func main() {
	ctx := gccjit.ContextAcquire()
	if ctx == nil {
		panic("no context")
	}

	defer ctx.Release()

	ctx.SetBoolOption(gccjit.BOOL_OPTION_DEBUGINFO, false)

	voidType := ctx.GetType(gccjit.TYPE_VOID)
	constCharType := ctx.GetType(gccjit.TYPE_CONST_CHAR_PTR)

	paramName := ctx.NewParam(nil, constCharType, "param")
	fn := ctx.NewFunction(nil, gccjit.FUNCTION_EXPORTED, voidType, "greet", []*gccjit.Param{paramName}, false)

	paramFormat := ctx.NewParam(nil, constCharType, "format")
	printfFunc := ctx.NewFunction(nil, gccjit.FUNCTION_IMPORTED, voidType, "printf", []*gccjit.Param{paramFormat}, true)

	block := ctx.NewBlock(fn, "entry")
	block.AddEval(
		nil,
		ctx.NewCall(
			nil,
			printfFunc,
			[]*gccjit.Rvalue{
				ctx.NewStringLiteral("Hello %s from GO!\n"),
				paramName.AsRvalue(),
			},
		),
	)

	block.EndWithVoidReturn(nil)

	res := ctx.Compile()
	if res == nil {
		panic("res is nil")
	}

	defer res.Release()

	var greet func(name string)
	res.RegisterFunc("greet", &greet)

	greet("world")
}
