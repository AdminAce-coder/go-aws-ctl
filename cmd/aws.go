package cmd

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/gogf/gf/v2/os/glog"
)

// LoadOptions 定义选项
type LoadOptions struct {
	Region     string // 区域
	ClientType string // 客户端类型
}

// GetClient 获取 AWS 客户端
func GetClient[T any](optFns ...func(*LoadOptions) error) T {
	// 创建上下文
	ctx := context.TODO()

	// 加载选项
	var options LoadOptions
	for _, fn := range optFns {
		if err := fn(&options); err != nil {
			glog.New().Error(ctx, "Failed to apply option:", err)
			var zero T // 返回零值
			return zero
		}
	}

	// 先从配置文件中加载 AWS 配置
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedCredentialsFiles(
			[]string{"../config/credentials", "data/credentials"},
		),
		config.WithSharedConfigFiles(
			[]string{"../config/config", "data/config"},
		),
		//config.WithCredentialsProvider(aws.AnonymousCredentials{}), // 禁用 IMDS
	)
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
		log.Fatalf("尝试从环境变量中获取 AWS 配置")
		cfg, err = GetAwsConfigFromEnv()
		if err != nil {
			log.Fatalf("Failed to load AWS configuration from environment variables: %v", err)
		}
	}

	// 设置区域（如果有）
	if options.Region != "" {
		cfg.Region = options.Region
	}

	// 客户端选择器
	clientMap := map[string]func(cfg aws.Config) interface{}{
		"lightsail": func(cfg aws.Config) interface{} {
			return lightsail.NewFromConfig(cfg)
		},
		"ec2": func(cfg aws.Config) interface{} {
			return ec2.NewFromConfig(cfg)
		},
		"backup": func(cfg aws.Config) interface{} { return backup.NewFromConfig(cfg) },
	}

	// 创建客户端
	if clientFunc, ok := clientMap[options.ClientType]; ok {
		client := clientFunc(cfg)

		// 断言为泛型类型
		if typedClient, ok := client.(T); ok {
			return typedClient
		}
		log.Fatalf("Failed to cast client to the expected type")
	}

	// 未知的客户端类型
	glog.New().Error(ctx, "Unknown client type:", options.ClientType)
	var zero T
	return zero
}

// LoadOptionsFunc 是 LoadOptions 的函数类型
type LoadOptionsFunc func(*LoadOptions) error

// WithRegion 设置区域
func WithRegion(region string) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.Region = region
		return nil
	}
}

// WithClientType 设置客户端类型
func WithClientType(clientType string) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.ClientType = clientType
		return nil
	}
}

// 获取一个默认区域的 Lightsail 客户端
func GetDefaultAwsLgClient() *lightsail.Client {
	// 明确指定返回类型为 *lightsail.Client
	lgClient := GetClient[*lightsail.Client](
		WithRegion("us-east-1"),
		WithClientType("lightsail"),
	)
	return lgClient
}

// 从环境变量中获取 AWS 配置
func GetAwsConfigFromEnv() (aws.Config, error) {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if accessKeyID == "" || secretAccessKey == "" {
		return aws.Config{}, errors.New("AWS_ACCESS_KEY_ID or AWS_SECRET_ACCESS_KEY is not set")
	}
	// 使用正确的凭证创建方法
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""))

	return aws.Config{
		Credentials: creds,
	}, nil
}
