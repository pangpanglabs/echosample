package controllers

import (
	"time"

	"github.com/pangpanglabs/echosample/models"
)

const (
	DefaultMaxResultCount = 30
)

type SearchInput struct {
	Sortby         []string `query:"sortby"`
	Order          []string `query:"order"`
	SkipCount      int      `query:"skipCount"`
	MaxResultCount int      `query:"maxResultCount"`
}
type DiscountInput struct {
	Name           string  `json:"name" valid:"required"`
	Desc           string  `json:"desc"`
	StartAt        string  `json:"startAt" valid:"required"`
	EndAt          string  `json:"endAt" valid:"required"`
	ActionType     string  `json:"actionType" valid:"required"`
	DiscountAmount float64 `json:"discountAmount" valid:"required"`
	Enable         bool    `json:"enable"`
}

func (d *DiscountInput) ToModel() (*models.Discount, error) {
	startAt, err := time.Parse("2006-01-02", d.StartAt)
	if err != nil {
		return nil, err
	}
	endAt, err := time.Parse("2006-01-02", d.EndAt)
	if err != nil {
		return nil, err
	}
	return &models.Discount{
		Name:           d.Name,
		Desc:           d.Desc,
		StartAt:        startAt,
		EndAt:          endAt,
		ActionType:     d.ActionType,
		DiscountAmount: d.DiscountAmount,
		Enable:         d.Enable,
	}, nil
}
