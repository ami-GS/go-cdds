package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"
import "unsafe"

type ArrayI interface {
	At(idx int) unsafe.Pointer
	ForEach(fn func(unsafe.Pointer))
}

type RawArray struct {
	head      unsafe.Pointer
	elmSize   uint32
	arraySize uint32
}

func (a RawArray) At(idx int) unsafe.Pointer {
	if uint32(idx) >= a.arraySize {
		panic("segmentation fault")
	}

	return unsafe.Pointer(uintptr(a.head) + uintptr(uint32(idx)*a.elmSize))
}

func (a RawArray) ForEach(fn func(unsafe.Pointer)) {
	for i := 0; uint32(i) < a.arraySize; i++ {
		fn(a.At(i))
	}
}

type SampleInfo C.dds_sample_info_t

func (info *SampleInfo) IsValid() bool {
	return bool((*C.dds_sample_info_t)(info).valid_data)
}

type Array struct {
	*RawArray
	infoElmSize uint32
	infos       unsafe.Pointer
}

// Should be called from allocator
func NewArray(arrayHead unsafe.Pointer, infoHead unsafe.Pointer, bufSize uint32, elmSize uint32) *Array {
	var info C.dds_sample_info_t
	return &Array{
		RawArray: &RawArray{
			head:      arrayHead,
			elmSize:   elmSize,
			arraySize: bufSize,
		},
		infoElmSize: uint32(unsafe.Sizeof(info)),
		infos:       infoHead,
	}
}

func (a Array) InfoAt(idx int) *SampleInfo {
	if uint32(idx) >= a.arraySize {
		panic("segmentation fault")
	}
	return (*SampleInfo)(unsafe.Pointer(uintptr(a.infos) + uintptr(uint32(idx)*a.infoElmSize)))
}

func (a Array) IsValidAt(idx int) bool {
	return a.InfoAt(idx).IsValid()
}
