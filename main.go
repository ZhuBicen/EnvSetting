package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os/user"
	"sort"
)

type Variable struct {
	Index int
	Name  string
	Value string
}

type EnvModel struct {
	envType EnvType
	walk.TableModelBase
	items []*Variable
}

func NewEnvModel(env EnvType) *EnvModel {
	m := new(EnvModel)
	m.envType = env

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
		return item.Name

	case 1:
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
func (m *EnvModel) Less(i, j int) bool {
	return m.items[i].Name < m.items[j].Name
}

func (m *EnvModel) Swap(i, j int) {
	m.items[i], m.items[j] = m.items[j], m.items[i]
}
func (m *EnvModel) ResetRows() {
	if usrEnv, err := ReadVariables(m.envType); err != nil {
		panic("Fail to read the user env")
	} else {
		m.items = make([]*Variable, 0)
		for k, v := range usrEnv {
			m.items = append(m.items, &Variable{0, k, v})
			log.Println(k, "=>", v)
		}
	}
	sort.Sort(m)
	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()
}
func main() {
	usr, _ := user.Current()
	usrModel := NewEnvModel(0)
	sysModel := NewEnvModel(1)
	var usrTableView, sysTableView *walk.TableView
	var mw *walk.MainWindow
	MainWindow{
		Title:    "Enviroment Variable",
		Size:     Size{600, 700},
		Layout:   VBox{},
		AssignTo: &mw,
		Children: []Widget{
			VSplitter{
				Children: []Widget{
					GroupBox{
						Title:  "User variables for " + usr.Username,
						Layout: VBox{},
						Children: []Widget{
							TableView{
								AssignTo:              &usrTableView,
								AlternatingRowBGColor: walk.RGB(255, 255, 224),
								ColumnsOrderable:      true,
								Columns: []TableViewColumn{
									{Title: "Variable"},
									{Title: "Value"},
								},
								LastColumnStretched: true,
								Model:               usrModel,
								OnItemActivated: func() {
									index := usrTableView.CurrentIndex()
									ShowDialog(mw, usrModel.items[index].Name, usrModel.items[index].Value)
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
								AssignTo:              &sysTableView,
								AlternatingRowBGColor: walk.RGB(255, 255, 224),
								ColumnsOrderable:      true,
								Columns: []TableViewColumn{
									{Title: "Variable"},
									{Title: "Value"},
								},
								LastColumnStretched: true,
								Model:               sysModel,
								OnItemActivated: func() {
									index := sysTableView.CurrentIndex()
									ShowDialog(mw, sysModel.items[index].Name, sysModel.items[index].Value)
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
