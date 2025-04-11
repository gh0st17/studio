package basic_types

import (
	"fmt"
	"time"
)

type AccessLevel uint

const (
	CUSTOMER AccessLevel = iota + 1
	OPERATOR
)

type Entity interface {
	FirstLastName() string
	AccessLevel() AccessLevel
	GetId() uint
}

type Customer struct {
	Id        uint
	FirstName string
	LastName  string
	Login     string
}

func (c *Customer) FirstLastName() string {
	return fmt.Sprint(c.FirstName, " ", c.LastName)
}

func (*Customer) AccessLevel() AccessLevel { return CUSTOMER }
func (c Customer) GetId() uint             { return c.Id }

type Employee struct {
	Id        uint
	FirstName string
	LastName  string
	JobId     uint
	Login     string
}

func (e *Employee) FirstLastName() string {
	return fmt.Sprint(e.FirstName, " ", e.LastName)
}

func (e *Employee) AccessLevel() AccessLevel { return AccessLevel(e.JobId) }

func (e Employee) GetId() uint { return e.Id }

type OrderStatus uint

const (
	Pending OrderStatus = iota + 1
	Processing
	Released
	Canceled
)

func (stat OrderStatus) String() string {
	return [...]string{"Ожидает", "На исполнении", "Выдан", "Отменен"}[stat-1]
}

type Order struct {
	Id          uint
	C_id        uint
	E_id        uint
	Status      OrderStatus
	TotalPrice  float64
	CreateDate  int64
	ReleaseDate int64
}

func (o *Order) time(unixSec int64) time.Time {
	return time.Unix(unixSec, 0)
}

func (o *Order) LocalCreateDate() time.Time {
	return o.time(o.CreateDate)
}

func (o *Order) LocalReleaseDate() time.Time {
	return o.time(o.ReleaseDate)
}

type OrderItem struct {
	Id        uint
	O_id      uint
	Model     Model
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
	Price     float64
	Materials map[uint]Material
	MatLeng   map[uint]float64
}
