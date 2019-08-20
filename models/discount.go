package models

import (
	"context"
	"time"

	"github.com/pangpanglabs/echosample/factory"
)

type Discount struct {
	Id             int64     `json:"id"`
	Name           string    `json:"name"`
	Desc           string    `json:"desc"`
	StartAt        time.Time `json:"startAt"`
	EndAt          time.Time `json:"endAt"`
	ActionType     string    `json:"actionType"`
	DiscountAmount float64   `json:"discountAmount"`
	Enable         bool      `json:"enable"`
	CreatedAt      time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt      time.Time `json:"updatedAt" xorm:"updated"`
}

func (d *Discount) Create(ctx context.Context) (int64, error) {
	return factory.DB(ctx).Insert(d)
}
func (Discount) GetById(ctx context.Context, id int64) (*Discount, error) {
	var v Discount
	if has, err := factory.DB(ctx).ID(id).Get(&v); err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return &v, nil
}
func (Discount) GetAll(ctx context.Context, sortby, order []string, offset, limit int) (int64, []Discount, error) {
	q := factory.DB(ctx)
	if err := setSortOrder(q, sortby, order); err != nil {
		factory.Logger(ctx).Error(err)
	}

	var items []Discount
	totalCount, err := q.Limit(limit, offset).FindAndCount(&items)
	if err != nil {
		return 0, nil, err
	}
	return totalCount, items, nil
}
func (d *Discount) Update(ctx context.Context) (err error) {
	_, err = factory.DB(ctx).ID(d.Id).Update(d)
	return
}

func (Discount) Delete(ctx context.Context, id int64) (err error) {
	_, err = factory.DB(ctx).ID(id).Delete(&Discount{})
	return
}
