package controllers

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/echosample/factory"
	"github.com/pangpanglabs/goutils/behaviorlog"
)

const (
	FlashName      = "flash"
	FlashSeparator = ";"
)

type ApiResult struct {
	Result  interface{}   `json:"result"`
	Success bool          `json:"success"`
	Error   factory.Error `json:"error"`
}

type ArrayResult struct {
	Items      interface{} `json:"items"`
	TotalCount int64       `json:"totalCount"`
}

func renderFail(c echo.Context, err error) error {
	if err == nil {
		err = factory.ErrorSystem.New(nil)
	}
	behaviorlog.FromCtx(c.Request().Context()).WithError(err)
	var apiError *factory.Error
	if ok := errors.As(err, &apiError); ok {
		return c.JSON(apiError.Status(), ApiResult{
			Success: false,
			Error:   *apiError,
		})
	}
	return err
}

func renderSucc(c echo.Context, status int, result interface{}) error {
	req := c.Request()
	if req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE" {
		if session, ok := factory.DB(req.Context()).(*xorm.Session); ok {
			err := session.Commit()
			if err != nil {
				return renderFail(c, factory.ErrorDB.New(err))
			}
		}
	}

	return c.JSON(status, ApiResult{
		Success: true,
		Result:  result,
	})
}

func setFlashMessage(c echo.Context, m map[string]string) {
	var flashValue string
	for key, value := range m {
		flashValue += "\x00" + key + "\x23" + FlashSeparator + "\x23" + value + "\x00"
	}

	c.SetCookie(&http.Cookie{
		Name:  FlashName,
		Value: url.QueryEscape(flashValue),
	})
}
func getFlashMessage(c echo.Context) map[string]string {
	cookie, err := c.Cookie(FlashName)
	if err != nil {
		return nil
	}

	m := map[string]string{}

	v, _ := url.QueryUnescape(cookie.Value)
	vals := strings.Split(v, "\x00")
	for _, v := range vals {
		if len(v) > 0 {
			kv := strings.Split(v, "\x23"+FlashSeparator+"\x23")
			if len(kv) == 2 {
				m[kv[0]] = kv[1]
			}
		}
	}
	//read one time then delete it
	c.SetCookie(&http.Cookie{
		Name:   FlashName,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	return m
}
