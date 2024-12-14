package lightsail

import (
	"github.com/golifez/go-aws-ctl/service"
)

type LgOp struct {
	LgOpsvc service.LgOpsvc
}

func NewLgOp() *LgOp {
	return &LgOp{
		LgOpsvc: NewLgInstanceOpCommand(), // 实例操作
	}
}
