package models

import "github.com/go-xorm/xorm"

type DB struct{ *xorm.Session }
