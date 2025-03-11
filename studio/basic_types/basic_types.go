package basic_types

import (
	"time"

	"github.com/shopspring/decimal"
)

type AccessLevel uint

const (
	CUSTOMER AccessLevel = iota
	OPERATOR
	SYSADMIN
)

type Entity interface {
	GetAccessLevel() AccessLevel
	GetId() uint
}

type Customer struct {
	Id         uint
	First_name string
	Last_name  string
}

func (*Customer) GetAccessLevel() AccessLevel { return CUSTOMER }
func (c Customer) GetId() uint                { return c.Id }

type Operator struct {
	Id         uint
	First_name string
	Last_name  string
}

func (*Operator) GetAccessLevel() AccessLevel { return OPERATOR }
func (o Operator) GetId() uint                { return o.Id }

type SysAdmin struct{}

func (*SysAdmin) GetAccessLevel() AccessLevel { return SYSADMIN }
func (s SysAdmin) GetId() uint                { return 0 }

type OrderStatus uint

const (
	Pending OrderStatus = iota
	Processing
	Released
	Canceled
)

type Order struct {
	Id          uint
	Customer_id uint
	Operator_id uint
	Status      OrderStatus
	Items       []Model
	Total_price decimal.Decimal
	CreateDate  time.Time
	ReleaseDate time.Time
}

type Material struct {
	Id    uint
	Title string
	Price decimal.Decimal
}

type Model struct {
	Id    uint
	Title string

	// Матриалы и их длина
	Materials map[Material]decimal.Decimal
}
