package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"

type Subscriber Entity

func CreateSubscriber(p Participant, qos *QoS, listener *Listener) Subscriber {
	sub := C.dds_create_subscriber(p.GetEntity(), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(sub, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Subscriber(sub)
}
