package lightsail

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/golifez/go-aws-ctl/cmd"
	ctltypes "github.com/golifez/go-aws-ctl/model"
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
	//lgclint := cmd.GetDefaultAwsLgClient()q
	ctx := context.Background()
	lgctl := NewLgQuery()
	list, err := lgctl.GetInstanceList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range list {
		t.Log("实例信息是：", i)
	}
}

// 删除实例
func TestLgDeleteInstance(t *testing.T) {
	lgctl := NewLgOp()
	err := lgctl.LgOpsvc.DeleteInstance("WordPress-1", "ap-southeast-1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("删除成功")
}

// 获取快照
func TestLgGetSnapshot(t *testing.T) {
	lgctl := NewLgQuery()
	snapshotList, err := lgctl.GetSnapshotList(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("快照列表：%v", snapshotList)
}

// 查看捆绑包
func TestLgGetBundles(t *testing.T) {
	lgctl := NewLgQuery()
	bundles, err := lgctl.GetBundlesInput(context.Background(), cmd.GetDefaultAwsLgClient())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("捆绑包列表：%v", bundles)
}

// 创建实例
func TestLgCreateInstance(t *testing.T) {
	lgctl := NewLgOp()
	err := lgctl.LgOpsvc.CreateInstance(ctltypes.LgCreateInstance{
		InstanceName:     "test-instance",
		Region:           "ap-northeast-1",
		SnapshotName:     "Ubuntu-1-1734237832",
		AvailabilityZone: "ap-northeast-1a",
		BundleId:         "nano_3_0",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("创建实例成功")
}
