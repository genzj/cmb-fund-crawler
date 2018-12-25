package cmd

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"syscall"

	db2 "github.com/genzj/cmb-fund-crawler/db"

	"github.com/mkideal/cli"
	"golang.org/x/crypto/ssh/terminal"
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
			fmt.Print("Enter Password: ")
			bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			if err != nil {
				return err
			}
			if err := db2.CreateUser(db, arg.Username, string(bytePassword)); err != nil {
				return err
			}
		} else if arg.Show {
			if userInfo, err := db2.GetUser(db, arg.Username); err != nil {
				return err
			} else {
				hashedPassword := userInfo.Password
				userInfo.Password = "----[HIDDEN FOR SAFE]----"
				ctx.JSONIndentln(userInfo, "", "  ")
				fmt.Print("Enter Password (optional, for validation, press Ctrl-C to abort): ")
				bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
				fmt.Println()
				if err != nil {
					return err
				}
				if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), bytePassword); err != nil {
					return err
				} else {
					fmt.Println("Password verified.")
				}
			}
		}
		return nil
	},
}
