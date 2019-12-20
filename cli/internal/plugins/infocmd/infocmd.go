package infocmd

import (
	"context"
	"io"
	"io/ioutil"
	"time"

	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/info"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/clara/genny/rx"
	"github.com/gobuffalo/genny"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = &InfoCmd{}
var _ plugins.PluginNeeder = &InfoCmd{}
var _ plugins.PluginScoper = &InfoCmd{}
var _ plugprint.Describer = &InfoCmd{}
var _ plugprint.FlagPrinter = &InfoCmd{}

type InfoCmd struct {
	pluginsFn plugins.PluginFeeder
	help      bool
}

func (ic *InfoCmd) WithPlugins(f plugins.PluginFeeder) {
	ic.pluginsFn = f
}

func (ic *InfoCmd) PrintFlags(w io.Writer) error {
	flags := ic.flagSet()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (ic *InfoCmd) Name() string {
	return "info"
}

func (ic *InfoCmd) Description() string {
	return "Print diagnostic information (useful for debugging)"
}

func (i InfoCmd) String() string {
	return i.Name()
}

// Info runs all of the plugins that implement the
// `Informer` interface in order.
func (ic *InfoCmd) plugins(ctx context.Context, args []string) error {
	for _, p := range ic.ScopedPlugins() {
		i, ok := p.(Informer)
		if !ok {
			continue
		}
		if err := i.Info(ctx, args); err != nil {
			return err
		}
	}
	return nil
}

func (ic *InfoCmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin

	if ic.pluginsFn == nil {
		return plugs
	}
	for _, p := range ic.pluginsFn() {
		if i, ok := p.(Informer); ok {
			plugs = append(plugs, i)
		}
	}

	return plugs
}

func (ic *InfoCmd) flagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet(ic.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVarP(&ic.help, "help", "h", false, "print this help")
	return flags
}

// Main implements the `buffalo info` command. Buffalo's checks
// are run first, then any plugins that implement plugins.Informer
// will be run in order at the end.
func (ic *InfoCmd) Main(ctx context.Context, args []string) error {

	flags := ic.flagSet()
	if err := flags.Parse(args); err != nil {
		return err
	}

	ioe := plugins.CtxIO(ctx)
	out := ioe.Stdout()

	if ic.help {
		return plugprint.Print(out, ic)
	}

	args = flags.Args()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	run := genny.WetRunner(ctx)

	opts := &rx.Options{
		Out: rx.NewWriter(out),
	}
	if err := run.WithNew(rx.New(opts)); err != nil {
		return err
	}

	iopts := &info.Options{
		Out: rx.NewWriter(out),
	}

	if err := run.WithNew(info.New(iopts)); err != nil {
		return err
	}

	if err := run.Run(); err != nil {
		return err
	}
	return ic.plugins(ctx, args)
}
