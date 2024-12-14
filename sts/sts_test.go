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

func TestGetUserNameByAccessKeyId(t *testing.T) {
	userName, err := staop.IamQuery.GetUserNameByAccessKeyId(ctx, "AKIA5FTZAFBPZP6OLGET")
	if err != nil {
		t.Error(err)
	}
	t.Log(userName)
}

func TestGetAccountByAccessKeyId(t *testing.T) {
	account, err := staop.IamQuery.GetAccountByAccessKeyId(ctx, "AKIA5FTZAFBPZP6OLGET")
	if err != nil {
		t.Error(err)
	}
	t.Log(account)
}
