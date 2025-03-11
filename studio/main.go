package main

import (
	"studio/basic_types"
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

	studio := studio.Studio{}
	if err = studio.Run(ui, &basic_types.SysAdmin{}); err != nil {
		errtype.ErrorHandler(err)
	}
}
