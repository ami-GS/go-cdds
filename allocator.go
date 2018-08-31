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

type AllocatorI interface {
	Free(sample unsafe.Pointer)
	AllFree()
	alloc(num uint32) unsafe.Pointer
	AllocArray(num uint32) *RawArray
}

type RawAllocator struct {
	elmSize      uint32
	mut          *sync.Mutex
	allockedList map[unsafe.Pointer]unsafe.Pointer
}

func NewRawAllocator(elmSize uint32) *RawAllocator {
	return &RawAllocator{
		elmSize:      elmSize,
		mut:          new(sync.Mutex),
		allockedList: make(map[unsafe.Pointer]unsafe.Pointer),
	}
}

func (a *RawAllocator) Free(sample unsafe.Pointer) {
	a.mut.Lock()
	defer a.mut.Unlock()

	_, ok := a.allockedList[sample]
	if !ok {
		panic("unallocated location free")
	}
	delete(a.allockedList, sample)
	C.dds_free(sample)
}

func (a *RawAllocator) AllFree() {
	a.mut.Lock()
	defer a.mut.Unlock()

	for array, _ := range a.allockedList {
		C.dds_free(array)
	}
}

func (a RawAllocator) alloc(num uint32) unsafe.Pointer /*error*/ {
	allocked := unsafe.Pointer(C.dds_alloc(C.ulong(a.elmSize * num)))
	return allocked
}

func (a *RawAllocator) AllocArray(num uint32) *RawArray {
	a.mut.Lock()
	defer a.mut.Unlock()

	head := a.alloc(num)
	a.allockedList[head] = nil
	return &RawArray{
		head:      head,
		elmSize:   a.elmSize,
		arraySize: num,
	}
}

type SampleAllocator struct {
	*RawAllocator
	desc unsafe.Pointer
}

func NewSampleAllocator(desc unsafe.Pointer, elmSize uint32) *SampleAllocator {
	return &SampleAllocator{
		RawAllocator: &RawAllocator{
			elmSize:      elmSize,
			mut:          new(sync.Mutex),
			allockedList: make(map[unsafe.Pointer]unsafe.Pointer),
		},
		desc: desc,
	}
}

func (a *SampleAllocator) allocInfo(num uint32) unsafe.Pointer /*error*/ {
	var val C.dds_sample_info_t
	allocked := unsafe.Pointer(C.dds_alloc(C.ulong(unsafe.Sizeof(val) * uintptr(num))))
	return allocked
}

//override
func (a *SampleAllocator) AllocArray(num uint32) *Array {
	a.mut.Lock()
	defer a.mut.Unlock()
	sample := a.alloc(num)
	infos := a.allocInfo(num)
	a.allockedList[sample] = infos
	return NewArray(sample, infos, num, a.elmSize)
}

//override
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

//override
func (a *SampleAllocator) AllFree() {
	a.mut.Lock()
	defer a.mut.Unlock()

	for array, infos := range a.allockedList {
		C.dds_sample_free(array, (*C.dds_topic_descriptor_t)(a.desc), C.DDS_FREE_ALL)
		C.dds_free(infos)
	}
}
