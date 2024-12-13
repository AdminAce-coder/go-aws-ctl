package lightsail

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/golifez/go-aws-ctl/cmd"
	"testing"
)

var ld = cmd.LoadOptions{}

func TestCreateInstancesFromSnapshot(t *testing.T) {
	InstanceName()
}

func TestLightsailClient(t *testing.T) {
	//ctx := context.TODO()
	client := cmd.GetClient[*lightsail.Client](
		cmd.WithClientType("lightsail"),
		cmd.WithRegion("us-east-1"),
	)
	if client == nil {
		t.Fatal("Failed to create Lightsail client")
	}
	t.Log("Lightsail client created successfully")
}

//获取实例

func TestLgGetinfo(t *testing.T) {
	//lgclint := cmd.GetDefaultAwsLgClient()
	ctx := context.Background()
	lgctl := NewLgQuery()
	list, err := lgctl.GetInstanceList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range list {
		t.Log("实例名是：", i)
	}
}
