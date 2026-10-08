package main

import (
	"github.com/andre-felipe-wonsik-alves/lpORM/input"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/assembler"
)

func main() {
	modelsList := []any{
		input.Cliente{},
		input.Job{},
	}
	a := &assembler.Assembler{}
	a.New(modelsList)
	a.Assemble()
}
