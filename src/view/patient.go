package view

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/vlabo/vet_note/src/model"
	"github.com/vlabo/vet_note/src/storage"
)

type PatientView struct {
	theme *material.Theme

	name  widget.Editor
	owner widget.Editor

	createButton widget.Clickable

	layout layout.Flex
}

func NewPatientView(t *material.Theme) *PatientView {
	return &PatientView{
		theme: t,
		name: widget.Editor{
			SingleLine: true,
		},
		owner: widget.Editor{
			SingleLine: true,
		},
		layout: layout.Flex{Axis: layout.Vertical},
	}
}

func (pv *PatientView) Layout(gtx layout.Context) layout.Dimensions {

	if pv.createButton.Clicked(gtx) {
		patient := model.Patient{
			Name:  pv.name.Text(),
			Owner: pv.owner.Text(),
		}
		pv.name.SetText("")
		pv.owner.SetText("")

		storage.AddPatient(&patient)

		ViewEventChan <- OpenSelectView
	}

	return pv.layout.Layout(gtx,
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			e := material.Editor(pv.theme, &pv.name, "Name")
			border := widget.Border{Color: pv.theme.ContrastBg, CornerRadius: unit.Dp(2), Width: unit.Dp(1)}
			return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(5)).Layout(gtx, e.Layout)
			})
		}),
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			e := material.Editor(pv.theme, &pv.owner, "Owner")
			border := widget.Border{Color: pv.theme.ContrastBg, CornerRadius: unit.Dp(2), Width: unit.Dp(1)}
			return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(5)).Layout(gtx, e.Layout)
			})
		}),
		layout.Flexed(1, layout.Spacer{Width: 20}.Layout),
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			if pv.name.Len() == 0 || pv.owner.Len() == 0 {
				gtx.Disabled()
			}
			return material.Button(pv.theme, &pv.createButton, "Add").Layout(gtx)
		}),
	)
}
