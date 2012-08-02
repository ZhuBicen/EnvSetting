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


func main() {

	model := &EnvModel{}
	model.Init(0)


	model2 := &EnvModel{}
	model2.Init(1)

	walk.Initialize(walk.InitParams{PanicOnError: true})
	defer walk.Shutdown()

	myWindow, _ := walk.NewMainWindow()

	mwVbox := walk.NewVBoxLayout()
	mwVbox.SetMargins(walk.Margins{5, 5, 5, 5})
	myWindow.SetLayout(mwVbox)

	//user group
	myWindow.SetTitle("环境变量")
	grp1, _ := walk.NewGroupBox(myWindow)
	grp1.SetTitle("bzhu的用户变量")

	grp1Vbox := walk.NewVBoxLayout()
	grp1Vbox.SetMargins(walk.Margins{10, 10, 10, 20})
	grp1.SetLayout(grp1Vbox)
	tableView, _ := walk.NewTableView(grp1)
	tableView.SetModel(model)

	buttonWidgetsGrp1, _ :=  walk.NewComposite(grp1)
	btnWidgetsHboxGrp1 := walk.NewHBoxLayout()
	buttonWidgetsGrp1.SetLayout(btnWidgetsHboxGrp1)
	btnWidgetsHboxGrp1.SetSpacing(15)
	_, _ = walk.NewHSpacer(buttonWidgetsGrp1)
	newUsrEnvBtn, _ := walk.NewPushButton(buttonWidgetsGrp1)
	newUsrEnvBtn.SetText("新建")
	newEditUsrEnvBtn, _ := walk.NewPushButton(buttonWidgetsGrp1)
	newEditUsrEnvBtn.SetText("编辑")
	deleteEditUsrEnvBtn, _ := walk.NewPushButton(buttonWidgetsGrp1)
	deleteEditUsrEnvBtn.SetText("编辑")



	//sys group
	grp2, _ := walk.NewGroupBox(myWindow)
	grp2.SetTitle("系统变量")
	grp2Vbox := walk.NewVBoxLayout()
	grp2Vbox.SetMargins(walk.Margins{10, 10, 10, 20})
	grp2.SetLayout(grp2Vbox)
    tableView2, _ := walk.NewTableView(grp2)
	tableView2.SetModel(model2)	

	buttonWidgets, _ :=  walk.NewComposite(myWindow)
	btnWidgetsHbox := walk.NewHBoxLayout()
	btnWidgetsHbox.SetMargins(walk.Margins{0, 5, 0, 5})

	buttonWidgetsGrp2, _ :=  walk.NewComposite(grp2)	
	btnWidgetsHboxGrp2 := walk.NewHBoxLayout()
	buttonWidgetsGrp2.SetLayout(btnWidgetsHboxGrp2)
	btnWidgetsHboxGrp2.SetSpacing(15)
	_, _ = walk.NewHSpacer(buttonWidgetsGrp2)
	newSysEnvBtn, _ := walk.NewPushButton(buttonWidgetsGrp2)
	newSysEnvBtn.SetText("新建")
	newEditSysEnvBtn, _ := walk.NewPushButton(buttonWidgetsGrp2)
	newEditSysEnvBtn.SetText("编辑")
	deleteEditSysEnvBtn, _ := walk.NewPushButton(buttonWidgetsGrp2)
	deleteEditSysEnvBtn.SetText("编辑")

	//bottom button
	buttonWidgets.SetLayout(btnWidgetsHbox)
	btnWidgetsHbox.SetSpacing(15)
	_, _ = walk.NewHSpacer(buttonWidgets)
	okBtn, _ := walk.NewPushButton(buttonWidgets)
	okBtn.SetText("确定")
	cancelBtn, _ := walk.NewPushButton(buttonWidgets)
	cancelBtn.SetText("取消")
	

	myWindow.Show()
	myWindow.SetMinMaxSize(walk.Size{320, 240}, walk.Size{})
	myWindow.SetSize(walk.Size{400, 500})
	myWindow.Run()
}