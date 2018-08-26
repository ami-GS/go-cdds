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

//type Topic Entity
type ReadCondition Entity

//TODO: can be error?
type Return C.dds_return_t

type DomainID C.dds_domainid_t
type QoS C.dds_qos_t
type Listener C.dds_listener_t

//type TopicDescriptor C.dds_topic_descriptor_t
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
