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
	affected, err := factory.DB(ctx).Insert(d)
	if err != nil {
		return 0, factory.ErrorDB.New(err)
	}
	return affected, nil
}

func (Discount) GetById(ctx context.Context, id int64) (*Discount, error) {
	var v Discount
	if has, err := factory.DB(ctx).ID(id).Get(&v); err != nil {
		return nil, factory.ErrorDB.New(err)
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
		return 0, nil, factory.ErrorDB.New(err)
	}
	return totalCount, items, nil
}

func (d *Discount) Update(ctx context.Context) error {
	if origin, err := d.GetById(ctx, d.Id); err != nil {
		return factory.ErrorDB.New(err)
	} else if origin == nil {
		return factory.ErrorDiscountNotExists.New(err, d.Id)
	}

	if _, err := factory.DB(ctx).ID(d.Id).Update(d); err != nil {
		return factory.ErrorDB.New(err)
	}
	return nil
}

func (Discount) Delete(ctx context.Context, id int64) error {
	if origin, err := (Discount{}).GetById(ctx, id); err != nil {
		return factory.ErrorDB.New(err)
	} else if origin == nil {
		return factory.ErrorDiscountNotExists.New(err, id)
	}

	if _, err := factory.DB(ctx).ID(id).Delete(&Discount{}); err != nil {
		return factory.ErrorDB.New(err)
	}
	return nil
}
