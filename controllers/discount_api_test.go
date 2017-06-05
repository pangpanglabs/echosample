package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"offer/filters"
	"offer/models"
	"os"
	"pangpanglabs/goutils/test"
	"strings"
	"testing"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

var (
	echoApp          *echo.Echo
	handleWithFilter func(handlerFunc echo.HandlerFunc, c echo.Context) error
)

func init() {
	os.Remove("test.db")
	xormEngine, err := xorm.NewEngine("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	xormEngine.Sync(new(models.Discount))
	echoApp = echo.New()
	echoApp.Validator = &filters.Validator{}

	handleWithFilter = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return filters.SetDbContext(xormEngine)(handlerFunc)(c)
	}
}

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
}

func Test_DiscountApiController_GetOne(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/api/discounts/1", nil)
	rec := httptest.NewRecorder()
	c := echoApp.NewContext(req, rec)
	c.SetPath("/api/discounts/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	test.Ok(t, handleWithFilter(DiscountApiController{}.GetOne, c))
	test.Equals(t, http.StatusOK, rec.Code)

	var v struct {
		Result  models.Discount `json:"result"`
		Success bool            `json:"success"`
	}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	test.Equals(t, v.Result.Name, "discount name2")
	test.Equals(t, v.Result.StartAt.Format("2006-01-02"), "2017-01-02")
}

func Test_DiscountApiController_GetAll(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/api/discounts", nil)
	rec := httptest.NewRecorder()
	c := echoApp.NewContext(req, rec)
	test.Ok(t, handleWithFilter(DiscountApiController{}.GetAll, c))
	test.Equals(t, http.StatusOK, rec.Code)

	var v struct {
		Result struct {
			TotalCount int
			Items      []models.Discount
		} `json:"result"`
		Success bool `json:"success"`
	}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	test.Equals(t, v.Result.TotalCount, 1)
	test.Equals(t, v.Result.Items[0].StartAt.Format("2006-01-02"), "2017-01-02")
}
