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
	topicNameToEntity *map[string]*Topic
	topicInfos        *map[*Topic]*TopicAccessor

	// TODO: currently participant:pub/sub = 1:1, but should be 1:n/m
	Publishers  []Publisher
	Subscribers []Subscriber
	WaitSets    []WaitSet
}

func CreateParticipant(domainID DomainID, qos *QoS, listener *Listener) (*Participant, error) {
	tmp := C.dds_create_participant((C.dds_domainid_t)(domainID), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if tmp < 0 {
		return nil, CddsErrorType(tmp)
	}

	topicNameToEntity := make(map[string]*Topic)
	topicInfos := make(map[*Topic]*TopicAccessor)
	return &Participant{
		Entity:            Entity{ent: tmp, qos: qos},
		topicNameToEntity: &topicNameToEntity,
		topicInfos:        &topicInfos,
	}, nil
}

func (p *Participant) CreateTopic(desc unsafe.Pointer, name string, qos *QoS, listener *Listener) (*Topic, error) {
	if _, ok := (*p.topicNameToEntity)[name]; ok {
		// error or ignore if qos and listener is same?
	}

	tmp := C.dds_create_topic(p.GetEntity(), (*C.dds_topic_descriptor_t)(desc), C.CString(name), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if tmp < 0 {
		return nil, CddsErrorType(tmp)
	}

	topic := &Topic{
		Entity: Entity{ent: tmp, qos: qos},
		desc:   desc,
		name:   name,
	}

	(*p.topicNameToEntity)[name] = topic
	(*p.topicInfos)[topic] = &TopicAccessor{}
	return topic, nil
}

func (p *Participant) GetOrCreateTopic(desc unsafe.Pointer, name string, qos *QoS, listener *Listener) (*Topic, error) {
	if topic, ok := (*p.topicNameToEntity)[name]; ok {
		// TODO: check qos and listener whether these are same
		return topic, nil
	}
	return p.CreateTopic(desc, name, qos, listener)
}

func (p *Participant) CreateReader(topic interface{}, elmSize uint32, qos *QoS, listener *Listener) (*Reader, error) {
	var topicEntity *Topic
	switch t := topic.(type) {
	case string:
		topicEntity = (*p.topicNameToEntity)[t]
	case *Topic:
		topicEntity = t
	default:
		panic("1st argument of CreateReader need to be string or *cdds.Topic")
	}
	tmp := C.dds_create_reader(p.GetEntity(), topicEntity.GetEntity(), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if tmp < 0 {
		return nil, CddsErrorType(tmp)
	}

	if ac, ok := (*p.topicInfos)[topicEntity]; ok {
		ac.Reader = Reader{
			Entity:    Entity{ent: tmp, qos: qos},
			allocator: NewSampleAllocator(topicEntity.desc, elmSize),
		}
		return &ac.Reader, nil
	}
	panic("topic was not created")
}

func (p *Participant) GetOrCreateReader(topic interface{}, elmSize uint32, qos *QoS, listener *Listener) (*Reader, error) {
	var topicEntity *Topic
	var ok bool
	switch t := topic.(type) {
	case string:
		topicEntity, ok = (*p.topicNameToEntity)[t]
		if !ok {
			panic("topic was not created")
		}
	case *Topic:
		topicEntity = t
	default:
		panic("1st argument of CreateWriter need to be string or *cdds.Topic")
	}

	acc, ok := (*p.topicInfos)[topicEntity]
	if !ok {
		panic("topic was not created")
	}
	if !acc.Reader.IsInitialized() {
		return p.CreateReader(topic, elmSize, qos, listener)
	}

	return &acc.Reader, nil
}

func (p *Participant) CreateWriter(topic interface{}, qos *QoS, listener *Listener) (*Writer, error) {
	var topicEntity *Topic
	switch t := topic.(type) {
	case string:
		topicEntity = (*p.topicNameToEntity)[t]
	case *Topic:
		topicEntity = t
	default:
		panic("1st argument of CreateWriter need to be string or cdds.Topic")
	}
	tmp := C.dds_create_writer(p.GetEntity(), topicEntity.GetEntity(), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if tmp < 0 {
		return nil, CddsErrorType(tmp)
	}

	if ac, ok := (*p.topicInfos)[topicEntity]; ok {
		ac.Writer = Writer{Entity{ent: tmp, qos: qos}}
		return &ac.Writer, nil
	}
	panic("topic was not created")
}

func (p *Participant) GetOrCreateWriter(topic interface{}, qos *QoS, listener *Listener) (*Writer, error) {
	var topicEntity *Topic
	var ok bool
	switch t := topic.(type) {
	case string:
		topicEntity, ok = (*p.topicNameToEntity)[t]
		if !ok {
			panic("topic was not created")
		}
	case *Topic:
		topicEntity = t
	default:
		panic("1st argument of CreateWriter need to be string or *cdds.Topic")
	}

	acc, ok := (*p.topicInfos)[topicEntity]
	if !ok {
		panic("topic was not created")
	}
	if !acc.Writer.IsInitialized() {
		return p.CreateWriter(topic, qos, listener)
	}

	return &acc.Writer, nil
}

func (p *Participant) GeTopicWriter(topicString string) (*Writer, bool) {
	entity, ok := (*p.topicNameToEntity)[topicString]
	if !ok {
		return nil, false
	}
	return &(*p.topicInfos)[entity].Writer, true
}

func (p *Participant) GeTopicAccessor(topicString string) (*TopicAccessor, bool) {
	entity, ok := (*p.topicNameToEntity)[topicString]
	if !ok {
		return nil, false
	}
	return (*p.topicInfos)[entity], true
}

func (p *Participant) CreatePublisher(qos *QoS, listener *Listener) (*Publisher, error) {
	tmp := C.dds_create_publisher(p.GetEntity(), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if tmp < 0 {
		return nil, CddsErrorType(tmp)
	}

	pub := Publisher(Participant{
		Entity:            Entity{ent: tmp, qos: qos},
		topicNameToEntity: p.topicNameToEntity,
		topicInfos:        p.topicInfos,
	})

	p.Publishers = append(p.Publishers, pub)
	return &pub, nil
}

func (p *Participant) CreateSubscriber(qos *QoS, listener *Listener) (*Subscriber, error) {
	tmp := C.dds_create_subscriber(p.GetEntity(), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if tmp < 0 {
		return nil, CddsErrorType(tmp)
	}

	sub := Subscriber(Participant{
		Entity:            Entity{ent: tmp, qos: qos},
		topicNameToEntity: p.topicNameToEntity,
		topicInfos:        p.topicInfos,
	})
	p.Subscribers = append(p.Subscribers, sub)
	return &sub, nil
}

func (p *Participant) CreateWaitSet() (*WaitSet, error) {
	tmp := C.dds_create_waitset(p.GetEntity())
	if tmp < 0 {
		return nil, CddsErrorType(tmp)
	}
	var sample C.dds_attach_t
	wait := WaitSet{
		Entity:    Entity{ent: tmp, qos: nil},
		allocator: NewRawAllocator(uint32(unsafe.Sizeof(sample))),
	}
	p.WaitSets = append(p.WaitSets, wait)
	return &wait, nil
}

func (p *Participant) Delete() error {
	for topic, accessor := range *p.topicInfos {
		if accessor.Reader.IsInitialized() {
			accessor.Reader.delete()
		}
		if accessor.Writer.IsInitialized() {
			// only for qos.delete()
			err := accessor.Writer.delete()
			if err != nil {
				return err
			}
		}
		// only for qos.delete()
		topic.delete()
	}
	for _, waitSet := range p.WaitSets {
		err := waitSet.delete()
		if err != nil {
			return err
		}
	}

	// Delete of participant propagete writer/reader/pub/sub entity implicitly
	return p.Entity.delete()
}
