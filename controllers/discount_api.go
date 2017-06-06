package controllers

import (
	"net/http"
	"offer/factory"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
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
	var v SearchInput
	if err := c.Bind(&v); err != nil {
		return c.JSON(http.StatusBadRequest, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}

	factory.Logger(c.Request().Context()).WithFields(logrus.Fields{
		"sortby":         v.Sortby,
		"order":          v.Order,
		"maxResultCount": v.MaxResultCount,
		"skipCount":      v.SkipCount,
	}).Info("SearchInput")

	totalCount, items, err := factory.DB(c.Request().Context()).GetAllDiscount(nil, v.Sortby, v.Order, v.SkipCount, v.MaxResultCount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}
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
	if _, err := factory.DB(c.Request().Context()).AddDiscount(discount); err != nil {
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
	v, err := factory.DB(c.Request().Context()).GetDiscountById(id)
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
	if err := factory.DB(c.Request().Context()).UpdateDiscountById(discount); err != nil {
		return c.JSON(http.StatusInternalServerError, ApiResult{
			Error: ApiError{Message: err.Error()},
		})
	}
	return c.JSON(http.StatusOK, ApiResult{
		Success: true,
		Result:  discount,
	})
}
