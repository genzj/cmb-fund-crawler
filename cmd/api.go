package cmd

import (
	"github.com/genzj/cmb-fund-crawler/api"
	"github.com/mkideal/cli"
)

// crawl command
type apiT struct {
	cli.Helper
}

var apiCmd = &cli.Command{
	Name: "api",
	Desc: "run api server",
	Argv: func() interface{} { return new(apiT) },
	Fn: func(ctx *cli.Context) error {
		api.StartAPIServer()
		return nil
	},
}
