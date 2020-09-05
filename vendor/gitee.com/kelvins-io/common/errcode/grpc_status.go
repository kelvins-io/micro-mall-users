package errcode

import (
	"gitee.com/kelvins-io/common/proto/common"
	sts "google.golang.org/grpc/status"
)

type Status struct {
	*sts.Status
}

func ToRPCStatus(code int, msgs ...string) *Status {
	msg := assemblyMessage(code, msgs...)
	s, _ := sts.New(ToRPCCode(int32(code)), msg).WithDetails(&common.Error{Code: int32(code), Message: msg})
	return &Status{s}
}

func FromError(err error) *Status {
	s, _ := sts.FromError(err)
	return &Status{s}
}

func (s *Status) CommonError() *common.Error {
	details := s.Details()
	for _, detail := range details {
		if v, ok := detail.(*common.Error); ok {
			return v
		}
	}

	return nil
}
