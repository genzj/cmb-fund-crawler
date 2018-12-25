package cmd

import (
	"fmt"
	"github.com/genzj/cmb-fund-crawler/db"
	"github.com/labstack/gommon/log"
	"github.com/mkideal/cli"
	"os"
)

var help = cli.HelpCommand("display help information")

type RootT struct {
	cli.Helper
	Debug bool `cli:"debug,D" usage:"enable debug log output"`
}

var root = &cli.Command{
	Desc: "crawler for funds with web gui",
	Argv: func() interface{} { return new(RootT) },
	Fn: func(ctx *cli.Context) error {
		readGlobalFlag(ctx.Argv().(*RootT))
		fmt.Println(ctx.Usage())
		return nil
	},
}

func readGlobalFlag(t *RootT) {
	if t.Debug {
		log.SetLevel(log.DEBUG)
	}
}

func Run() {
	defer func() {
		if err := db.GetDatabaseInstance().Close(); err != nil {
			log.Errorf("cannot close database due to %s\n", err)
		} else {
			log.Debug("database closed")
		}
	}()
	if err := cli.Root(
		root,
		cli.Tree(help),
		cli.Tree(crawlCmd),
		cli.Tree(apiCmd),
		cli.Tree(userCmd),
	).Run(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
