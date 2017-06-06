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
		if len(sortby) != 0 {
			if len(sortby) == len(order) {
				// 1) for each sort field, there is an associated order
				for i, v := range sortby {
					if order[i] == "desc" {
						q.Desc(v)
					} else if order[i] == "asc" {
						q.Asc(v)
					} else {
						factory.Logger(ctx).Error("Invalid order. Must be either [asc|desc]")
					}
				}
			} else if len(sortby) != len(order) && len(order) == 1 {
				// 2) there is exactly one order, all the sorted fields will be sorted by this order
				for _, v := range sortby {
					if order[0] == "desc" {
						q.Desc(v)
					} else if order[0] == "asc" {
						q.Asc(v)
					} else {
						factory.Logger(ctx).Error("Invalid order. Must be either [asc|desc]")
					}
				}
			} else if len(sortby) != len(order) && len(order) != 1 {
				factory.Logger(ctx).Error("'sortby', 'order' sizes mismatch or 'order' size is not 1")
			}
		} else {
			if len(order) != 0 {
				factory.Logger(ctx).Error("unused 'order' fields")
			}
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
