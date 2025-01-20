package config

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var AWSClient *s3.Client

// InitializeAWSSession initializes the AWS session
func InitializeAWSSession() error {
	ctx := context.TODO()

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(
		ctx,
		// Add more custom config values if needed.
		// config.WithRegion("us-west-2"),
		// config.WithSharedConfigProfile("customProfile"),
	)
	if err != nil {
		log.Fatal(err)
	}

	AWSClient := sts.NewFromConfig(cfg)

	identity, err := AWSClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Account: %s, Arn: %s", aws.ToString(identity.Account), aws.ToString(identity.Arn))

	log.Println("AWS session initialized successfully")
	return nil
}
