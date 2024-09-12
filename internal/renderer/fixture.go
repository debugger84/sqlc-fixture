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
		"lowerTitle": sdk.LowerTitle,
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
		AddSqlDriver().
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
			Add(pkField.Type().Import()).
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
