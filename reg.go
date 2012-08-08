package main
import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
	"strings"
)
import 	. "github.com/ZhuBicen/go-winapi"

type EnvType int
const (
	USR_SUBKEY = "Environment"
	SYS_SUBKEY = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
)

func CreateVariable(etype EnvType, varName string, varValue string) error {
	var rootkey HKEY
	var subkey string

	if etype == 0 {
		rootkey = HKEY_CURRENT_USER
		subkey = USR_SUBKEY
	} else {
		rootkey = HKEY_LOCAL_MACHINE
		subkey = SYS_SUBKEY
	}

	var mykey HKEY

	if ret := RegOpenKeyEx(rootkey, syscall.StringToUTF16Ptr(subkey), 0, KEY_WRITE, &mykey); ret != ERROR_SUCCESS {
		return errors.New(fmt.Sprintf("CreateEnvVar error, RegOpenKeyEx = %d", ret))
	}
	
	dataType := REG_SZ
	if strings.Index(varValue, "%") != -1 {
		dataType = REG_EXPAND_SZ
	}
	if ret := RegSetValueEx(mykey, 
		syscall.StringToUTF16Ptr(varName),
		0, 
		dataType, 
		(*byte)(unsafe.Pointer(syscall.StringToUTF16Ptr(varValue))), 
		uint32(len(syscall.StringToUTF16(varValue)))); ret != ERROR_SUCCESS {
		return errors.New(fmt.Sprintf("CreateEnvVar error, RegSetValueEx = %d", ret))
	}
	return nil
}

func ReadVariable(etype EnvType, varName string) (string, error){
	var rootkey HKEY
	var subkey string

	if etype == 0 {
		rootkey = HKEY_CURRENT_USER
		subkey = USR_SUBKEY
	} else {
		rootkey = HKEY_LOCAL_MACHINE
		subkey = SYS_SUBKEY
	}

	var mykey HKEY

	if ret := RegOpenKeyEx(rootkey, syscall.StringToUTF16Ptr(subkey), 0, KEY_READ, &mykey); ret != ERROR_SUCCESS {
		return "", errors.New(fmt.Sprintf("ReadVariable error, RegOpenKeyEx = %d", ret))
	}
	//get the data length
	var dataLen uint32 = 0
	var dataType uint32 = 0

	if ret := RegQueryValueEx(mykey, syscall.StringToUTF16Ptr(varName), 
		nil, &dataType, nil, &dataLen); ret != ERROR_SUCCESS{
		return "", errors.New(fmt.Sprintf("ReadVariable error1, RegQueryValueEx = %d", ret))
	}

	dataBuffer := make([]uint16, dataLen/2 + 1)

	if ret := RegQueryValueEx(mykey, syscall.StringToUTF16Ptr(varName), 
		nil, &dataType, (*byte)(unsafe.Pointer(&dataBuffer[0])), &dataLen); ret != ERROR_SUCCESS{
		return "", errors.New(fmt.Sprintf("ReadVariable error2, RegQueryValueEx = %d", ret))
	}
	return syscall.UTF16ToString(dataBuffer), nil

}

func ReadVariables(etype EnvType) (map[string]string, error) {
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