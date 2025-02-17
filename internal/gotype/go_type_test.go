package gotype_test

import (
	"github.com/debugger84/sqlc-fixture/internal/gotype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewGoType(t *testing.T) {
	t.Run(
		"base type", func(t *testing.T) {
			gt := gotype.NewGoType("int")
			require.NotNil(t, gt)
			assert.Equal(t, "int", gt.TypeName())
			assert.Equal(t, "", gt.PackageName())
			assert.False(t, gt.IsPointer())
		},
	)

	t.Run(
		"pointer type", func(t *testing.T) {
			gt := gotype.NewGoType("*int")
			require.NotNil(t, gt)
			assert.Equal(t, "int", gt.TypeName())
			assert.Equal(t, "", gt.PackageName())
			assert.True(t, gt.IsPointer())
		},
	)

	t.Run(
		"custom type without package", func(t *testing.T) {
			gt := gotype.NewGoType("Identity")
			require.NotNil(t, gt)
			assert.Equal(t, "Identity", gt.TypeName())
			assert.Equal(t, "", gt.PackageName())
			assert.False(t, gt.IsPointer())
		},
	)

	t.Run(
		"custom type with package", func(t *testing.T) {
			gt := gotype.NewGoType("test.Identity")
			require.NotNil(t, gt)
			assert.Equal(t, "Identity", gt.TypeName())
			assert.Equal(t, "test", gt.PackageName())
			assert.Equal(t, "", gt.Import().Path)
			assert.False(t, gt.IsPointer())
		},
	)

	t.Run(
		"custom type with full path of the local package", func(t *testing.T) {
			gt := gotype.NewGoType("myproj/test.Identity")
			require.NotNil(t, gt)
			assert.Equal(t, "Identity", gt.TypeName())
			assert.Equal(t, "test", gt.PackageName())
			assert.Equal(t, "myproj/test", gt.Import().Path)
			assert.False(t, gt.IsPointer())
		},
	)

	t.Run(
		"custom type with full path of the remote package", func(t *testing.T) {
			gt := gotype.NewGoType("github.com/myproj/test.Identity")
			require.NotNil(t, gt)
			assert.Equal(t, "Identity", gt.TypeName())
			assert.Equal(t, "test", gt.PackageName())
			assert.Equal(t, "github.com/myproj/test", gt.Import().Path)
			assert.Equal(t, "", gt.Import().Alias)
			assert.False(t, gt.IsPointer())
		},
	)

	t.Run(
		"array type", func(t *testing.T) {
			gt := gotype.NewGoType("[]int")
			require.NotNil(t, gt)
			assert.Equal(t, "int", gt.TypeName())
			assert.Equal(t, "", gt.PackageName())
			assert.False(t, gt.IsPointer())
			assert.True(t, gt.IsArray())
		},
	)
}
