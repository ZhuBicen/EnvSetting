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
	//RegOpenKeyEx(HKEY_CURRENT_USER, syscall.StringToUTF16Ptr("Environment"), 0, KEY_READ, &hkey)
	RegOpenKeyEx(HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`), 0, KEY_READ, &hkey)
	for i := 0; i < 10; i++ {
		valueBuffer := make([]uint16, 256)
		var valueLen uint32 = 256

		dataBuffer := make([]uint16, 1024)
		var dataLen uint32 = 1024
		var dataType  uint32 = 0
		RegEnumValue(hkey, uint32(i), &valueBuffer[0], &valueLen, nil, &dataType, uintptr(unsafe.Pointer(&dataBuffer[0])), &dataLen)

		fmt.Println("DataType = ", dataType, " DataLen = ", dataLen)			
		if dataType == 1 {
			fmt.Println(syscall.UTF16ToString(valueBuffer), "=", syscall.UTF16ToString(dataBuffer))
		}else{
			fmt.Println(syscall.UTF16ToString(valueBuffer), "=", syscall.UTF16ToString(dataBuffer))
		}
		
	}

}