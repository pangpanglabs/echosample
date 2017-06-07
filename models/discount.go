package models

import (
	"context"
	"offer/factory"
	"time"

	"github.com/go-xorm/xorm"
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
func (Discount) GetAll(ctx context.Context, sortby, order []string, offset, limit int) (totalCount int64, items []Discount, err error) {
	queryBuilder := func() *xorm.Session {
		q := factory.DB(ctx).Table("discount")
		if err := setSortOrder(q, sortby, order); err != nil {
			factory.Logger(ctx).Error(err)
		}
		return q
	}

	errc := make(chan error)
	go func() {
		v, err := queryBuilder().Count(&Discount{})
		if err != nil {
			errc <- err
			return
		}
		totalCount = v
		errc <- nil

	}()

	go func() {
		if err := queryBuilder().Limit(limit, offset).Find(&items); err != nil {
			errc <- err
			return
		}
		errc <- nil
	}()

	if err := <-errc; err != nil {
		return 0, nil, err
	}
	if err := <-errc; err != nil {
		return 0, nil, err
	}
	return
}
func (d *Discount) Update(ctx context.Context) (err error) {
	_, err = factory.DB(ctx).ID(d.Id).Update(d)
	return
}

func (Discount) Delete(ctx context.Context, id int64) (err error) {
	_, err = factory.DB(ctx).ID(id).Delete(&Discount{})
	return
}
