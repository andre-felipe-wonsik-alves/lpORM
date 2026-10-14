package dialect

import (
	"fmt"
	"reflect"

	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
)

type Postgres struct{}

func (Postgres) Quote(ident string) string {
	return fmt.Sprintf(`"%s"`, ident)
}

func (Postgres) SQLType(f schema.Field) string {
	if f.Auto {
		return "BIGSERIAL"
	}

	switch f.GoType.Kind() {
	case reflect.Int, reflect.Int32:
		return "INTEGER"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Int8, reflect.Int16:
		return "SMALLINT"
	case reflect.Float32:
		return "REAL"
	case reflect.Float64:
		return "DOUBLE PRECISION"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.String:
		if f.Size > 0 {
			return fmt.Sprintf("VARCHAR(%d)", f.Size)
		}
		return "TEXT"
	default:
		return "TEXT"
	}
}
