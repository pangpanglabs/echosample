package models

import (
	"context"
	"time"

	"github.com/pangpanglabs/echosample/factory"

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
	queryBuilder := func() xorm.Interface {
		q := factory.DB(ctx)
		if err := setSortOrder(q, sortby, order); err != nil {
			factory.Logger(ctx).Error(err)
		}
		return q
	}

	errc, totalCountc, discountc := make(chan error), make(chan int64, 1), make(chan []Discount, 1)
	go func() {
		v, err := queryBuilder().Count(&Discount{})
		totalCountc <- v
		errc <- err
	}()

	go func() {
		var v []Discount
		err := queryBuilder().Limit(limit, offset).Find(&v)
		discountc <- v
		errc <- err
	}()
	for i := 0; i < 2; i++ {
		if err := <-errc; err != nil {
			return 0, nil, err
		}
	}
	totalCount = <-totalCountc
	items = <-discountc
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
