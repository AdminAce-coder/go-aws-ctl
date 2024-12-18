package main

//
//func main() {
//	// 加载 SSO 登录后的临时凭证配置
//	cfg, err := config.LoadDefaultConfig(context.TODO(),
//		config.WithSharedConfigProfile("sso-user-profile"), // 使用 SSO 配置的 profile
//	)
//	if err != nil {
//		log.Fatalf("无法加载 AWS 配置文件: %v", err)
//	}
//
//	// 初始化 S3 客户端
//	s3Client := s3.NewFromConfig(cfg)
//
//	// 调用 S3 API：列出存储桶
//	result, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
//	if err != nil {
//		log.Fatalf("无法列出 S3 存储桶: %v", err)
//	}
//
//	// 输出存储桶名称
//	fmt.Println("S3 Buckets:")
//	for _, bucket := range result.Buckets {
//		fmt.Println("- ", aws.ToString(bucket.Name))
//	}
//}
