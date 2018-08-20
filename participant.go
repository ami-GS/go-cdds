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
