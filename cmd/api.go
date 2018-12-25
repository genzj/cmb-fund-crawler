package cmd

import (
	"github.com/genzj/cmb-fund-crawler/api"
	"github.com/mkideal/cli"
)

// crawl command
type apiT struct {
	RootT
}

var apiCmd = &cli.Command{
	Name: "api",
	Desc: "run api server",
	Argv: func() interface{} { return new(apiT) },
	Fn: func(ctx *cli.Context) error {
		readGlobalFlag(&ctx.Argv().(*apiT).RootT)
		api.StartAPIServer()
		return nil
	},
}
