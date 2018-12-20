package crawl

import (
	"log"

	"github.com/andybalholm/cascadia"
	"github.com/genzj/cmb-fund-crawler/util"
	"github.com/levigross/grequests"
	"golang.org/x/net/html"
)

type CmbFundCrawlOption struct {
	baseURL string
}

type CmbFundCrawl struct {
	option CmbFundCrawlOption
}

var defaultCmbFundOption = CmbFundCrawlOption{
	baseURL: "http://fund.cmbchina.com/FundPages/OpenFund/OpenFundDetail.aspx",
}

var (
	nameSelector               = cascadia.MustCompile(".fundDetail_NameRow")
	netValueSelector           = cascadia.MustCompile("#ctl00_FundContentPlace_ucFundSummary_NetValue")
	cumulativeNetValueSelector = cascadia.MustCompile("#ctl00_FundContentPlace_ucFundSummary_CumulativeNetValue")
	changePCTSelector          = cascadia.MustCompile("#ctl00_FundContentPlace_ucFundSummary_ChangePCT")
	returnOneDaySelector       = cascadia.MustCompile("#ctl00_FundContentPlace_ucFundSummary_RETURN1DAY")
	updateTimeSelector         = cascadia.MustCompile("#ctl00_FundContentPlace_ucFundSummary_UpdateTime")
)

func NewCmbFundCrawl(option *CmbFundCrawlOption) *CmbFundCrawl {
	if option == nil {
		option = &defaultCmbFundOption
	}
	return &CmbFundCrawl{*option}
}

func (c CmbFundCrawl) Crawl(id string) (*FundDetail, error) {
	resp, err := grequests.Get(c.option.baseURL, &grequests.RequestOptions{
		Params: map[string]string{
			"FundID": id,
		},
		Headers: map[string]string{
			"Refer": "http://fund.cmbchina.com/FundPages/OpenFund/FundNetValue.aspx?Channel=OpenFund",
		},
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36",
	})

	if err != nil {
		log.Printf("ERROR Unable to make request: %s\n", err)
		return nil, err
	}

	if !resp.Ok {
		log.Printf("ERROR HTTP request failed: %s\n", resp.StatusCode)
		return nil, err
	}

	defer resp.Close()
	doc, err := html.Parse(resp.RawResponse.Body)
	if err != nil {
		log.Printf("ERROR Unable to parse HTML document: %s\n", err)
	}

	name := nameSelector.MatchFirst(doc)
	netValue := netValueSelector.MatchFirst(doc)
	cumulativeNetValue := cumulativeNetValueSelector.MatchFirst(doc)
	changePCT := changePCTSelector.MatchFirst(doc)
	returnOneDay := returnOneDaySelector.MatchFirst(doc)
	updateTime := updateTimeSelector.MatchFirst(doc)

	if !util.All(func(i interface{}) bool {
		d, ok := i.(*html.Node)
		if !ok {
			log.Println("DEBUG type mismatch")
			return false
		}
		if d == nil {
			log.Println("DEBUG nil node")
			return false
		}
		return true
	}, name, netValue, cumulativeNetValue, changePCT, returnOneDay, updateTime) {
		log.Printf("ERROR not all selector match content\n")
		log.Printf("DEBUG %v, %v, %v, %v, %v, %v\n", name, netValue, cumulativeNetValue, changePCT, returnOneDay, updateTime)

		return nil, nil
	}

	ans := NewFundDetailFromString(
		name.FirstChild.Data, netValue.FirstChild.Data, cumulativeNetValue.FirstChild.Data,
		changePCT.FirstChild.Data, returnOneDay.FirstChild.Data, updateTime.FirstChild.Data,
	)

	log.Printf("DEBUG crawled fund detail: %v", ans)

	return ans, nil
}
