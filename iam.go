package main

import (
    "context"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/sts"
)

func WhoAmI() error {
    cfg, _ := config.LoadDefaultConfig(context.TODO())
    client := sts.NewFromConfig(cfg)
    resp, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
    if err != nil {
        return err
    }
    fmt.Printf("Account: %s\nUserID: %s\nARN: %s\n", *resp.Account, *resp.UserId, *resp.Arn)
    return nil
}
