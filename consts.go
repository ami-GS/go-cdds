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
