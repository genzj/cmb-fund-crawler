package cmd

import (
	"fmt"

	db2 "github.com/genzj/cmb-fund-crawler/db"

	"github.com/mkideal/cli"
)

// crawl command
type userT struct {
	RootT
	Username string `cli:"*username,u" usage:"username to manipulated (required)"`
	Add      bool   `cli:"add,a" usage:"add the specified user"`
	Show     bool   `cli:"show,s" usage:"show user information"`
}

var userCmd = &cli.Command{
	Name: "user",
	Desc: "offline user manipulation commands",
	Argv: func() interface{} { return new(userT) },
	Fn: func(ctx *cli.Context) error {
		arg := (ctx.Argv()).(*userT)
		readGlobalFlag(&arg.RootT)

		fmt.Printf("username %s\n", arg.Username)
		db := db2.GetDatabaseInstance()

		if arg.Add {
			if err := db2.CreateUser(db, arg.Username); err != nil {
				return err
			}
		} else if arg.Show {
			if userInfo, err := db2.GetUser(db, arg.Username); err != nil {
				return err
			} else {
				ctx.JSONIndentln(userInfo, "", "  ")
			}
		}
		return nil
	},
}
