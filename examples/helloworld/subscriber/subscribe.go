package main

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc ${SRCDIR}/../HelloWorldData.o
#include "ddsc/dds.h"
#include "../HelloWorldData.h"
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"

	cdds "github.com/ami-GS/go-cdds"
)

const MAX_SAMPLES = 1

func main() {
	var samples [MAX_SAMPLES]unsafe.Pointer
	var infos [MAX_SAMPLES]cdds.SampleInfo
	var msg *C.HelloWorldData_Msg
	participant := cdds.CreateParticipant(cdds.DomainDefault, nil, nil)
	participant.CreateTopic(unsafe.Pointer(&C.HelloWorldData_Msg_desc), "HelloWorldData_Msg", nil, nil)
	qos := cdds.CreateQoS()
	qos.SetReliability(cdds.Reliability(C.DDS_RELIABILITY_RELIABLE), time.Second*10)
	participant.CreateReader(qos, nil)
	qos.Delete()
	fmt.Println("=== [Subscriber] Waiting for sample ...")

	allocator := cdds.NewSampleAllocator(unsafe.Pointer(&C.HelloWorldData_Msg_desc), unsafe.Sizeof(*msg))
	samples[0] = allocator.Alloc()

	for {
		_ = participant.Reader.Read(&samples[0], &infos[0], MAX_SAMPLES, MAX_SAMPLES)
		if infos[0].IsValid() {
			/* Print Message. */
			msg = (*C.HelloWorldData_Msg)(samples[0])
			fmt.Print("=== [Subscriber] Received : ")
			fmt.Printf("Message (%d, %s)\n", msg.userID, C.GoString(msg.message))
			break
		}
		cdds.SleepFor(time.Millisecond * 20)
	}

	allocator.Free(samples[0])
	participant.Delete()

}
