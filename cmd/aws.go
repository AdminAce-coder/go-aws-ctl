package cmd

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/gogf/gf/v2/os/glog"
)

// LoadOptions 定义选项
type LoadOptions struct {
	Region     string // 区域
	ClientType string // 客户端类型
	cfg        *aws.Config
}

// GetClient 获取 AWS 客户端
func GetClient[T any](optFns ...func(*LoadOptions) error) T {
	// 创建上下文
	ctx := context.TODO()
	var cfg aws.Config
	var err error // 声明 err 变量

	// 加载选项
	var options LoadOptions
	for _, fn := range optFns {
		if err := fn(&options); err != nil {
			glog.New().Error(ctx, "Failed to apply option:", err)
			var zero T
			return zero
		}
	}

	// 先判断CFG是否存在
	if options.cfg != nil {
		cfg = *options.cfg
	} else {
		cfg, err = GetAwsConfigFromConfigFile()
		if err != nil {
			glog.New().Warning(ctx, "Failed to load AWS configuration from config file:", err)
			// 尝试从环境变量加载
			cfg, err = GetAwsConfigFromEnv()
			if err != nil {
				log.Fatalf("Failed to load AWS configuration: %v", err)
			}
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
		"sts":    func(cfg aws.Config) interface{} { return sts.NewFromConfig(cfg) },
		"iam":    func(cfg aws.Config) interface{} { return iam.NewFromConfig(cfg) },
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

// WithConfig 设置配置
func WithConfig(cfg *aws.Config) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.cfg = cfg
		return nil
	}
}

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

// 获取一个带区域的 Lightsail 客户端
func GetAwsLgClient(region string) *lightsail.Client {
	lgClient := GetClient[*lightsail.Client](
		WithRegion(region),
		WithClientType("lightsail"),
	)
	return lgClient
}

// 获取一个默认区域的 EC2 客户端
func GetDefaultAwsEc2Client() *ec2.Client {
	// 明确指定返回类型为 *ec2.Client
	ec2Client := GetClient[*ec2.Client](
		WithRegion("us-east-1"),
		WithClientType("ec2"),
	)
	return ec2Client
}

// 获取一个带区域的 EC2 客户端
func GetAwsEc2ClientWithRegion(region string) *ec2.Client {
	ec2Client := GetClient[*ec2.Client](
		WithRegion(region),
		WithClientType("ec2"),
	)
	return ec2Client
}

// 返回一个STS客户端
func GetDefaultAwsStsClient() *sts.Client {
	// 明确指定返回类型为 *sts.Client
	stsClient := GetClient[*sts.Client](
		WithRegion("us-east-1"),
		WithClientType("sts"),
	)
	return stsClient
}

// 添加获取 IAM 客户端的函数
func GetDefaultAwsIamClient() *iam.Client {
	iamClient := GetClient[*iam.Client](
		WithRegion("us-east-1"),
		WithClientType("iam"),
	)
	return iamClient
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

// 从.aws/config文件中获取 AWS 配置
func GetAwsConfigFromConfigFile() (aws.Config, error) {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return aws.Config{}, err
	}

	// 构建标准的 AWS 配置文件路径
	credentialsPath := filepath.Join(homeDir, ".aws", "credentials")
	configPath := filepath.Join(homeDir, ".aws", "config")

	// 加载 AWS 配置
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedCredentialsFiles([]string{credentialsPath}),
		config.WithSharedConfigFiles([]string{configPath}),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}

// 直接传入id和key,获取 AWS 配置
func GetAwsConfigFromIdAndKey(accessKeyId, secretAccessKey string) (aws.Config, error) {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, ""))
	return aws.Config{
		Credentials: creds,
	}, nil
}

// GetAwsConfigFromIdAndKeyWithRegion 获取一个带区域的 AWS 配置
func GetAwsConfigFromIdAndKeyWithRegion(accessKeyId, secretAccessKey, region string) (aws.Config, error) {
	cfg, err := GetAwsConfigFromIdAndKey(accessKeyId, secretAccessKey)
	if err != nil {
		return aws.Config{}, err
	}
	cfg.Region = region
	return cfg, nil
}
