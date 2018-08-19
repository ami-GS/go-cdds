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

type Participant C.dds_entity_t
type Topic C.dds_entity_t
type Writer C.dds_entity_t
type Reader C.dds_entity_t

//TODO: can be error?
type Return C.dds_return_t

type DomainID C.dds_domainid_t
type QoS C.dds_qos_t
type Reliability C.dds_reliability_kind_t
type Listener C.dds_listener_t
type TopicDescriptor C.dds_topic_descriptor_t
type SampleInfo C.dds_sample_info_t
type Sample unsafe.Pointer

func CreateParticipant(domainID DomainID, qos *QoS, listener *Listener) Participant {
	tmp := C.dds_create_participant((C.dds_domainid_t)(domainID), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(tmp, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Participant(tmp)
}

func (p Participant) Delete() {
	ret := C.dds_delete(C.dds_entity_t(p))
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
}

func CreateTopic(participant Participant, desc unsafe.Pointer, name string, qos *QoS, listener *Listener) Topic {
	tmp := C.dds_create_topic(C.dds_entity_t(participant), (*C.dds_topic_descriptor_t)(desc), C.CString(name), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))

	ErrorCheck(tmp, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Topic(tmp)
}

func CreateWriter(participant Participant, topic Topic, qos *QoS, listener *Listener) Writer {
	tmp := C.dds_create_writer((C.dds_entity_t)(participant), (C.dds_entity_t)(topic), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(tmp, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Writer(tmp)
}

func (w Writer) Write(data unsafe.Pointer) {
	ret := C.dds_write(C.dds_entity_t(w), data)
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
}

func CreateReader(participant Participant, topic Topic, qos *QoS, listener *Listener) Reader {
	tmp := C.dds_create_reader((C.dds_entity_t)(participant), (C.dds_entity_t)(topic), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(tmp, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Reader(tmp)
}

func (r Reader) Read(samples *unsafe.Pointer, info *SampleInfo, bufsz int, maxsz uint32) Return {
	ret := C.dds_read(C.dds_entity_t(r), samples, (*C.dds_sample_info_t)(info), C.size_t(bufsz), C.uint32_t(maxsz))
	return Return(ret)
}

func CreateQoS() *QoS {
	return (*QoS)(C.dds_qos_create())
}

func (qos *QoS) SetReliability(rel Reliability, n time.Duration) {
	C.dds_qset_reliability((*C.dds_qos_t)(qos), C.dds_reliability_kind_t(rel), C.int64_t(int64(n)))
}

func (qos *QoS) Delete() {
	C.dds_qos_delete((*C.dds_qos_t)(qos))
}

func (info *SampleInfo) IsValid() bool {
	return bool((*C.dds_sample_info_t)(info).valid_data)
}

// need class which has alocater/free for specific desc?
type SampleAllocator struct {
	size uintptr
	desc unsafe.Pointer
}

func NewSampleAllocator(desc unsafe.Pointer, size uintptr) *SampleAllocator {
	return &SampleAllocator{
		size: size,
		desc: desc,
	}
}

func (a *SampleAllocator) Alloc() unsafe.Pointer /*error*/ {
	return unsafe.Pointer(C.dds_alloc(C.ulong(a.size)))
}

func (a *SampleAllocator) Free(sample unsafe.Pointer) /*error*/ {
	C.dds_sample_free(sample, (*C.dds_topic_descriptor_t)(a.desc), C.DDS_FREE_ALL)
}

func SetEnabledStatus(writer Writer, mask uint32) Return {
	ret := C.dds_set_enabled_status(C.dds_entity_t(writer), C.uint32_t(mask))
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Return(ret)
}

func GetStatusChanges(writer Writer) uint32 {
	var status uint32
	ret := C.dds_get_status_changes(C.dds_entity_t(writer), (*C.uint32_t)(&status))
	ErrorCheck(ret, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return status
}

func SleepFor(n time.Duration) {
	C.dds_sleepfor(C.int64_t(int64(n)))
}
