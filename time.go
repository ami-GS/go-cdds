package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"
import "time"

type Time C.dds_time_t
type Duration C.dds_duration_t

func SleepFor(n time.Duration) {
	C.dds_sleepfor(C.int64_t(int64(n)))
}

/*
func SleepUntil(n Time) {

}
*/

func DdsTime() Time {
	return Time(C.dds_time_t(C.dds_time()))
}
