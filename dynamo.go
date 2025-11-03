package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func DynamoClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg), nil
}

func RecordFileToDynamo(fileName, hash string, size int64, mode, location string) error {
	client, err := DynamoClient()
	if err != nil {
		return err
	}
	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("VaultMetadata"),
		Item: map[string]dynamodb.AttributeValue{
			"FileName":   &dynamodb.AttributeValueMemberS{Value: fileName},
			"UploadedAt": &dynamodb.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
			"Hash":       &dynamodb.AttributeValueMemberS{Value: hash},
			"Size":       &dynamodb.AttributeValueMemberN{Value: fmt.Sprintf("%d", size)},
			"Mode":       &dynamodb.AttributeValueMemberS{Value: mode},
			"Location":   &dynamodb.AttributeValueMemberS{Value: location},
		},
	})
	return err
}

func RecordAuditToDynamo(action, filename, target string, success bool, errMsg string) error {
	client, err := DynamoClient()
	if err != nil {
		return err
	}
	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("VaultAudit"),
		Item: map[string]dynamodb.AttributeValue{
			"ActionID": &dynamodb.AttributeValueMemberS{Value: fmt.Sprintf("%s-%d", action, time.Now().UnixNano())},
			"Action":   &dynamodb.AttributeValueMemberS{Value: action},
			"FileName": &dynamodb.AttributeValueMemberS{Value: filename},
			"Target":   &dynamodb.AttributeValueMemberS{Value: target},
			"Success":  &dynamodb.AttributeValueMemberBOOL{Value: success},
			"Error":    &dynamodb.AttributeValueMemberS{Value: errMsg},
			"TS":       &dynamodb.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
		},
	})
	return err
}
