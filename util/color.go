package util

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"sync"
)

type ColorManager struct {
	sync.RWMutex
	colors []color.Attribute
}

func (cm *ColorManager) GetColor() color.Attribute {
	cm.Lock()
	defer cm.Unlock()

	if len(cm.colors) == 0 {
		cm.setColors()
	}

	colAttr := cm.colors[0]
	cm.colors = append(cm.colors[:0], cm.colors[1:]...)
	return colAttr
}

func (cm *ColorManager) setColors() {
	cm.colors = []color.Attribute{
		color.FgHiGreen,
		color.FgHiYellow,
		color.FgHiBlue,
		color.FgHiCyan,
		color.FgHiRed,
		color.FgHiMagenta,
		color.FgGreen,
		color.FgBlue,
		color.FgYellow,
		color.FgCyan,
		color.FgRed,
		color.FgMagenta,
		color.BgHiRed,
		color.BgHiGreen,
		color.BgHiYellow,
		color.BgHiBlue,
		color.BgHiMagenta,
		color.BgHiCyan,
	}
}

func (cm *ColorManager) GetPrefix(prefix string) string {
	beginning := strings.Split(prefix, "-")[0]
	bytes := []byte(prefix)
	lastFive := bytes[len(bytes)-5:]
	return cm.Colorize(fmt.Sprintf("[%v...%v]", beginning, string(lastFive)))
}

func (cm *ColorManager) Colorize(s string) string {
	colAttribute := cm.GetColor()
	col := color.New(colAttribute).Add(color.Bold).SprintfFunc()
	return col(s)
}
