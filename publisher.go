package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"

type Publisher C.dds_entity_t

func CreatePublisher(p Participant, qos *QoS, listener *Listener) Publisher {
	pub := C.dds_create_publisher(C.dds_entity_t(p), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(pub, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Publisher(pub)
}

func (p Publisher) GetEntity() C.dds_entity_t {
	return C.dds_entity_t(p)
}
