package main

import (
	"errors"
	"fmt"
	. "github.com/ZhuBicen/go-winapi" 
	//"github.com/AllenDang/w32"
	"syscall"
	//"bytes"
	//"encoding/binary"
	"unsafe"
)

var usrEnv map[string]string
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
func main() {
	if usrMap, err := ReadEnvMap(0); err == nil {
		fmt.Println(usrMap)
	}
	if sysMap, err := ReadEnvMap(1); err == nil {
		fmt.Println(sysMap)
	}
}