package main

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc ${SRCDIR}/../RoundTrip.c.o
#include "ddsc/dds.h"
#include "../RoundTrip.h"
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"

	cdds "github.com/ami-GS/go-cdds"
)

const MAX_SAMPLES = 10

func main() {
	part, err := cdds.CreateParticipant(cdds.DomainDefault, nil, nil)
	if err != nil {
		panic(err)
	}
	useListener := false
	var listener *cdds.Listener
	if useListener {
		listener = cdds.CreateListener(nil)
		//listener.SetDataAvailable()
	}

	// begin prepare_dds()
	_, err = part.CreateTopic(unsafe.Pointer(&C.RoundTripModule_DataType_desc), "RoundTrip", nil, nil)
	if err != nil {
		panic(err)
	}
	pubQos := cdds.CreateQoS()
	pubPartition := [1]string{"pong"}
	pubQos.SetPartition(1, &pubPartition[0])
	pub, err := part.CreatePublisher(pubQos, nil)
	if err != nil {
		panic(err)
	}

	dwQoS := cdds.CreateQoS()
	dwQoS.SetReliability(cdds.Reliable, time.Second*10)
	dwQoS.SetWriterDataLifecycle(false)
	wr, err := pub.CreateWriter("RoundTrip", dwQoS, nil)
	if err != nil {
		panic(err)
	}

	subQos := cdds.CreateQoS()
	subPartition := [1]string{"ping"}
	subQos.SetPartition(1, &subPartition[0])
	sub, err := part.CreateSubscriber(subQos, nil)
	if err != nil {
		panic(err)
	}

	drQos := cdds.CreateQoS()
	drQos.SetReliability(cdds.Reliable, time.Second*10)
	var msg *C.RoundTripModule_DataType
	rd, err := sub.CreateReader("RoundTrip", uint32(unsafe.Sizeof(*msg)), drQos, listener)
	if err != nil {
		panic(err)
	}

	waitSet, err := part.CreateWaitSet()
	if err != nil {
		panic(err)
	}

	var rdcond cdds.ReadCondition
	if listener == nil {
		rdcond = *rd.CreateReadCondition(cdds.AnyState)
		err := waitSet.Attach(rdcond, rd)
		if err != nil {
			panic(err)
		}
	} else {
		rdcond.SetEntity(0)
	}
	err = waitSet.Attach(waitSet, waitSet)
	if err != nil {
		panic(err)
	}

	fmt.Println("Waiting for samples from ping to send back...")
	// fflush();
	// end prepare_dds()

	wsresultsize := 1
	for waitSet.Triggered() == nil {
		wsResults, err := waitSet.Wait(wsresultsize, time.Duration(C.DDS_INFINITY))
		fmt.Println("78:wsResults", wsResults)
		if err != nil {
			panic(err)
		}
		if listener == nil {
			// BEGIN data_available()
			data, count, err := rd.AllocRead(MAX_SAMPLES, MAX_SAMPLES, true)
			fmt.Println(data, data.IsValidAt(0), count, "WriteTimeStamp")
			if err != nil {
				panic(err)
			}
			for i := 0; waitSet.Triggered() == nil && i < count; i++ {
				infoOne := data.InfoAt(i)
				if infoOne.GetInstanceState() == cdds.IstNotAliveDisposed {
					fmt.Println("Received termination request. Teminating.")
					waitSet.SetTrigger(true)
					break
				} else if infoOne.IsValid() {
					// no need for casting?
					//validSample := (*C.RoundTripModule_DataType)(data.At(0))
					err := wr.WriteTimeStamp(data.At(i), infoOne.GetSrcTimeStamp())
					if err != nil {
						panic(err)
					}
				}
			}
			// END data_available()
		}
		time.Sleep(time.Millisecond * 500)
	}

	// begin finalize_dds
	// remove callback
	err = rd.SetEnabledStatus(0)
	if err != nil {
		panic(err)
	}
	err = waitSet.Detach(rdcond)
	if err != nil {
		panic(err)
	}
	part.Delete()
	// end finalize_dds
}
