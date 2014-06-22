package main

import (
	"./env"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os/user"
)

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
									ShowDialog(mw, usrModel.GetVariable(index).Name, usrModel.GetVariable(index).Value)
								},
							},
							Composite{
								Layout: HBox{},
								Children: []Widget{
									HSpacer{},
									PushButton{
										Text: "New...",
										OnClicked: func() {
											if ret, name, value := ShowDialog(mw, "", ""); ret == 0 {
												log.Println("You will creating var name = ", name, ", value = ", value)
											}
										},
									},
									PushButton{
										Text: "Edit...",
										OnClicked: func() {
											index := usrTableView.CurrentIndex()
											ShowDialog(mw, usrModel.GetVariable(index).Name, usrModel.GetVariable(index).Value)
										},
									},
									PushButton{
										Text: "Delete",
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
									ShowDialog(mw, sysModel.GetVariable(index).Name, sysModel.GetVariable(index).Value)
								},
							},
							Composite{
								Layout: HBox{},
								Children: []Widget{
									HSpacer{},
									PushButton{
										Text: "New...",
									},
									PushButton{
										Text: "Edit...",
										OnClicked: func() {
											index := usrTableView.CurrentIndex()
											ShowDialog(mw, sysModel.GetVariable(index).Name, sysModel.GetVariable(index).Value)
										},
									},
									PushButton{
										Text: "Delete",
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
						Text: "Apply",
					},
					PushButton{
						Text: "Cancel",
					},
				},
			},
		},
	}.Run()

}
