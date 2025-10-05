package code

import (
	"encoding/binary"
	"fmt"
)

// Instruction 은 컴파일된 바이트코드 명령어의 기본 단위입니다.
// 실제로는 바이트 슬라이스입니다.
type Instruction []byte

// Opcode 는 명령어의 종류를 나타내는 한 바이트 숫자입니다.
type Opcode byte

// 컴파일러가 지원하는 Opcode 목록입니다.
const (
	// OpConstant 는 상수를 스택에 푸시하는 명령어입니다.
	OpConstant Opcode = iota
)

// Definition 은 각 Opcode에 대한 명세입니다.
type Definition struct {
	Name          string // 연산자 이름 (디버깅용)
	OperandWidths []int  // 각 피연산자가 차지하는 바이트 크기 배열
}

// Definitions 는 Opcode를 해당 Definition에 매핑합니다.
var Definitions = map[Opcode]*Definition{
	// OpConstant 명령어는 2바이트 크기의 피연산자 하나를 가집니다.
	// 이 피연산자는 상수 풀(constant pool)의 인덱스를 나타냅니다.
	OpConstant: {"OpConstant", []int{2}},
}

// Lookup 함수는 주어진 Opcode(바이트)에 해당하는 Definition을 찾습니다.
// 만약 정의되지 않은 Opcode라면 에러를 반환합니다.
func Lookup(op byte) (*Definition, error) {
	def, ok := Definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

// Make 함수는 Opcode와 피연산자들을 이용해 바이트코드 Instruction을 생성합니다.
func Make(op Opcode, operands ...int) []byte {
	def, ok := Definitions[op]
	if !ok {
		return []byte{}
	}

	// 명령어의 전체 길이를 계산합니다. (Opcode 1바이트 + 모든 피연산자의 길이)
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	// 계산된 길이만큼의 바이트 슬라이스를 생성합니다.
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op) // 첫 바이트에는 Opcode를 저장합니다.

	offset := 1 // 실제 피연산자 데이터가 시작될 위치
	// 각 피연산자를 순회하며 바이트 슬라이스에 씁니다.
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			// 2바이트 피연산자의 경우, Big-Endian 순서로 바이트를 씁니다.
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}
