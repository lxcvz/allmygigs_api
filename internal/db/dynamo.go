package db

import (
	"allmygigs/config"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBClient struct {
	Client *dynamodb.DynamoDB
}

func NewDynamoDBClient(cfg *config.Config) (*DynamoDBClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AwsRegion),
		Credentials: credentials.NewStaticCredentials(cfg.AwsDynamoId, cfg.AwsDynamoSecret, ""),
	})
	if err != nil {
		return nil, err
	}

	client := dynamodb.New(sess)

	log.Println("Conexão com o DynamoDB estabelecida com sucesso.")

	return &DynamoDBClient{Client: client}, nil
}
