package awsom

import (
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// Service

func SecretsManagerService() (*secretsmanager.SecretsManager, error) {
	sess, err := CreateSession()
	if err != nil {
		return nil, err
	}
	return secretsmanager.New(sess), err
}
