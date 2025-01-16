package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
)

type secretString struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	Engine               string `json:"engine"`
	Host                 string `json:"host"`
	Port                 int    `json:"port"`
	DBName               string `json:"dbname"`
	DBInstanceIdentifier string `json:"dbInstanceIdentifier"`
}

func GetSecrets(ctx context.Context, secretsClient *sm.Client, secretName string) (string, string, error) {
	in := sm.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	out, err := secretsClient.GetSecretValue(ctx, &in)
	if err != nil {
		log.Printf("failed to get secret: %v", err)
		return "", "", fmt.Errorf("failed to get secret: %v", err)
	}

	ss := secretString{}
	err = json.Unmarshal([]byte(*out.SecretString), &ss)
	if err != nil {
		log.Printf("failed to unmarshal secret string: %v", err)
		return "", "", fmt.Errorf("failed to unmarshal secret string: %v", err)
	}

	return ss.Username, ss.Password, nil
}
