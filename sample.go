package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"
import "unsafe"

type SampleInfo C.dds_sample_info_t

func (info *SampleInfo) IsValid() bool {
	return bool((*C.dds_sample_info_t)(info).valid_data)
}

type Array struct {
	arr     unsafe.Pointer
	infos   unsafe.Pointer
	elmSize uint32
	bufSize uint32
}

// Should be called from allocator
func NewArray(arrayHead unsafe.Pointer, infoHead unsafe.Pointer, bufSize uint32, elmSize uint32) *Array {
	return &Array{
		arr:     arrayHead,
		infos:   infoHead,
		elmSize: elmSize,
		bufSize: bufSize,
	}
}

func (a Array) At(idx int) unsafe.Pointer {
	if uint32(idx) >= a.bufSize {
		panic("segmentation fault")
	}

	return unsafe.Pointer(uintptr(a.arr) + uintptr(uint32(idx)*a.bufSize))
}

func (a Array) InfoAt(idx int) *SampleInfo {
	if uint32(idx) >= a.bufSize {
		panic("segmentation fault")
	}

	return (*SampleInfo)(unsafe.Pointer(uintptr(a.infos) + uintptr(uint32(idx)*a.bufSize)))
}

func (a Array) IsValidAt(idx int) bool {
	return a.InfoAt(idx).IsValid()
}

func (a Array) ForEach(fn func(unsafe.Pointer)) {
	for i := 0; uint32(i) < a.bufSize; i++ {
		fn(a.At(i))
	}
}
