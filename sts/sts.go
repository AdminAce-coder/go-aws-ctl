package sts

import (
	"github.com/golifez/go-aws-ctl/service"
)

type StsOp struct {
	IamQuery service.IamQuery
}

func NewStsOp() *StsOp {
	return &StsOp{
		IamQuery: NewIamQueryCommand(),
	}
}
