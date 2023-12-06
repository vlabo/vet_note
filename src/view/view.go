package view

import "gioui.org/layout"

type View interface {
	Layout(gtx layout.Context) layout.Dimensions
}
