package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"
import (
	"unsafe"
)

type Participant struct {
	Entity
	Topic      Topic
	Reader     Reader
	Writer     Writer
	Publisher  Publisher
	Subscriber Subscriber
}

func CreateParticipant(domainID DomainID, qos *QoS, listener *Listener) Participant {
	tmp := C.dds_create_participant((C.dds_domainid_t)(domainID), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(tmp, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return Participant{
		Entity: Entity(tmp),
	}
}

func (p *Participant) CreateTopic(desc unsafe.Pointer, name string, qos *QoS, listener *Listener) {
	tmp := C.dds_create_topic(p.GetEntity(), (*C.dds_topic_descriptor_t)(desc), C.CString(name), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))

	ErrorCheck(tmp, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	p.Topic = Topic(tmp)
}

func (p *Participant) CreateReader(qos *QoS, listener *Listener) {
	tmp := C.dds_create_reader(p.GetEntity(), (C.dds_entity_t)(p.Topic), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(tmp, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	p.Reader = Reader{Entity(tmp)}
}

func (p *Participant) CreateWriter(qos *QoS, listener *Listener) {
	tmp := C.dds_create_writer(p.GetEntity(), (C.dds_entity_t)(p.Topic), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(tmp, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	p.Writer = Writer{Entity(tmp)}
}

func (p *Participant) CreatePublisher(qos *QoS, listener *Listener) {
	pub := C.dds_create_publisher(p.GetEntity(), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(pub, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	p.Publisher = Publisher(pub)
}

func (p *Participant) CreateSubscriber(qos *QoS, listener *Listener) {
	sub := C.dds_create_subscriber(p.GetEntity(), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	ErrorCheck(sub, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	p.Subscriber = Subscriber(sub)
}
