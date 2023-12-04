package view

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/sahilm/fuzzy"
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

type PatientEntries []PatientEntry

// String returns search string for the specific entry. Part of fuzzy.Source interface
func (pe PatientEntries) String(i int) string {
	return pe[i].searchText
}

// Len returns number of entries in the list. Part of fuzzy.Source interface
func (pe PatientEntries) Len() int {
	return len(pe)
}

type SelectView struct {
	theme  *material.Theme
	layout widget.List

	search      widget.Editor
	patientList PatientEntries
}

// NewSelectView creates a new select view object
func NewSelectView(theme *material.Theme) *SelectView {
	m := SelectView{
		theme: theme,
		layout: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		search: widget.Editor{
			SingleLine: true,
		},
	}

	m.search.SetText("")

	// for _, patient := range patients {
	// 	m.patientList = append(m.patientList, newPatientEntry(theme, patient))
	// }

	// dataDir, err := app.DataDir()
	// if err == nil {
	// 	m.patientList = append(m.patientList, newPatientEntry(theme, Patient{name: dataDir}))
	// }
	return &m
}

func (m *SelectView) Layout(gtx layout.Context) {
	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			e := material.Editor(m.theme, &m.search, "Hint")
			border := widget.Border{Color: m.theme.ContrastBg, CornerRadius: unit.Dp(2), Width: unit.Dp(1)}
			return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(10)).Layout(gtx, e.Layout)
			})
		},
	}

	if m.search.Len() > 0 {
		matches := fuzzy.FindFrom(m.search.Text(), m.patientList)

		for _, match := range matches {
			patient := m.patientList[match.Index]
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return patient.asWidget(gtx)
			})
		}
	} else {

		for _, patient := range m.patientList {
			patient := patient
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return patient.asWidget(gtx)
			})
		}
	}

	material.List(m.theme, &m.layout).Layout(gtx, len(widgets), func(gtx layout.Context, i int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(2)).Layout(gtx, widgets[i])
	})
}
