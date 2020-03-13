package controllers

import (
	"net/http"
	"strconv"

	"github.com/pangpanglabs/goutils/behaviorlog"

	"github.com/pangpanglabs/echosample/factory"
	"github.com/pangpanglabs/echosample/models"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
	"github.com/sirupsen/logrus"
)

type DiscountApiController struct {
}

func (c DiscountApiController) Init(g echoswagger.ApiGroup) {
	g.GET("", c.GetAll).AddParamQueryNested(SearchInput{})
	g.POST("", c.Create).AddParamBody(DiscountInput{}, "body", "", true)
	g.GET("/:id", c.GetOne).AddParamPath(0, "id", "")
	g.PUT("/:id", c.Update).AddParamPath(0, "id", "").AddParamBody(DiscountInput{}, "body", "", true)
}

func (DiscountApiController) GetAll(c echo.Context) error {
	var v SearchInput
	if err := c.Bind(&v); err != nil {
		return renderFail(c, factory.ErrorParameter.New(err))
	}
	if v.MaxResultCount == 0 {
		v.MaxResultCount = DefaultMaxResultCount
	}

	// behavior log
	behaviorlog.FromCtx(c.Request().Context()).WithBizAttr("maxResultCount", v.MaxResultCount).Log("SearchDiscount")

	// console log
	factory.Logger(c.Request().Context()).WithFields(logrus.Fields{
		"sortby":         v.Sortby,
		"order":          v.Order,
		"maxResultCount": v.MaxResultCount,
		"skipCount":      v.SkipCount,
	}).Info("SearchStart")

	totalCount, items, err := models.Discount{}.GetAll(c.Request().Context(), v.Sortby, v.Order, v.SkipCount, v.MaxResultCount)
	if err != nil {
		return renderFail(c, err)
	}

	// behavior log
	behaviorlog.FromCtx(c.Request().Context()).
		WithCallURLInfo(http.MethodGet, "https://play.google.com/books", nil, 200).
		WithBizAttrs(map[string]interface{}{
			"totalCount": totalCount,
			"itemCount":  len(items),
		}).
		Log("SearchComplete")

	return renderSucc(c, http.StatusOK, ArrayResult{
		TotalCount: totalCount,
		Items:      items,
	})
}

func (DiscountApiController) Create(c echo.Context) error {
	var v DiscountInput
	if err := c.Bind(&v); err != nil {
		return renderFail(c, factory.ErrorParameter.New(err))
	}
	if err := c.Validate(&v); err != nil {
		return renderFail(c, factory.ErrorParameter.New(err))
	}
	discount, err := v.ToModel()
	if err != nil {
		return renderFail(c, factory.ErrorParameter.New(err))
	}
	if _, err := discount.Create(c.Request().Context()); err != nil {
		return renderFail(c, err)
	}
	return renderSucc(c, http.StatusOK, discount)
}

func (DiscountApiController) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return renderFail(c, factory.ErrorParameter.New(err))
	}
	v, err := models.Discount{}.GetById(c.Request().Context(), id)
	if err != nil {
		return renderFail(c, err)
	}
	if v == nil {
		return renderFail(c, factory.ErrorNotFound.New(nil))
	}
	return renderSucc(c, http.StatusOK, v)
}

func (DiscountApiController) Update(c echo.Context) error {
	var v DiscountInput
	if err := c.Bind(&v); err != nil {
		return renderFail(c, factory.ErrorParameter.New(err))
	}
	if err := c.Validate(&v); err != nil {
		return renderFail(c, factory.ErrorParameter.New(err))
	}
	discount, err := v.ToModel()
	if err != nil {
		return renderFail(c, factory.ErrorParameter.New(err))
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return renderFail(c, factory.ErrorParameter.New(err))
	}
	discount.Id = id
	if err := discount.Update(c.Request().Context()); err != nil {
		return renderFail(c, err)
	}
	return renderSucc(c, http.StatusOK, v)
}
