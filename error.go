package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"

func ErrorCheck(err C.dds_entity_t, flags uint8, where string) {
	C.dds_err_check(err, C.uint(flags), C.CString(where))
}
