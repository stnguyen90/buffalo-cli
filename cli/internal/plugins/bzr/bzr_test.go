package bzr

import (
	"context"
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Bzr_Generalities(t *testing.T) {
	r := require.New(t)
	b := Versioner{}

	r.Equal("bzr", b.PluginName(), "Name should be bzr")
	r.Equal("Provides bzr related hooks to Buffalo applications.", b.Description(), "Description does not match")
}

func Test_Bzr_BuildVersion(t *testing.T) {

	r := require.New(t)
	vr := &commandRunner{
		stdout: "123",
	}

	v := &Versioner{
		pluginsFn: func() []plugins.Plugin {
			return []plugins.Plugin{vr}
		},
	}

	s, err := v.BuildVersion(context.Background(), "")
	r.NoError(err)
	r.Equal(vr.stdout, s)
	r.NotNil(vr.cmd)

	r.Equal([]string{"bzr", "revno"}, vr.cmd.Args)
}
