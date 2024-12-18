package main

//func main() {
//	// 配置临时凭证
//	cfg, err := config.LoadDefaultConfig(context.TODO(),
//		config.WithCredentialsProvider(
//			aws.NewCredentialsCache(
//				stscreds.NewTemporaryCredentialsProvider("ASIA123EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", "IQoJb3JpZ2luX2VjEOr//////////wEaCWV1LXdlc3QtMSJHMEUCIQDPFJFlxQpN..."),
//			),
//		),
//	)
//	if err != nil {
//		log.Fatalf("无法加载配置: %v", err)
//	}
//
//	// 使用临时凭证调用 AWS 服务
//	client := s3.NewFromConfig(cfg)
//
//	// 示例：列出 S3 存储桶
//	result, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
//	if err != nil {
//		log.Fatalf("无法列出存储桶: %v", err)
//	}
//
//	fmt.Println("S3 Buckets:")
//	for _, bucket := range result.Buckets {
//		fmt.Printf("- %s\n", aws.ToString(bucket.Name))
//	}
//}
