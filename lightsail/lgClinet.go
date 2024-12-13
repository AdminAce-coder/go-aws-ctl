package lightsail

import (
	"context"
	"github.com/golifez/go-aws-ctl/cmd"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/gogf/gf/v2/os/glog"
)

func lgClinet(ctx context.Context) *lightsail.Client {
	client := cmd.GetClient[*lightsail.Client](
		cmd.WithClientType("lightsail"),
		cmd.WithRegion("us-east-1"),
	)
	if client == nil {
		glog.New().Error(ctx, "Failed to create Lightsail client")
	}
	glog.New().Info(ctx, "Lightsail client created successfully")
	return client
}
