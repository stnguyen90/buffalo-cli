package grifts

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/build"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/here"
)

var _ build.Importer = Importer{}
var _ plugins.Plugin = Importer{}

type Importer struct{}

func (Importer) PluginName() string {
	return "grifts/importer"
}

func (Importer) BuildImports(ctx context.Context, root string) ([]string, error) {
	dir := filepath.Join(root, "grifts")
	info, err := here.Dir(dir)
	if err != nil {
		return nil, nil
	}
	return []string{info.ImportPath}, nil
}
