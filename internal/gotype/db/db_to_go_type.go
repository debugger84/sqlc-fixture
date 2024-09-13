package db

import (
	"fmt"
	"github.com/debugger84/sqlc-fixture/internal/gotype"
	"github.com/debugger84/sqlc-fixture/internal/opts"
	"github.com/debugger84/sqlc-fixture/internal/sqltype"
)

func NewDbTOGoTypeTransformer(
	engine opts.SQLEngine,
	customTypes []sqltype.CustomType,
	options *opts.Options,
) (gotype.DbTOGoTypeTransformer, error) {
	var typeTransformer gotype.DbTOGoTypeTransformer
	switch engine {
	case opts.SQLEngineMySQL:
		return NewMysqlTypeTransformer(customTypes), nil
	case opts.SQLEngineSQLite:
		return NewSqlLiteTypeTransformer(options, customTypes), nil
	case opts.SQLEnginePostgresql:
		return NewPostgresqlTypeTransformer(options, customTypes), nil
	}
	return typeTransformer, fmt.Errorf("unsupported sql engine %s", engine)
}
