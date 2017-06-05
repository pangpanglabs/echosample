package models

import "time"

type DiscountInput struct {
	Name           string `valid:"required"`
	Desc           string
	StartAt        string  `valid:"required"`
	EndAt          string  `valid:"required"`
	ActionType     string  `valid:"required"`
	DiscountAmount float64 `valid:"required"`
	Enable         bool
}

func (d *DiscountInput) ToModel() (*Discount, error) {
	startAt, err := time.Parse("2006-01-02", d.StartAt)
	if err != nil {
		return nil, err
	}
	endAt, err := time.Parse("2006-01-02", d.EndAt)
	if err != nil {
		return nil, err
	}
	return &Discount{
		Name:           d.Name,
		Desc:           d.Desc,
		StartAt:        startAt,
		EndAt:          endAt,
		ActionType:     d.ActionType,
		DiscountAmount: d.DiscountAmount,
		Enable:         d.Enable,
	}, nil
}
