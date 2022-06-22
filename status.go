package napi

/*
#include <node/node_api.h>
*/
import "C"

type Status int

const (
	StatusOK                            Status = C.napi_ok
	StatusInvalidArg                    Status = C.napi_invalid_arg
	StatusObjectExpected                Status = C.napi_object_expected
	StatusStringExpected                Status = C.napi_string_expected
	StatusNameExpected                  Status = C.napi_name_expected
	StatusFunctionExpected              Status = C.napi_function_expected
	StatusNumberExpected                Status = C.napi_number_expected
	StatusBooleanExpected               Status = C.napi_boolean_expected
	StatusArrayExpected                 Status = C.napi_array_expected
	StatusGenericFailure                Status = C.napi_generic_failure
	StatusPendingException              Status = C.napi_pending_exception
	StatusCancelled                     Status = C.napi_cancelled
	StatusEscapeCalledTwice             Status = C.napi_escape_called_twice
	StatusHandleScopeMismatch           Status = C.napi_handle_scope_mismatch
	StatusCallbackScopeMismatch         Status = C.napi_callback_scope_mismatch
	StatusQueueFull                     Status = C.napi_queue_full
	StatusClosing                       Status = C.napi_closing
	StatusBigintExpected                Status = C.napi_bigint_expected
	StatusDateExpected                  Status = C.napi_date_expected
	StatusArraybufferExpected           Status = C.napi_arraybuffer_expected
	StatusDetachableArraybufferExpected Status = C.napi_detachable_arraybuffer_expected
	StatusWouldDeadlock                 Status = C.napi_would_deadlock
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "napi_ok"
	case StatusInvalidArg:
		return "napi_invalid_arg"
	case StatusObjectExpected:
		return "napi_object_expected"
	case StatusStringExpected:
		return "napi_string_expected"
	case StatusNameExpected:
		return "napi_name_expected"
	case StatusFunctionExpected:
		return "napi_function_expected"
	case StatusNumberExpected:
		return "napi_number_expected"
	case StatusBooleanExpected:
		return "napi_boolean_expected"
	case StatusArrayExpected:
		return "napi_array_expected"
	case StatusGenericFailure:
		return "napi_generic_failure"
	case StatusPendingException:
		return "napi_pending_exception"
	case StatusCancelled:
		return "napi_cancelled"
	case StatusEscapeCalledTwice:
		return "napi_escape_called_twice"
	case StatusHandleScopeMismatch:
		return "napi_handle_scope_mismatch"
	case StatusCallbackScopeMismatch:
		return "napi_callback_scope_mismatch"
	case StatusQueueFull:
		return "napi_queue_full"
	case StatusClosing:
		return "napi_closing"
	case StatusBigintExpected:
		return "napi_bigint_expected"
	case StatusDateExpected:
		return "napi_date_expected"
	case StatusArraybufferExpected:
		return "napi_arraybuffer_expected"
	case StatusDetachableArraybufferExpected:
		return "napi_detachable_arraybuffer_expected"
	case StatusWouldDeadlock:
		return "napi_would_deadlock"
	}

	return "napi_go_status_unknown"
}
