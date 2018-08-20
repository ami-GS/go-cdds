package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"

type Subscriber C.dds_entity_t

func CreateSubscriber(p Participant, qos *QoS, listener *Listener) Subscriber {
	sub := C.dds_create_subscriber(C.dds_entity_t(p), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(sub, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Subscriber(sub)
}
func (s Subscriber) GetEntity() C.dds_entity_t {
	return C.dds_entity_t(s)
}
