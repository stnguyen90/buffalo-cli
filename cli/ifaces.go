package cli

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
)

type Aliases = plugins.Aliases

// Command represents a plugin that can be
// used as a full sub-command. Like Go program's the
// `Main` method is called to run that command.
type Command interface {
	plugins.Plugin
	Main(ctx context.Context, root string, args []string) error
}