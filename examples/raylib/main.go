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

	voidPtrType := ctx.GetType(gccjit.TYPE_VOID_PTR)
	voidType := ctx.GetType(gccjit.TYPE_VOID)
	intType := ctx.GetType(gccjit.TYPE_INT)
	constCharType := ctx.GetType(gccjit.TYPE_CONST_CHAR_PTR)
	boolType := ctx.GetType(gccjit.TYPE_BOOL)
	charType := ctx.GetType(gccjit.TYPE_CHAR)
	charPtrPtrType := ctx.GetType(gccjit.TYPE_CONST_CHAR_PTR).GetPointer()

	printfFunc := ctx.NewFunction(
		nil,
		gccjit.FUNCTION_IMPORTED,
		intType,
		"printf",
		[]*gccjit.Param{
			ctx.NewParam(nil, constCharType, "format"),
		},
		true,
	)

	dlOpenFunc := ctx.NewFunction(
		nil,
		gccjit.FUNCTION_IMPORTED,
		voidPtrType,
		"dlopen",
		[]*gccjit.Param{
			ctx.NewParam(nil, constCharType, "filename"),
			ctx.NewParam(nil, intType, "flags"),
		},
		false,
	)

	dlCloseFunc := ctx.NewFunction(
		nil,
		gccjit.FUNCTION_IMPORTED,
		voidType,
		"dlclose",
		[]*gccjit.Param{
			ctx.NewParam(nil, voidPtrType, "handle"),
		},
		false,
	)

	dlsymFunc := ctx.NewFunction(
		nil,
		gccjit.FUNCTION_IMPORTED,
		voidPtrType,
		"dlsym",
		[]*gccjit.Param{
			ctx.NewParam(nil, voidPtrType, "handle"),
			ctx.NewParam(nil, constCharType, "symbol"),
		},
		false,
	)

	dlerrorFunc := ctx.NewFunction(
		nil,
		gccjit.FUNCTION_IMPORTED,
		constCharType,
		"dlerror",
		[]*gccjit.Param{},
		false,
	)

	rField := ctx.NewField(nil, charType, "r")
	gField := ctx.NewField(nil, charType, "g")
	bField := ctx.NewField(nil, charType, "b")
	aField := ctx.NewField(nil, charType, "a")
	colorType := ctx.NewStructType(
		nil,
		"Color",
		[]*gccjit.Field{
			rField,
			gField,
			bField,
			aField,
		},
	)

	initWindowFuncType := ctx.NewFunctionPtrType(
		nil,
		voidType,
		[]*gccjit.Type{
			intType,
			intType,
			constCharType,
		},
		false,
	)

	setTargetFPSFuncType := ctx.NewFunctionPtrType(
		nil,
		voidType,
		[]*gccjit.Type{
			intType,
		},
		false,
	)

	windowShouldCloseFuncType := ctx.NewFunctionPtrType(
		nil,
		boolType,
		[]*gccjit.Type{},
		false,
	)

	beginDrawingFuncType := ctx.NewFunctionPtrType(
		nil,
		voidType,
		[]*gccjit.Type{},
		false,
	)

	endDrawingFuncType := ctx.NewFunctionPtrType(
		nil,
		voidType,
		[]*gccjit.Type{},
		false,
	)

	clearBackgroundFuncType := ctx.NewFunctionPtrType(
		nil,
		voidType,
		[]*gccjit.Type{
			colorType.AsType(),
		},
		false,
	)

	drawTextFuncType := ctx.NewFunctionPtrType(
		nil,
		voidType,
		[]*gccjit.Type{
			constCharType,
			intType,
			intType,
			intType,
			colorType.AsType(),
		},
		false,
	)

	closeWindowFuncType := ctx.NewFunctionPtrType(
		nil,
		voidType,
		[]*gccjit.Type{},
		false,
	)

	mainFunc := ctx.NewFunction(
		nil,
		gccjit.FUNCTION_EXPORTED,
		intType,
		"main",
		[]*gccjit.Param{
			ctx.NewParam(nil, intType, "argc"),
			ctx.NewParam(nil, charPtrPtrType, "argv"),
		},
		false,
	)

	block := mainFunc.NewBlock("entry")

	handle := mainFunc.NewLocal(
		nil,
		voidPtrType,
		"handle",
	)

	block.AddAssignment(
		nil,
		handle,
		ctx.NewCall(
			nil,
			dlOpenFunc,
			[]*gccjit.Rvalue{
				ctx.NewStringLiteral("libraylib.so"),
				ctx.NewRValueFromLong(intType, 0x00001),
			},
		),
	)

	handleLoaded := mainFunc.NewBlock("handleLoaded")
	handleNotLoaded := mainFunc.NewBlock("handleNotLoaded")

	block.EndWithConditional(
		nil,
		ctx.NewNewComparison(
			nil,
			gccjit.COMPARISON_NE,
			handle.AsRvalue(),
			ctx.NewRvalueFromPtr(voidPtrType, 0),
		),
		handleLoaded,
		handleNotLoaded,
	)

	handleNotLoaded.AddEval(
		nil,
		ctx.NewCall(
			nil,
			printfFunc,
			[]*gccjit.Rvalue{
				ctx.NewStringLiteral("Failed to load library libraylib.so: %s\n"),
				ctx.NewCall(nil, dlerrorFunc, []*gccjit.Rvalue{}),
			},
		),
	)

	handleNotLoaded.EndWithReturn(nil, ctx.NewRValueFromInt(intType, 1))

	initWindowFunc := mainFunc.NewLocal(nil, initWindowFuncType, "InitWindow")
	handleLoaded.AddAssignment(
		nil,
		initWindowFunc,
		ctx.NewCast(
			nil,
			ctx.NewCall(
				nil,
				dlsymFunc,
				[]*gccjit.Rvalue{
					handle.AsRvalue(),
					ctx.NewStringLiteral("InitWindow"),
				},
			),
			initWindowFuncType,
		),
	)

	setTargetFPSFunc := mainFunc.NewLocal(nil, setTargetFPSFuncType, "SetTargetFPS")
	handleLoaded.AddAssignment(
		nil,
		setTargetFPSFunc,
		ctx.NewCast(
			nil,
			ctx.NewCall(
				nil,
				dlsymFunc,
				[]*gccjit.Rvalue{
					handle.AsRvalue(),
					ctx.NewStringLiteral("SetTargetFPS"),
				},
			),
			setTargetFPSFuncType,
		),
	)

	beginDrawingFunc := mainFunc.NewLocal(nil, beginDrawingFuncType, "BeginDrawing")
	handleLoaded.AddAssignment(
		nil,
		beginDrawingFunc,
		ctx.NewCast(
			nil,
			ctx.NewCall(
				nil,
				dlsymFunc,
				[]*gccjit.Rvalue{
					handle.AsRvalue(),
					ctx.NewStringLiteral("BeginDrawing"),
				},
			),
			beginDrawingFuncType,
		),
	)

	endDrawingFunc := mainFunc.NewLocal(nil, endDrawingFuncType, "EndDrawing")
	handleLoaded.AddAssignment(
		nil,
		endDrawingFunc,
		ctx.NewCast(
			nil,
			ctx.NewCall(
				nil,
				dlsymFunc,
				[]*gccjit.Rvalue{
					handle.AsRvalue(),
					ctx.NewStringLiteral("EndDrawing"),
				},
			),
			endDrawingFuncType,
		),
	)

	clearBackgroundFunc := mainFunc.NewLocal(nil, clearBackgroundFuncType, "ClearBackground")
	handleLoaded.AddAssignment(
		nil,
		clearBackgroundFunc,
		ctx.NewCast(
			nil,
			ctx.NewCall(
				nil,
				dlsymFunc,
				[]*gccjit.Rvalue{
					handle.AsRvalue(),
					ctx.NewStringLiteral("ClearBackground"),
				},
			),
			clearBackgroundFuncType,
		),
	)

	drawTextFunc := mainFunc.NewLocal(nil, drawTextFuncType, "DrawText")
	handleLoaded.AddAssignment(
		nil,
		drawTextFunc,
		ctx.NewCast(
			nil,
			ctx.NewCall(
				nil,
				dlsymFunc,
				[]*gccjit.Rvalue{
					handle.AsRvalue(),
					ctx.NewStringLiteral("DrawText"),
				},
			),
			drawTextFuncType,
		),
	)

	closeWindowFunc := mainFunc.NewLocal(nil, closeWindowFuncType, "CloseWindow")
	handleLoaded.AddAssignment(
		nil,
		closeWindowFunc,
		ctx.NewCast(
			nil,
			ctx.NewCall(
				nil,
				dlsymFunc,
				[]*gccjit.Rvalue{
					handle.AsRvalue(),
					ctx.NewStringLiteral("CloseWindow"),
				},
			),
			closeWindowFuncType,
		),
	)

	windowShouldCloseFunc := mainFunc.NewLocal(nil, windowShouldCloseFuncType, "WindowShouldClose")
	handleLoaded.AddAssignment(
		nil,
		windowShouldCloseFunc,
		ctx.NewCast(
			nil,
			ctx.NewCall(
				nil,
				dlsymFunc,
				[]*gccjit.Rvalue{
					handle.AsRvalue(),
					ctx.NewStringLiteral("WindowShouldClose"),
				},
			),
			windowShouldCloseFuncType,
		),
	)

	lightGray := mainFunc.NewLocal(nil, colorType.AsType(), "lightGray")
	handleLoaded.AddAssignment(
		nil,
		lightGray.AccessField(nil, rField),
		ctx.NewRValueFromInt(charType, 200),
	)

	handleLoaded.AddAssignment(
		nil,
		lightGray.AccessField(nil, gField),
		ctx.NewRValueFromInt(charType, 200),
	)

	handleLoaded.AddAssignment(
		nil,
		lightGray.AccessField(nil, bField),
		ctx.NewRValueFromInt(charType, 200),
	)

	handleLoaded.AddAssignment(
		nil,
		lightGray.AccessField(nil, aField),
		ctx.NewRValueFromInt(charType, 255),
	)

	rayWhite := mainFunc.NewLocal(nil, colorType.AsType(), "rayWhite")
	handleLoaded.AddAssignment(
		nil,
		rayWhite.AccessField(nil, rField),
		ctx.NewRValueFromInt(charType, 245),
	)

	handleLoaded.AddAssignment(
		nil,
		rayWhite.AccessField(nil, gField),
		ctx.NewRValueFromInt(charType, 245),
	)

	handleLoaded.AddAssignment(
		nil,
		rayWhite.AccessField(nil, bField),
		ctx.NewRValueFromInt(charType, 245),
	)

	handleLoaded.AddAssignment(
		nil,
		rayWhite.AccessField(nil, aField),
		ctx.NewRValueFromInt(charType, 255),
	)

	handleLoaded.AddEval(
		nil,
		ctx.NewCallThroughPtr(
			nil,
			initWindowFunc.AsRvalue(),
			[]*gccjit.Rvalue{
				ctx.NewRValueFromInt(intType, 800),
				ctx.NewRValueFromInt(intType, 450),
				ctx.NewStringLiteral("raylib [core] example - basic window"),
			},
		),
	)

	handleLoaded.AddEval(
		nil,
		ctx.NewCallThroughPtr(
			nil,
			setTargetFPSFunc.AsRvalue(),
			[]*gccjit.Rvalue{
				ctx.NewRValueFromInt(intType, 60),
			},
		),
	)

	gameLoopCheck := mainFunc.NewBlock("gameLoopCheck")
	gameLoop := mainFunc.NewBlock("gameLoop")
	gameExit := mainFunc.NewBlock("gameExit")

	handleLoaded.EndWithJump(nil, gameLoopCheck)
	gameLoopCheck.EndWithConditional(
		nil,
		ctx.NewCallThroughPtr(
			nil,
			windowShouldCloseFunc.AsRvalue(),
			[]*gccjit.Rvalue{},
		),
		gameExit,
		gameLoop,
	)

	gameLoop.AddEval(
		nil,
		ctx.NewCallThroughPtr(
			nil,
			beginDrawingFunc.AsRvalue(),
			[]*gccjit.Rvalue{},
		),
	)

	gameLoop.AddEval(
		nil,
		ctx.NewCallThroughPtr(
			nil,
			clearBackgroundFunc.AsRvalue(),
			[]*gccjit.Rvalue{
				rayWhite.AsRvalue(),
			},
		),
	)
	gameLoop.AddEval(
		nil,
		ctx.NewCallThroughPtr(
			nil,
			drawTextFunc.AsRvalue(),
			[]*gccjit.Rvalue{
				ctx.NewStringLiteral("Congrats! You created your first window!"),
				ctx.NewRValueFromInt(intType, 190),
				ctx.NewRValueFromInt(intType, 200),
				ctx.NewRValueFromInt(intType, 20),
				lightGray.AsRvalue(),
			},
		),
	)

	gameLoop.AddEval(
		nil,
		ctx.NewCallThroughPtr(
			nil,
			endDrawingFunc.AsRvalue(),
			[]*gccjit.Rvalue{},
		),
	)

	gameLoop.EndWithJump(nil, gameLoopCheck)

	gameExit.AddEval(
		nil,
		ctx.NewCallThroughPtr(
			nil,
			closeWindowFunc.AsRvalue(),
			[]*gccjit.Rvalue{},
		),
	)

	gameExit.AddEval(
		nil,
		ctx.NewCall(
			nil,
			dlCloseFunc,
			[]*gccjit.Rvalue{
				handle.AsRvalue(),
			},
		),
	)

	gameExit.EndWithReturn(nil, ctx.NewRValueFromInt(intType, 0))

	// ctx.DumpToFile("./test.txt", false)
	ctx.CompileToFile(gccjit.OUTPUT_KIND_EXECUTABLE, "./a.out")
}
