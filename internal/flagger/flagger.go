package flagger

import (
	"flag"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugflag"
	"github.com/spf13/pflag"
)

func CleanPflags(p plugins.Plugin, pflags []*pflag.Flag) []*flag.Flag {
	flags := make([]*flag.Flag, len(pflags))
	for i, f := range pflags {
		flags[i] = &flag.Flag{
			Name:  f.Name,
			Usage: f.Usage,
			Value: f.Value,
		}
	}
	return plugflag.Clean(p, flags)
}

// SetToSlice takes a flag set and returns a slice
// of the flags
func SetToSlice(set *pflag.FlagSet) []*pflag.Flag {
	var flags []*pflag.Flag
	set.VisitAll(func(f *pflag.Flag) {
		flags = append(flags, f)
	})
	return flags
}
