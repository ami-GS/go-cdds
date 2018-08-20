package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"

type Publisher Entity

func CreatePublisher(p Participant, qos *QoS, listener *Listener) Publisher {
	pub := C.dds_create_publisher(p.GetEntity(), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(pub, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Publisher(pub)
}
