package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"
import (
	"sync"
	"unsafe"
)

type SampleAllocator struct {
	elmSize      uint32
	desc         unsafe.Pointer
	mut          *sync.Mutex
	allockedList map[unsafe.Pointer]unsafe.Pointer
}

func NewSampleAllocator(desc unsafe.Pointer, elmSize uint32) *SampleAllocator {
	return &SampleAllocator{
		elmSize:      elmSize,
		desc:         desc,
		mut:          new(sync.Mutex),
		allockedList: make(map[unsafe.Pointer]unsafe.Pointer),
	}
}

func (a *SampleAllocator) alloc(num uint32) unsafe.Pointer /*error*/ {
	allocked := unsafe.Pointer(C.dds_alloc(C.ulong(a.elmSize * num)))
	return allocked
}
func (a *SampleAllocator) allocInfo(num uint32) unsafe.Pointer /*error*/ {
	var val C.dds_sample_info_t
	allocked := unsafe.Pointer(C.dds_alloc(C.ulong(unsafe.Sizeof(val) * uintptr(num))))
	return allocked
}

func (a *SampleAllocator) AllocArray(num uint32) *Array {
	a.mut.Lock()
	defer a.mut.Unlock()
	sample := a.alloc(num)
	infos := a.allocInfo(num)
	a.allockedList[sample] = infos
	return NewArray(sample, infos, num, a.elmSize)
}

func (a *SampleAllocator) Free(sample unsafe.Pointer) /*error*/ {
	a.mut.Lock()
	defer a.mut.Unlock()

	infos, ok := a.allockedList[sample]
	if !ok {
		panic("unallocated location free")
	}
	delete(a.allockedList, sample)
	C.dds_sample_free(sample, (*C.dds_topic_descriptor_t)(a.desc), C.DDS_FREE_ALL)
	C.dds_free(infos)
}

func (a *SampleAllocator) AllFree() {
	a.mut.Lock()
	defer a.mut.Unlock()

	for array, infos := range a.allockedList {
		C.dds_sample_free(array, (*C.dds_topic_descriptor_t)(a.desc), C.DDS_FREE_ALL)
		C.dds_free(infos)
	}
}
