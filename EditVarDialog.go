package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strings"
)

func ShowDialog(parent walk.Form, name string, value string) (int, error) {
	minSize := Size{400, 300}
	if name == "Path" {
		value = strings.Replace(value, ";", "\r\n", -1)
		minSize = Size{450, 600}
	}
	var okButton, cancelButton *walk.PushButton
	return Dialog{
		Title:         "Edit Variable",
		Layout:        VBox{},
		MinSize:       minSize,
		DefaultButton: &okButton,
		CancelButton:  &cancelButton,
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
					},
					PushButton{
						Text:     "Cancel",
						AssignTo: &cancelButton,
					},
				},
			},
		},
	}.Run(parent)
}
