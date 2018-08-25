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

	participant := cdds.CreateParticipant(cdds.DomainDefault, nil, nil)
	participant.CreateTopic(unsafe.Pointer(&C.HelloWorldData_Msg_desc), "HelloWorldData_Msg", nil, nil)
	writer := participant.CreateWriter("HelloWorldData_Msg", nil, nil)
	fmt.Println("=== [Publisher] Waiting for a reader to be discovered ...")

	writer.SetEnabledStatus(cdds.PublicationMatched)

	for {
		status := writer.GetStatusChanges()
		if status == cdds.PublicationMatched {
			break
		}
		cdds.SleepFor(time.Millisecond * 20)
	}
	msg.userID = 1
	msg.message = C.CString("Hello World!")
	fmt.Println("=== [Publisher] Writing : ")
	fmt.Printf("Message (%d, %s)\n", msg.userID, C.GoString(msg.message))
	writer.Write(unsafe.Pointer(&msg))

}
