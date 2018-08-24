package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"
import (
	"time"
	"unsafe"
)

type Reader struct {
	Entity
	allocator *SampleAllocator
}

func (r Reader) Read(samples *unsafe.Pointer, info *SampleInfo, bufsz int, maxsz uint32) Return {
	ret := C.dds_read(r.GetEntity(), samples, (*C.dds_sample_info_t)(info), C.size_t(bufsz), C.uint32_t(maxsz))
	return Return(ret)
}

func (r Reader) CreateReadCondition(mask uint32) ReadCondition {
	return ReadCondition(C.dds_create_readcondition(r.GetEntity(), C.uint32_t(mask)))
}

func (r Reader) Delete() {
	if r.allocator != nil {
		r.allocator.AllFree()
	}
	// reader entity will be deleted by participant, no need to call from here
	//r.Entity.Delete()
}
