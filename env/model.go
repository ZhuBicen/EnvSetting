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

type Model struct {
	envType EnvType
	walk.TableModelBase
	items []*Variable

	deletedItems map[string]bool
}

func NewModel(env EnvType) *Model {
	m := new(Model)
	m.envType = env
	m.deletedItems = make(map[string]bool)
	m.ResetRows()
	return m

}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *Model) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *Model) Value(row, col int) interface{} {
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
func (m *Model) Checked(row int) bool {
	return false
}

func (m *Model) Len() int {
	return len(m.items)
}

// Called by the TableView to retrieve an item image.
func (m *Model) Image(row int) interface{} {
	return nil
}
func (m *Model) Less(i, j int) bool {
	return m.items[i].Name < m.items[j].Name
}

func (m *Model) Swap(i, j int) {
	m.items[i], m.items[j] = m.items[j], m.items[i]
}
func (m *Model) ResetRows() {
	if usrEnv, err := LoadVariables(m.envType); err != nil {
		panic("Fail to read the user env")
	} else {
		m.items = make([]*Variable, 0)
		for k, v := range usrEnv {
			m.items = append(m.items, &Variable{k, v})
		}
	}
	sort.Sort(m)
	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()
}
func (m *Model) AddVariable(name string, value string) bool {
	if e, _ := m.exists(name); e {
		return false
	}
	m.items = append(m.items, &Variable{name, value})
	sort.Sort(m)
	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()
	return true
}
func (m *Model) GetVariable(index int) *Variable {
	return m.items[index]
}
func (m *Model) EditVariable(name string, value string) bool {
	if e, v := m.exists(name); e {
		v.Value = value
		sort.Sort(m)
		// Notify TableView and other interested parties about the reset.
		m.PublishRowsReset()
		return true
	}
	return false
}

func (m *Model) DeleteVariable(name string) bool {
	for i, variable := range m.items {
		if variable.Name == name {
			m.deletedItems[name] = true
			m.items = append(m.items[:i], m.items[i+1:]...)
			sort.Sort(m)
			// Notify TableView and other interested parties about the reset.
			m.PublishRowsReset()
			return true
		}
	}
	return false
}
func (m *Model) exists(name string) (bool, *Variable) {
	for i, variable := range m.items {
		if variable.Name == name {
			return true, m.items[i]
		}
	}
	return false, nil
}

func (m *Model) Apply() error {
	for _, variable := range m.items {
		value, err := ReadVariable(m.envType, variable.Name)
		if err != nil {
			// this is a new variable
			log.Println("New Variable:", variable.Name, "=>", variable.Value)
			goto editOrNew
		}
		if variable.Value == value {
			continue
		}
	editOrNew:
		// user changed the variable
		log.Println("Change Variable:", variable.Name, "=>", variable.Value)
		if err := EditVariable(m.envType, variable.Name, variable.Value); err != nil {
			return err
		}
	}
	for name, _ := range m.deletedItems {
		log.Println("Delete Variable:", name)
		if err := DeleteVariable(m.envType, name); err != nil {
			return err
		}
	}
	return nil
}
