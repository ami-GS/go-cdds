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

//type Entity C.dds_entity_t
type Entity struct {
	ent C.dds_entity_t
	qos *QoS //participantI and Topic and ? have qos
}

func (e Entity) GetEntity() C.dds_entity_t {
	return e.ent
}
func (e *Entity) delete() error {
	ret := C.dds_delete(e.GetEntity())
	if ret < 0 {
		return CddsErrorType(ret)
	}
	if e.qos != nil {
		e.qos.delete()
	}
	return nil
}

func (e Entity) GetStatusChanges() (CommunicationStatus, error) {
	var status CommunicationStatus
	ret := C.dds_get_status_changes(e.GetEntity(), (*C.uint32_t)(&status))
	if ret < 0 {
		return CommunicationNil, CddsErrorType(ret)
	}
	return status, nil
}

func (e *Entity) SetEnabledStatus(comStatusMask CommunicationStatus) error {
	ret := C.dds_set_enabled_status(e.GetEntity(), C.uint32_t(comStatusMask))
	if ret < 0 {
		return CddsErrorType(ret)
	}
	return nil
}

func (e Entity) IsInitialized() bool {
	return e.GetEntity() > 0
}

func (e Entity) Triggered() error {
	ret := C.dds_triggered(e.GetEntity())
	if ret < 0 {
		return CddsErrorType(ret)
	}
	return nil
}
