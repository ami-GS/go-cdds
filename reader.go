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
	allocator      *SampleAllocator
	readConditions []ReadCondition
}

func (r Reader) Read(samples *unsafe.Pointer, info *SampleInfo, bufsz int, maxsz uint32) Return {
	ret := C.dds_read(r.GetEntity(), samples, (*C.dds_sample_info_t)(info), C.size_t(bufsz), C.uint32_t(maxsz))
	return Return(ret)
}

func (r Reader) ReadWithCallback(bufsz int, maxsz uint32, finCh *chan error, callback func(*Array)) {
	// WARN: currently this might have issue when participant.Delete()
	// TODO: allock first, then use with loop
	// TODO: need choise this to run forever
	samples, err := r.BlockAllocRead(bufsz, maxsz)
	if err != nil {
		*finCh <- err
	}
	callback(samples)
	r.allocator.Free(unsafe.Pointer(samples.arr))
	*finCh <- nil

}

func (r Reader) BlockAllocRead(bufsz int, maxsz uint32) (*Array, error) {
	// this is not GCed by Golang, maybe
	samples := r.allocator.AllocArray(maxsz)

	var ret C.dds_return_t
	for i := 0; i < bufsz; {
		loc := samples.At(i)
		info := (*C.dds_sample_info_t)(samples.InfoAt(i))

		ret = C.dds_read(r.GetEntity(), &loc, info, C.size_t(bufsz), C.uint32_t(maxsz))
		if ret < 0 {
			return nil, CddsErrorType(ret)
		}

		if info.valid_data {
			i++
			break
		}
		time.Sleep(time.Millisecond * 20)
	}

	return samples, nil
}

func (r Reader) AllocRead(bufsz int, maxsz uint32) (*Array, error) {
	// this is not GCed by Golang, maybe
	samples := r.allocator.AllocArray(maxsz)
	loc := samples.At(0)

	ret := C.dds_read(r.GetEntity(), &loc, (*C.dds_sample_info_t)(samples.InfoAt(0)), C.size_t(bufsz), C.uint32_t(maxsz))
	if ret < 0 {
		return nil, CddsErrorType(ret)
	}
	return samples, nil
}

func (r *Reader) CreateReadCondition(mask uint32) *ReadCondition {
	rd := ReadCondition(C.dds_create_readcondition(r.GetEntity(), C.uint32_t(mask)))
	r.readConditions = append(r.readConditions, rd)
	return &rd
}

func (r Reader) Delete() {
	if r.allocator != nil {
		r.allocator.AllFree()
	}
	// reader entity will be deleted by participant, no need to call from here
	//r.Entity.Delete()
}
