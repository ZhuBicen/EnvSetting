package main

import (
	"./env"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"os/user"
)

func NewVariable(mw *walk.MainWindow, m *env.Model) {
	if ret, name, value := ShowDialog(mw, "", ""); ret == 0 {
		if !m.AddVariable(name, value) {
			walk.MsgBox(mw, "Error", "The variable has already existed.", walk.MsgBoxOK)
		}
	}
}
func EditVariable(mw *walk.MainWindow, m *env.Model, name string, value string) {
	if ret, newName, newValue := ShowDialog(mw, name, value); ret == 0 {
		if newName != name {
			walk.MsgBox(mw, "Error", "You can't change the variable name when editing.", walk.MsgBoxOK)
			return
		}
		if !m.EditVariable(name, newValue) {
			walk.MsgBox(mw, "Error", "Please ensure the variable has already existed.", walk.MsgBoxOK)
			return
		}
	}
}

func DeleteVariable(mw *walk.MainWindow, m *env.Model, name string) {
	if win.IDYES == walk.MsgBox(mw, "Information", "You want to delete variable "+name, walk.MsgBoxYesNo) {
		if !m.DeleteVariable(name) {
			walk.MsgBox(mw, "Error", "Please ensure the variable has already existed.", walk.MsgBoxOK)
		}
	}
}

func ApplyEnv(mw *walk.MainWindow, usrModel *env.Model, sysModel *env.Model) bool {
	if err := usrModel.Apply(); err != nil {
		walk.MsgBox(mw, "Error", fmt.Sprintf("%s", err), walk.MsgBoxOK)
		return false
	}
	if err := sysModel.Apply(); err != nil {
		walk.MsgBox(mw, "Error", fmt.Sprintf("%s", err), walk.MsgBoxOK)
		return false
	}
	return true
}
func main() {
	font := Font{
		Family:    "Times New Roman",
		PointSize: 13,
		Bold:      true,
	}
	usr, _ := user.Current()
	usrModel := env.NewModel(0)
	sysModel := env.NewModel(1)
	var usrTableView, sysTableView *walk.TableView
	var mw *walk.MainWindow
	MainWindow{
		Title:    "Enviroment Variable",
		Size:     Size{600, 700},
		Layout:   VBox{},
		AssignTo: &mw,
		Font:     font,
		Children: []Widget{
			VSplitter{
				Children: []Widget{
					GroupBox{
						Title:  "User variables for " + usr.Username,
						Font:   font,
						Layout: VBox{},
						Children: []Widget{
							TableView{
								AssignTo:              &usrTableView,
								AlternatingRowBGColor: walk.RGB(255, 255, 224),
								ColumnsOrderable:      true,
								Columns: []TableViewColumn{
									{Title: "name", Width: 200},
									{Title: "value"},
								},
								LastColumnStretched: true,
								Model:               usrModel,
								OnItemActivated: func() {
									index := usrTableView.CurrentIndex()
									if index != -1 {
										EditVariable(mw, usrModel, usrModel.GetVariable(index).Name, usrModel.GetVariable(index).Value)
									}
								},
							},
							Composite{
								Layout: HBox{},
								Children: []Widget{
									HSpacer{},
									PushButton{
										Text: "New...",
										OnClicked: func() {
											NewVariable(mw, usrModel)
										},
									},
									PushButton{
										Text: "Edit...",
										OnClicked: func() {
											index := usrTableView.CurrentIndex()
											if index != -1 {
												EditVariable(mw, usrModel, usrModel.GetVariable(index).Name, usrModel.GetVariable(index).Value)
											}
										},
									},
									PushButton{
										Text: "Delete",
										OnClicked: func() {
											index := usrTableView.CurrentIndex()
											if index != -1 {
												DeleteVariable(mw, usrModel, usrModel.GetVariable(index).Name)
											}
										},
									},
								},
							},
						},
					},
					GroupBox{
						Title:  "System variables",
						Font:   font,
						Layout: VBox{},
						Children: []Widget{
							TableView{
								AssignTo:              &sysTableView,
								AlternatingRowBGColor: walk.RGB(255, 255, 224),
								ColumnsOrderable:      true,
								Columns: []TableViewColumn{
									{Title: "Variable", Width: 200},
									{Title: "Value"},
								},
								LastColumnStretched: true,
								Model:               sysModel,
								OnItemActivated: func() {
									index := sysTableView.CurrentIndex()
									if index != -1 {
										EditVariable(mw, sysModel, sysModel.GetVariable(index).Name, sysModel.GetVariable(index).Value)
									}
								},
							},
							Composite{
								Layout: HBox{},
								Children: []Widget{
									HSpacer{},
									PushButton{
										Text: "New...",
										OnClicked: func() {
											NewVariable(mw, sysModel)
										},
									},
									PushButton{
										Text: "Edit...",
										OnClicked: func() {
											index := sysTableView.CurrentIndex()
											if index != -1 {
												EditVariable(mw, sysModel, sysModel.GetVariable(index).Name, sysModel.GetVariable(index).Value)
											}
										},
									},
									PushButton{
										Text: "Delete",
										OnClicked: func() {
											index := sysTableView.CurrentIndex()
											if index != -1 {
												DeleteVariable(mw, sysModel, sysModel.GetVariable(index).Name)
											}
										},
									},
								},
							},
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text: "OK",
						OnClicked: func() {
							if ApplyEnv(mw, usrModel, sysModel) {
								mw.Close()
							}
						},
					},
					PushButton{
						Text: "Cancel",
						OnClicked: func() {
							mw.Close()
						},
					},
				},
			},
		},
	}.Run()

}
