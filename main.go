package main

import (
	"fmt"
	. "github.com/ZhuBicen/go-winapi" 
	//"github.com/AllenDang/w32"
	"syscall"
	//"bytes"
	//"encoding/binary"
	"unsafe"
)

func main() {
	var  hkey HKEY
	RegOpenKeyEx(HKEY_CURRENT_USER, syscall.StringToUTF16Ptr(`Environment`), 0, KEY_READ, &hkey)
	//RegOpenKeyEx(HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`), 0, KEY_READ, &hkey)
	for i := 0; i < 10; i++ {

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
			fmt.Println("ERROR2")
		}
		fmt.Println(syscall.UTF16ToString(valueBuffer),"=",  syscall.UTF16ToString(dataBuffer))
		//fmt.Println(syscall.UTF16ToString(valueBuffer), "=", syscall.UTF16ToString(dataBuffer), dataType, valueLen, dataLen)
	}

}