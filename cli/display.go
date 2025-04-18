package cli

import (
	"fmt"
	bt "studio/basic_types"
	"time"
)

type customerOrders []bt.Order

func (orders customerOrders) String() (s string) {
	var ctime, rtime string

	s = fmt.Sprintf(
		"  # Статус заказа %9s %19s %19s\n",
		"Сумма", "Создан", "Выдан",
	)

	for _, o := range orders {
		ctime = o.CreateDate.Format(bt.DateFormat)

		if o.ReleaseDate.Equal(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)) {
			rtime = "---"
		} else {
			rtime = o.ReleaseDate.Format(bt.DateFormat)
		}

		s += fmt.Sprintf("%3d %13s %9.2f %19s %19s\n",
			o.Id, o.Status, o.TotalPrice, ctime, rtime,
		)
	}

	return s
}

type employeeOrders []bt.Order

func (orders employeeOrders) String() (s string) {
	var ctime, rtime string

	s = fmt.Sprintf(
		"  # Статус заказа %9s %19s %19s %s\n",
		"Сумма", "Создан", "Выдан", "Имя клиента",
	)

	for _, o := range orders {
		ctime = o.CreateDate.Format(bt.DateFormat)

		if o.ReleaseDate.Equal(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)) {
			rtime = "---"
		} else {
			rtime = o.ReleaseDate.Format(bt.DateFormat)
		}

		s += fmt.Sprintf("%3d %13s %9.2f %19s %19s %s\n",
			o.Id, o.Status, o.TotalPrice, ctime, rtime, o.CustomerName,
		)
	}

	return s
}

type (
	OrderItems []bt.OrderItem
	Model      bt.Model
	Models     []bt.Model
)

func (orderItems OrderItems) String() (s string) {
	var sum float64 = 0.0

	for i, oi := range orderItems {
		s += fmt.Sprintln("Позиция:", i+1)
		s += Model(oi.Model).String()
		sum += oi.UnitPrice
	}
	s += fmt.Sprintln("Общая стоимость заказа:", sum)

	return s
}

func (model Model) String() (s string) {
	s += fmt.Sprintf("%s (Артикул %d):\n", model.Title, model.Id)

	for _, mat := range model.Materials {
		s += fmt.Sprintf("\t%s стоимостью %2.2f за погонный метр длиной %2.2f метра\n",
			mat.Title, mat.Price, model.MatLeng[mat.Id],
		)
	}
	s += fmt.Sprintf("\tCтоимость изготовления %2.2f\n\n", model.Price)

	return s
}

func (models Models) String() (s string) {
	for _, m := range models {
		s += Model(m).String()
	}

	return s
}
