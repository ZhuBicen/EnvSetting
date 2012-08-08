package main

import (
	"github.com/lxn/walk" 
)


var sysEnv map[string]string



type Env struct {
	Name     string
	Value     string
}

type EnvModel struct {
	items               []*Env
	rowsResetPublisher  walk.EventPublisher
	rowChangedPublisher walk.IntEventPublisher
}

// Make sure we implement all required interfaces.
var _ walk.TableModel = &EnvModel{}


func (m *EnvModel) Init(env EnvType) {
	if usrEnv, err := ReadVariables(env); err != nil {
		panic("Fail to read the user env")
	}else {
		m.items = make([]*Env, 0)
		for k, v := range usrEnv {
			m.items  = append(m.items, &Env{k, v})
		}
	}
	// if sysMap, err := ReadEnvMap(1); err != nil {
	// 	panic("Fail to read the sys env")
	// }


}

// Called by the TableView from SetModel to retrieve column information. 
func (m *EnvModel) Columns() []walk.TableColumn {
	return []walk.TableColumn{
		{Title: "变量"},
		{Title: "值"},
	}
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

// The TableView attaches to this event to synchronize its internal item count.
func (m *EnvModel) RowsReset() *walk.Event {
	return m.rowsResetPublisher.Event()
}

// The TableView attaches to this event to get notified when a row changed and
// needs to be repainted.
func (m *EnvModel) RowChanged() *walk.IntEvent {
	return m.rowChangedPublisher.Event()
}

func (m *EnvModel) ResetRows() {
	// Notify TableView and other interested parties about the reset.
	m.rowsResetPublisher.Publish()
}

type EnvSettingDialog struct{
	*walk.Dialog
	ui envSettingDialogUI
}

func main(){

	dlg := new(EnvSettingDialog)
	if err := dlg.init(nil); err != nil {
		panic("cna't init the dialogUI")
	}
	dlg.ui.usrEnvTableView.SetLastColumnStretched(true)
	dlg.ui.sysEnvTableView.SetLastColumnStretched(true)

	usrEnv := &EnvModel{}
	usrEnv.Init(0)
	sysEnv := &EnvModel{}
	sysEnv.Init(1)

	dlg.ui.usrEnvTableView.SetModel(usrEnv)
	dlg.ui.sysEnvTableView.SetModel(sysEnv)

	dlg.Run()
}
