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
	"unsafe"

	cdds "github.com/ami-GS/go-cdds"
)

func main() {
	var participant cdds.Participant
	var topic cdds.Topic
	var writer cdds.Writer
	var msg C.HelloWorldData_Msg

	participant = cdds.CreateParticipant(C.DDS_DOMAIN_DEFAULT, nil, nil)
	topic = cdds.CreateTopic(participant, unsafe.Pointer(&C.HelloWorldData_Msg_desc), "HelloWorldData_Msg", nil, nil)
	writer = cdds.CreateWriter(participant, topic, nil, nil)
	fmt.Println("=== [Publisher] Waiting for a reader to be discovered ...")

	cdds.SetEnabledStatus(writer, C.DDS_PUBLICATION_MATCHED_STATUS)

	for {
		status := cdds.GetStatusChanges(writer)
		if status == C.DDS_PUBLICATION_MATCHED_STATUS {
			fmt.Println(2, status)
			break
		}
		C.dds_sleepfor(1000000 * 20)
	}
	msg.userID = 1
	msg.message = C.CString("Hello World!")
	fmt.Println("=== [Publisher] Writing : ")
	fmt.Printf("Message (%d, %s)\n", msg.userID, C.GoString(msg.message))
	writer.Write(unsafe.Pointer(&msg))
	participant.Delete()

}
