package renderer

import (
	"fmt"
	"github.com/debugger84/sqlc-fixture/internal/model"
	"github.com/debugger84/sqlc-fixture/internal/opts"
	"strings"
)

type StructHelper struct {
	s      model.Struct
	driver opts.SQLDriver
}

func (h *StructHelper) ColumnNames() string {
	out := ""
	for i, f := range h.s.Fields() {
		if i > 0 {
			out += ", "
		}
		out += f.DBName()
	}
	return out
}

func (h *StructHelper) ColumnPlaceholders() string {
	out := ""
	for i, _ := range h.s.Fields() {
		if i > 0 {
			out += ", "
		}
		if h.driver.IsPGX() || h.driver.IsLibPQ() {
			out += fmt.Sprintf("$%d", i+1)
		} else {
			out += "?"
		}
	}
	return out
}

func (h *StructHelper) UpdateSql() string {
	out := ""
	updatedFields := make([]string, 0, len(h.s.Fields()))
	whereClause := ""
	for i, field := range h.s.Fields() {
		if field.IsPrimaryKey() {
			whereClause = fmt.Sprintf("%s = ?", field.DBName())
			if h.driver.IsPGX() || h.driver.IsLibPQ() {
				whereClause = fmt.Sprintf("%s = $%d", field.DBName(), i+1)
			}
		} else {
			updatedField := fmt.Sprintf("%s = ?", field.DBName())
			if h.driver.IsPGX() || h.driver.IsLibPQ() {
				updatedField = fmt.Sprintf("%s = $%d", field.DBName(), i+1)
			}
			updatedFields = append(updatedFields, updatedField)
		}
	}
	out = fmt.Sprintf(
		"UPDATE %s SET \n            %s\n        WHERE %s",
		h.s.FullTableName(),
		strings.Join(updatedFields, ",\n            "),
		whereClause,
	)
	return out
}

func NewStructHelper(s model.Struct, driver opts.SQLDriver) *StructHelper {
	return &StructHelper{s: s, driver: driver}
}
