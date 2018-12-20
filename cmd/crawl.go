package cmd

import (
	crawl2 "github.com/genzj/cmb-fund-crawler/crawl"
	"github.com/mkideal/cli"
)

func doCmbFundCrawl(id string) (*crawl2.FundDetail, error) {
	crawl := crawl2.NewCmbFundCrawl(nil)
	return crawl.Crawl(id)
	//crawl.Crawl("217011")
	//crawl.Crawl("485011")
	//crawl.Crawl("001868")
}

type crawlT struct {
	cli.Helper
	ID []string `cli:"*id,i" usage:"ID of fund to crawl"`
}

var crawlCmd = &cli.Command{
	Name: "crawl",
	Desc: "run crawler and exit",
	Argv: func() interface{} { return new(crawlT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*crawlT)

		for _, id := range argv.ID {
			if detail, err := doCmbFundCrawl(id); err != nil {
				return err
			} else {
				ctx.JSONIndentln(detail, "", "  ")
			}
		}

		return nil
	},
}
