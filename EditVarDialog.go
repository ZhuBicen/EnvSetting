package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strings"
)

func ShowDialog(parent walk.Form, name string, value string) (int, error) {
	minSize := Size{400, 300}
	if name == "Path" {

		minSize = Size{450, 600}
	}
	value = strings.Replace(value, ";", "\r\n", -1)
	var okButton, cancelButton *walk.PushButton
	var dialog *walk.Dialog
	return Dialog{
		Title:         "Edit Variable",
		Layout:        VBox{},
		MinSize:       minSize,
		DefaultButton: &okButton,
		CancelButton:  &cancelButton,
		AssignTo:      &dialog,
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "Name:",
					},
					LineEdit{
						Text: name,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "Value:",
					},
					TextEdit{
						Text: value,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text:     "OK",
						AssignTo: &okButton,
						OnClicked: func() {
							dialog.Close(1)
						},
					},
					PushButton{
						Text:     "Cancel",
						AssignTo: &cancelButton,
						OnClicked: func() {
							dialog.Close(1)
						},
					},
				},
			},
		},
	}.Run(parent)
}
