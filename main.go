package main

import (
	"studio/cli"
	"studio/errtype"
	"studio/params"
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
		if cli, err := cli.New(); err != nil {
			errtype.ErrorHandler(err)
		} else {
			if p.Reg {
				login := cli.Login()
				cli.Registration(login)
			}
			ui = cli
		}
	case params.Web:
		if web, err := web.New(); err != nil {
			errtype.ErrorHandler(err)
		} else {
			ui = web
		}
	}

	if err = ui.Run(); err != nil {
		errtype.ErrorHandler(err)
	}
}
