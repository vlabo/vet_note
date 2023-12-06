package view

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/sahilm/fuzzy"
	"github.com/vlabo/vet_note/src/storage"
)

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

	newEntryButton widget.Clickable
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
	m.ReloadPatientList()

	return &m
}

func (m *SelectView) ReloadPatientList() error {
	patientList, err := storage.GetAllPatients()
	if err != nil {
		return err
	}
	m.patientList = nil
	for _, patient := range patientList {
		m.patientList = append(m.patientList, newPatientEntry(m.theme, patient))
	}

	return nil
}

func (m *SelectView) Layout(gtx layout.Context) layout.Dimensions {
	if m.newEntryButton.Clicked(gtx) {
		ViewEventChan <- OpenPatientView
	}

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

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return material.List(m.theme, &m.layout).Layout(gtx, len(widgets), func(gtx layout.Context, i int) layout.Dimensions {
				return layout.UniformInset(unit.Dp(2)).Layout(gtx, widgets[i])
			})
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(0.5, layout.Spacer{Width: 20}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Button(m.theme, &m.newEntryButton, "New").Layout(gtx)
				}),
			)
		}),
	)

}
