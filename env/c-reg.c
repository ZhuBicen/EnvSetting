#include <Windows.h>
#include "reg.h"
#include "_cgo_export.h"

int deleteVariable(void* varName)
{
	HKEY hkey1 = NULL;

	if (ERROR_SUCCESS != RegOpenKeyEx( HKEY_CURRENT_USER, "Environment", 
		0, KEY_ALL_ACCESS, &hkey1)) {
			return -1;
	}
	LONG result = RegDeleteValueW(hkey1, (wchar_t*)(varName));
	if (ERROR_SUCCESS != result) {
		return result;
	}

	RegCloseKey(hkey1);
	return 0;
}