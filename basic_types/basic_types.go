package basic_types

import (
	"fmt"
	"time"
)

const DateFormat string = "02.01.2006 15:04:05"

type AccessLevel uint

const (
	CUSTOMER AccessLevel = iota + 1
	OPERATOR
)

type Entity interface {
	FullName() string
	AccessLevel() AccessLevel
	GetId() uint
	GetLogin() string
}

type Customer struct {
	Id        uint
	FirstName string
	LastName  string
	Login     string
}

func (c *Customer) FullName() string {
	return fmt.Sprint(c.FirstName, " ", c.LastName)
}

func (*Customer) AccessLevel() AccessLevel { return CUSTOMER }
func (c Customer) GetId() uint             { return c.Id }
func (c Customer) GetLogin() string        { return c.Login }

type Employee struct {
	Id        uint
	FirstName string
	LastName  string
	JobId     uint
	Login     string
}

func (e *Employee) FullName() string {
	return fmt.Sprint(e.FirstName, " ", e.LastName)
}

func (e *Employee) AccessLevel() AccessLevel { return AccessLevel(e.JobId) }
func (e Employee) GetId() uint               { return e.Id }
func (e Employee) GetLogin() string          { return e.Login }

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
	Id           uint
	CustomerName string
	EmployeeName string
	Status       OrderStatus
	TotalPrice   float64
	CreateDate   int64
	ReleaseDate  int64
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
