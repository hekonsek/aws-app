package awsom

import (
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	awsom_session "github.com/hekonsek/awsom-session"
)

// Service

func SecretsManagerService() (*secretsmanager.SecretsManager, error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return nil, err
	}
	return secretsmanager.New(sess), err
}
