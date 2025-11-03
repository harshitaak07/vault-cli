package main

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func GenerateDataKey(kmsKeyID string) ([]byte, []byte, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, nil, err
	}
	client := kms.NewFromConfig(cfg)
	out, err := client.GenerateDataKey(context.TODO(), &kms.GenerateDataKeyInput{
		KeyId:   aws.String(kmsKeyID),
		KeySpec: types.DataKeySpecAes256,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("generate data key: %w", err)
	}
	return out.Plaintext, out.CiphertextBlob, nil
}

func DecryptDataKey(encryptedKey []byte) ([]byte, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := kms.NewFromConfig(cfg)
	out, err := client.Decrypt(context.TODO(), &kms.DecryptInput{
		CiphertextBlob: encryptedKey,
	})
	if err != nil {
		return nil, fmt.Errorf("decrypt data key: %w", err)
	}
	return out.Plaintext, nil
}

func LocalKey() ([]byte, []byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, nil, err
	}
	return key, key, nil
}
