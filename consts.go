package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"

const DomainDefault = C.DDS_DOMAIN_DEFAULT

type Reliability C.dds_reliability_kind_t

const (
	BestEffort Reliability = C.DDS_RELIABILITY_BEST_EFFORT
	Reliable   Reliability = C.DDS_RELIABILITY_RELIABLE
)

type CommunicationStatus C.uint32_t

const (
	CommunicationNil   CommunicationStatus = 0
	PublicationMatched CommunicationStatus = C.DDS_PUBLICATION_MATCHED_STATUS
)

type ReadConditionState uint32

const (
	ReadSampleState    ReadConditionState = C.DDS_READ_SAMPLE_STATE
	NotReadSampleState ReadConditionState = C.DDS_NOT_READ_SAMPLE_STATE
	AnySampleState     ReadConditionState = C.DDS_ANY_SAMPLE_STATE

	NewViewState    ReadConditionState = C.DDS_NEW_VIEW_STATE
	NotNewViewState ReadConditionState = C.DDS_NOT_NEW_VIEW_STATE
	AnyViewState    ReadConditionState = C.DDS_ANY_VIEW_STATE

	AliveInstanceState             ReadConditionState = C.DDS_ALIVE_INSTANCE_STATE
	NotAliveDisposedInstanceState  ReadConditionState = C.DDS_NOT_ALIVE_DISPOSED_INSTANCE_STATE
	NotAliveNoWritersInstanceState ReadConditionState = C.DDS_NOT_ALIVE_NO_WRITERS_INSTANCE_STATE
	AnyInstanceState               ReadConditionState = C.DDS_ANY_INSTANCE_STATE

	AnyState ReadConditionState = C.DDS_ANY_STATE
)

type InstanceState uint32

const (
	IstAlive             InstanceState = C.DDS_ALIVE_INSTANCE_STATE
	IstNotAliveDisposed  InstanceState = C.DDS_NOT_ALIVE_DISPOSED_INSTANCE_STATE
	IstNotAliceNoWriters InstanceState = C.DDS_NOT_ALIVE_NO_WRITERS_INSTANCE_STATE
)
