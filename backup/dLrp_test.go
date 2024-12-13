package backup

import (
	"context"
	"goAwsCtrl/cmd"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/backup"
)

func TestDlRP(t *testing.T) {
	client := cmd.GetClient[*backup.Client](
		cmd.WithClientType("backup"),
		cmd.WithRegion("us-east-1"),
	)
	ctx := context.Background()
	deleteRP(ctx, client)
}
