package env

import (
	"github.com/lxn/walk"
	"log"
	"sort"
)

type EnvType int

type Variable struct {
	Name  string
	Value string
}

type EnvModel struct {
	envType EnvType
	walk.TableModelBase
	items []*Variable
}

func NewModel(env EnvType) *EnvModel {
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
			m.items = append(m.items, &Variable{k, v})
			log.Println(k, "=>", v)
		}
	}
	sort.Sort(m)
	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()
}
func (m *EnvModel) AddVariable(name string, value string) bool {
	if e, _ := m.exists(name); e {
		return false
	}
	m.items = append(m.items, &Variable{name, value})
	return true
}
func (m *EnvModel) GetVariable(index int) *Variable {
	return m.items[index]
}
func (m *EnvModel) EditVariable(name string, value string) bool {
	if e, v := m.exists(name); e {
		v.Value = value
		return true
	}
	return false
}

func (m *EnvModel) DeleteVariable(name string) bool {
	for i, variable := range m.items {
		if variable.Name == name {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return true
		}
	}
	return false
}
func (m *EnvModel) exists(name string) (bool, *Variable) {
	for i, variable := range m.items {
		if variable.Name == name {
			return true, m.items[i]
		}
	}
	return false, nil
}
