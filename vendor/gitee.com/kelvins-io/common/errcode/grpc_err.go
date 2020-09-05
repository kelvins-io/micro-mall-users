// 通用错误码组件-grpc错误返回定义
package errcode

import (
	"gitee.com/kelvins-io/common/proto/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

func TogRPCError(code int, msgs ...string) error {
	msg := assemblyMessage(code, msgs...)
	s, _ := status.New(ToRPCCode(int32(code)), msg).WithDetails(&common.Error{Code: int32(code), Message: msg})
	return s.Err()
}

func ToRPCCode(code int32) codes.Code {
	var statusCode codes.Code
	switch code {
	case FAIL:
		statusCode = codes.Internal
	case INVALID_PARAMS:
		statusCode = codes.InvalidArgument
	case UNAUTH:
		statusCode = codes.Unauthenticated
	case ACCESS_DENIED:
		statusCode = codes.PermissionDenied
	case DEADLINE_EXCEEDED:
		statusCode = codes.DeadlineExceeded
	case NOT_FOUND:
		statusCode = codes.NotFound
	case LIMIT_EXCEED:
		statusCode = codes.ResourceExhausted
	case METHOD_NOT_ALLOWED:
		statusCode = codes.Unimplemented
	default:
		statusCode = codes.Unknown
	}

	return statusCode
}

func ToHttpStatusCode(code codes.Code) int {
	var statusCode int
	switch code {
	case codes.Unknown:
		statusCode = http.StatusInternalServerError
	case codes.Internal:
		statusCode = http.StatusInternalServerError
	case codes.InvalidArgument:
		statusCode = http.StatusBadRequest
	case codes.Unauthenticated:
		statusCode = http.StatusUnauthorized
	case codes.PermissionDenied:
		statusCode = http.StatusUnauthorized
	case codes.DeadlineExceeded:
		statusCode = http.StatusRequestTimeout
	case codes.NotFound:
		statusCode = http.StatusNotFound
	case codes.ResourceExhausted:
		statusCode = http.StatusTooManyRequests
	case codes.Unimplemented:
		statusCode = http.StatusMethodNotAllowed
	default:
		statusCode = http.StatusOK
	}

	return statusCode
}

func assemblyMessage(code int, msgs ...string) string {
	var msg string
	if len(msgs) == 0 {
		msg = GetErrMsg(code)
	} else {
		msg = strings.Join(msgs, ",")
	}

	return msg
}
