package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"
import (
	"sync"
	"time"
	"unsafe"
)

//type Topic Entity
type ReadCondition Entity

//TODO: can be error?
type Return C.dds_return_t

type DomainID C.dds_domainid_t
type QoS C.dds_qos_t
type Listener C.dds_listener_t

//type TopicDescriptor C.dds_topic_descriptor_t
type SampleInfo C.dds_sample_info_t
type Sample unsafe.Pointer
type Attach C.dds_attach_t

// originally argument is void* arg
func CreateListener(arg unsafe.Pointer) *Listener {
	return (*Listener)(C.dds_listener_create(arg))
}

func CreateQoS() *QoS {
	return (*QoS)(C.dds_qos_create())
}

func (qos *QoS) SetReliability(rel Reliability, n time.Duration) {
	C.dds_qset_reliability((*C.dds_qos_t)(qos), C.dds_reliability_kind_t(rel), C.int64_t(int64(n)))
}

func (qos *QoS) SetWriterDataLifecycle(autoDispose bool) {
	C.dds_qset_writer_data_lifecycle((*C.dds_qos_t)(qos), C.bool(autoDispose))
}

func (qos *QoS) SetPartition(num int, partitions *string) {
	C.dds_qset_partition((*C.dds_qos_t)(qos), C.uint32_t(num), (**C.char)(unsafe.Pointer(partitions)))

}

func (qos *QoS) Delete() {
	C.dds_qos_delete((*C.dds_qos_t)(qos))
}

func (info *SampleInfo) IsValid() bool {
	return bool((*C.dds_sample_info_t)(info).valid_data)
}

// need class which has alocater/free for specific desc?
type SampleAllocator struct {
	size         uintptr
	desc         unsafe.Pointer
	mut          *sync.Mutex
	allockedList []unsafe.Pointer
}

func NewSampleAllocator(desc unsafe.Pointer, size uintptr) *SampleAllocator {
	return &SampleAllocator{
		size:         size,
		desc:         desc,
		mut:          new(sync.Mutex),
		allockedList: make([]unsafe.Pointer, 0),
	}
}

func (a *SampleAllocator) Alloc(num uint32) unsafe.Pointer /*error*/ {
	a.mut.Lock()
	defer a.mut.Unlock()
	allocked := unsafe.Pointer(C.dds_alloc(C.ulong(a.size * uintptr(num))))
	a.allockedList = append(a.allockedList, allocked)
	return allocked
}

func (a *SampleAllocator) Free(sample unsafe.Pointer) /*error*/ {
	a.mut.Lock()
	defer a.mut.Unlock()

	C.dds_sample_free(sample, (*C.dds_topic_descriptor_t)(a.desc), C.DDS_FREE_ALL)
	var i int
	var pointer unsafe.Pointer
	for i, pointer = range a.allockedList {
		if pointer == sample {
			break
		}
	}
	// remove entry (change order)
	lastIdx := len(a.allockedList) - 1
	a.allockedList[i] = a.allockedList[lastIdx]
	a.allockedList[lastIdx] = nil
	a.allockedList = a.allockedList[:lastIdx]

}

func (a *SampleAllocator) AllFree() {
	a.mut.Lock()
	defer a.mut.Unlock()

	for _, allocked := range a.allockedList {
		C.dds_sample_free(allocked, (*C.dds_topic_descriptor_t)(a.desc), C.DDS_FREE_ALL)
	}
}
