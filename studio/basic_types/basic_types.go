package basic_types

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type AccessLevel uint

const (
	CUSTOMER AccessLevel = iota + 1
	OPERATOR
	SYSADMIN
)

type Entity interface {
	GetFirstLastName() string
	GetAccessLevel() AccessLevel
	GetId() uint
}

type Customer struct {
	Id         uint
	First_name string
	Last_name  string
}

func (c *Customer) GetFirstLastName() string {
	return fmt.Sprint(c.First_name, " ", c.Last_name)
}

func (*Customer) GetAccessLevel() AccessLevel { return CUSTOMER }
func (c Customer) GetId() uint                { return c.Id }

type Operator struct {
	Id         uint
	First_name string
	Last_name  string
}

func (o *Operator) GetFirstLastName() string {
	return fmt.Sprint(o.First_name, " ", o.Last_name)
}

func (*Operator) GetAccessLevel() AccessLevel { return OPERATOR }
func (o Operator) GetId() uint                { return o.Id }

type SysAdmin struct{}

func (*SysAdmin) GetFirstLastName() string {
	return "Системный администратор"
}

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
