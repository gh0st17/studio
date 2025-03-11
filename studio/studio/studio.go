package studio

import "studio/ui"

type Studio struct {
	customers []Customer
	orders    []Order
	materials []Material
	models    []Model
}

func (s *Studio) Run(u ui.UI) error {
	u.Run()
	return nil
}
