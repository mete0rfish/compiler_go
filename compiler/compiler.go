package compiler

import (
	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

type Compiler struct {
	instructions code.Instructions // 컴파일된 바이트코드 명령어
	constants    []object.Object    // 상수 풀
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

// Compile 메서드는 주어진 AST 노드를 컴파일합니다.
// 현재는 아무 작업도 하지 않고 nil을 반환합니다.
func (c *Compiler) Compile(node ast.Node) error {
	return nil
}

// Bytecode 메서드는 컴파일된 바이트코드 명령어와 상수 풀을 포함하는 Bytecode 객체를 반환합니다.
// 이렇게 만들어진 Bytecode 객체는 가상 머신에서 실행할 수 있습니다.
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
