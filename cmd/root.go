package cmd

import (
	"fmt"
	"github.com/genzj/cmb-fund-crawler/db"
	"github.com/mkideal/cli"
	"os"
)

var help = cli.HelpCommand("display help information")

type rootT struct {
	cli.Helper
}

var root = &cli.Command{
	Desc: "crawler for funds with web gui",
	Argv: func() interface{} { return new(rootT) },
	Fn: func(ctx *cli.Context) error {
		fmt.Println(ctx.Usage())
		return nil
	},
}

func Run() {
	defer func() {
		db.GetDatabaseInstance().Close()
		fmt.Println("database closed")
	}()
	if err := cli.Root(
		root,
		cli.Tree(help),
		cli.Tree(crawlCmd),
		cli.Tree(apiCmd),
		cli.Tree(userCmd),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
