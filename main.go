package main

import (
	"github.com/ZhuBicen/walk"
	. "github.com/ZhuBicen/walk/declarative"
	"math/rand"
	"strings"
)

type Variable struct {
	Index int
	Name  string
	Value string
}

type EnvModel struct {
	walk.TableModelBase
	items []*Variable
}

func NewEnvModel() *EnvModel {
	m := new(EnvModel)

	m.ResetRows()
	return m
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *EnvModel) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *EnvModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index

	case 1:
		return item.Name

	case 2:
		return item.Value
	}

	panic("unexpected col")
}

// Called by the TableView to retrieve if a given row is checked.
func (m *EnvModel) Checked(row int) bool {
	return false
}

func (m *EnvModel) Len() int {
	return len(m.items)
}

// Called by the TableView to retrieve an item image.
func (m *EnvModel) Image(row int) interface{} {
	return nil
}

func (m *EnvModel) ResetRows() {
	// Create some random data.
	m.items = make([]*Variable, rand.Intn(50000))

	for i := range m.items {
		m.items[i] = &Variable{
			Index: i,
			Name:  strings.Repeat("*", rand.Intn(5)+1),
			Value: strings.Repeat("*", rand.Intn(5)+1),
		}
	}

	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()
}
func main() {
	model := NewEnvModel()
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
							{Title: "Variable"},
							{Title: "Value"},
						},
						LastColumnStretched: true,
						Model:               model,
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
							{Title: "Variable"},
							{Title: "Value"},
						},
						LastColumnStretched: true,
						Model:               model,
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
