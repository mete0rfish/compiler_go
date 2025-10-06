package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

// testIntegerObject는 주어진 object.Object가 특정 정수 값을 가진
// *object.Integer 타입인지 확인하는 헬퍼 함수입니다.
func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}

	return nil
}

// testConstants는 컴파일러가 생성한 상수 풀(actual)이
// 테스트 케이스에서 기대하는 상수(expected)와 일치하는지 확인하는 헬퍼 함수입니다.
func testConstants(
	t *testing.T,
	expected []interface{},
	actual []object.Object,
) error {
	// 상수 풀의 개수가 일치하는지 먼저 확인합니다.
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d",
			len(actual), len(expected))
	}

	// 각 상수를 순회하며 값을 비교합니다.
	for i, constant := range expected {
		switch constant := constant.(type) {
		case int: // 기대값이 int일 경우
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}

	return nil
}

// parse는 입력 문자열을 받아 어휘 분석과 파싱을 거쳐
// AST의 루트 노드(*ast.Program)를 반환하는 유틸리티 함수입니다.
func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

// concatInstructions는 여러 개의 명령어 슬라이스(Instructions)를
// 하나의 단일 슬라이스로 합치는 유틸리티 함수입니다.
func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}

// compilerTestCase는 컴파일러 테스트를 위한 단일 테스트 케이스의 구조체입니다.
type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}       // 컴파일 후 기대되는 상수 풀
	expectedInstructions []code.Instructions // 컴파일 후 기대되는 명령어들
}

// runCompilerTests는 compilerTestCase 슬라이스를 받아
// 각 케이스에 대해 컴파일러를 실행하고 결과를 검증하는 메인 테스트 러너 함수입니다.
func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper() // 이 함수가 테스트 헬퍼 함수임을 명시합니다.

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()                // 새로운 컴파일러 인스턴스를 생성합니다.
		err := compiler.Compile(program) // AST를 컴파일합니다.
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode() // 컴파일된 바이트코드를 가져옵니다.

		// 생성된 명령어가 기대값과 일치하는지 확인합니다.
		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		// 생성된 상수 풀이 기대값과 일치하는지 확인합니다.
		err = testConstants(tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

// testInstructions는 컴파일러가 생성한 실제 명령어(actual)가
// 테스트 케이스에서 기대하는 명령어(expected)와 정확히 일치하는지 확인하는 헬퍼 함수입니다.
func testInstructions(
	expected []code.Instructions,
	actual code.Instructions,
) error {
	concatted := concatInstructions(expected) // 기대 명령어들을 하나로 합칩니다.

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot =%q",
			concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot =%q",
				i, concatted, actual)
		}
	}

	return nil
}

// TestIngegerArithmetic는 정수 산술 연산에 대한 컴파일러의 동작을 테스트합니다.
func TestIngegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0), // 상수 1 (인덱스 0)
				code.Make(code.OpConstant, 1), // 상수 2 (인덱스 1)
				code.Make(code.OpAdd),         // 덧셈 명령어
			},
		},
	}

	runCompilerTests(t, tests)
}
