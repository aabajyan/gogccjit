package main

import "github.com/ebitengine/purego"

func greet() {
	ctx := contextAcquire()
	if ctx == nil {
		panic("ctx is null")
	}

	defer contextRelease(ctx)

	contextSetBoolOption(ctx, BOOL_OPTION_DEBUGINFO, false)

	void_type := contextGetType(ctx, TYPE_VOID)
	const_char_type := contextGetType(ctx, TYPE_CONST_CHAR_PTR)
	param_name := contextNewParam(ctx, nil, const_char_type, "param")
	fn := contextNewFunction(
		ctx,
		nil,
		FUNCTION_EXPORTED,
		void_type,
		"greet",
		1,
		[]*Param{param_name},
		0,
	)

	param_format := contextNewParam(ctx, nil, const_char_type, "format")
	printf_func := contextNewFunction(
		ctx,
		nil,
		FUNCTION_IMPORTED,
		contextGetType(ctx, TYPE_INT),
		"printf",
		1,
		[]*Param{param_format},
		1,
	)

	block := functionNewBlock(fn, "entry")
	blockAddEval(block, nil, contextNewCall(
		ctx,
		nil,
		printf_func,
		2,
		[]*Rvalue{contextNewStringLiteral(ctx, "Hello %s from GO!\n"), paramAsRvalue(param_name)},
	))

	blockEndWithVoidReturn(block, nil)

	res := contextCompile(ctx)
	if res == nil {
		panic("res is null")
	}
	defer resultRelease(res)

	var greet func(name string)
	ptr := resultGetCode(res, "greet")
	purego.RegisterFunc(&greet, ptr)

	greet("world")
}
