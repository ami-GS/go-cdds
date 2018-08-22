package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"
import "unsafe"

type Writer struct {
	Entity
}

func CreateWriter(participant EntityI, topic Topic, qos *QoS, listener *Listener) Writer {
	tmp := C.dds_create_writer(participant.GetEntity(), (C.dds_entity_t)(topic), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(tmp, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Writer{Entity(tmp)}
}

func (w Writer) Write(data unsafe.Pointer) {
	ret := C.dds_write(w.GetEntity(), data)
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
}

func (w Writer) WriteTimeStampe(data unsafe.Pointer, ts Time) {
	ret := C.dds_write_ts(w.GetEntity(), data, C.dds_time_t(ts))
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
}

func (w Writer) WriteDispose(data unsafe.Pointer) {
	ret := C.dds_writedispose(w.GetEntity(), data)
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
}
