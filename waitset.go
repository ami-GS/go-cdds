package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"
import "time"

type WaitSet struct {
	Entity
	allocator *RawAllocator
}

func (w *WaitSet) Wait(size int, d time.Duration) (*RawArray, error) {
	// TODO: to return RawArray in stead of array of Attach is not convenient?
	wsresults := w.allocator.AllocArray(uint32(size))

	ret := C.dds_waitset_wait(w.GetEntity(), (*C.dds_attach_t)(wsresults.head), C.size_t(size), C.dds_duration_t(int64(d)))
	if ret < 0 {
		return nil, CddsErrorType(ret)
	}
	return wsresults, nil

}

func (w *WaitSet) Attach(entity EntityI, arg EntityI) error {
	ret := C.dds_waitset_attach(w.GetEntity(), entity.GetEntity(), C.dds_attach_t(arg.GetEntity()))
	if ret < 0 {
		return CddsErrorType(ret)
	}
	return nil
}

func (w *WaitSet) Detach(entity EntityI) error {
	ret := C.dds_waitset_detach(w.GetEntity(), entity.GetEntity())
	if ret < 0 {
		return CddsErrorType(ret)
	}
	return nil
}

func (w *WaitSet) delete() error {
	if w.allocator != nil {
		w.allocator.AllFree()
	}
	return w.Entity.delete()
}
