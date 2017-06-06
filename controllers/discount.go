package controllers

import (
	"fmt"
	"net/http"
	"offer/factory"
	"offer/models"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type DiscountController struct {
}

func (c DiscountController) Init(g *echo.Group) {
	g.GET("", c.GetAll)
	g.GET("/new", c.New)
	g.POST("", c.Create)
	g.GET("/:id", c.GetOne)
	g.GET("/:id/edit", c.Edit)
	g.POST("/:id", c.Update)
}

func (DiscountController) GetAll(c echo.Context) error {
	var v SearchInput
	if err := c.Bind(&v); err != nil {
		setFlashMessage(c, map[string]string{"warning": err.Error()})
	}
	if v.MaxResultCount == 0 {
		v.MaxResultCount = DefaultMaxResultCount
	}

	factory.Logger(c.Request().Context()).WithFields(logrus.Fields{
		"sortby":         v.Sortby,
		"order":          v.Order,
		"maxResultCount": v.MaxResultCount,
		"skipCount":      v.SkipCount,
	}).Info("SearchInput")

	totalCount, items, err := models.Discount{}.GetAll(c.Request().Context(), v.Sortby, v.Order, v.SkipCount, v.MaxResultCount)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "discount/index", map[string]interface{}{
		"TotalCount":     totalCount,
		"Discounts":      items,
		"MaxResultCount": v.MaxResultCount,
	})
}
func (DiscountController) New(c echo.Context) error {
	return c.Render(http.StatusOK, "discount/new", map[string]interface{}{
		FlashName: getFlashMessage(c),
		"Form":    &models.Discount{},
	})
}
func (DiscountController) Create(c echo.Context) error {
	var v DiscountInput
	if err := c.Bind(&v); err != nil {
		setFlashMessage(c, map[string]string{"error": err.Error()})
		return c.Redirect(http.StatusFound, "/discount/new")
	}
	if err := c.Validate(&v); err != nil {
		setFlashMessage(c, map[string]string{"error": err.Error()})
		return c.Redirect(http.StatusFound, "/discounts/new")
	}
	discount, err := v.ToModel()
	if err != nil {
		setFlashMessage(c, map[string]string{"error": err.Error()})
		return c.Redirect(http.StatusFound, "/discounts/new")
	}
	if _, err := discount.Create(c.Request().Context()); err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, fmt.Sprintf("/discounts/%d", discount.Id))
}
func (DiscountController) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.Render(http.StatusNotFound, "404", nil)
	}
	v, err := models.Discount{}.GetById(c.Request().Context(), id)
	if err != nil {
		return err
	}
	if v == nil {
		return c.Render(http.StatusNotFound, "404", nil)
	}
	return c.Render(http.StatusOK, "discount/show", map[string]interface{}{"Discount": v})
}

func (DiscountController) Edit(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.Render(http.StatusNotFound, "404", nil)
	}
	v, err := models.Discount{}.GetById(c.Request().Context(), id)
	if err != nil {
		return err
	}
	if v == nil {
		return c.Render(http.StatusNotFound, "404", nil)
	}
	return c.Render(http.StatusOK, "discount/edit", map[string]interface{}{
		FlashName: getFlashMessage(c),
		"Form":    v,
	})
}
func (DiscountController) Update(c echo.Context) error {
	var v DiscountInput
	if err := c.Bind(&v); err != nil {
		setFlashMessage(c, map[string]string{"error": err.Error()})
		return c.Redirect(http.StatusFound, "/discount/new")
	}
	if err := c.Validate(&v); err != nil {
		setFlashMessage(c, map[string]string{"error": err.Error()})
		return c.Redirect(http.StatusFound, "/discounts/new")
	}
	discount, err := v.ToModel()
	if err != nil {
		setFlashMessage(c, map[string]string{"error": err.Error()})
		return c.Redirect(http.StatusFound, "/discounts/new")
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.Render(http.StatusNotFound, "404", nil)
	}
	discount.Id = id
	if err := discount.Update(c.Request().Context()); err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, fmt.Sprintf("/discounts/%d", discount.Id))
}
