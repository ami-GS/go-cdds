package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"

type EntityI interface {
	GetEntity() C.dds_entity_t
	IsInitialized() bool
}

type Entity C.dds_entity_t

func (e Entity) GetEntity() C.dds_entity_t {
	return C.dds_entity_t(e)
}
func (e Entity) Delete() {
	ret := C.dds_delete(e.GetEntity())
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
}

func (e Entity) GetStatusChanges() CommunicationStatus {
	var status CommunicationStatus
	ret := C.dds_get_status_changes(e.GetEntity(), (*C.uint32_t)(&status))
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return status
}

func (e Entity) SetEnabledStatus(comStatusMask CommunicationStatus) Return {
	ret := C.dds_set_enabled_status(e.GetEntity(), C.uint32_t(comStatusMask))
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Return(ret)
}

func (e Entity) IsInitialized() bool {
	return e.GetEntity() > 0
}
