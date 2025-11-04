package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
		Item: map[string]types.AttributeValue{
			"FileName":   &types.AttributeValueMemberS{Value: fileName},
			"UploadedAt": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
			"Hash":       &types.AttributeValueMemberS{Value: hash},
			"Size":       &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", size)},
			"Mode":       &types.AttributeValueMemberS{Value: mode},
			"Location":   &types.AttributeValueMemberS{Value: location},
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
		Item: map[string]types.AttributeValue{
			"ActionID": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s-%d", action, time.Now().UnixNano()),
			},
			"Action":   &types.AttributeValueMemberS{Value: action},
			"FileName": &types.AttributeValueMemberS{Value: filename},
			"Target":   &types.AttributeValueMemberS{Value: target},
			"Success":  &types.AttributeValueMemberBOOL{Value: success},
			"Error":    &types.AttributeValueMemberS{Value: errMsg},
			"TS":       &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
		},
	})
	return err
}

