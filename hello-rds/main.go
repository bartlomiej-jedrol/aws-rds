package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/bartlomiej-jedrol/aws-rds/amzn/rdspsql"
)

var (
	ctx       context.Context
	rdsClient *rds.Client
)

func init() {
	ctx = context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load sdk config")
		return
	}

	rdsClient = rds.NewFromConfig(cfg)
}

func prepareInput() rdspsql.Instance {
	opts := rdspsql.NewOpts()
	creds := rdspsql.NewCreds()
	return rdspsql.Instance{
		Opts:  opts,
		Creds: creds,
	}
}

func main() {
	cfg, err := rdspsql.ParseConfig("cfg/instance_config.yaml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("instance cfg: %+v", cfg)
	// ic := rdspsql.NewInstanceController(rdsClient)
	// input := prepareInput()
	// ic.CreateInstance(ctx, input)
	// ic.DescribeDBInstances(ctx)
	// createDBInstance()
	// deleteDBInstance()
	// describeDBInstances()
}
