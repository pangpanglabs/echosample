package controllers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo"
)

const (
	FlashName      = "flash"
	FlashSeparator = ";"
)

type ApiResult struct {
	Result  interface{} `json:"result"`
	Success bool        `json:"success"`
	Error   ApiError    `json:"error"`
}

type ApiError struct {
	Code    int         `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ArrayResult struct {
	Items      interface{} `json:"items"`
	TotalCount int64       `json:"totalCount"`
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
