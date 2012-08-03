package main

import (
	"errors"
	//"fmt"
	. "github.com/ZhuBicen/go-winapi"
	"github.com/lxn/walk" 
	"syscall"
	//"bytes"
	//"encoding/binary"
	"unsafe"
)


var sysEnv map[string]string

type EnvType int

func ReadEnvMap(etype EnvType) (map[string]string, error) {
	var  hkey HKEY
	envMap := make(map[string]string)
	if etype == 0 {
		RegOpenKeyEx(HKEY_CURRENT_USER, syscall.StringToUTF16Ptr(`Environment`), 0, KEY_READ, &hkey)
	}else {
		RegOpenKeyEx(HKEY_LOCAL_MACHINE, 
			syscall.StringToUTF16Ptr(`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`), 
			0, KEY_READ, &hkey)	
	}

	for i := 0; ; i++ {

		//http://msdn.microsoft.com/en-us/library/windows/desktop/ms724872(v=vs.85).aspx
		var valueLen uint32 = 256
		valueBuffer := make([]uint16, 256)

		var dataLen uint32 = 0
		var dataType uint32 = 0
		
		if ERROR_NO_MORE_ITEMS == RegEnumValue( hkey, uint32(i), &valueBuffer[0], &valueLen, 	nil, &dataType, nil, &dataLen) {
			break
		}

		dataBuffer := make([]uint16, dataLen/2 + 1)

		if ERROR_SUCCESS != RegQueryValueEx(hkey, &valueBuffer[0], nil, &dataType, (*byte)(unsafe.Pointer(&dataBuffer[0])), &dataLen){
			return nil, errors.New("ERROR2")
		}
		envMap[syscall.UTF16ToString(valueBuffer)] = syscall.UTF16ToString(dataBuffer)

	}
	return envMap, nil
}

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
	if usrEnv, err := ReadEnvMap(env); err != nil {
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
	ui dialogUI
}

func main(){

	dlg := new(EnvSettingDialog)
	if err := dlg.init(nil); err != nil {
		panic("cna't init the dialogUI")
	}
	dlg.Run()
}
