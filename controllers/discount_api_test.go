package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"

	"github.com/pangpanglabs/echosample/models"
	"github.com/pangpanglabs/goutils/test"
)

func Test_DiscountApiController_Create(t *testing.T) {
	req := httptest.NewRequest(echo.POST, "/api/discounts", strings.NewReader(`{"name":"discount name", "desc":"discount desc", "startAt":"2017-01-01","endAt":"2017-12-31","actionType":"Percentage","discountAmount":10,"enable":true}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	test.Ok(t, handleWithFilter(DiscountApiController{}.Create, echoApp.NewContext(req, rec)))
	test.Equals(t, http.StatusOK, rec.Code)

	var v struct {
		Result  models.Discount `json:"result"`
		Success bool            `json:"success"`
	}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	test.Equals(t, v.Result.Name, "discount name")
	test.Equals(t, v.Result.StartAt.Format("2006-01-02"), "2017-01-01")
}

func Test_DiscountApiController_Create2(t *testing.T) {
	req := httptest.NewRequest(echo.POST, "/api/discounts", strings.NewReader(`{"name":"discount name#2", "desc":"discount desc#2", "startAt":"2017-02-01","endAt":"2017-11-30","actionType":"Percentage","discountAmount":20,"enable":true}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	test.Ok(t, handleWithFilter(DiscountApiController{}.Create, echoApp.NewContext(req, rec)))
	test.Equals(t, http.StatusOK, rec.Code)

	var v struct {
		Result  models.Discount `json:"result"`
		Success bool            `json:"success"`
	}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	test.Equals(t, v.Result.Name, "discount name#2")
	test.Equals(t, v.Result.StartAt.Format("2006-01-02"), "2017-02-01")
}

func Test_DiscountApiController_Update(t *testing.T) {
	req := httptest.NewRequest(echo.PUT, "/api/discounts/1", strings.NewReader(`{"name":"discount name2", "desc":"discount desc2", "startAt":"2017-01-02","endAt":"2017-12-30","actionType":"Percentage","discountAmount":10,"enable":true}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echoApp.NewContext(req, rec)
	c.SetPath("/api/discounts/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	test.Ok(t, handleWithFilter(DiscountApiController{}.Update, c))
	test.Equals(t, http.StatusOK, rec.Code)

	var v struct {
		Result  models.Discount `json:"result"`
		Success bool            `json:"success"`
	}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	test.Equals(t, v.Result.Name, "discount name2")
	test.Equals(t, v.Result.StartAt.Format("2006-01-02"), "2017-01-02")
	test.Equals(t, v.Result.UpdatedAt.Format("2006-01-02"), time.Now().Format("2006-01-02"))
}

func Test_DiscountApiController_GetOne(t *testing.T) {
	// given
	req := httptest.NewRequest(echo.GET, "/api/discounts/1", nil)
	rec := httptest.NewRecorder()
	c := echoApp.NewContext(req, rec)
	c.SetPath("/api/discounts/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	// when
	test.Ok(t, handleWithFilter(DiscountApiController{}.GetOne, c))
	test.Equals(t, http.StatusOK, rec.Code)

	// then
	var v struct {
		Result  map[string]interface{} `json:"result"`
		Success bool                   `json:"success"`
	}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	test.Equals(t, v.Result["name"], "discount name2")
	test.Equals(t, strings.HasPrefix(v.Result["startAt"].(string), "2017-01-02"), true)
	test.Equals(t, strings.HasPrefix(v.Result["endAt"].(string), "2017-02-02"), true)
}

func Test_DiscountApiController_GetAll_SortByAsc(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/api/discounts?sortby=discount_amount&order=asc", nil)
	rec := httptest.NewRecorder()
	c := echoApp.NewContext(req, rec)
	test.Ok(t, handleWithFilter(DiscountApiController{}.GetAll, c))
	if rec.Code != http.StatusOK {
		fmt.Println(rec.Body.String())
	}
	test.Equals(t, http.StatusOK, rec.Code)

	var v struct {
		Result struct {
			TotalCount int
			Items      []models.Discount
		} `json:"result"`
		Success bool `json:"success"`
	}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	test.Equals(t, v.Result.TotalCount, 2)
	test.Equals(t, v.Result.Items[0].DiscountAmount, float64(10))
}

func Test_DiscountApiController_GetAll_SortByDesc(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/api/discounts?sortby=discount_amount&order=desc", nil)
	rec := httptest.NewRecorder()
	c := echoApp.NewContext(req, rec)
	test.Ok(t, handleWithFilter(DiscountApiController{}.GetAll, c))
	if rec.Code != http.StatusOK {
		fmt.Println(rec.Body.String())
	}
	test.Equals(t, http.StatusOK, rec.Code)

	var v struct {
		Result struct {
			TotalCount int
			Items      []models.Discount
		} `json:"result"`
		Success bool `json:"success"`
	}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	test.Equals(t, v.Result.TotalCount, 2)
	test.Equals(t, v.Result.Items[0].DiscountAmount, float64(20))
}
