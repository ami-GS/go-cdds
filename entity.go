package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"

type EntityI interface {
	GetEntity() C.dds_entity_t
}

type Entity C.dds_entity_t

func (e Entity) GetEntity() C.dds_entity_t {
	return C.dds_entity_t(e)
}
func (e Entity) Delete() {
	ret := C.dds_delete(e.GetEntity())
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
}

func (e Entity) SetEnabledStatus(mask uint32) Return {
	ret := C.dds_set_enabled_status(e.GetEntity(), C.uint32_t(mask))
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Return(ret)
}
