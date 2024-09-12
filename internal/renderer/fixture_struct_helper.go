package renderer

import (
	"fmt"
	"github.com/debugger84/sqlc-fixture/internal/model"
	"github.com/debugger84/sqlc-fixture/internal/opts"
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

func NewStructHelper(s model.Struct, driver opts.SQLDriver) *StructHelper {
	return &StructHelper{s: s, driver: driver}
}
