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
	GetFirstLastName() string
	GetAccessLevel() AccessLevel
	GetId() uint
}

type Customer struct {
	Id        uint
	FirstName string
	LastName  string
	Login     string
}

func (c *Customer) GetFirstLastName() string {
	return fmt.Sprint(c.FirstName, " ", c.LastName)
}

func (*Customer) GetAccessLevel() AccessLevel { return CUSTOMER }
func (c Customer) GetId() uint                { return c.Id }

type Employee struct {
	Id        uint
	FirstName string
	LastName  string
	JobId     uint
	Login     string
}

func (e *Employee) GetFirstLastName() string {
	return fmt.Sprint(e.FirstName, " ", e.LastName)
}

func (e *Employee) GetAccessLevel() AccessLevel { return AccessLevel(e.JobId) }

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
	Price     float64
	Materials map[uint]Material
	MatLeng   map[uint]float64
}
