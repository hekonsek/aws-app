package awsom

import "github.com/aws/aws-sdk-go/aws/session"

func CreateSession() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	return sess, nil
}
