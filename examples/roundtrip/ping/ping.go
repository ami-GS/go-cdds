package main

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc ${SRCDIR}/../RoundTrip.c.o
#include "ddsc/dds.h"
#include "../RoundTrip.h"
#define US_IN_ONE_SEC 1000000LL
*/
import "C"
import (
	"fmt"
	"sort"
	"time"
	"unsafe"

	cdds "github.com/ami-GS/go-cdds"
)

const MAX_SAMPLES = 100

type ExampleTimeStats struct {
	values     []cdds.Time
	valuesSize uint32
	valuesMax  uint32
	avg        float64
	min        cdds.Time
	max        cdds.Time
	count      uint32
	/*
	  dds_time_t * values;
	  unsigned long valuesSize;
	  unsigned long valuesMax;
	  double average;
	  dds_time_t min;
	  dds_time_t max;
	  unsigned long count;
	*/
}

const TIME_STATS_SIZE_INCREMENT = 50000

func NewExampleTimeStats() *ExampleTimeStats {
	return &ExampleTimeStats{
		values:    make([]cdds.Time, TIME_STATS_SIZE_INCREMENT),
		valuesMax: TIME_STATS_SIZE_INCREMENT,
	}
}

func (stats *ExampleTimeStats) Reset() {
	stats.values = make([]cdds.Time, stats.valuesMax)
	stats.valuesSize = 0
	stats.avg = 0
	stats.min = 0
	stats.max = 0
	stats.count = 0
}

func (stats *ExampleTimeStats) Delete() {
	// not implemented
}

func (stats *ExampleTimeStats) AddTiming(timing cdds.Time) {
	if stats.valuesSize > stats.valuesMax {
		stats.values = make([]cdds.Time, stats.valuesMax+TIME_STATS_SIZE_INCREMENT)
		stats.valuesMax += TIME_STATS_SIZE_INCREMENT
	}
	if stats.values != nil && stats.valuesSize < stats.valuesMax {
		stats.values[stats.valuesSize] = timing
		stats.valuesSize++
	}
	stats.avg = float64(cdds.Time(float64(stats.count)*stats.avg)+timing) / float64(stats.count+1)
	if stats.count == 0 || timing < stats.min {
		stats.min = timing
	}
	if stats.count == 0 || timing > stats.max {
		stats.max = timing
	}
	stats.count++
}

func (stats *ExampleTimeStats) GetMedian() float64 {
	mid := stats.valuesSize / 2
	sort.Slice(stats.values, func(i, j int) bool { return stats.values[i] < stats.values[j] })

	if stats.valuesSize%2 == 0 {
		return float64(stats.values[mid]+stats.values[stats.valuesSize/2-1]) / 2
	}
	return float64(stats.values[mid])
}

func (stats *ExampleTimeStats) Get99Percentile() int64 {
	sort.Slice(stats.values, func(i, j int) bool { return stats.values[i] < stats.values[j] })
	return int64(stats.values[stats.valuesSize-stats.valuesSize/100])
}

func main() {
	roundTrip := NewExampleTimeStats()
	writeAccess := NewExampleTimeStats()
	readAccess := NewExampleTimeStats()
	roundTripOverall := NewExampleTimeStats()
	writeAccessOverall := NewExampleTimeStats()
	readAccessOverall := NewExampleTimeStats()

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
	pubPartition := [1]string{"ping"}
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
	subPartition := [1]string{"pong"}
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

	var rdcond *cdds.ReadCondition
	if listener == nil {
		rdcond = rd.CreateReadCondition(cdds.AnyState)
		err = waitSet.Attach(rdcond, rd)
		if err != nil {
			panic(err)
		}
	} else {
		rdcond = nil
	}
	err = waitSet.Attach(waitSet, waitSet)
	if err != nil {
		panic(err)
	}
	// end prepare_dds()

	payloadSize := uint32(0)
	numSamples := uint64(0)
	timeOut := C.dds_time_t(0)
	fmt.Println("# payloadSize: ", payloadSize, " | numSamples: ", numSamples, " | timeOut: ", timeOut, "\n", payloadSize, numSamples, timeOut)

	var pubData C.RoundTripModule_DataType
	pubData.payload._length = C.uint(payloadSize)
	if payloadSize != 0 {
		pubData.payload._buffer = (*C.uchar)(C.dds_alloc(C.ulong(payloadSize)))
	} else {
		pubData.payload._buffer = nil
	}
	pubData.payload._release = true
	pubData.payload._maximum = 0
	for i := uint32(0); i < payloadSize; i++ {
		*(*C.uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(pubData.payload._buffer)) + uintptr(i*4))) = C.uchar('a')
	}
	fmt.Println("# Waiting for startup jitter to stabilise")

	wsresultsize := 1
	startTime := cdds.DdsTime() //time.Now()
	difference := cdds.Time(0)
	for err := waitSet.Triggered(); err == nil && difference < C.DDS_NSECS_IN_SEC*5; err = waitSet.Triggered() {
		//wsResults[0], err = waitSet.Wait(wsresultsize, time.Second*1)
		// wsResults is *RawArray, need to be converted to Attach when to use
		wsResults, err := waitSet.Wait(wsresultsize, time.Second*1)
		if err != nil {
			panic(err)
		}

		fmt.Println("186:wsResults", wsResults)
		if err == nil && listener == nil {
			samples, _, err := rd.AllocRead(MAX_SAMPLES, MAX_SAMPLES, true)
			fmt.Println("189:samples", samples, samples.IsValidAt(0))
			if err != nil {
				panic(err)
			}
		}
		difference = cdds.DdsTime() - startTime
	}

	warmUp := true
	if err := waitSet.Triggered(); err == nil {
		warmUp = false
		fmt.Println(`# Warm up complete.
# Round trip measurements (in us)
#             Round trip time [us]                           Write-access time [us]       Read-access time [us]
# Seconds     Count   median      min      99%%      max      Count   median      min      Count   median      min`)
	}

	roundTrip.Reset()
	writeAccess.Reset()
	readAccess.Reset()

	preWriteTime := cdds.DdsTime()
	err = wr.WriteTimeStamp(unsafe.Pointer(&pubData), preWriteTime)
	if err != nil {
		panic(err)
	}
	postWriteTime := cdds.DdsTime()

	elapsed := 0
	for i := uint64(0); waitSet.Triggered() == nil && (numSamples == 0 || i < numSamples); i++ {
		wsResults, err := waitSet.Wait(wsresultsize, time.Second*1)
		fmt.Println("220:wsResults", wsResults)
		if err == nil && listener == nil {
			// BEGIN data_available(reader, NULL);
			preTakeTime := cdds.DdsTime()
			data, num, err := rd.AllocRead(MAX_SAMPLES, MAX_SAMPLES, true)
			if err != nil {
				panic(err)
			}
			postTakeTime := cdds.DdsTime()

			diff := (postWriteTime - preWriteTime) / C.DDS_NSECS_IN_USEC
			writeAccess.AddTiming(diff)
			writeAccessOverall.AddTiming(diff)

			diff = (postTakeTime - preTakeTime) / C.DDS_NSECS_IN_USEC
			readAccess.AddTiming(diff)
			readAccessOverall.AddTiming(diff)

			infoOne := data.InfoAt(0)
			fmt.Println(num, data.IsValidAt(0), postTakeTime, infoOne.GetSrcTimeStamp())
			diff = (postTakeTime - infoOne.GetSrcTimeStamp()) / C.DDS_NSECS_IN_USEC
			roundTrip.AddTiming(diff)
			roundTripOverall.AddTiming(diff)
			if !warmUp {
				diff = (postTakeTime - startTime) / C.DDS_NSECS_IN_USEC
				if diff > C.US_IN_ONE_SEC {
					fmt.Printf("%9d %9d %8.0f %8d %8d %8d %10d %8.0f %8d %10d %8.0f %8d\n",
						elapsed+1,
						roundTrip.count,
						roundTrip.GetMedian(),
						roundTrip.min,
						roundTrip.Get99Percentile(),
						roundTrip.max,
						writeAccess.count,
						writeAccess.GetMedian(),
						writeAccess.min,
						readAccess.count,
						readAccess.GetMedian(),
						readAccess.min,
					)
					roundTrip.Reset()
					writeAccess.Reset()
					readAccess.Reset()
					startTime = cdds.DdsTime()
					elapsed++
				}
			}
			preWriteTime := cdds.DdsTime()
			err = wr.WriteTimeStamp(unsafe.Pointer(&pubData), preWriteTime)
			if err != nil {
				panic(err)
			}
			postWriteTime = cdds.DdsTime()
			// END data_available(reader, NULL);
		}
	}
	if !warmUp {
		fmt.Printf("\n%9s %9d %8.0f %8d %10d %8.0f %8d %10d %8.0f %8d\n",
			"# Overall",
			roundTripOverall.count,
			roundTrip.GetMedian(),
			roundTrip.min,
			writeAccessOverall.count,
			writeAccessOverall.GetMedian(),
			writeAccessOverall.min,
			readAccessOverall.count,
			readAccessOverall.GetMedian(),
			readAccessOverall.min,
		)
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
	err = part.Delete()
	if err != nil {
		panic(err)
	}
	// end finalize_dds

	// skipping clean up variables

	roundTrip.Delete()
	writeAccess.Delete()
	readAccess.Delete()
	roundTripOverall.Delete()
	writeAccessOverall.Delete()
	readAccessOverall.Delete()

	C.dds_sample_free(unsafe.Pointer(&pubData), &C.RoundTripModule_DataType_desc, C.DDS_FREE_CONTENTS)

}
