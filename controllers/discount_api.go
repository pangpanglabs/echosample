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
		return c.JSON(http.StatusBadRequest, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
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
		return c.JSON(http.StatusInternalServerError, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}
	tracer.LogFields(
		log.String("Action", "Search From DB"),
		log.Int64("TotalCount", totalCount),
	)
	return c.JSON(http.StatusOK, ApiResult{
		Success: true,
		Result: ArrayResult{
			TotalCount: totalCount,
			Items:      items,
		},
	})
}

func (DiscountApiController) Create(c echo.Context) error {
	var v DiscountInput
	if err := c.Bind(&v); err != nil {
		return c.JSON(http.StatusBadRequest, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}
	if err := c.Validate(&v); err != nil {
		return c.JSON(http.StatusBadRequest, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}
	discount, err := v.ToModel()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}
	if _, err := discount.Create(c.Request().Context()); err != nil {
		return c.JSON(http.StatusInternalServerError, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}
	return c.JSON(http.StatusOK, ApiResult{
		Success: true,
		Result:  discount,
	})
}

func (DiscountApiController) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}
	v, err := models.Discount{}.GetById(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}
	if v == nil {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, ApiResult{
		Success: true,
		Result:  v,
	})
}

func (DiscountApiController) Update(c echo.Context) error {
	var v DiscountInput
	if err := c.Bind(&v); err != nil {
		return c.JSON(http.StatusBadRequest, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}
	if err := c.Validate(&v); err != nil {
		return c.JSON(http.StatusBadRequest, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}
	discount, err := v.ToModel()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}
	discount.Id = id
	if err := discount.Update(c.Request().Context()); err != nil {
		return c.JSON(http.StatusInternalServerError, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}
	return c.JSON(http.StatusOK, ApiResult{
		Success: true,
		Result:  discount,
	})
}
