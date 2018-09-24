package main

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc  ${SRCDIR}/../HelloWorldData.o
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

func main() {
	var msg C.HelloWorldData_Msg

	participant, err := cdds.CreateParticipant(cdds.DomainDefault, nil, nil)
	defer participant.Delete()
	if err != nil {
		panic(err)
	}

	_, err = participant.CreateTopic(unsafe.Pointer(&C.HelloWorldData_Msg_desc), "HelloWorldData_Msg", nil, nil)
	if err != nil {
		panic(err)
	}
	writer, err := participant.CreateWriter("HelloWorldData_Msg", nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("=== [Publisher] Waiting for a reader to be discovered ...")

	// err = writer.SearchTopic(time.Millisecond * 20)
	// if err != nil {
	// 	panic(err)
	// }

	err = writer.SetEnabledStatus(cdds.PublicationMatched)
	if err != nil {
		panic(err)
	}
	var status cdds.CommunicationStatus
	for status != cdds.PublicationMatched {
		status, err = writer.GetStatusChanges()
		if err != nil {
			panic(err)
		}
		cdds.SleepFor(time.Millisecond * 20)
	}

	msg.userID = 12343

	jsonStr := "{\"Name\":\"cyclone\", \"Age\":22}"

	msg.message = C.CString(jsonStr)

	fmt.Println("=== [Publisher] Writing : ")
	fmt.Printf("Message (%d, %s)\n", msg.userID, C.GoString(msg.message))
	writer.Write(unsafe.Pointer(&msg))
}
