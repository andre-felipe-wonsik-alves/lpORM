package dialect

import "github.com/andre-felipe-wonsik-alves/lpORM/src/schema"

type Dialect interface {
	SQLType(f schema.Field) string
	Quote(ident string) string
}
