package main

import "github.com/ebitengine/purego"

func greet() {
	ctx := ContextAcquire()
	if ctx == nil {
		panic("no context")
	}

	defer ctx.Release()

	ctx.SetBoolOption(BOOL_OPTION_DEBUGINFO, false)

	voidType := ctx.GetType(TYPE_VOID)
	constCharType := ctx.GetType(TYPE_CONST_CHAR_PTR)

	paramName := ctx.NewParam(constCharType, "param")
	fn := ctx.NewFunction(FUNCTION_EXPORTED, voidType, "greet", []*Param{paramName}, false)

	paramFormat := ctx.NewParam(constCharType, "format")
	printfFunc := ctx.NewFunction(FUNCTION_IMPORTED, voidType, "printf", []*Param{paramFormat}, true)

	block := ctx.NewBlock(fn, "entry")
	block.AddEval(
		ctx.NewCall(
			printfFunc,
			[]*Rvalue{
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
	ptr := res.GetCode("greet")
	purego.RegisterFunc(&greet, ptr)
	greet("world")
}
