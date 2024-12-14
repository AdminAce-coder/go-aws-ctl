package lightsail

import (
	"github.com/golifez/go-aws-ctl/service"
)

type LgOp struct {
	LgOpsvc service.LgOpsvc
	lgQuery service.LgQuerysvc
}

func NewLgOp() *LgOp {
	return &LgOp{
		LgOpsvc: NewLgInstanceOpCommand(), // 实例操作
		lgQuery: NewLgQuery(),             // 查询
	}
}
