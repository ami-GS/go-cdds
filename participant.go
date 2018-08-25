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

// TODO: participant to be interface? ParticipantI
type Participant struct {
	Entity
	topicEntityToName map[Topic]string
	topicNameToEntity map[string]Topic
	topicInfos        map[Topic]*TopicAccessor

	// TODO: currently participant:pub/sub = 1:1, but should be 1:n/m
	Publisher  Publisher
	Subscriber Subscriber
}

func CreateParticipant(domainID DomainID, qos *QoS, listener *Listener) (*Participant, error) {
	tmp := C.dds_create_participant((C.dds_domainid_t)(domainID), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if tmp < 0 {
		return nil, CddsErrorType(tmp)
	}

	ErrorCheck(tmp, C.DDS_CHECK_REPORT|C.DDS_CHECK_EXIT, "tmp where")
	return &Participant{
		Entity:            Entity(tmp),
		topicEntityToName: make(map[Topic]string),
		topicNameToEntity: make(map[string]Topic),
		topicInfos:        make(map[Topic]*TopicAccessor),
	}, nil
}

func (p *Participant) CreateTopic(desc unsafe.Pointer, name string, qos *QoS, listener *Listener) (*Topic, error) {
	if _, ok := p.topicNameToEntity[name]; ok {
		// error or ignore if qos and listener is same?
	}

	tmp := C.dds_create_topic(p.GetEntity(), (*C.dds_topic_descriptor_t)(desc), C.CString(name), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if tmp < 0 {
		return nil, CddsErrorType(tmp)
	}

	topic := &Topic{
		Entity: Entity(tmp),
		desc:   desc,
	}

	p.topicEntityToName[topic] = name
	p.topicNameToEntity[name] = topic
	p.topicInfos[topic] = &TopicAccessor{}
	return topic, nil
}

func (p *Participant) GetOrCreateTopic(desc unsafe.Pointer, name string, qos *QoS, listener *Listener) (*Topic, error) {
	if topic, ok := p.topicNameToEntity[name]; ok {
		// TODO: check qos and listener whether these are same
		return topic, nil
	}
	return p.CreateTopic(desc, name, qos, listener)
}

func (p *Participant) CreateReader(topic interface{}, elmSize uint32, qos *QoS, listener *Listener) (*Reader, error) {
	var topicEntity *Topic
	switch t := topic.(type) {
	case string:
		topicEntity = p.topicNameToEntity[t]
	case *Topic:
		topicEntity = t
	default:
		panic("1st argument of CreateReader need to be string or *cdds.Topic")
	}
	tmp := C.dds_create_reader(p.GetEntity(), topicEntity.GetEntity(), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if tmp < 0 {
		return nil, CddsErrorType(tmp)
	}

	if ac, ok := p.topicInfos[topicEntity]; ok {
		ac.Reader = Reader{
			Entity:    Entity(tmp),
			allocator: NewSampleAllocator(topicEntity.desc, size),
		}
		return &ac.Reader, nil
	}
	panic("topic was not created")
}

func (p *Participant) CreateWriter(topic interface{}, qos *QoS, listener *Listener) (*Writer, error) {
	var topicEntity *Topic
	switch t := topic.(type) {
	case string:
		topicEntity = p.topicNameToEntity[t]
	case *Topic:
		topicEntity = t
	default:
		panic("1st argument of CreateWriter need to be string or cdds.Topic")
	}
	tmp := C.dds_create_writer(p.GetEntity(), topicEntity.GetEntity(), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if tmp < 0 {
		return nil, CddsErrorType(tmp)
	}

	if ac, ok := p.topicInfos[topicEntity]; ok {
		ac.Writer = Writer{Entity(tmp)}
		return &ac.Writer, nil
	}
	panic("topic was not created")
}

func (p *Participant) GetOrCreateWriter(topic interface{}, qos *QoS, listener *Listener) (*Writer, error) {
	var topicEntity *Topic
	var ok bool
	switch t := topic.(type) {
	case string:
		topicEntity, ok = p.topicNameToEntity[t]
		if !ok {
			panic("topic was not created")
		}
	case *Topic:
		topicEntity = t
	default:
		panic("1st argument of CreateWriter need to be string or *cdds.Topic")
	}

	acc, ok := p.topicInfos[topicEntity]
	if !ok {
		panic("topic was not created")
	}
	if !acc.Writer.IsInitialized() {
		return p.CreateWriter(topic, qos, listener)
	}

	return &acc.Writer, nil
}

func (p *Participant) GeTopicWriter(topicString string) (*Writer, bool) {
	entity, ok := p.topicNameToEntity[topicString]
	if !ok {
		return nil, false
	}
	return &p.topicInfos[entity].Writer, true
}

func (p *Participant) GeTopicAccessor(topicString string) (*TopicAccessor, bool) {
	entity, ok := p.topicNameToEntity[topicString]
	if !ok {
		return nil, false
	}
	return p.topicInfos[entity], true
}

func (p *Participant) CreatePublisher(qos *QoS, listener *Listener) error {
	pub := C.dds_create_publisher(p.GetEntity(), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if pub < 0 {
		return CddsErrorType(pub)
	}
	p.Publisher = Publisher(pub)
	return nil
}

func (p *Participant) CreateSubscriber(qos *QoS, listener *Listener) error {
	sub := C.dds_create_subscriber(p.GetEntity(), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if sub < 0 {
		return CddsErrorType(sub)
	}
	p.Subscriber = Subscriber(sub)
	return nil
}

func (p *Participant) Delete() {
	for _, accessor := range p.topicInfos {
		if accessor.Reader.IsInitialized() {
			accessor.Reader.Delete()
		}
	}
	p.Entity.Delete()
}
