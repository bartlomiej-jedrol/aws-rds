package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/bartlomiej-jedrol/aws-rds/amzn/rdspsql"
	"github.com/bartlomiej-jedrol/aws-rds/cfg"
)

var (
	ctx           context.Context
	rdsClient     *rds.Client
	secretsClient *sm.Client
)

func init() {
	ctx = context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load sdk config: %v", err)
		return
	}

	rdsClient = rds.NewFromConfig(cfg)
	secretsClient = sm.NewFromConfig(cfg)
}

func main() {
	instance, err := cfg.NewInstanceFromCfg(ctx, secretsClient)
	if err != nil {
		return
	}

	ic := rdspsql.NewInstanceController(rdsClient)
	err = ic.CreateDBInstance(ctx, instance)
	if err != nil {
		return
	}

	// ic.DescribeDBInstances(ctx)
	// createDBInstance()
	// deleteDBInstance()
	// describeDBInstances()
}
