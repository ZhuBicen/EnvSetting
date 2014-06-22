package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strings"
)

func ShowDialog(parent walk.Form, name string, value string) (int, string, string) {
	minSize := Size{400, 300}
	if name == "Path" {

		minSize = Size{450, 600}
	}
	value = strings.Replace(value, ";", "\r\n", -1)
	var okButton, cancelButton *walk.PushButton
	var nameLineEdit *walk.LineEdit
	var valueTextEdit *walk.TextEdit
	var dialog *walk.Dialog
	ret, _ := Dialog{
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
						Text: "Variable name:",
					},
					LineEdit{
						Text:     name,
						AssignTo: &nameLineEdit,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "Variable value:",
					},
					TextEdit{
						Text:     value,
						AssignTo: &valueTextEdit,
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
							name = nameLineEdit.Text()
							value = valueTextEdit.Text()
							dialog.Close(0)
						},
					},
					PushButton{
						Text:     "Cancel",
						AssignTo: &cancelButton,
						OnClicked: func() {
							name = ""
							value = ""
							dialog.Close(1)
						},
					},
				},
			},
		},
	}.Run(parent)
	return ret, name, value
}
