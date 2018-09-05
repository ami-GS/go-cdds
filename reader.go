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

type ReadCondition struct {
	Entity
}

type Reader struct {
	Entity
	allocator      *SampleAllocator
	readConditions []ReadCondition
}

// take == false just return copy of data, take == true removes data after reading
// https://github.com/eclipse/cyclonedds/issues/17
func (r Reader) Read(samples *unsafe.Pointer, info *SampleInfo, bufsz int, maxsz uint32, take bool) Return {
	var ret C.dds_entity_t
	if take {
		ret = C.dds_take(r.GetEntity(), samples, (*C.dds_sample_info_t)(info), C.size_t(bufsz), C.uint32_t(maxsz))
	} else {
		ret = C.dds_read(r.GetEntity(), samples, (*C.dds_sample_info_t)(info), C.size_t(bufsz), C.uint32_t(maxsz))
	}
	return Return(ret)
}

func (r Reader) ReadWithCallback(bufsz int, maxsz uint32, take bool, finCh *chan error, callback func(*Array)) {
	// WARN: currently this might have issue when participant.Delete()
	// TODO: allock first, then use with loop
	// TODO: need choise this to run forever
	samples, err := r.BlockAllocRead(bufsz, maxsz, take)
	if err != nil {
		*finCh <- err
	}
	callback(samples)
	r.allocator.Free(unsafe.Pointer(samples.head))
	*finCh <- nil

}

func (r Reader) BlockAllocRead(bufsz int, maxsz uint32, take bool) (*Array, error) {
	// this is not GCed by Golang, maybe
	samples := r.allocator.AllocArray(maxsz)

	var ret C.dds_return_t
	for i := 0; i < bufsz; {
		loc := samples.At(i)
		info := (*C.dds_sample_info_t)(samples.InfoAt(i))
		if take {
			ret = C.dds_take(r.GetEntity(), &loc, info, C.size_t(bufsz), C.uint32_t(maxsz))
		} else {
			ret = C.dds_read(r.GetEntity(), &loc, info, C.size_t(bufsz), C.uint32_t(maxsz))
		}
		fmt.Println(ret, info.valid_data)
		if ret < 0 {
			return nil, CddsErrorType(ret)
		}

		if info.valid_data {
			i++
		}
		time.Sleep(time.Millisecond * 20)
	}

	return samples, nil
}

func (r Reader) AllocRead(bufsz int, maxsz uint32, take bool) (*Array, error) {
	// this is not GCed by Golang, maybe
	samples := r.allocator.AllocArray(maxsz)
	loc := samples.At(0)
	var ret C.dds_entity_t
	if take {
		ret = C.dds_take(r.GetEntity(), &loc, (*C.dds_sample_info_t)(samples.InfoAt(0)), C.size_t(bufsz), C.uint32_t(maxsz))
	} else {
		ret = C.dds_read(r.GetEntity(), &loc, (*C.dds_sample_info_t)(samples.InfoAt(0)), C.size_t(bufsz), C.uint32_t(maxsz))
	}
	if ret < 0 {
		return nil, CddsErrorType(ret)
	}
	return samples, nil
}

func (r Reader) Alloc(bufsz int) *Array {
	// this is not GCed by Golang, maybe
	return r.allocator.AllocArray(uint32(bufsz))
}

func (r Reader) ReadWithBuff(samples *Array, take bool) error {
	// this is not GCed by Golang, maybe
	if samples == nil {
		panic("buffer was not allocated")
	}

	loc := samples.At(0)
	var ret C.dds_entity_t
	if take {
		ret = C.dds_take(r.GetEntity(), &loc, (*C.dds_sample_info_t)(samples.InfoAt(0)), C.size_t(samples.elmSize), C.uint32_t(samples.elmSize))
	} else {
		ret = C.dds_read(r.GetEntity(), &loc, (*C.dds_sample_info_t)(samples.InfoAt(0)), C.size_t(samples.elmSize), C.uint32_t(samples.elmSize))
	}
	if ret < 0 {
		return CddsErrorType(ret)
	}
	return nil
}

func (r *Reader) CreateReadCondition(mask ReadConditionState) *ReadCondition {
	rd := ReadCondition{
		Entity: Entity{ent: C.dds_create_readcondition(r.GetEntity(), C.uint32_t(mask)), qos: nil},
	}
	r.readConditions = append(r.readConditions, rd)
	return &rd
}

func (r Reader) delete() error {
	if r.allocator != nil {
		r.allocator.AllFree()
	}
	if r.qos != nil {
		r.qos.delete()
	}
	for _, rdcond := range r.readConditions {
		// TODO: be careful, this might be deleted via participant.Delete(), need to check in the future
		err := rdcond.delete()
		if err != nil {
			return err
		}
	}
	return nil
	// reader entity will be deleted by participant, no need to call from here
	//r.Entity.Delete()
}
