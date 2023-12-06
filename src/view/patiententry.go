package view

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/vlabo/vet_note/src/model"
)

type PatientEntry struct {
	theme      *material.Theme
	name       material.LabelStyle
	owner      material.LabelStyle
	searchText string
}

func (e *PatientEntry) asWidget(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		// Draw entry background
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			// Apply round corners
			const r = 4 // roundness
			bounds := image.Rectangle{Min: image.Point{0, 0}, Max: gtx.Constraints.Min}
			stackClip := clip.RRect{Rect: bounds, SE: r, SW: r, NW: r, NE: r}.Push(gtx.Ops)
			defer stackClip.Pop()
			// Draw rect, with previously applied round corners
			defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: e.theme.ContrastBg}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		// Draw entry elements (name, owner...)
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			margins := layout.Inset{
				Top:    unit.Dp(10),
				Bottom: unit.Dp(5),
				Right:  unit.Dp(3),
				Left:   unit.Dp(3),
			}

			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{}.Layout(gtx,
					// Left
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return e.name.Layout(gtx)
					}),
					// Middle spacer
					layout.Flexed(0.5, layout.Spacer{Width: 20}.Layout),
					// Right
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return e.owner.Layout(gtx)
					}),
				)

			})
		}),
	)
}
func newPatientEntry(theme *material.Theme, p model.Patient) PatientEntry {
	return PatientEntry{
		theme:      theme,
		name:       material.Body1(theme, p.Name),
		owner:      material.Body1(theme, p.Owner),
		searchText: p.Name + p.Owner,
	}
}
