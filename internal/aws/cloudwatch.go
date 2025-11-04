package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func LogToCloudWatch(message string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	svc := cloudwatchlogs.NewFromConfig(cfg)
	logGroup := "/vault/logs"
	logStream := "vault-cli"

	_, _ = svc.CreateLogStream(context.TODO(), &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  &logGroup,
		LogStreamName: &logStream,
	})

	ts := time.Now().UnixMilli()
	_, err = svc.PutLogEvents(context.TODO(), &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  &logGroup,
		LogStreamName: &logStream,
		LogEvents: []types.InputLogEvent{
			{
				Message:   aws.String(message),
				Timestamp: aws.Int64(ts),
			},
		},
	})
	return err
}
