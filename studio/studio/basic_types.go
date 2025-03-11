package studio

import (
	"github.com/shopspring/decimal"
)

type Entity interface {
	GetAccessLevel() uint
}

type Customer struct {
	id         uint
	first_name string
	last_name  string
	tel        string
}

type Operator struct {
	id         uint
	first_name string
	last_name  string
	tel        string
}

type SysAdmin struct{}

func (*SysAdmin) GetAccessLevel() uint { return 3 }

type Order struct {
	id          uint
	customer_id uint
	items       []Model
	total_price decimal.Decimal
}

type Material struct {
	id    uint
	title string
	price decimal.Decimal
}

type Model struct {
	id        uint
	title     string
	materials []Material
	price     decimal.Decimal
}
