package cmd

import (
	"fmt"
	crawl2 "github.com/genzj/cmb-fund-crawler/crawl"
	db2 "github.com/genzj/cmb-fund-crawler/db"
	"github.com/mkideal/cli"
)

func doCmbFundCrawl(id string) (*crawl2.FundDetail, error) {
	crawl := crawl2.NewCmbFundCrawl(nil)
	return crawl.Crawl(id)
}

type crawlT struct {
	RootT
	Save bool     `cli:"save,s" usage:"save crawled record to database"`
	ID   []string `cli:"*id,i" usage:"ID of fund to crawl"`
}

var crawlCmd = &cli.Command{
	Name: "crawl",
	Desc: "run crawler and exit",
	Argv: func() interface{} { return new(crawlT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*crawlT)
		readGlobalFlag(&argv.RootT)

		for _, id := range argv.ID {
			if detail, err := doCmbFundCrawl(id); err != nil {
				return err
			} else if argv.Save {
				ctx.JSONIndentln(detail, "", "  ")
				db := db2.GetDatabaseInstance()
				err := db2.SaveFundValueRecord(db, "cmb", id, *detail)
				if err != nil {
					fmt.Printf("ERROR %s\n", err)
				}
				_ = db2.IterateFundValue(db, "cmb", id, func(value db2.FundValue) error {
					ctx.JSONIndentln(value, "", "  ")
					return nil
				})
			} else {
				ctx.JSONIndentln(detail, "", "  ")
			}
		}

		return nil
	},
}
