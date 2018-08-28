package cdds

import "unsafe"

type Topic struct {
	Entity
	desc unsafe.Pointer // message descriptor
	name string
}

type TopicAccessor struct {
	Reader Reader
	Writer Writer
}
