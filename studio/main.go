package main

import (
	"studio/cli"
	"studio/errtype"
	"studio/params"
	"studio/web"
)

func main() {
	p, err := params.ParseParams()
	if err != nil {
		errtype.ErrorHandler(err)
	}

	switch p.IType {
	case params.CLI:
		err = cli.New().Run(p.DBPath)
	case params.Web:
		err = web.New().Run(p.DBPath)
	}

	if err != nil {
		errtype.ErrorHandler(err)
	}
}
