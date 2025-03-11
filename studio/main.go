package main

import (
	"studio/cli"
	"studio/errtype"
	"studio/params"
	"studio/studio"
	"studio/ui"
	"studio/web"
)

func main() {
	p, err := params.ParseParams()
	if err != nil {
		errtype.ErrorHandler(err)
	}

	var ui ui.UI
	switch p.IType {
	case params.CLI:
		ui = &cli.CLI{}
	case params.Web:
		ui = &web.Web{}
	}

	studio, err := studio.New(ui)
	if err != nil {
		errtype.ErrorHandler(err)
	}

	if err = studio.Run(p.DBPath); err != nil {
		errtype.ErrorHandler(err)
	}
}
