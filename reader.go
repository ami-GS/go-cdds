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

func (r Reader) ReadWithCallback(info *SampleInfo, bufsz int, maxsz uint32, finCh *chan struct{}, callback func(unsafe.Pointer)) {
	// WARN: currently this might have issue when participant.Delete()
	samples := r.BlockAllocRead(info, bufsz, maxsz)
	callback(samples)
	r.allocator.Free(samples)
	*finCh <- struct{}{}

}

func (r Reader) BlockAllocRead(info *SampleInfo, bufsz int, maxsz uint32) unsafe.Pointer {
	// this is not GCed by Golang, maybe
	samples := r.allocator.Alloc(maxsz)

	var ret C.dds_return_t
	for {
		ret = C.dds_read(r.GetEntity(), &samples, (*C.dds_sample_info_t)(info), C.size_t(bufsz), C.uint32_t(maxsz))
		if info.IsValid() {
			break
		}
		time.Sleep(time.Millisecond * 20)
	}
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return samples
}

func (r Reader) AllocRead(info *SampleInfo, bufsz int, maxsz uint32) unsafe.Pointer {
	// this is not GCed by Golang, maybe
	samples := r.allocator.Alloc(maxsz)

	ret := C.dds_read(r.GetEntity(), &samples, (*C.dds_sample_info_t)(info), C.size_t(bufsz), C.uint32_t(maxsz))
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return samples
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
