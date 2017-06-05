package controllers

import (
	"fmt"
	"net/http"
	"offer/models"
	"strconv"

	"github.com/labstack/echo"
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
	totalCount, items, err := c.Get("DB").(*models.DB).GetAllDiscount(nil, nil, nil, 0, 30)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "discount/index", map[string]interface{}{
		"TotalCount": totalCount,
		"Discounts":  items,
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
	if _, err := c.Get("DB").(*models.DB).AddDiscount(discount); err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, fmt.Sprintf("/discounts/%d", discount.Id))
}
func (DiscountController) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.Render(http.StatusNotFound, "404", nil)
	}
	v, err := c.Get("DB").(*models.DB).GetDiscountById(id)
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
	v, err := c.Get("DB").(*models.DB).GetDiscountById(id)
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
	if err := c.Get("DB").(*models.DB).UpdateDiscountById(discount); err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, fmt.Sprintf("/discounts/%d", discount.Id))
}
