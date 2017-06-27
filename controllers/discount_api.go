package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"

	"github.com/pangpanglabs/echosample/factory"
	"github.com/pangpanglabs/echosample/models"
)

type DiscountApiController struct {
}

func (c DiscountApiController) Init(g *echo.Group) {
	g.GET("", c.GetAll)
	g.POST("", c.Create)
	g.GET("/:id", c.GetOne)
	g.PUT("/:id", c.Update)
}
func (DiscountApiController) GetAll(c echo.Context) error {
	tracer := factory.Tracer(c.Request().Context())
	tracer.LogEvent("Start GetAll")

	var v SearchInput
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if v.MaxResultCount == 0 {
		v.MaxResultCount = DefaultMaxResultCount
	}
	tracer.LogFields(
		log.String("Action", "Bind Request"),
		log.Int("MaxResultCount", v.MaxResultCount),
		log.Int("SkipCount", v.SkipCount),
	)

	factory.Logger(c.Request().Context()).WithFields(logrus.Fields{
		"sortby":         v.Sortby,
		"order":          v.Order,
		"maxResultCount": v.MaxResultCount,
		"skipCount":      v.SkipCount,
	}).Info("SearchInput")

	totalCount, items, err := models.Discount{}.GetAll(c.Request().Context(), v.Sortby, v.Order, v.SkipCount, v.MaxResultCount)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	tracer.LogFields(
		log.String("Action", "Search From DB"),
		log.Int64("TotalCount", totalCount),
	)
	return ReturnApiSucc(c, http.StatusOK, ArrayResult{
		TotalCount: totalCount,
		Items:      items,
	})
}

func (DiscountApiController) Create(c echo.Context) error {
	var v DiscountInput
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if err := c.Validate(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	discount, err := v.ToModel()
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if _, err := discount.Create(c.Request().Context()); err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	return ReturnApiSucc(c, http.StatusOK, discount)
}

func (DiscountApiController) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	v, err := models.Discount{}.GetById(c.Request().Context(), id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if v == nil {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}
	return ReturnApiSucc(c, http.StatusOK, v)
}

func (DiscountApiController) Update(c echo.Context) error {
	var v DiscountInput
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if err := c.Validate(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	discount, err := v.ToModel()
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	discount.Id = id
	if err := discount.Update(c.Request().Context()); err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	return ReturnApiSucc(c, http.StatusOK, discount)
}
