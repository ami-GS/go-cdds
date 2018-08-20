package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"
import "unsafe"

type Reader struct {
	Entity
}

func CreateReader(participant EntityI, topic Topic, qos *QoS, listener *Listener) Reader {
	tmp := C.dds_create_reader(participant.GetEntity(), (C.dds_entity_t)(topic), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(tmp, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Reader{Entity(tmp)}
}

func (r Reader) Read(samples *unsafe.Pointer, info *SampleInfo, bufsz int, maxsz uint32) Return {
	ret := C.dds_read(r.GetEntity(), samples, (*C.dds_sample_info_t)(info), C.size_t(bufsz), C.uint32_t(maxsz))
	return Return(ret)
}

func (r Reader) CreateReadCondition(mask uint32) ReadCondition {
	return ReadCondition(C.dds_create_readcondition(r.GetEntity(), C.uint32_t(mask)))
}
