package main

import (
	"github.com/ZhuBicen/walk"
	. "github.com/ZhuBicen/walk/declarative"
)

func main() {
	MainWindow{
		Title:  "Enviroment Variable",
		Size:   Size{600, 700},
		Layout: VBox{},
		Children: []Widget{
			GroupBox{
				Title:  "User variables for ",
				Layout: VBox{},
				Children: []Widget{
					TableView{
						AlternatingRowBGColor: walk.RGB(255, 255, 224),
						ColumnsOrderable:      true,
						Columns: []TableViewColumn{
							{Title: "#"},
							{Title: "Bar"},
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
				Layout: VBox{},
				Children: []Widget{
					TableView{
						AlternatingRowBGColor: walk.RGB(255, 255, 224),
						ColumnsOrderable:      true,
						Columns: []TableViewColumn{
							{Title: "#"},
							{Title: "Bar"},
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
							},
							PushButton{
								Text: "Delete",
							},
						},
					},
				},
			},
		},
	}.Run()

}
