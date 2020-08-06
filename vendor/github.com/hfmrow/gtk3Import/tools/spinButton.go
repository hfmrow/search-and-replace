// spinButton.go

/// +build ignore

/*
	©2019 H.F.M. MIT license
*/

package gtk3Import

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

// SpinbuttonSetValues: Configure spin button
// min, max, value, stepIncrement, pageIncrement, pageSize
func SpinbuttonSetValues(sb *gtk.SpinButton, min, max, value int, step ...int) {
	incStep, pageIncrement, pageSize := 1, 0, 0
	switch len(step) {
	case 1:
		incStep = step[0]
	case 2:
		incStep = step[0]
		pageIncrement = step[1]
	case 3:
		incStep = step[0]
		pageIncrement = step[1]
		pageSize = step[2]
	}
	if ad, err := gtk.AdjustmentNew(float64(value), float64(min), float64(max),
		float64(incStep), float64(pageIncrement), float64(pageSize)); err == nil {
		sb.Configure(ad, 1, 0)
	} else {
		fmt.Errorf("SpinbuttonSetValues: %s", err.Error())
	}
}
