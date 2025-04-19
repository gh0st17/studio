package main

import (
	"github.com/gh0st17/studio/cli"
	"github.com/gh0st17/studio/errtype"
	"github.com/gh0st17/studio/params"
	"github.com/gh0st17/studio/ui"
	"github.com/gh0st17/studio/web"
)

func main() {
	p, err := params.ParseParams()
	if err != nil {
		errtype.ErrorHandler(err)
	}

	var ui ui.UI

	switch p.IType {
	case params.CLI:
		if cli, err := cli.New(p.PgSqlSoc); err != nil {
			errtype.ErrorHandler(err)
		} else {
			if p.Reg {
				login := cli.Login()
				cli.Registration(login)
			}
			ui = cli
		}
	case params.Web:
		if web, err := web.New(p.PgSqlSoc, p.RedisSoc, p.HttpSoc); err != nil {
			errtype.ErrorHandler(err)
		} else {
			ui = web
		}
	}

	if err = ui.Run(); err != nil {
		errtype.ErrorHandler(err)
	}
}
