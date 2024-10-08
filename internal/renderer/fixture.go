package renderer

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/debugger84/sqlc-fixture/internal/imports"
	"github.com/debugger84/sqlc-fixture/internal/model"
	"github.com/debugger84/sqlc-fixture/internal/opts"
	"github.com/iancoleman/strcase"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
	"go/format"
	"text/template"
)

type FixtureTplData struct {
	Struct               model.Struct
	Helper               *StructHelper
	Package              string
	PrimaryKeyColumnName string
	PrimaryKeyFieldType  string
	PrimaryKeyFieldName  string
	Imports              []imports.Import
}

type FixtureFactoryTplData struct {
	Structs      []model.Struct
	Package      string
	Imports      []imports.Import
	ModelPackage string
}

type FixtureRenderer struct {
	structs       []model.Struct
	loaderPackage string
	importer      *imports.ImportBuilder
	driver        opts.SQLDriver
}

func NewFixtureRenderer(
	structs []model.Struct,
	options *opts.Options,
	importer *imports.ImportBuilder,
) *FixtureRenderer {
	return &FixtureRenderer{
		structs:       structs,
		loaderPackage: options.Package,
		importer:      importer,
		driver:        options.Driver(),
	}
}

func (r *FixtureRenderer) Render() ([]*plugin.File, error) {
	if len(r.structs) == 0 {
		return nil, nil
	}
	funcMap := template.FuncMap{
		"lowerTitle": func(s string) string {
			title := sdk.LowerTitle(s)
			title = r.renameReservedWords(title)

			return title
		},
	}
	tmpl := template.Must(
		template.New("fixture.tmpl").
			Funcs(funcMap).
			ParseFS(
				templates,
				"templates/fixture.tmpl",
			),
	)
	files := make([]*plugin.File, 0)
	loaderImporter := r.importer.
		AddWithoutAlias("testing").
		AddWithoutAlias("context")

	for _, s := range r.structs {
		if !s.HasPrimaryKey() {
			continue
		}
		file, err := r.renderFixture(tmpl, s, loaderImporter)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (r *FixtureRenderer) renameReservedWords(title string) string {
	if title == "type" {
		return "typ"
	}
	if title == "range" {
		return "rng"
	}
	if title == "map" {
		return "mp"
	}
	if title == "string" {
		return "str"
	}
	if title == "interface" {
		return "iface"
	}
	if title == "select" {
		return "sel"
	}
	if title == "default" {
		return "def"
	}
	if title == "case" {
		return "c"
	}
	if title == "switch" {
		return "sw"
	}
	if title == "for" {
		return "f"
	}
	if title == "func" {
		return "fn"
	}
	if title == "return" {
		return "ret"
	}
	if title == "package" {
		return "pkg"
	}
	if title == "import" {
		return "imp"
	}
	if title == "var" {
		return "v"
	}
	if title == "const" {
		return "cst"
	}
	if title == "struct" {
		return "st"
	}
	return title
}

func (r *FixtureRenderer) renderFixture(
	tmpl *template.Template,
	s model.Struct,
	importer *imports.ImportBuilder,
) (*plugin.File, error) {
	var pkField model.Field
	for _, f := range s.Fields() {
		if f.IsPrimaryKey() {
			pkField = f
			break
		}
	}

	tctx := FixtureTplData{
		Struct:               s,
		Helper:               NewStructHelper(s, r.driver),
		Package:              r.loaderPackage,
		PrimaryKeyColumnName: pkField.DBName(),
		PrimaryKeyFieldType:  pkField.Type().TypeWithPackage(),
		PrimaryKeyFieldName:  pkField.Name(),
		Imports: importer.
			ImportContainer(&s).
			Build(),
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err := tmpl.ExecuteTemplate(w, "fixture.tmpl", &tctx)
	w.Flush()
	if err != nil {
		return nil, err
	}
	code, err := format.Source(b.Bytes())
	if err != nil {
		fmt.Println(b.String())
		return nil, fmt.Errorf("source error: %w", err)
	}
	filename := fmt.Sprintf("%s_loader.go", strcase.ToSnake(s.Type().TypeName()))
	if r.loaderPackage != s.Type().PackageName() {
		filename = fmt.Sprintf("%s/%s.go", r.loaderPackage, strcase.ToSnake(s.Type().TypeName()))
	}
	file := &plugin.File{
		Name:     filename,
		Contents: code,
	}
	return file, nil
}
