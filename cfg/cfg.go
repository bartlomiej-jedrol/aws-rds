package cfg

import (
	"context"
	"log"
	"os"

	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/bartlomiej-jedrol/aws-rds/amzn/secrets"
	"gopkg.in/yaml.v2"
)

var (
	cfgPath    string = "cfg/instance_config.yaml"
	secretName string = "rds/psql"
)

type ConfigCreds struct {
	SecretName string `yaml:"secret_name"`
}

type ConfigInstannce struct {
	Opts  Opts        `yaml:"opts"`
	Creds ConfigCreds `yaml:"creds"`
}

type Config struct {
	Instance ConfigInstannce `yaml:"instance"`
}

type Opts struct {
	AllocatedStorage            int32
	AutoMinorVersionUpgrade     bool
	BackupRetentionPeriod       int32
	DBInstanceClass             string
	DBInstanceIdentifier        string
	DBName                      string
	DBParameterGroup            string
	DBSubnetGroupName           string
	DeletionProtection          bool
	EnableCloudwatchLogsExports []string
	Engine                      string
	EngineVersion               string
	LicenseModel                string
	PubliclyAccessible          bool
	SkipFinalSnapshot           bool
}

type Creds struct {
	MasterUserName     string
	MasterUserPassword string
}

type Instance struct {
	Opts
	Creds
}

// ParseConfig
func parseConfig(filePath string) (*ConfigInstannce, error) {
	cfg := Config{}

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("failed to read instance config: %v", err)
		return nil, err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Printf("failed to unmarshal instance config: %v", err)
		return nil, err
	}

	return &cfg.Instance, nil
}

func NewInstanceInputFromCfg(ctx context.Context, secretsClient *sm.Client) (*Instance, error) {
	user, pass, err := secrets.GetSecrets(ctx, secretsClient, secretName)
	if err != nil {
		return nil, err
	}

	cfg, err := parseConfig(cfgPath)
	if err != nil {
		return nil, err
	}

	opts := Opts{
		AllocatedStorage:            cfg.Opts.AllocatedStorage,
		AutoMinorVersionUpgrade:     cfg.Opts.AutoMinorVersionUpgrade,
		BackupRetentionPeriod:       cfg.Opts.BackupRetentionPeriod,
		DBInstanceClass:             cfg.Opts.DBInstanceClass,
		DBInstanceIdentifier:        cfg.Opts.DBInstanceIdentifier,
		DBName:                      cfg.Opts.DBName,
		DBParameterGroup:            cfg.Opts.DBParameterGroup,
		DBSubnetGroupName:           cfg.Opts.DBSubnetGroupName,
		DeletionProtection:          cfg.Opts.DeletionProtection,
		EnableCloudwatchLogsExports: cfg.Opts.EnableCloudwatchLogsExports,
		Engine:                      cfg.Opts.Engine,
		EngineVersion:               cfg.Opts.EngineVersion,
		LicenseModel:                cfg.Opts.LicenseModel,
		PubliclyAccessible:          cfg.Opts.PubliclyAccessible,
		SkipFinalSnapshot:           cfg.Opts.SkipFinalSnapshot,
	}

	creds := Creds{
		MasterUserName:     user,
		MasterUserPassword: pass,
	}

	return &Instance{
		Opts:  opts,
		Creds: creds,
	}, nil
}
