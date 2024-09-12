package internal

import (
	"context"
	"github.com/debugger84/sqlc-fixture/internal/imports"
	"github.com/debugger84/sqlc-fixture/internal/model"
	"github.com/debugger84/sqlc-fixture/internal/opts"
	"github.com/debugger84/sqlc-fixture/internal/renderer"
	"github.com/debugger84/sqlc-fixture/internal/sqltype"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

func Generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	options, err := opts.Parse(req)
	if err != nil {
		return nil, err
	}

	if err := opts.ValidateOpts(options); err != nil {
		return nil, err
	}

	if options.DefaultSchema != "" {
		req.Catalog.DefaultSchema = options.DefaultSchema
	}
	customTypes := sqltype.NewCustomTypes(req.Catalog.Schemas, options)
	structs := model.BuildStructs(req, options, customTypes)

	importer := imports.NewImportBuilder(options)

	loaderRendered := renderer.NewFixtureRenderer(structs, options, importer)

	files := make([]*plugin.File, 0)
	loaderFiles, err := loaderRendered.Render()
	if err != nil {
		return nil, err
	}
	files = append(files, loaderFiles...)

	return &plugin.GenerateResponse{
		Files: files,
	}, nil
}
