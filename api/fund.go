package api

import (
	"github.com/genzj/cmb-fund-crawler/db"
	"github.com/labstack/echo"
	"net/http"
)

func registerFundRoutes(e *echo.Group) {
	e.GET("/fund/:org/:id", fundList)
}

func fundList(c echo.Context) error {
	var ans []db.FundValue
	dbi := db.GetDatabaseInstance()
	err := db.IterateFundValue(
		dbi,
		c.Param("org"),
		c.Param("id"),
		func(value db.FundValue) error {
			ans = append(ans, value)
			return nil
		},
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&map[string]string{
				"error": err.Error(),
			},
		)
	} else {
		return c.JSON(
			http.StatusOK,
			&map[string]interface{}{
				"funds": ans,
			},
		)
	}
}
