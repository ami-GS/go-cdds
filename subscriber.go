package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"

type Subscriber Participant

func (p *Subscriber) CreateReader(topic interface{}, elmSize uint32, qos *QoS, listener *Listener) (*Reader, error) {
	return (*Participant)(p).CreateReader(topic, elmSize, qos, listener)
}
