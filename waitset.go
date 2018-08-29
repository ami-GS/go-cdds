package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"
import "time"

type WaitSet Entity

func (w WaitSet) Wait(wsresults *Attach, size int, d time.Duration) {
	ret := C.dds_waitset_wait(C.dds_entity_t(w), (*C.dds_attach_t)(wsresults), C.size_t(size), C.dds_duration_t(int64(d)))
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
}

func (w WaitSet) Attach(entity EntityI, arg Entity) error {
	ret := C.dds_waitset_attach(C.dds_entity_t(w), entity.GetEntity(), C.dds_attach_t(arg))
	if ret < 0 {
		return CddsErrorType(ret)
	}
	return nil
}

func (w WaitSet) Detach(entity EntityI) {
	C.dds_waitset_detach(C.dds_entity_t(w), entity.GetEntity())
}
