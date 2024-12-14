package sts

import (
	"context"
	"testing"
)

var (
	staop = NewStsOp()
	ctx   = context.Background()
)

func TestNewIamQueryCommand(t *testing.T) {
	IamInfo, err := staop.IamQuery.GetIamInfo(ctx)
	if err != nil {
		t.Error(err)
	}
	t.Log(IamInfo)
}
