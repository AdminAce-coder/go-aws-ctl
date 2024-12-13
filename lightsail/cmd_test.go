package lightsail

import (
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
