package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const STACK_CAPACITY = 1024

type Word int64

var (
	ErrStackOverflow  = errors.New("stack overflow")
	ErrStackUnderflow = errors.New("stack underflow")
	ErrIllegalInst    = errors.New("illegal instruction")
	ErrDivideByZero   = errors.New("divide by zero")
)

type Bm struct {
	stack     []Word
	stackSize int
	ip Word
	halted bool
}

func (bm *Bm) Init() {
	bm.stack = make([]Word, STACK_CAPACITY)
	bm.stackSize = 0
	bm.ip = 0
	bm.halted = false
}

func (bm *Bm) ExecuteInst(inst Inst) error {
	switch inst.typ {
	case INST_PUSH:
		if bm.stackSize >= STACK_CAPACITY {
			return ErrStackOverflow
		}
		bm.stack[bm.stackSize] = inst.operand
		bm.stackSize += 1
	case INST_PLUS:
		if bm.stackSize < 2 {
			return ErrStackUnderflow
		}
		bm.stack[bm.stackSize-2] += bm.stack[bm.stackSize-1]
		bm.stackSize -= 1
	case INST_MINUS:
		if bm.stackSize < 2 {
			return ErrStackUnderflow
		}
		bm.stack[bm.stackSize-2] -= bm.stack[bm.stackSize-1]
		bm.stackSize -= 1
	case INST_MULT:
		if bm.stackSize < 2 {
			return ErrStackUnderflow
		}
		bm.stack[bm.stackSize-2] *= bm.stack[bm.stackSize-1]
		bm.stackSize -= 1
	case INST_DIV:
		if bm.stackSize < 2 {
			return ErrStackUnderflow
		} else if bm.stack[bm.stackSize-1] == 0 {
			return ErrDivideByZero
		}
		bm.stack[bm.stackSize-2] /= bm.stack[bm.stackSize-1]
		bm.stackSize -= 1
	case INST_HALT:
		bm.halted = true
	default:
		return ErrIllegalInst
	}
	return nil
}

func (bm *Bm) Dump(f io.Writer) {
	fmt.Fprintln(f, "Stack:")
	if bm.stackSize > 0 {
		for _, v := range bm.stack[:bm.stackSize] {
			fmt.Fprintf(f, "  %d\n", v)
		}
	} else {
		fmt.Fprintln(f, "  [empty]")
	}
}

type InstType int

const (
	INST_PUSH InstType = iota
	INST_PLUS
	INST_MINUS
	INST_MULT
	INST_DIV
	INST_JMP
	INST_HALT
)

var instTypeNames = [...]string{
	INST_PUSH:  "INST_PUSH",
	INST_PLUS:  "INST_PLUS",
	INST_MINUS: "INST_MINUS",
	INST_MULT:  "INST_MULT",
	INST_DIV:   "INST_DIV",
	INST_JMP:  "INST_JMP",
	INST_HALT:   "INST_HALT",
}

func InstTypeName(typ InstType) string {
	return instTypeNames[typ]
}

type Inst struct {
	typ     InstType
	operand Word
}

var program = []Inst{
	{typ: INST_PUSH, operand: 69},
	{typ: INST_PUSH, operand: 420},
	{typ: INST_PLUS},
	{typ: INST_PUSH, operand: 42},
	{typ: INST_MINUS},
	{typ: INST_PUSH, operand: 2},
	{typ: INST_MULT},
	{typ: INST_PUSH, operand: 4},
	{typ: INST_DIV},
	{typ: INST_HALT},
}

func main() {
	var bm Bm
	bm.Init()
	bm.Dump(os.Stdout)
	for !bm.halted {
		fmt.Fprintln(os.Stdout, InstTypeName(program[bm.ip].typ))
		err := bm.ExecuteInst(program[bm.ip])
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			bm.Dump(os.Stderr)
			os.Exit(1)
		}
		bm.Dump(os.Stdout)
	}
	bm.Dump(os.Stdout)
}
