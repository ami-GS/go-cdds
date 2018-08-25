package cdds

/*
#cgo CFLAGS: -I/usr/local/include/ddsc
#cgo LDFLAGS: -lddsc
#include "ddsc/dds.h"
*/
import "C"

type CddsErrorType uint16

const (
	Ok CddsErrorType = iota
	Error
	Unsupported
	BadParameter
	PreconditionNotMet
	OutOfResource
	NotEnabled
	ImmutablePolicy
	InconsistencyPolicy
	AlreadyDeleted
	TimeOut
	NoData
	IllegalOperation
	NotAllowedBySecurity
)

func (e CddsErrorType) Error() string {
	return []string{
		"Success",
		"Non specific error",
		"Feature unsupported",
		"Bad parameter value",
		"Precondition for operation not met",
		"Out of resources",
		"Configurable feature is not enabled",
		"Attempt is made to modify an immutable policy",
		"Policy is used with inconsistent values",
		"Attempt is made to delete something more than once",
		"Timeout",
		"Expected data is not provided",
		"Function is called when it should not be",
		"credentials are not enough to use the function",
	}[int(e)]
}

func ErrorCheck(err C.dds_entity_t, flags uint8, where string) {
	C.dds_err_check(err, C.uint(flags), C.CString(where))
}
