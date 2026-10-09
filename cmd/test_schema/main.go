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

	//cfg := pq.Config{
	//	Host:           "localhost",
	//	Port:           5432,
	//	User:           "pqgo",
	//	ConnectTimeout: 5 * time.Second,
	//}

	a := &assembler.Assembler{}
	a.New(modelsList, "localhost", "appdb")
	a.Assemble()
}
