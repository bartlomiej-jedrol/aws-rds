package rdspsql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/bartlomiej-jedrol/aws-rds/cfg"
)

type InstanceController struct {
	Client *rds.Client
}

// NewInstanceController
func NewInstanceController(rdsClient *rds.Client) *InstanceController {
	return &InstanceController{
		Client: rdsClient,
	}
}

// CreateInstance
func (ic *InstanceController) CreateDBInstance(ctx context.Context, in *cfg.Instance) {
	startTime := time.Now()

	input := rds.CreateDBInstanceInput{
		AllocatedStorage:            aws.Int32(in.Opts.AllocatedStorage),
		AutoMinorVersionUpgrade:     aws.Bool(in.Opts.AutoMinorVersionUpgrade),
		BackupRetentionPeriod:       aws.Int32(in.Opts.BackupRetentionPeriod),
		DBInstanceClass:             aws.String(in.Opts.DBInstanceClass),
		DBInstanceIdentifier:        aws.String(in.Opts.DBInstanceIdentifier),
		DBName:                      aws.String(in.Opts.DBName),
		DBParameterGroupName:        aws.String(in.Opts.DBParameterGroup),
		DBSubnetGroupName:           aws.String(in.Opts.DBSubnetGroupName),
		DeletionProtection:          aws.Bool(in.Opts.DeletionProtection),
		EnableCloudwatchLogsExports: in.Opts.EnableCloudwatchLogsExports,
		Engine:                      aws.String(in.Opts.Engine),
		EngineVersion:               aws.String(in.Opts.EngineVersion),
		LicenseModel:                aws.String(in.Opts.LicenseModel),
		MasterUsername:              aws.String(in.Creds.MasterUserName),
		MasterUserPassword:          aws.String(in.Creds.MasterUserPassword),
		PubliclyAccessible:          aws.Bool(in.Opts.PubliclyAccessible),
	}

	_, err := ic.Client.CreateDBInstance(ctx, &input)
	if err != nil {
		log.Printf("failed to create db instance: %+v", err)
		return
	}
	log.Printf("started creation of instance: %s", in.Opts.DBInstanceIdentifier)

	waiter := rds.NewDBInstanceAvailableWaiter(ic.Client)
	dscInput := rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(in.Opts.DBInstanceIdentifier),
	}
	maxWaitTime := 20 * time.Minute
	output, err := waiter.WaitForOutput(ctx, &dscInput, maxWaitTime)
	if err != nil {
		log.Printf("error waiting for instance creation: %v", err)
		return
	}
	duration := time.Since(startTime)
	log.Printf("instance: %s successfully created in: %v", in.Opts.DBInstanceIdentifier, duration)
	log.Printf("created instance arn: %s", *output.DBInstances[0].DBInstanceArn)
}

// DescribeDBInstances
func (ic *InstanceController) DescribeDBInstances(ctx context.Context) {
	filter := types.Filter{
		Name:   aws.String("db-instance-id"),
		Values: []string{"db-instance-1"}}

	input := rds.DescribeDBInstancesInput{
		// DBInstanceIdentifier: aws.String("database-1"),
		MaxRecords: aws.Int32(20),
		Filters:    []types.Filter{filter},
	}

	rdsInstances, err := ic.Client.DescribeDBInstances(ctx, &input)
	if err != nil {
		log.Printf("failed to describe rds instances: %v", err)
	}

	// Print instances in a more readable format
	fmt.Println("\nRDS Instances:")
	fmt.Println("==============")

	for _, instance := range rdsInstances.DBInstances {
		fmt.Printf("\nInstance: %s\n", *instance.DBInstanceIdentifier)
		if instance.DBName != nil {
			fmt.Printf("  DB Name: %s\n", *instance.DBName)
		}
		fmt.Printf("  Engine: %s %s\n", *instance.Engine, *instance.EngineVersion)
		fmt.Printf("  Allocated Storage: %d GB\n", *instance.AllocatedStorage)
		fmt.Printf("  Auto Minor Version Upgrade: %v\n", *instance.AutoMinorVersionUpgrade)
		fmt.Printf("  Availability Zone: %s\n", *instance.AvailabilityZone)
		fmt.Printf("  Backup Retention Period: %d\n", *instance.BackupRetentionPeriod)
		fmt.Printf("  DB Instance ARN: %s\n", *instance.DBInstanceArn)
		fmt.Printf("  DB Instance Class: %s\n", *instance.DBInstanceClass)
		fmt.Printf("  DB Instance Status: %s\n", *instance.DBInstanceStatus)
		if instance.Endpoint != nil {
			fmt.Printf("  Endpoint: %s:%d\n", *instance.Endpoint.Address, *instance.Endpoint.Port)
		}
		if instance.MasterUsername != nil {
			fmt.Printf("  Master Username: %s\n", *instance.MasterUsername)
		}
		fmt.Printf("  Multi-AZ: %v\n", *instance.MultiAZ)
		fmt.Printf("\n  DB Parameter Group Name: %v\n", *instance.DBParameterGroups[0].DBParameterGroupName)
		fmt.Printf("  DB Subnet Group Name: %v\n", *instance.DBSubnetGroup.DBSubnetGroupName)
		fmt.Printf("\n  Instance Create Time: %v\n", *instance.InstanceCreateTime)
	}
}

// DeleteDBInstance
func (ic *InstanceController) DeleteDBInstance(ctx context.Context, in *cfg.Instance) {
	startTime := time.Now()

	input := rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: aws.String(in.Opts.DBInstanceIdentifier),
		SkipFinalSnapshot:    aws.Bool(in.Opts.SkipFinalSnapshot),
	}
	_, err := ic.Client.DeleteDBInstance(ctx, &input)
	if err != nil {
		log.Printf("failed to delete db instance: %v", err)
		return
	}
	log.Printf("started deletion of instance: %s", in.Opts.DBInstanceIdentifier)

	waiter := rds.NewDBInstanceDeletedWaiter(ic.Client)
	dscInput := rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(in.Opts.DBInstanceIdentifier),
	}

	maxWaitTime := 20 * time.Minute
	_, err = waiter.WaitForOutput(ctx, &dscInput, maxWaitTime)
	if err != nil {
		log.Printf("error waiting for instance deletion: %v", err)
		return
	}
	duration := time.Since(startTime)
	log.Printf("instance: %s has been successfully deleted in: %v", in.Opts.DBInstanceIdentifier, duration)
}
