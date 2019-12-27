package buildcmd

const mainBuildTmpl = `
package main

import (
	"context"
	"log"
	"os"

  {{if .WithFallthrough -}}
	appcli "{{.Info.Module.Path}}/cli"
	{{end -}}

	{{range $imp := .Imports -}}
	_ "{{$imp}}"
	{{end -}}

	"github.com/gobuffalo/buffalo-cli/built"
	"github.com/gobuffalo/buffalo-cli/cli"
	"github.com/gobuffalo/buffalo/runtime"
)

func main() {
	err := func() error {
		cb, err := cli.New()
		if err != nil {
			return err
		}

		b := &built.App{
			Plugger:      cb,
			BuildTime:    {{.BuildTime}},
			BuildVersion: {{.BuildVersion}},
		  {{if .WithFallthrough -}}
			Fallthrough:  appcli.Buffalo,
			{{else -}}
			Fallthrough:  cb.Main,
	    {{end -}}
			OriginalMain: originalMain,
		}

		if err := setBuildInfo(b); err != nil {
			return err
		}

		ctx := context.Background()
		return b.Main(ctx, os.Args[1:])
	}()

	if err != nil {
		log.Fatal(err)
	}
}

func setBuildInfo(b *built.App) error {
	t, err := time.Parse(time.RFC3339, b.BuildTime)
	if err != nil {
		t = time.Now()
	}
	runtime.SetBuild(runtime.BuildInfo{
		Version: b.BuildVersion,
		Time:    t,
	})
	return nil
}
`
