package basic_types

import (
	"fmt"
	"time"
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

type Employee struct {
	Id         uint
	First_name string
	Last_name  string
	JobId      uint
}

func (o *Employee) GetFirstLastName() string {
	return fmt.Sprint(o.First_name, " ", o.Last_name)
}

func (*Employee) GetAccessLevel() AccessLevel { return OPERATOR }
func (o Employee) GetId() uint                { return o.Id }

type OrderStatus uint

const (
	Pending OrderStatus = iota
	Processing
	Released
	Canceled
)

func (stat OrderStatus) String() string {
	return [...]string{"Ожидает", "На исполнении", "Выдан", "Отменен"}[stat]
}

type Order struct {
	Id          uint
	C_id        uint
	E_id        uint
	Status      OrderStatus
	TotalPrice  float64
	CreateDate  time.Time
	ReleaseDate time.Time
}

type RawOrder struct {
	Id          uint
	C_id        uint
	E_id        uint
	Status      OrderStatus
	TotalPrice  float64
	CreateDate  int64
	ReleaseDate int64
}

type OrderItem struct {
	Id        uint
	O_id      uint
	Model     []Model
	UnitPrice float64
}

type RawOrderItem struct {
	Id        uint
	O_id      uint
	Model     uint
	UnitPrice float64
}

type Material struct {
	Id    uint
	Title string
	Price float64
}

type Model struct {
	Id        uint
	Title     string
	Materials map[uint]Material
	MatLeng   map[uint]float64
}
