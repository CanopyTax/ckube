package util

import (
	"fmt"
	"time"
)

type Age struct {
	Time time.Time
}

func (a *Age) Relative() string {
	var relativeAge string
	switch d := time.Since(a.Time); true {
	case d.Minutes() < 1:
		relativeAge = fmt.Sprintf("%vs", int(d.Seconds()))
	case d.Hours() < 1:
		relativeAge = fmt.Sprintf("%vm", int(d.Minutes()))
	case d.Hours() < 24:
		relativeAge = fmt.Sprintf("%vh", int(d.Hours()))
	default:
		relativeAge = fmt.Sprintf("%vd", int(d.Hours()/24))
	}
	return relativeAge
}
