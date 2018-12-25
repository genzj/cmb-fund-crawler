package crawl

import (
	"github.com/labstack/gommon/log"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type FundDetail struct {
	Name               string
	NetValue           decimal.NullDecimal
	CumulativeNetValue decimal.NullDecimal
	ChangePCT          decimal.NullDecimal
	ReturnOneDay       decimal.NullDecimal
	UpdateTime         time.Time
}

func NewFundDetailFromString(name, netValue, cumulativeNetValue, changePCT, returnOneDay, updateTime string) *FundDetail {
	ans := &FundDetail{}
	ans.Name = strings.TrimSpace(name)
	if d, err := decimal.NewFromString(netValue); err != nil {
		log.Debugf("invalid netValue data: %s\n", netValue)
		return nil
	} else {
		ans.NetValue.Decimal = d
		ans.NetValue.Valid = true
	}
	if d, err := decimal.NewFromString(cumulativeNetValue); err != nil {
		log.Debugf("invalid cumulativeNetValue data: %s\n", cumulativeNetValue)
		return nil
	} else {
		ans.CumulativeNetValue.Decimal = d
		ans.CumulativeNetValue.Valid = true
	}
	if d, err := decimal.NewFromString(changePCT); err != nil {
		log.Debugf("invalid changePCT data: %s\n", changePCT)
		ans.ChangePCT.Valid = false
	} else {
		ans.ChangePCT.Decimal = d
		ans.ChangePCT.Valid = true
	}
	if d, err := decimal.NewFromString(returnOneDay); err != nil {
		log.Debugf("invalid returnOneDay data: %s\n", returnOneDay)
		ans.ReturnOneDay.Valid = false
	} else {
		ans.ReturnOneDay.Decimal = d
		ans.ReturnOneDay.Valid = true
	}
	if t, err := time.Parse("2006/01/02", updateTime); err != nil {
		log.Debugf("invalid updateTime data: %s\n", updateTime)

		return nil
	} else {
		ans.UpdateTime = t
	}
	return ans
}

type Crawl interface {
	Crawl(id string) (*FundDetail, error)
}
