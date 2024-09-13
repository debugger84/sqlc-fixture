package sqltype

import (
	"github.com/debugger84/sqlc-fixture/internal/gotype"
	"github.com/debugger84/sqlc-fixture/internal/imports"
	"github.com/debugger84/sqlc-fixture/internal/naming"
	"github.com/debugger84/sqlc-fixture/internal/opts"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

type CustomTypeKind string

const (
	EnumType      CustomTypeKind = "enum"
	CompositeType CustomTypeKind = "composite"
)

type CustomType struct {
	GoType      gotype.GoType
	SqlTypeName string
	Schema      string
	Kind        CustomTypeKind
	IsNullable  bool
}

func NewCustomTypes(
	schemas []*plugin.Schema,
	options *opts.Options,
) []CustomType {
	normalizer := naming.NewNameNormalizer(options)
	driver := options.Driver()
	customTypes := make([]CustomType, 0)
	for _, schema := range schemas {
		if schema.Name == "pg_catalog" || schema.Name == "information_schema" {
			continue
		}

		for _, enum := range schema.Enums {
			baseName := normalizer.NormalizeSqlName(schema.Name, enum.Name)
			pckg := ""
			if options.ModelImport != "" {
				pckg = options.ModelImport + "."
			}
			nullGoType := gotype.NewGoType(pckg + "Null" + normalizer.NormalizeGoType(baseName))
			goType := gotype.NewGoType(pckg + normalizer.NormalizeGoType(baseName))
			customTypes = append(
				customTypes, CustomType{
					GoType:      *nullGoType,
					SqlTypeName: enum.Name,
					Schema:      schema.Name,
					Kind:        EnumType,
					IsNullable:  true,
				}, CustomType{
					GoType:      *goType,
					SqlTypeName: enum.Name,
					Schema:      schema.Name,
					Kind:        EnumType,
					IsNullable:  false,
				},
			)
		}

		emitPointersForNull := driver.IsPGX() && options.EmitPointersForNullTypes

		for _, ct := range schema.CompositeTypes {
			goType := gotype.NewGoType("string")
			nullGoType := gotype.NewGoType("sql.NullString")
			if emitPointersForNull {
				nullGoType = gotype.NewGoType("*string")
			} else {
				nullGoType.SetImport(imports.Import{Path: "database/sql"})
			}

			customTypes = append(
				customTypes, CustomType{
					GoType:      *goType,
					SqlTypeName: ct.Name,
					Schema:      schema.Name,
					Kind:        CompositeType,
					IsNullable:  false,
				}, CustomType{
					GoType:      *nullGoType,
					SqlTypeName: ct.Name,
					Schema:      schema.Name,
					Kind:        CompositeType,
					IsNullable:  true,
				},
			)
		}
	}

	return customTypes
}
